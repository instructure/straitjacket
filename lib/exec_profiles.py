#!/usr/bin/env python
#
# Copyright (C) 2011 Instructure, Inc.
#
# This file is part of StraitJacket.
#
# StraitJacket is free software: you can redistribute it and/or modify it under
# the terms of the GNU Affero General Public License as published by the Free
# Software Foundation, version 3 of the License.
#
# StraitJacket is distributed in the hope that it will be useful, but WITHOUT ANY
# WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
# A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
# details.
#
# You should have received a copy of the GNU Affero General Public License along
# with this program. If not, see <http://www.gnu.org/licenses/>.
#

import time, os, subprocess, threading, base64, shutil
try:
  import LibAppArmor
except ImportError:
  LibAppArmor = None

__author__ = "JT Olds"
__copyright__ = "Copyright 2011 Instructure, Inc."
__license__ = "AGPLv3"
__email__ = "jt@instructure.com"

class Error_(Exception): pass
class AppArmorProtectionFailure(Error_): pass

def aa_change_onexec(profile):
  if LibAppArmor is None or LibAppArmor.aa_change_onexec(profile) != 0:
    raise AppArmorProtectionFailure, ("failed to switch to apparmor profile %s"
        % profile)

class BaseProfile(object):

  def __init__(self, config):
    self.config = config
    self.max_runtime = config.getint("general", "max_runtime")

  def _kill(self, pid, completed, max_runtime=None):
    if max_runtime is None: max_runtime = self.max_runtime
    for _ in xrange(max(int(max_runtime), 1)):
      if completed: return
      time.sleep(1)
    if not completed:
      os.kill(pid, 9)
      completed.append("killed")

  def _run_user_program(self, user_program, stdin, aa_profile, time_used=0,
      executable=None, chdir=None, custom_timelimit=None):
    if custom_timelimit == None: custom_timelimit = float('inf')
    completed = []
    runtime = None
    start_time = time.time()

    def preexec_fn():
      if chdir: os.chdir(chdir)
      aa_change_onexec(aa_profile)

    proc = subprocess.Popen(user_program, executable=executable,
        stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
        close_fds=True, preexec_fn=preexec_fn)
    kill_thread = threading.Thread(target=self._kill, args=(proc.pid, completed,
        min(self.max_runtime - time_used, custom_timelimit)))

    kill_thread.start()
    returncode = None
    try:
      stdout, stderr = proc.communicate(stdin)
    except Exception, e:
      stdout, stderr = "", str(e)
      returncode = -9
    runtime = time.time() - start_time
    completed.append(True)
    kill_thread.join()

    if returncode is None: returncode = proc.returncode

    if "killed" in completed or runtime > custom_timelimit:
      error = "runtime_timelimit"
    elif returncode != 0:
      error = "runtime_error"
    else:
      error = ""

    return stdout, stderr, returncode, runtime, error

  def _filename_gen(self): return base64.b64encode(os.urandom(42), "_-")

  def run(self, lang_conf, source, stdin, custom_timelimit=None):
    raise NotImplementedError

  def apparmor_profile(self, lang_conf, profile_name="apparmor_profile"):
    profile = lang_conf.get(profile_name, "").strip()
    if len(profile) > 0: return profile
    return self.default_apparmor_profile()

  def default_apparmor_profile(self):
    raise NotImplementedError


class CompilerProfile(BaseProfile):

  def __init__(self, config): BaseProfile.__init__(self, config)

  def run(self, lang_conf, source, stdin, custom_timelimit=None):
    source_dir = os.path.join(self.config.get("directories", "source"),
        self._filename_gen())
    source_file = os.path.join(source_dir, lang_conf["filename"])
    compiler_file = os.path.join(self.config.get("directories", "compiler"),
        self._filename_gen())
    executable_file = os.path.join(self.config.get("directories", "execution"),
        self._filename_gen())
    try:
      os.mkdir(source_dir)
      f = file(source_file, "w")
      try:
        f.write(source)
      finally:
        f.close()

      completed = []
      compile_start_time = time.time()

      def compiler_preexec():
        os.environ["TMPDIR"] = self.config.get("directories", "compiler")
        aa_change_onexec(self.apparmor_profile(lang_conf))

      if lang_conf.has_key("compilation_command"):
        command = eval(lang_conf["compilation_command"])(source_file,
            compiler_file)
      else:
        command = [lang_conf["binary"], "-o", compiler_file, source_file]
      proc = subprocess.Popen(command, stdin=None, stdout=subprocess.PIPE,
          stderr=subprocess.STDOUT, close_fds=True, preexec_fn=compiler_preexec)
      kill_thread = threading.Thread(target=self._kill, args=(proc.pid,
          completed))

      kill_thread.start()
      returncode = None
      try:
        compile_out = proc.communicate()[0]
      except Exception, e:
        compile_out = str(e)
        returncode = -9
      completed.append(True)
      kill_thread.join()

      if returncode is None: returncode = proc.returncode

      if returncode != 0:
        if "killed" in completed:
          error = "compilation_timelimit"
        else:
          error = "compilation_error"
        return "", compile_out, returncode, 0.0, error

      os.rename(compiler_file, executable_file)

      if lang_conf.has_key("compiled_apparmor_profile"):
        compiled_profile = lang_conf["compiled_apparmor_profile"]
      else:
        compiled_profile = self.config.get("default-apparmor-profiles",
            "compiled")

      return self._run_user_program(["straitjacket-binary"], stdin,
          compiled_profile, time.time() - compile_start_time, executable_file,
          custom_timelimit=custom_timelimit)
    finally:
      shutil.rmtree(source_dir)
      if os.path.exists(compiler_file): os.unlink(compiler_file)
      if os.path.exists(executable_file): os.unlink(executable_file)

  def default_apparmor_profile(self):
    return self.config.get("default-apparmor-profiles", "compiler")

class InterpreterProfile(BaseProfile):

  def __init__(self, config): BaseProfile.__init__(self, config)

  def run(self, lang_conf, source, stdin, custom_timelimit=None):
    dirname = os.path.join(self.config.get("directories", "source"),
        self._filename_gen())
    filename = os.path.join(dirname, lang_conf["filename"])
    try:
      os.mkdir(dirname)
      f = file(filename, "w")
      try:
        f.write(source)
      finally:
        f.close()
      if lang_conf.has_key("interpretation_command"):
        command = eval(lang_conf["interpretation_command"])(filename)
      else:
        command = [lang_conf["binary"], filename]
      return self._run_user_program(command, stdin,
          self.apparmor_profile(lang_conf), custom_timelimit=custom_timelimit)

    finally:
      shutil.rmtree(dirname)

  def default_apparmor_profile(self):
    return self.config.get("default-apparmor-profiles", "interpreter")


class VMProfile(BaseProfile):

  def __init__(self, config): BaseProfile.__init__(self, config)

  def run(self, lang_conf, source, stdin, custom_timelimit=None):
    source_dir = os.path.join(self.config.get("directories", "source"),
        self._filename_gen())
    source_file = os.path.join(source_dir, lang_conf["filename"])
    try:
      os.mkdir(source_dir)
      f = file(source_file, "w")
      try:
        f.write(source)
      finally:
        f.close()

      completed = []
      compile_start_time = time.time()

      def compiler_preexec():
        os.environ["TMPDIR"] = self.config.get("directories", "compiler")
        os.chdir(source_dir)
        aa_change_onexec(self.apparmor_profile(lang_conf,
            "compiler_apparmor_profile"))

      if lang_conf.has_key("compilation_command"):
        command = eval(lang_conf["compilation_command"])(source_file)
      else:
        command = [lang_conf["binary"], source_file]
      proc = subprocess.Popen(command, stdin=None, stdout=subprocess.PIPE,
          stderr=subprocess.STDOUT, close_fds=True, preexec_fn=compiler_preexec)
      kill_thread = threading.Thread(target=self._kill, args=(proc.pid,
          completed))

      kill_thread.start()
      returncode = None
      try:
        compile_out = proc.communicate()[0]
      except Exception, e:
        compile_out = str(e)
        returncode = -9
      completed.append(True)
      kill_thread.join()

      if returncode is None: returncode = proc.returncode

      if returncode != 0:
        if "killed" in completed:
          error = "compilation_timelimit"
        else:
          error = "compilation_error"
        return "", compile_out, returncode, 0.0, error

      return self._run_user_program(eval(lang_conf["vm_command"])(source_file),
          stdin, self.apparmor_profile(lang_conf, "vm_apparmor_profile"),
          time.time() - compile_start_time, chdir=source_dir,
          custom_timelimit=custom_timelimit)
    finally:
      shutil.rmtree(source_dir)

  def default_apparmor_profile(self):
    return self.config.get("default-apparmor-profiles", "compiled")

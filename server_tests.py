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

import server, unittest, urllib, json, re

__author__ = "JT Olds"
__copyright__ = "Copyright 2011 Instructure, Inc."
__license__ = "AGPLv3"
__email__ = "jt@instructure.com"

wrapper = server.straitjacket.StraitJacket(server.DEFAULT_CONFIG_DIR)
webapp = server.webapp(wrapper)

class TestCase(unittest.TestCase): pass

class WebServerTests(TestCase):

  def execute(self, data, expected_stdout, expected_stderr,
      expected_exitstatus, expected_error):
    response = webapp.request("/execute", "POST",
        urllib.urlencode(data))
    self.assertEquals(response.status, "200 OK")
    r = json.loads(response.data)
    try:
      if type(expected_stdout) in (str, unicode):
        if r["stdout"] != expected_stdout: raise Exception
      else:
        if not expected_stdout.search(r["stdout"]): raise Exception
      if type(expected_stderr) in (str, unicode):
        if r["stderr"] != expected_stderr: raise Exception
      else:
        if not expected_stderr.search(r["stderr"]): raise Exception
      if r["exitstatus"] != expected_exitstatus: raise Exception
    except:
      raise Exception, "unexpected output: %s" % r

  def testBadLanguage(self):
    self.assertEquals(webapp.request("/execute", "POST",
        urllib.urlencode({
      "language": "non-existent-language",
      "source": "",
      "stdin": ""})).status, "400 Bad Request")

  def testTooLongExecution(self):
    self.execute({
        "language": "ruby1.8",
        "source": "sleep 20\n",
        "stdin": ""},
      "",
      "",
      -9,
      "runtime_timelimit")

  def testOkayExecution(self):
    self.execute({
        "language": "ruby1.8",
        "source": "sleep 10\n",
        "stdin": ""},
      "",
      "",
      0,
      "")

  def testLimitedExecution(self):
    self.execute({
        "language": "ruby1.8",
        "source": "sleep 10\n",
        "stdin": "",
        "timelimit": "1.5"},
      "",
      "",
      -9,
      "runtime_timelimit")

  def testOkayLimitedExecution(self):
    self.execute({
        "language": "ruby1.8",
        "source": "sleep 1\n",
        "stdin": "",
        "timelimit": "2.5"},
      "",
      "",
      0,
      "")


class StraitJacketTests(TestCase):

  def sj_run(self, lang, source):
    stdout, stderr, exitstatus, runtime, error = wrapper.run(lang, source, "",
        custom_timelimit=1.0)
    return (exitstatus, error)

  def testCompilerProfileErrors(self):
    self.assertEquals(wrapper.enabled_languages["c"][
        "execution_profile"].__class__.__name__, "CompilerProfile")
    self.assertEquals(self.sj_run("c",
        "int main(int argc, char** argv) { return 0; }"), (0, ""))
    self.assertEquals(self.sj_run("c",
        "(int argc, char** argv) { }"), (1, "compilation_error"))
    self.assertEquals(self.sj_run("c",
        "int main(int argc, char** argv) { return 1; }"), (1, "runtime_error"))
    self.assertEquals(self.sj_run("c",
        "#include <time.h>\nint main(int argc, char** argv) { "
        "sleep(20); return 0; }"), (-9, "runtime_timelimit"))

  def testInterpreterProfileErrors(self):
    self.assertEquals(wrapper.enabled_languages["python"][
        "execution_profile"].__class__.__name__, "InterpreterProfile")
    self.assertEquals(self.sj_run("python",
        "print 'hi'\n"), (0, ""))
    self.assertEquals(self.sj_run("python",
        "raise 'error'\n"), (1, "runtime_error"))
    self.assertEquals(self.sj_run("python",
        "import time\ntime.sleep(20)"), (-9, "runtime_timelimit"))

  def testVMProfileErrors(self):
    self.assertEquals(wrapper.enabled_languages["java"][
        "execution_profile"].__class__.__name__, "VMProfile")
    self.assertEquals(self.sj_run("java",
        "class Main { public static void main(String[] args) { } }"), (0, ""))
    self.assertEquals(self.sj_run("java",
        "object extends { 0 }"), (1, "compilation_error"))
    self.assertEquals(self.sj_run("java",
        "class Main { public static void main(String[] args) throws Exception "
        "{ throw new Exception(); } }"), (1, "runtime_error"))
    self.assertEquals(self.sj_run("java",
        "class Main { public static void main(String[] args) throws Exception "
        "{ Thread.sleep(20000); } }"), (-9, "runtime_timelimit"))


if __name__ == "__main__": unittest.main()

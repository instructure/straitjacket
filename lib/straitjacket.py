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

import os, ConfigParser, re, sys, exec_profiles, subprocess

__author__ = "JT Olds"
__copyright__ = "Copyright 2011 Instructure, Inc."
__license__ = "AGPLv3"
__email__ = "jt@instructure.com"

class Error_(Exception): pass
class InputError(Error_): pass
class ConfigError(Error_): pass
class LanguageInitError(ConfigError): pass

_LANG_MATCH = re.compile(r'^lang-(.*).conf$')
_TEST_MATCH = re.compile(r'^test_')
REQUIRED_TESTS = ["test-simple", "test-rlimit", "test-apparmor"]
def stderr_log(msg): print >>sys.stderr, msg


def safe_language_check(config_file, language, remote_call, log_method):
  passed_tests = 0
  total_tests = 0
  test_names = set()
  for test in config_file.sections():
    if test[:5] != "test-": continue
    test_names.add(test)
    total_tests += 1
    config = dict(config_file.items(test))
    try:
      stdout, stderr, exitstatus, runtime, error = remote_call(config["source"],
          config.get("stdin", ""))
      if not re.compile(config.get("stderr", "")).search(stderr):
        raise LanguageInitError, "stderr: %s" % stderr
      if not re.compile(config.get("stdout", "")).search(stdout):
        raise LanguageInitError, "stdout: %s" % stdout
      if int(config["exitstatus"]) != exitstatus:
        raise LanguageInitError, "exitstatus: %s" % exitstatus
      if config.has_key("error") and config["error"] != error:
        raise LanguageInitError, "error: %s" % error
      passed_tests += 1
    except LanguageInitError, e:
      log_method("%s:%s unexpected %s" % (language, test, e.args[0]))
    except Exception, e:
      log_method("%s:%s unexpected error: %s" % (language, test, e))
  all_required_tests = True
  for required_test in REQUIRED_TESTS:
    if required_test not in test_names:
      all_required_tests = False
      log_method("%s: missing required test '%s'" % (language, required_test))
      break
  log_method("%s passed %d/%d tests: %s" % (language, passed_tests, total_tests,
      passed_tests == total_tests and all_required_tests and "enabled." or
      "DISABLED!"))
  return passed_tests == total_tests and all_required_tests


class StraitJacket(object):

  def _read_config_file(self, file_):
    parsed = ConfigParser.SafeConfigParser()
    f = file(os.path.join(self.config_dir, file_))
    parsed.readfp(f)
    f.close()
    return parsed

  def _get_all_languages(self):
    return sorted((match.group(1) for match in (
              _LANG_MATCH.search(dir_) for dir_ in os.listdir(self.config_dir))
            if match))

  def __init__(self, config_dir, log_method=None, skip_language_checks=False):
    self.log_method = log_method or stderr_log
    self.config_dir = config_dir

    self.config = self._read_config_file("global.conf")
    if os.path.exists(os.path.join(self.config_dir, "cached_versions.conf")):
      self.cached_versions = self._read_config_file("cached_versions.conf")
      self.cached_versions_needs_flush = False
    else:
      self.cached_versions = ConfigParser.SafeConfigParser()

    if not self.cached_versions.has_section("versions"):
      self.cached_versions.add_section("versions")
      self.cached_versions_needs_flush = True

    self.enabled_languages = {}
    self.exec_profiles = {}

    languages = self._get_all_languages()

    if not skip_language_checks:
      self.log_method("Initializing %d languages." % len(languages))

    for language in languages:
      self._initialize_language(language, skip_language_checks)

    if self.cached_versions_needs_flush:
      try:
        f = file(os.path.join(self.config_dir, "cached_versions.conf"), "w")
        self.cached_versions.write(f)
        f.close()
      except:
        self.log_method("Failed to save cached language version data")

    if not skip_language_checks:
      self.log_method("Initialized %d languages, %d enabled." % (
          len(languages), len(self.enabled_languages)))

  def _initialize_language(self, language, skip_language_checks):
    config_file = self._read_config_file("lang-%s.conf" % language)
    language_config = dict(config_file.items("general"))

    language_config["execution_profile"] = self._get_exec_profile(
        language_config["execution_profile"])

    options = language_config.get("options", "").strip()
    if len(options) > 0: language_config["options"] = options.split(" ")
    else: language_config["options"] = []

    language_config["visible_name"] = "%s (%s)" % (
        language_config["name"], " ".join([language_config["binary"]] +
            language_config["options"]))

    if skip_language_checks:
      if self.cached_versions.has_option("versions", language):
        language_config["version"] = self.cached_versions.get("versions",
            language)
        self.enabled_languages[language] = language_config
      return

    try:
      language_config["version"] = [x for x in (x.strip()
          for x in subprocess.Popen([language_config["binary"],
          language_config.get("version_switch", "--version")],
          stdout=subprocess.PIPE, stderr=subprocess.STDOUT).communicate()[0
          ].split("\n")) if len(x) > 0][0]
    except:
      language_config["version"] = "Unknown version"

    self.log_method("Initializing %s - %s" % (language_config["visible_name"],
        language_config["version"]))

    if not self._safe_language(config_file, language_config, language):
      if self.cached_versions.has_option("versions", language):
        self.cached_versions.remove_option("versions", language)
        self.cached_versions_needs_flush = True
      return

    if not self.cached_versions_needs_flush:
      if self.cached_versions.has_option("versions", language):
        if self.cached_versions.get("versions", language) != language_config[
            "version"]:
          self.cached_versions_needs_flush = True
      else:
        self.cached_versions_needs_flush = True

    if self.cached_versions_needs_flush:
      self.cached_versions.set("versions", language,
          language_config["version"])

    self.enabled_languages[language] = language_config

  def _safe_language(self, config_file, lang_config, language):
    return safe_language_check(config_file, language,
        lambda source, stdin: self._real_run(lang_config, source, stdin),
        self.log_method)

  def _get_exec_profile(self, profile_name):
    if self.exec_profiles.has_key(profile_name):
      return self.exec_profiles[profile_name]
    if not hasattr(exec_profiles, profile_name):
      raise ConfigError, "invalid execution profile"
    profile_class = getattr(exec_profiles, profile_name)
    if not issubclass(profile_class, exec_profiles.BaseProfile):
      raise ConfigError, "invalid execution profile"
    self.exec_profiles[profile_name] = profile_class(self.config)
    return self.exec_profiles[profile_name]

  def run(self, language, source, stdin, custom_timelimit=None):
    if language not in self.enabled_languages:
      raise InputError, "invalid language"
    return self._real_run(self.enabled_languages[language], source, stdin,
        custom_timelimit)

  def _real_run(self, conf, source, stdin, custom_timelimit=None):
    return conf["execution_profile"].run(conf, source, stdin, custom_timelimit)

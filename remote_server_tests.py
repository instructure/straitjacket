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

import sys, urllib2, json, ConfigParser, urllib
from lib.straitjacket import safe_language_check

__author__ = "JT Olds"
__copyright__ = "Copyright 2011 Instructure, Inc."
__license__ = "AGPLv3"
__email__ = "jt@instructure.com"

def test_language(language, server_info, remote_server):
  config_file = ConfigParser.SafeConfigParser()
  config_file.readfp(file("config/lang-%s.conf" % language))
  def remote_call(source, stdin):
    response = json.loads(urllib2.urlopen("%s/execute" % remote_server,
        urllib.urlencode({
          "language": language,
          "source": source,
          "stdin": stdin})).read())
    return response["stdout"], response["stderr"], response["exitstatus"], \
        response["time"], response["error"]
  def log_method(message):
    print >>sys.stderr, message

  return safe_language_check(config_file, language, remote_call,
      log_method)

def test(remote_server):
  server_info = json.loads(urllib2.urlopen("%s/info" % remote_server).read())
  total_languages = 0
  working_languages = 0
  for language in server_info["languages"]:
    total_languages += 1
    if test_language(language, server_info["languages"][language],
        remote_server):
      working_languages += 1
  print "%d/%d languages working" % (working_languages, total_languages)

def main(argv):
  try:
    remote_server_base_url = argv[1]
  except:
    print "usage: %s <remote_server_base_url>" % argv[0]
    return 1

  test(remote_server_base_url)

  return 0

if __name__ == "__main__": sys.exit(main(sys.argv))

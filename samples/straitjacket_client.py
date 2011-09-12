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

import sys, urllib2, json, urllib, time

__author__ = "JT Olds"
__copyright__ = "Copyright 2011 Instructure, Inc."
__license__ = "AGPLv3"
__email__ = "jt@instructure.com"

SERVER_INFO_EXPIRY = 1800

class StraitJacketClient(object):

  def __init__(self, base_url):
    self.base_url = base_url
    self._server_info = {}
    self._server_info_expiration = 0

  def _ensure_server_info(self):
    now = time.time()
    if self._server_info_expiration < now:
      self._server_info = json.loads(urllib2.urlopen("%s/info" % self.base_url
          ).read())
      self._server_info_expiration = now + SERVER_INFO_EXPIRY

  def enabled_languages(self):
    self._ensure_server_info()
    return self._server_info["languages"]

  def run(self, language, source, stdin, custom_timelimit=None):
    post_data = { "language": language,
                  "source": source,
                  "stdin": stdin }
    if custom_timelimit is not None:
      post_data["timelimit"] = custom_timelimit
    response = json.loads(urllib2.urlopen("%s/execute" % self.base_url,
        urllib.urlencode(post_data)).read())
    return response["stdout"], response["stderr"], response["exitstatus"], \
        response["time"], response["error"]

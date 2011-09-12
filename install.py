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

import os, subprocess, ConfigParser

CONFIG_PATH = "config/global.conf"
config = ConfigParser.SafeConfigParser()
config.readfp(file(CONFIG_PATH))

def mkdir_p(path):
  try: os.makedirs(path)
  except: pass

mkdir_p(config.get("directories", "temp_root"))
mkdir_p(config.get("directories", "source"))
mkdir_p(config.get("directories", "compiler"))
mkdir_p(config.get("directories", "execution"))

print "Make sure to make your source, compiler, and execution directories are "\
    "writable by the user that runs straitjacket."

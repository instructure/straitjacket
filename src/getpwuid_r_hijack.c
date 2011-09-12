/*
 * Copyright (C) 2011 Instructure, Inc.
 *
 * This file is part of StraitJacket.
 *
 * StraitJacket is free software: you can redistribute it and/or modify it under
 * the terms of the GNU Affero General Public License as published by the Free
 * Software Foundation, version 3 of the License.
 *
 * StraitJacket is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
 * A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
 * details.
 *
 * You should have received a copy of the GNU Affero General Public License along
 * with this program. If not, see <http://www.gnu.org/licenses/>.
 */

#include <stdlib.h>
#include <sys/types.h>
#include <pwd.h>
#include <errno.h>
#include <string.h>

// gcc -fPIC -c getpwuid_r_hijack.c -o getpwuid_r_hijack.o
// gcc -shared -o getpwuid_r_hijack.so getpwuid_r_hijack.o

int getpwuid_r(uid_t uid, struct passwd *pwbuf, char *buf, size_t buflen,
    struct passwd **pwbufp) {
  *pwbufp = NULL;
  if(buflen < 31) return ERANGE;
  memcpy(buf, "user\0x\0/nonexistent\0/bin/false\0", 31);
  pwbuf->pw_name = buf;
  pwbuf->pw_passwd = buf + 5;
  pwbuf->pw_uid = uid;
  pwbuf->pw_gid = 65534;
  pwbuf->pw_gecos = buf;
  pwbuf->pw_dir = buf + 7;
  pwbuf->pw_shell = buf + 20;
  *pwbufp = pwbuf;
  return 0;
}

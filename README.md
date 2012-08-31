StraitJacket 0.1
=====

This web application is a (hopefully) safe and secure remote execution
environment framework. It builds on top of Linux' AppArmor system calls and as
such won't be able to run on any other operating system.

The end goal is to be able to run someone else's source code in any (configured)
language automatically and not worry about hax.

Design
=====

StraitJacket comes with a number of predetermined AppArmor profiles. When
StraitJacket gets an incoming request to run some code, it will, after calling
fork, but before exec, tell AppArmor that on exec, it wants to switch that
process into a specific profile permanently. AppArmor profiles also provide
standard resource-limit style constraints.

AppArmor really does all the heavy lifting. For more information please see
[AppArmor's wiki](http://wiki.apparmor.net/). A big thanks to Immunix and the
subsequent AppArmor team!

API
===

The API has two calls.

```
GET /info
 * No arguments.
 * This will return, in JSON format, server info, such as what languages are
   currently supported.

POST /execute
 * Takes parameters: language (required), stdin (required, but can be empty),
   source (required), and timelimit (optional, in seconds).
 * Returns, in JSON format, stdout, stderr, exitstatus, time, and error.
```

A sample client library is provided in the samples directory (it's what
[CodeWarden](https://github.com/instructure/codewarden/) uses).

Installation
=====

AMI
-----

If you'd rather just ignore all of the following and use a pre-existing Amazon
Web Services AMI, try looking for ami-8de553e4

Dependencies
-----

You will need to install all of the appropriate files for each language you want
to run. On Ubuntu 12.04, suggested packages include but are not limited to:

gcc, mono-gmcs, g++, guile-1.8, ghc, lua5.1, ocaml, php5, python, ruby,
ruby1.9, scala, racket, golang, openjdk-6-jdk

You'll probably want to get nodejs and DMD from elsewhere.
http://dlang.org/download.html
https://github.com/joyent/node/wiki/Installing-Node.js-via-package-manager

Dependencies include:

python-webpy, python-libapparmor, apparmor

AppArmor
------

There are a number of AppArmor profiles provided in files/etc/apparmor.d.
You should transfer these to wherever your AppArmor profiles are stored.
Additionally, you need to transfer the AppArmor profile abstractions provided in
files/etc/apparmor.d/abstractions similarly.

Once you have successfully installed your AppArmor profiles, make sure to force
AppArmor to reload its configuration.

System directories
-------

There are a number of system directories StraitJacket uses for intermediate
stages of execution, all (configurably) prefixed by /var/local/straitjacket.
Please take a look at both config/global.conf and install.py (which currently
only can be relied upon to make these directories for you, unfortunately).

LD_PRELOAD hacks
-------

Some languages (c#) require access to the getpwuid_r system call, which reads
/etc/passwd, which is disallowed by AppArmor, which promptly fails, causing
the runtime to bail. To counteract this without actually just adding /etc/passwd
read access, there is a getpwuid_r LD_PRELOAD library in the src/ directory.

The current config/lang-c#.conf file expects the getpwuid_r_hijack.so module
that can be built in the src/ directory to be in /var/local/straitjacket/lib/

Web
-----

This application is (mostly) a standard web.py WSGI-capable web app. A sample
Apache configuration is provided in files/etc/apache2/sites-available.

It is recommended that you verify that your server is properly and
safely configured before full use. The only thing to know here is that by
default, StraitJacket will not enable a language unless it passes all of that
language's specific tests, UNLESS you are running in WSGI mode. If you are
running in WSGI mode, this preventative step is disabled.

You can both run tests locally (using server_tests.py) to ensure your system
is correctly set up, or remotely (using remote_server_tests.py).

License
=====

StraitJacket is released under the AGPLv3. Please see COPYRIGHT and LICENSE.

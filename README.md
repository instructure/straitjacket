StraitJacket 2.0
=====

This web application is a (hopefully) safe and secure remote execution
environment framework. It builds on top of Docker and Linux' AppArmor system
calls and as such won't be able to run on any other operating system.

The end goal is to be able to run someone else's source code in any (configured)
language automatically and not worry about hax.

Design
=====

StraitJacket comes with a number of predetermined AppArmor profiles, and docker
containers built for each supported language. When StraitJacket gets an incoming
request to run some code, it will launch that container with the AppArmor
profile applied.

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

A sample client library is provided in the samples directory.

Installation
=====

Development
-----------

There is a Vagrantfile provided for devloping Straitjacket. Run `vagrant up` to
build the image.

To run straitjacket locally for development, ssh into the VM with `vagrant ssh` and run:

    cd straitjacket
    ./straitjacket-setup.sh
    ./run-dev.sh

This will listen on port 8081, which is forwarded to the host machine.

You'll need to re-run `straitjacket-setup.sh` any time you add/modify a language
apparmor profile or docker image. New docker images need to be added there, as well.

AMI
-----

You can build an AWS AMI using [Packer](https://packer.io/) by calling the
`build_ami.sh` build script. You'll need to modify `packer.json` for your VPC
and subnet IDs.

A pre-built AMI may be made public later.

License
=====

StraitJacket is released under the AGPLv3. Please see COPYRIGHT and LICENSE.

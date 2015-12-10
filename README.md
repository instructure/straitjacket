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

The API has two calls:

```
GET /info
POST /execute
```

There is also a more advanced websocket API at `GET /executews`.

You can view the API documentation directly from github at
http://petstore.swagger.io/?url=https://raw.githubusercontent.com/instructure/straitjacket/master/public/api/2015-07-14.yml
though you'll need to spin up an instance of straitjacket to actually perform
API calls from that page.

Installation
=====

Development
-----------

There is a Vagrantfile provided for developing Straitjacket. Run `vagrant up` to
build the image.

To run straitjacket locally for development, ssh into the VM with `vagrant ssh` and run:

    cd straitjacket
    sudo ./straitjacket-setup.sh
    ./run-dev.sh

This will listen on port 8081, which is forwarded to the host machine.

You'll need to re-run `straitjacket-setup.sh` any time you add/modify a language
apparmor profile or docker image. New docker images need to be added there, as well.

To run the language tests (sanity checks) defined in the config .yml files, run:

    ./run-dev.sh --test

AMI
-----

You can build an AWS AMI using [Packer](https://packer.io/) by calling the
`build_ami.sh` build script. You'll need to modify `packer.json` for your VPC
and subnet IDs.

A pre-built AMI may be made public later.

License
=====

StraitJacket is released under the AGPLv3. Please see COPYRIGHT and LICENSE.

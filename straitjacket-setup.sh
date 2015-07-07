#!/bin/bash
set -e

cp -R files/etc/apparmor.d/* /etc/apparmor.d/
service apparmor reload

docker build -t straitjacket-csharp languages/csharp
docker build -t straitjacket-nodejs languages/nodejs
docker build -t straitjacket-ruby languages/ruby
docker build -t straitjacket-d languages/d
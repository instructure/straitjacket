#!/bin/bash
set -e

docker build -t straitjacket .
docker run --rm --entrypoint go straitjacket test -v ./handlers
deploy/install-containers
./run-dev.sh --test --disable-apparmor

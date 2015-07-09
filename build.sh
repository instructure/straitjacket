#!/bin/bash
set -e

deploy/build-containers
./run-dev.sh --test --disable-apparmor

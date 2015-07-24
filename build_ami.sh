#!/bin/bash

set -e

git archive -o straitjacket.tar HEAD
packer build deploy/packer.json

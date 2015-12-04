#!/bin/bash

docker build -t straitjacket . && docker run -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp -v $(pwd):/go/src/straitjacket --rm -p 8081:8081 straitjacket $@

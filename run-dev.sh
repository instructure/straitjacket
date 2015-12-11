#!/bin/bash

docker build -t straitjacket . && exec docker run -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/go/src/straitjacket --rm -p 8081:8081 -it straitjacket $@

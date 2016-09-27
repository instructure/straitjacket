FROM golang:1.7.1

WORKDIR /go/src/straitjacket

RUN go get github.com/tools/godep
ADD . /go/src/straitjacket
RUN go install straitjacket

ENTRYPOINT ["/go/bin/straitjacket"]

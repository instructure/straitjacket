FROM golang:1.5

ENV GO15VENDOREXPERIMENT 1
WORKDIR /go/src/straitjacket

RUN go get github.com/tools/godep
ADD . /go/src/straitjacket
RUN go-wrapper install

ENTRYPOINT ["go-wrapper", "run"]
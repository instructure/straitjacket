FROM golang:1.4.2

WORKDIR /go/src/straitjacket

RUN go get github.com/tools/godep
ADD Godeps /go/src/straitjacket/Godeps
RUN godep restore
ADD . /go/src/straitjacket
RUN go-wrapper install

ENTRYPOINT ["go-wrapper", "run"]
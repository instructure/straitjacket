FROM golang:1.7.1

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

WORKDIR /src

ENTRYPOINT ["/build-run"]

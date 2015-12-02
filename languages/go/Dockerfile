FROM golang:1.4.2

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

WORKDIR /src

ENTRYPOINT ["/build-run"]

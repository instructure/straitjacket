FROM buildpack-deps:wheezy

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["/build-run"]

FROM buildpack-deps:jessie

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

ENTRYPOINT ["/build-run"]

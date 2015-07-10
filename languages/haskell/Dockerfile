FROM buildpack-deps:wheezy

RUN         apt-get update && \
            apt-get install -y ghc && \
            apt-get clean && \
            rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["/build-run"]

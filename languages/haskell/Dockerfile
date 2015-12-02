FROM buildpack-deps:jessie

RUN         apt-get update && \
            apt-get install -y ca-certificates libtinfo-dev ca-certificates g++ libgmp10 libgmp-dev libffi-dev zlib1g-dev && \
            apt-get clean && \
            cd /tmp && \
            wget -nv https://haskell.org/platform/download/7.10.2/haskell-platform-7.10.2-a-unknown-linux-deb7.tar.gz && \
            tar xf haskell-platform-7.10.2-a-unknown-linux-deb7.tar.gz && \
            ./install-haskell-platform.sh && \
            rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

ENTRYPOINT ["/build-run"]

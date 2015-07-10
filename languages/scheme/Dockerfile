FROM buildpack-deps:wheezy

RUN         apt-get update && \
            apt-get install -y racket && \
            apt-get clean && \
            rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN useradd docker
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["plt-r5rs"]

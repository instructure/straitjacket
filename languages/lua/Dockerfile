FROM buildpack-deps:jessie

RUN apt-get update && \
    apt-get install -y lua5.2 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN useradd docker
USER docker

ENTRYPOINT ["lua"]

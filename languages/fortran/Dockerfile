FROM buildpack-deps:jessie

RUN apt-get update && \
    apt-get install -y gfortran && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

ENTRYPOINT ["/build-run"]

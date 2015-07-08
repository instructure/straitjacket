FROM buildpack-deps:wheezy

RUN         apt-get update && \
            apt-get install -y ocaml-nox && \
            apt-get clean && \
            rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN useradd docker
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["ocaml"]

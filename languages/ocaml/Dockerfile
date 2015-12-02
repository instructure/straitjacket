FROM buildpack-deps:jessie

RUN         apt-get update && \
            apt-get install -y ocaml-nox opam camlp4-extra && \
            apt-get clean && \
            rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN useradd -m docker
USER docker

RUN opam init -y && \
    opam switch 4.02.3 && \
    opam install -y core

USER root
ADD build-run /build-run
RUN chmod +x /build-run
USER docker

WORKDIR /src

ENTRYPOINT ["/build-run"]

FROM buildpack-deps:jessie

RUN curl -f -L https://static.rust-lang.org/rustup.sh -O && \
    SHELL=sh sh rustup.sh --yes --verbose --revision=1.1.0 --disable-sudo && \
    rm rustup.sh

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

WORKDIR /src

ENTRYPOINT ["/build-run"]

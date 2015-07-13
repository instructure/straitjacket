FROM buildpack-deps:jessie

RUN useradd docker
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["perl"]

FROM buildpack-deps:wheezy

RUN useradd docker
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["perl"]

FROM mono:4.0

RUN mkdir /build
RUN useradd docker && chown docker /build
USER docker

ADD build-run /build-run

VOLUME /src
WORKDIR /src

ENTRYPOINT ["/build-run"]

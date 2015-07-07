FROM ibuclaw/gdc:4.9-v2.066

ADD build-run /build-run
RUN chmod +x /build-run

RUN mkdir /build
RUN useradd docker && chown docker /build
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["/build-run"]

FROM java:7

ADD build-run /build-run
RUN chmod +x /build-run

RUN mkdir /build
RUN useradd docker && chown docker /build
USER docker

VOLUME /src
WORKDIR /src

ENV JAVA_TOOL_OPTIONS "-Xmx256m -Xms256m -Xss256k"
ENTRYPOINT ["/build-run"]

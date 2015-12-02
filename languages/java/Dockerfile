FROM java:7

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

ENV JAVA_TOOL_OPTIONS "-Xmx256m -Xms256m -Xss256k"
ENTRYPOINT ["/build-run"]

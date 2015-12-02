FROM williamyeh/scala:2.11.6

ADD build-run /build-run
RUN chmod +x /build-run

RUN useradd docker
USER docker

# ENV SCALA_RUNNER_DEBUG true
ENTRYPOINT ["/build-run"]

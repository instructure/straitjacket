FROM java:7

RUN useradd docker

RUN mkdir /clojure
RUN curl -o /clojure/clojure-1.7.0-slim.jar https://oss.sonatype.org/content/repositories/snapshots/org/clojure/clojure/1.7.0-master-SNAPSHOT/clojure-1.7.0-master-20150617.180331-46-slim.jar
RUN chown -R docker /clojure

USER docker

VOLUME /src
WORKDIR /src

ENV JAVA_TOOL_OPTIONS "-Xmx256m -Xms256m -Xss256k"
ENTRYPOINT ["java", "-cp", "/clojure/clojure-1.7.0-slim.jar", "clojure.main"]

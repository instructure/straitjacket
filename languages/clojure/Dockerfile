FROM java:8

RUN useradd docker

RUN mkdir /clojure
RUN curl -o /clojure/clojure.zip http://repo1.maven.org/maven2/org/clojure/clojure/1.8.0/clojure-1.8.0.zip
RUN unzip /clojure/clojure.zip 'clojure*.jar' -d /clojure && chmod -R go+rx /clojure
RUN chown -R docker /clojure

USER docker

ENV JAVA_TOOL_OPTIONS "-Xmx256m -Xms256m -Xss256k -XX:-UsePerfData"
ENTRYPOINT ["java", "-cp", "/clojure/clojure-1.8.0/clojure-1.8.0-slim.jar", "clojure.main"]

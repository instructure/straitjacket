FROM buildpack-deps:jessie

RUN echo "deb http://packages.erlang-solutions.com/debian jessie contrib" >> /etc/apt/sources.list && \
    wget -qO - http://packages.erlang-solutions.com/debian/erlang_solutions.asc | apt-key add - && \
    apt-get update && \
    apt-get install -y --no-install-recommends erlang=1:18.* elixir=1.0.* && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# elixir relies on the locale being set to UTF-8
ENV LANG C.UTF-8

RUN useradd docker
USER docker

WORKDIR /src

ENTRYPOINT ["elixir"]

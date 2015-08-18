FROM node:0.12

RUN useradd docker
USER docker

VOLUME /src
WORKDIR /src

ENTRYPOINT ["node", "--harmony"]
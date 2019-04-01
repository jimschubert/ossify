ARG GO_VERSION=1.11.6

FROM golang:${GO_VERSION}

LABEL "name"="ossify github action"
LABEL "maintainer"="Jim Schubert <james.schubert@gmail.com>"
LABEL "version"="1.0.0"

LABEL "com.github.actions.name"="ossify github action"
LABEL "com.github.actions.description"="Building a golang app."
# See https://feathericons.com/
LABEL "com.github.actions.icon"="github"
LABEL "com.github.actions.color"="gray-dark"

ENV DEP_VERSION=0.5.0
ENV SHELL=bash

RUN apt-get update \
 && apt-get --no-install-recommends install zip \
 && wget -nv https://github.com/golang/dep/releases/download/v$DEP_VERSION/dep-linux-amd64 -O /usr/local/bin/dep \
 && chmod +x /usr/local/bin/dep

VOLUME $GOPATH/src

COPY ./entrypoint.sh /

ENTRYPOINT /entrypoint.sh

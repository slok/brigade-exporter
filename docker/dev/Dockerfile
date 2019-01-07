FROM golang:1.11-alpine

RUN apk --no-cache add \
    g++ \
    git \
    bash

RUN mkdir /src

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid app && \
    adduser -D -u $uid -G app app && \
    chown app:app -R /go && \
    chown app:app -R /src

# Fill go mod cache.
RUN mkdir /tmp/cache
COPY go.mod /tmp/cache
COPY go.sum /tmp/cache
RUN cd /tmp/cache && \
    go mod download

USER app

WORKDIR /src

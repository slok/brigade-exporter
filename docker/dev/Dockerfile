FROM golang:1.11-alpine

RUN apk --no-cache add \
    g++ \
    git \
    bash

# Mock creator
RUN go get -u github.com/vektra/mockery/.../

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid app && \
    adduser -D -u $uid -G app app && \
    chown app:app -R /go


USER app
# Although we are using go 1.11 we set here in order to allow vendor directory
# work as a cache.
WORKDIR /go/src/github.com/slok/brigade-exporter

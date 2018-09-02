FROM golang:1.11-alpine AS build-stage

RUN apk --no-cache add \
    g++ \
    git \
    make

ARG VERSION
ENV VERSION=${VERSION}
WORKDIR /go/src/github.com/slok/brigade-exporter
COPY . .
RUN ./hack/scripts/build-binary.sh

# Final image.
FROM alpine:latest
RUN apk --no-cache add \
  ca-certificates
COPY --from=build-stage /go/src/github.com/slok/brigade-exporter/bin/brigade-exporter /usr/local/bin/brigade-exporter
ENTRYPOINT ["/usr/local/bin/brigade-exporter"]


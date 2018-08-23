#!/usr/bin/env sh

set -e

if [ -z ${VERSION} ]; then
    echo "VERSION env var needs to be set"
    exit 1
fi

REPOSITORY="slok/"
IMAGE="brigade-exporter"


docker build \
    --build-arg VERSION=${VERSION} \
    -t ${REPOSITORY}${IMAGE}:${VERSION} \
    -f ./docker/prod/Dockerfile .
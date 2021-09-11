FROM golang:1.17

RUN apt-get update && apt-get install --no-install-recommends -y \
    make \
    libmapnik-dev \
    fonts-noto-cjk \
    fonts-noto-hinted \
    fonts-noto-unhinted \
    fonts-hanazono \
    ttf-unifont \
    fonts-dejavu-core

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1

ENV WORKDIR=/src

RUN mkdir -p ${WORKDIR}

WORKDIR ${WORKDIR}

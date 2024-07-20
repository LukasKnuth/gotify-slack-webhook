ARG GO_VERSION
FROM golang:$GO_VERSION-bullseye
RUN git config --global --add safe.directory '*'

RUN \
  apt-get update && \
  apt-get install -y gcc-aarch64-linux-gnu g++-aarch64-linux-gnu

ENV CGO_ENABLED=1
ENV CC=aarch64-linux-gnu-gcc
ENV CXX=aarch64-linux-gnu-g++
ENV GOOS=linux
ENV GOARCH=arm64

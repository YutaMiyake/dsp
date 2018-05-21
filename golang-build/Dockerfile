FROM golang:1.9.0-stretch

RUN apt-get update && \
    apt-get install -y build-essential --no-install-recommends

RUN go get -u github.com/golang/dep/cmd/dep

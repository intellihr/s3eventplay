FROM golang:1.9.2 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/intellihr/s3eventplay
WORKDIR /go/src/github.com/intellihr/s3eventplay

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

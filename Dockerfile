# syntax=docker/dockerfile:1

## Certs
FROM alpine:3.16.0 as certs
RUN apk add --update --no-cache ca-certificates

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /github-client

## Deploy
#FROM gcr.io/distroless/base-debian10
FROM debian:10
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /github-client /bin/github-client
ENTRYPOINT ["/bin/github-client"]
WORKDIR /workdir
ENV LISTEN :2701
EXPOSE 2701
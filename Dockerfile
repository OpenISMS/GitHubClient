# syntax=docker/dockerfile:1

FROM docker.io/library/alpine:3.16.2 as certs
RUN apk add --update --no-cache ca-certificates

FROM docker.io/library/busybox:1.35.0
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY githubclient ./githubclient
COPY templates ./templates
ENTRYPOINT ["./githubclient"]

WORKDIR /workdir
ENV LISTEN :2701
EXPOSE 2701
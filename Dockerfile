# syntax=docker/dockerfile:1

FROM docker.io/library/alpine:3.16.2 as certs
RUN apk add --update --no-cache ca-certificates

FROM docker.io/library/busybox:1.35.0
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY githubclient /bin/githubclient
ENTRYPOINT ["/bin/githubclient"]

WORKDIR /workdir
ENV LISTEN :2701
EXPOSE 2701

LABEL org.opencontainers.image.source=https://github.com/OpenISMS/GitHubClient
LABEL org.opencontainers.image.vendor=OpenISMS
LABEL org.opencontainers.image.title="GitHub Client for OpenISMS"
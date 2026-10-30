FROM golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb AS builder

COPY . /src/webtts
WORKDIR /src/webtts

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

LABEL maintainer="Knut Ahlers <knut@ahlers.me>"
LABEL org.opencontainers.image.authors="Knut Ahlers <knut@ahlers.me>" \
      org.opencontainers.image.url="https://github.com/users/Luzifer/packages/container/package/webtts" \
      org.opencontainers.image.source="https://github.com/Luzifer/webtts" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.title="Simple wrapper around the Google Cloud Text-To-Speech and Azure Text-To-Speech API to output OGG Vorbis Audio to be used with OBS overlays"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/webtts /usr/local/bin/webtts

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/webtts"]
CMD ["--"]

# vim: set ft=Dockerfile:

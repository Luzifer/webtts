FROM golang:1.26.4-alpine@sha256:7a3e50096189ad57c9f9f865e7e4aa8585ed1585248513dc5cda498e2f41812c AS builder

COPY . /src/webtts
WORKDIR /src/webtts

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.24@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4

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

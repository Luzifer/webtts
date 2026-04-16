FROM golang:1.26-alpine@sha256:27f829349da645e287cb195a9921c106fc224eeebbdc33aeb0f4fca2382befa6 AS builder

COPY . /src/webtts
WORKDIR /src/webtts

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

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

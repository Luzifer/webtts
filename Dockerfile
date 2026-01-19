FROM golang:1.25-alpine@sha256:d9b2e14101f27ec8d09674cd01186798d227bb0daec90e032aeb1cd22ac0f029 AS builder

COPY . /src/webtts
WORKDIR /src/webtts

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.23@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62

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

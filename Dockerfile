FROM golang:1.24-alpine@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS builder

COPY . /src/webtts
WORKDIR /src/webtts

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715

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

FROM golang:1.18.3-stretch as builder

LABEL maintainer="Adriano M. La Selva <adrianolaselva@gmail.com>"

ARG VERSION
ENV VERSION=$VERSION

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o actions .

FROM debian:stable-slim

RUN apt-get update \
     && apt-get install -y --no-install-recommends ca-certificates \
     && update-ca-certificates

WORKDIR /app

COPY --from=builder ./app/actions .

ENTRYPOINT ["/app/actions"]
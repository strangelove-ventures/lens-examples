# syntax=docker/dockerfile:latest
FROM golang:1.17-alpine as build
WORKDIR /lens
RUN --mount=type=bind,source=.,ro \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  go build -o /usr/local/bin/indexer ./...

FROM alpine
COPY --link --from=build /usr/local/bin/indexer /usr/local/bin/indexer

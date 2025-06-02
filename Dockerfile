FROM golang:1.21-bookworm as build

COPY . /go/src/bete
WORKDIR /go/src/bete

RUN bin/build.sh

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /go/src/bete/dist /bete

WORKDIR /bete
ENTRYPOINT ["/bete/bin/bete"]

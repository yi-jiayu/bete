FROM golang:1.21-bookworm as build

COPY . /go/src/bete
WORKDIR /go/src/bete

RUN bin/build.sh

FROM debian:bookworm-slim

COPY --from=build /go/src/bete/dist /bete

WORKDIR /bete
ENTRYPOINT ["/bete/bin/bete"]

FROM golang:1.14-buster as build

COPY . /go/src/bete
WORKDIR /go/src/bete

RUN bin/build.sh

FROM gcr.io/distroless/base-debian10

COPY --from=build /go/src/bete/dist /bete

WORKDIR /bete
ENTRYPOINT ["/bete/bin/bete"]
FROM golang:1.14-buster as build

COPY . /go/src/bete

WORKDIR /go/src/bete/cmd/bete
ARG commit
RUN go build -o /go/bin/bete -ldflags "-X main.commit=$commit"

WORKDIR /go/src/bete
RUN go get -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate

FROM gcr.io/distroless/base-debian10

COPY --from=build /go/bin/bete /bete/bin/bete

COPY --from=build /go/src/bete/migrations /bete/migrations

COPY --from=build /go/bin/migrate /bete/bin/migrate

WORKDIR /bete
ENTRYPOINT ["/bete/bin/bete"]
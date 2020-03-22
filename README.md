# Bete
Bus Eta Bot's secret identity

## Developing

### Generating mocks

* Install `mockgen` binary to `bin/mockgen`: `GOBIN=$PWD/bin go get github.com/golang/mock/mockgen`
* Generate mocks: `go generate`

### Database migrations

* Install `migrate` binary to `bin/migrate`: `GOBIN=$PWD/bin go get -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate`
* The `postgres` tag is needed to run migrations against Postgres.
* To generate new migrations: `bin/migrate create -dir migrations -ext sql MIGRATION_NAME`

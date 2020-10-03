# Bete
Bus Eta Bot's secret identity

## Developing

### Running tests

* The `DATABASE_URL` should be set appropriately.
* For a local database: `env DATABASE_URL='postgres://localhost/bete_test?sslmode=disable' go test ./...`

### Generating mocks

* Install `mockgen` binary to `bin/mockgen`: `GOBIN=$PWD/bin go get github.com/golang/mock/mockgen`
* Generate mocks: `go generate`

### Database migrations

* Install `migrate` binary to `bin/migrate`: `GOBIN=$PWD/bin go get -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate`
* The `postgres` tag is needed to run migrations against Postgres.
* To generate new migrations: `bin/migrate create -dir migrations -ext sql MIGRATION_NAME`
* To run migrations: `bin/migrate -path migrations -database 'postgres://localhost/bete_test?sslmode=disable' up`

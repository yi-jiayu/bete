# Bete
Bus Eta Bot's secret identity

## Developing

### Running tests

* The `DATABASE_URL` should be set appropriately.
* For a local database: `env DATABASE_URL='postgres://localhost/bete_test?sslmode=disable' go test ./...`

### Testing with a real bot

* Create a new bot with [@BotFather](https://t.me/BotFather) on Telegram and
  obtain a bot token.
* Request for API access to [LTA
  DataMall](https://www.mytransport.sg/content/mytransport/home/dataMall.html)
  and obtain an account key.
* Build the server binary: `make`
* Run the server with the `DATAMALL_ACCOUNT_KEY`, `TELEGRAM_BOT_TOKEN` and
  `DATABASE_URL` environment variables:

```
env DATABASE_URL='postgres://localhost/bete_test?sslmode=disable' DATAMALL_ACCOUNT_KEY='xxx' TELEGRAM_BOT_TOKEN='xxx' bin/bete
```

* Start a separate ngrok process to expose the server: `ngrok http 8080`
* Set the bot webhook to the public ngrok URL (use the HTTPS one). Note that
  the path should be `/telegram/updates`

```
curl https://api.telegram.org/bot$bot_token/setWebhook?url=https://xxx/telegram/updates
```

* Find your bot on Telegram and chat with it, you should see HTTP requests
  logged by ngrok and bete.
* For convenience, the Makefile contains a `start` target to build and run the
  server, and a `webhook` target to automatically set the webhook. These both
  need the appropriate environment variables to be set.

### Generating mocks

* Install `mockgen` binary to `bin/mockgen`: `GOBIN=$PWD/bin go get github.com/golang/mock/mockgen`
* Generate mocks: `go generate`

### Database migrations

* Install `migrate` binary to `bin/migrate`: `GOBIN=$PWD/bin go get -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate`
* The `postgres` tag is needed to run migrations against Postgres.
* To generate new migrations: `bin/migrate create -dir migrations -ext sql MIGRATION_NAME`
* To run migrations: `bin/migrate -path migrations -database 'postgres://localhost/bete_test?sslmode=disable' up`

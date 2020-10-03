.PHONY: build start webhook

build:
	cd cmd/bete && go build -o ../../bin/bete

start: build
	bin/bete	

webhook:
	curl "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook?url=$(shell curl localhost:4040/api/tunnels/command_line | jq -r '.public_url')/telegram/updates"


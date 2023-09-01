#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
set -x

commit=$(git rev-parse --short --verify HEAD)
mkdir dist

(cd cmd/bete && go build -o ../../dist/bin/bete -ldflags "-X main.commit=$commit")
(cd cmd/seed && go build -o ../../dist/bin/seed)
GOBIN=$PWD/dist/bin go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate

cp -r migrations dist/migrations
cp tour.yaml dist/

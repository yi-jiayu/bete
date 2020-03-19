#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

commit=${1:-}
environment=${2:-}
values=${3:-}

if [[ -z "$commit" || -z "$environment" || -z "$values" ]]; then
  echo "usage: deploy.sh <commit> <environment> <values>"
  exit 1
fi

if [[ -z "${SENTRY_AUTH_TOKEN:-}" ]]; then
  echo "SENTRY_AUTH_TOKEN environment variable not set"
  exit 1
fi

if [[ -z "${SENTRY_ORG:-}" ]]; then
  echo "SENTRY_ORG environment variable not set"
  exit 1
fi

helm upgrade --install \
  bete charts/server \
  -n bete \
  -f "$values" \
  --set image.tag="$commit" \
  --set environment="$environment"

sentry-cli releases deploys "$commit" new -e "$environment"

#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

commit=${1:-}
environment=${2:-}

if [[ -z "$commit" || -z "$environment" ]]; then
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
  bete deploy/charts/server \
  -n bete \
  -f deploy/values.yaml \
  --set image.tag="$commit" \
  --set environment="$environment"

sentry-cli releases deploys "$commit" new -e "$environment"

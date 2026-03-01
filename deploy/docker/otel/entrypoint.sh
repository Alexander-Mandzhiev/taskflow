#!/bin/sh
set -e

if [ -n "${ELASTIC_PASSWORD}" ] && command -v base64 >/dev/null 2>&1; then
  needs_regen="false"
  if [ -z "${ELASTICSEARCH_BASIC_AUTH}" ]; then
    needs_regen="true"
  else
    decoded="$(printf "%s" "$ELASTICSEARCH_BASIC_AUTH" | base64 -d 2>/dev/null || true)"
    echo "$decoded" | grep -q "^elastic:" || needs_regen="true"
  fi

  if [ "$needs_regen" = "true" ]; then
    ELASTICSEARCH_BASIC_AUTH="$(printf "elastic:%s" "$ELASTIC_PASSWORD" | base64 | tr -d '\n')"
    export ELASTICSEARCH_BASIC_AUTH
  fi
fi

exec /otelcol-contrib --config "${OTEL_CONFIG:-/etc/otel-collector-config.yaml}"

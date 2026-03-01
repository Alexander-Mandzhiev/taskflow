FROM otel/opentelemetry-collector-contrib:0.133.0 AS base

FROM alpine:3.18
COPY --from=base /otelcol-contrib /otelcol-contrib

USER root
RUN apk add --no-cache coreutils
COPY otel-collector-config.yaml /etc/otel-collector-config.yaml
COPY otel-collector-config.secure.yaml /etc/otel-collector-config.secure.yaml
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Создаем директорию для сертификатов (если используются)
RUN mkdir -p /etc/otel/certs

# Выбор конфига через env:
# - OTEL_CONFIG=/etc/otel-collector-config.yaml (по умолчанию, без auth)
# - OTEL_CONFIG=/etc/otel-collector-config.secure.yaml (с API key к Elasticsearch)
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 13133
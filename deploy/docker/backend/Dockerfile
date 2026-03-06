# =============================================================================
# Stage 1: Build
# =============================================================================
FROM golang:1.25.7-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Модуль один — backend (без go.work)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

# Статический бинарник для linux/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/backend-app ./cmd/main.go

# =============================================================================
# Stage 2: Runtime
# =============================================================================
FROM alpine:3.21

RUN apk add --no-cache wget \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/backend-app ./
COPY --from=builder /app/config ./config
COPY --from=builder /app/db ./db
RUN chown -R appuser:appgroup /app

USER appuser

ENV CONFIG_PATH=/app/config/production.yaml

EXPOSE 4000

ENTRYPOINT ["./backend-app"]

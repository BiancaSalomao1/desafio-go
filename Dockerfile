# ==========================
# Build Stage
# ==========================
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o api \
    ./cmd/api

# ==========================
# Runtime Stage
# ==========================
FROM alpine:3.24

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./api"]
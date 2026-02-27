# ── Build stage ──
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api-quest .

# ── Run stage ──
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/api-quest .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./api-quest"]

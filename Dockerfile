# ── Builder stage ──────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api cmd/api/main.go


# ── Final stage ───────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /app/api .

# 🔥 COPY MIGRATIONS (INI KUNCI UTAMA)
COPY --from=builder /src/migrations /migrations

EXPOSE 8080

ENTRYPOINT ["./api"]
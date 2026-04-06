# ── Builder stage ──────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api cmd/api/main.go


# ── Final stage ───────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

# ✅ copy binary dari builder
COPY --from=builder /app/api .

# ✅ copy migrations dari builder
COPY --from=builder /src/migrations /migrations

# ✅ install migrate binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
  | tar xvz \
  && mv migrate /usr/local/bin/migrate

EXPOSE 8080

ENTRYPOINT ["./api"]
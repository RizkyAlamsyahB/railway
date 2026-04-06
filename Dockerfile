FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /app/api .

# ✅ copy migrations
COPY --from=builder /src/migrations /migrations

# ✅ install migrate binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
  | tar xvz \
  && mv migrate /usr/local/bin/migrate

EXPOSE 8080

ENTRYPOINT ["./api"]
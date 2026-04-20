FROM golang:1.23.4-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations
COPY seeds ./seeds
COPY docker ./docker

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/cee4ai-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/cee4ai-migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/cee4ai-seed ./cmd/seed

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S cee4ai && \
    adduser -S -G cee4ai cee4ai

COPY --from=builder /out/cee4ai-api /app/cee4ai-api
COPY --from=builder /out/cee4ai-migrate /app/cee4ai-migrate
COPY --from=builder /out/cee4ai-seed /app/cee4ai-seed
COPY --from=builder /src/migrations /app/migrations
COPY --from=builder /src/seeds /app/seeds
COPY --from=builder /src/docker/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh /app/cee4ai-api /app/cee4ai-migrate /app/cee4ai-seed && \
    chown -R cee4ai:cee4ai /app

USER cee4ai

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]

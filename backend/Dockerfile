FROM golang:1.23.3-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/config.json .
COPY --from=builder /app/migrations /app/migrations

CMD ["./server", "--config_path=/root/config.json", "--logger_level=debug"]

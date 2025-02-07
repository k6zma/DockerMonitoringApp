FROM golang:1.23.3-alpine AS builder
RUN apk add --no-cache git build-base

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o pinger ./cmd/pinger/main.go

FROM alpine:latest
RUN apk --no-cache add iputils

WORKDIR /root/
COPY --from=builder /app/pinger .
COPY --from=builder /app/config.json .

CMD ["./pinger", "--config_path=/root/config.json"]
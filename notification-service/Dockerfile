FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN go build -o notify-service ./cmd/server/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/notify-service .

CMD ["./notify-service"]

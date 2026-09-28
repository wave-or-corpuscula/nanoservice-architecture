FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY ./proto ./proto
COPY ./order-service/go.mod ./order-service/go.sum ./order-service/

WORKDIR /app/order-service
RUN go mod download

COPY ./order-service/ .

RUN go build -o order-service ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/order-service .

EXPOSE 9999

CMD ["./order-service"]

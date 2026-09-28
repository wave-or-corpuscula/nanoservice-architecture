FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY ./proto ./proto
COPY ./user-service/go.mod ./user-service/go.sum ./user-service/

WORKDIR /app/user-service
RUN go mod download

COPY ./user-service/ .

RUN go build -o user-service ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user-service .

EXPOSE 8888

CMD ["./user-service"]

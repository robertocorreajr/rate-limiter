# Etapa 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/rate-limiter/main.go

# Etapa 2: Execução
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
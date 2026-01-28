FROM golang:1.25-alpine AS builder

WORKDIR /app

# копируем зависимости + go.sum
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o autoonline .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/autoonline .

RUN apk --no-cache add bash

WORKDIR /app

CMD ["./autoonline"]
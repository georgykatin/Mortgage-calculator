FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myapp .

COPY --from=builder /app/internal/config /app/internal/config

EXPOSE 8080

CMD ["./myapp"]
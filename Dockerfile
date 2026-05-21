# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o product-service .

# Run stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/product-service .
COPY .env .env

EXPOSE 8080
CMD ["./product-service"]

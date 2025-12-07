# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main server.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/public ./public

EXPOSE 8080
CMD ["./main"]

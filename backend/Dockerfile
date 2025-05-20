# Build Stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Runtime Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./main"]
# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Stage 2: Minimal final image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["./main"]
# Start from a small base image with Go installed
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of your code
COPY . .

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
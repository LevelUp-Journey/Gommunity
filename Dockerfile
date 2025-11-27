# Start from the official Go image to build the binary.
FROM golang:1.23-alpine AS builder

# Install git and make (if needed for dependencies or build scripts)
RUN apk add --no-cache git make

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache for dependencies.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the application.
# We disable CGO for a static binary and target the api entry point.
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Start a new stage from a small Alpine image.
FROM alpine:latest

# Install ca-certificates for HTTPS requests.
RUN apk --no-cache add ca-certificates

# Set the working directory.
WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/main .

# Copy the .env file if it exists (optional, usually better to pass env vars at runtime)
# COPY .env .

# Expose the port the app runs on.
EXPOSE 8080

# Command to run the executable.
CMD ["./main"]

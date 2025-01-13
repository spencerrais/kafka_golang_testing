# Use the latest official Go image as a base, can hardcode the version when moving to production
FROM golang:latest AS builder

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kafka-consumer ./cmd/consumer

# Use a minimal image for runtime, can hardcode the version when moving to production
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder
COPY --from=builder /app/kafka-consumer .

# Expose the port for debugging or monitoring (optional)
EXPOSE 9094

# Run the binary
CMD ["./kafka-consumer"]

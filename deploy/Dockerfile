# Build stage
FROM golang:1.22-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o npm-download-spammer .

# Runtime stage
FROM alpine:latest

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/npm-download-spammer .
# Copy the config file
COPY --from=builder /app/npm-downloads-increaser.json .

# Set environment variables with default values
ENV NPM_PACKAGE_NAME=""
ENV NPM_NUM_DOWNLOADS=1000
ENV NPM_MAX_CONCURRENT_DOWNLOAD=300
ENV NPM_DOWNLOAD_TIMEOUT=3000

# Run the binary
ENTRYPOINT ["/app/npm-download-spammer"] 
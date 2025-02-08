# syntax=docker/dockerfile:1

# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Install dependencies required for the build
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /opt/wa_bot_service

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build binary
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /wa_bot_service_build

# Stage 2: Minimal runtime container
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /opt/wa_bot_service

# Copy the compiled binary from the builder stage
COPY --from=builder /wa_bot_service_build /wa_bot_service_build

# RUN go install github.com/air-verse/air@latest
# CMD ["air", "-d"]
CMD ["/wa_bot_service_build"]
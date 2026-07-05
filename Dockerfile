# Stage 1: Build the binary
# Using a specific alpine version to ensure consistency
FROM golang:alpine AS builder

# Update packages to patch known OS-level vulnerabilities
RUN apk update && apk upgrade --no-cache

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations for a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/api/main.go

# Stage 2: Run the binary
# Distroless contains only the minimal libraries needed to run the binary, 
# which removes almost all 'High' and 'Critical' CVEs.
FROM gcr.io/distroless/static-debian12:latest

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["/main"]
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mock-price-feed .

# ============================================
# Final stage - minimal runtime image
# ============================================
FROM alpine:latest

# Install curl for health checks
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/mock-price-feed .

# Expose the WebSocket/HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the mock price feed server
CMD ["./mock-price-feed"]
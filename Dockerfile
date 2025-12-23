FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/

# Download dependencies
RUN go mod download

# Copy source code
COPY services/api-gateway/ ./services/api-gateway/
COPY pkg/ ./pkg/

# Build the application
WORKDIR /app/services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/api-gateway/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]

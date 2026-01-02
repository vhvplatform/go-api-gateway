# Stage 1: Build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Cài đặt git nếu cần tải các module từ private repo
RUN apk add --no-cache git

# 1. Copy file go.mod/sum của cả hai để cache layer (Tăng tốc build)
# Lưu ý đường dẫn từ Context là thư mục gốc ~/workspace/go
COPY go-shared/go.mod go-shared/go.sum ./go-shared/
COPY go-api-gateway/go.mod go-api-gateway/go.sum ./go-api-gateway/

# 2. Download dependencies
# Đứng từ folder gateway để download vì nó chứa các logic dependency
WORKDIR /app/go-api-gateway
RUN go mod download

# 3. Copy mã nguồn của cả hai
WORKDIR /app
COPY go-shared/ ./go-shared/
COPY go-api-gateway/ ./go-api-gateway/

# 4. Build ứng dụng
WORKDIR /app/go-api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/bin/main ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata || true
WORKDIR /app
RUN chown 1000:1000 /app
COPY --from=builder /app/bin/main .
USER 1000
EXPOSE 8080
CMD ["./main"]

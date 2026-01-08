# Stage 1: Build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Cài đặt git nếu cần tải các module từ private repo
RUN apk add --no-cache git

# 1. Copy go.mod và go.sum để cache dependencies layer (Tăng tốc build)
COPY go.mod go.sum ./

# 2. Download dependencies
RUN go mod download

# 3. Copy toàn bộ mã nguồn
COPY . .

# 4. Build ứng dụng
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

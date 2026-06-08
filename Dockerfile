# ===== STAGE 1: BUILD =====
# Dùng image golang:1.26 để compile code
FROM golang:1.26-alpine AS builder

# Cài gcc để hỗ trợ race detector và CGO
RUN apk add --no-cache gcc musl-dev

# Tạo thư mục làm việc trong container
WORKDIR /app

# Copy go.mod và go.sum trước để cache dependencies
COPY go.mod go.sum ./

# Tải dependencies
RUN go mod download

# Copy toàn bộ source code vào container
COPY . .

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api/main.go

# Build Worker binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker ./cmd/worker/main.go


# ===== STAGE 2: API RUNNER =====
# Dùng image nhỏ alpine để chạy API
FROM alpine:latest AS api

WORKDIR /app

# Copy binary API từ stage builder
COPY --from=builder /app/bin/api .

# Copy file .env nếu có
COPY --from=builder /app/.env* ./

EXPOSE 8080

CMD ["./api"]


# ===== STAGE 3: WORKER RUNNER =====
# Dùng image nhỏ alpine để chạy Worker
FROM alpine:latest AS worker

WORKDIR /app

# Copy binary Worker từ stage builder
COPY --from=builder /app/bin/worker .

COPY --from=builder /app/.env* ./

CMD ["./worker"]
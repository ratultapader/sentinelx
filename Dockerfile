# -------- BUILD STAGE --------
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 🔥 REQUIRED FOR CGO (SQLite)
RUN apk add --no-cache git gcc musl-dev

# Enable CGO
ENV CGO_ENABLED=1

# Copy go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o sentinelx ./scripts/sentinelx_main.go


# -------- RUNTIME STAGE --------
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/sentinelx .

# Copy env + required files
COPY .env .
COPY data ./data

COPY rules ./rules

# 🔥 Create logs directory
RUN mkdir -p logs

# Expose API port
EXPOSE 9090

# Run app
CMD ["./sentinelx"]
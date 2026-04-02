# -------- BUILD STAGE --------
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o sentinelx ./scripts/sentinelx_main.go


# -------- RUNTIME STAGE --------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sentinelx .

COPY .env .

EXPOSE 9090

CMD ["./sentinelx"]
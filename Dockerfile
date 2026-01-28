# Multi-stage build for Alpine-compatible binary

# 1) Builder stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

# Enable Go modules explicitly (useful in CI)
ENV GO111MODULE=on

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build statically linked binary for Linux (Alpine compatible)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/json2xls ./cmd/json2xls

# 2) Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/json2xls /app/json2xls

EXPOSE 8080

ENTRYPOINT ["/app/json2xls"]

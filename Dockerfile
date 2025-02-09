# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /app

# Install required packages
RUN apk add --no-cache git

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate swagger docs
RUN /go/bin/swag init --parseDependency --parseInternal

# Build binary
RUN go build -o main .

# Runtime stage
FROM alpine:latest
WORKDIR /app

# Copy binary and swagger docs
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# เปิด port 8080 สำหรับรับ request
EXPOSE 8080

# คำสั่งเริ่มโปรแกรม
ENTRYPOINT ["./main"]

# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime Stage
FROM alpine:latest

RUN apk add --no-cache \
    libreoffice \
    poppler-utils \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /app/storage/uploads /app/storage/temp /app/storage/output && \
    chmod -R 755 /app/storage

EXPOSE 8080

ENV PORT=8080 \
    STORAGE_PATH=/app/storage \
    LOG_LEVEL=info \
    LOG_FORMAT=json

CMD ["./server"]

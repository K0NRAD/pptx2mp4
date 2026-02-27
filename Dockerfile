# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

ARG VITE_BASE_PATH=/pptx2mp4/

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN VITE_BASE_PATH=$VITE_BASE_PATH npm run build

# Stage 2: Build backend
FROM golang:alpine AS backend-builder

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
COPY --from=frontend-builder /app/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 3: Runtime
FROM alpine:latest

RUN apk add --no-cache \
    libreoffice \
    poppler-utils \
    ffmpeg \
    ca-certificates \
    font-liberation \
    font-dejavu \
    font-noto \
    font-noto-cjk \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=backend-builder /app/server .

RUN mkdir -p /app/storage/uploads /app/storage/temp /app/storage/output && \
    chmod -R 755 /app/storage

EXPOSE 8080

ENV PORT=8080 \
    STORAGE_PATH=/app/storage \
    LOG_LEVEL=info \
    LOG_FORMAT=json

CMD ["./server"]

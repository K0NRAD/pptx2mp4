# Stage 1: Build frontend (immer nativ auf dem Build-Host, kein QEMU)
FROM --platform=$BUILDPLATFORM node:20-alpine AS frontend-builder

ARG VITE_BASE_PATH=/pptx2mp4/

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN VITE_BASE_PATH=$VITE_BASE_PATH npm run build

# Stage 2: Build backend (nativ mit Go Cross-Compilation)
FROM --platform=$BUILDPLATFORM golang:alpine AS backend-builder

ARG TARGETARCH

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
COPY --from=frontend-builder /app/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o server ./cmd/server

# Stage 3: Runtime (Zielplattform)
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

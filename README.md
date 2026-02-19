# PPTX to MP4 Converter

Web-Anwendung zur Konvertierung von PowerPoint-Präsentationen (PPTX) in MP4-Videos.

## Features

- Web-basierter Upload von PPTX-Dateien
- Konvertierung in MP4-Video (statische Slides)
- Konfigurierbare Parameter:
  - FPS (Frames per Second): 1-60
  - Auflösung: 720p, 1080p, 1440p, 2160p
  - Dauer pro Slide: 1-60 Sekunden
- Asynchrone Verarbeitung mit Status-Tracking
- Download-Link für fertige Videos

## Technologie-Stack

### Backend
- **Go 1.23** - HTTP Server und Business-Logik
- **Gin** - Web Framework
- **LibreOffice** - PPTX zu PDF Konvertierung
- **Poppler** - PDF zu Bild-Sequenz
- **FFmpeg** - Video-Encoding

### Frontend
- **Svelte 5** - UI Framework
- **Vite** - Build Tool
- **TailwindCSS** - Styling

### Infrastructure
- **Docker** - Containerisierung
- **Docker Compose** - Multi-Container Orchestrierung
- **Nginx** - Frontend Serving

## Konvertierungs-Pipeline

```
PPTX Upload → LibreOffice → PDF → Poppler → PNG-Sequenz → FFmpeg → MP4
```

## Voraussetzungen

- Docker & Docker Compose
- (Für lokale Entwicklung: Go 1.23+, Node.js 20+)

## Installation & Start

### Mit Docker (Empfohlen)

```bash
# Repository klonen
git clone <repository-url>
cd pptx2mp4

# Container bauen und starten
docker-compose up --build

# Zugriff:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
```

### Lokale Entwicklung

#### Backend

```bash
cd backend
go mod download
go run cmd/server/main.go
```

**Benötigt:**
- LibreOffice: `brew install libreoffice` (macOS) / `apt install libreoffice` (Linux)
- Poppler: `brew install poppler` (macOS) / `apt install poppler-utils` (Linux)
- FFmpeg: `brew install ffmpeg` (macOS) / `apt install ffmpeg` (Linux)

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

## API-Dokumentation

### POST /api/v1/convert

Upload einer PPTX-Datei und Start der Konvertierung.

**Request:**
```
Content-Type: multipart/form-data

file: <PPTX-Datei>
fps: 24
resolution: 1080
duration: 5
```

**Response:**
```json
{
  "jobId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending"
}
```

### GET /api/v1/jobs/{jobId}/status

Status einer Konvertierung abfragen.

**Response:**
```json
{
  "jobId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing",
  "progress": 45,
  "error": null
}
```

Status-Werte: `pending`, `processing`, `completed`, `failed`

### GET /api/v1/jobs/{jobId}/download

Herunterladen des fertigen MP4-Videos.

**Response:** Binary MP4-Datei

### GET /api/v1/health

Health Check des Backend-Services.

**Response:**
```json
{
  "status": "ok",
  "libreoffice": true,
  "ffmpeg": true
}
```

## Projektstruktur

```
pptx2mp4/
├── backend/              # Go Backend
│   ├── cmd/
│   │   └── server/      # Entry Point
│   ├── internal/
│   │   ├── api/         # HTTP Layer
│   │   ├── domain/      # Business Entities
│   │   ├── service/     # Business Logic
│   │   ├── repository/  # Data Access
│   │   ├── converter/   # External Tool Integration
│   │   └── config/      # Configuration
│   └── storage/         # Temporäre Dateien
├── frontend/            # Svelte Frontend
│   └── src/
│       ├── lib/
│       │   ├── components/
│       │   ├── stores/
│       │   └── api/
│       └── routes/
└── docker-compose.yml
```

## Architektur

Das Backend folgt einer klassischen Layered Architecture mit strikter Separation of Concerns:

- **Domain Layer**: Pure Business Entities ohne externe Abhängigkeiten
- **Repository Layer**: Abstraktion der Datenpersistenz
- **Service Layer**: Business-Logik und Pipeline-Orchestrierung
- **API Layer**: HTTP Interface (Handlers, Middleware, Router)
- **Converter Layer**: Integration externer Tools (LibreOffice, FFmpeg)

Die Konvertierung läuft asynchron in Goroutines, Status-Updates erfolgen über Polling vom Frontend.

## Sicherheit

- MIME Type Validierung für Uploads
- File Extension Check (nur `.pptx`)
- File Size Limit (100MB)
- Filename Sanitization gegen Path Traversal
- Input Validation für alle Config-Parameter
- CORS Configuration
- Automatisches Cleanup alter Dateien nach 1 Stunde

## Lizenz

MIT

## Autor

K0NRAD

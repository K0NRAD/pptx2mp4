package api

import (
	"io/fs"
	"net/http"
	"strings"

	"pptx2mp4/backend/internal/api/handlers"
	"pptx2mp4/backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Router struct {
	engine          *gin.Engine
	uploadHandler   *handlers.UploadHandler
	statusHandler   *handlers.StatusHandler
	downloadHandler *handlers.DownloadHandler
	healthHandler   *handlers.HealthHandler
	logger          *logrus.Logger
	allowedOrigins  []string
	staticFiles     fs.FS
	basePath        string
}

func NewRouter(
	uploadHandler *handlers.UploadHandler,
	statusHandler *handlers.StatusHandler,
	downloadHandler *handlers.DownloadHandler,
	healthHandler *handlers.HealthHandler,
	logger *logrus.Logger,
	allowedOrigins []string,
	staticFiles fs.FS,
	basePath string,
) *Router {
	return &Router{
		uploadHandler:   uploadHandler,
		statusHandler:   statusHandler,
		downloadHandler: downloadHandler,
		healthHandler:   healthHandler,
		logger:          logger,
		allowedOrigins:  allowedOrigins,
		staticFiles:     staticFiles,
		basePath:        strings.TrimRight(basePath, "/"),
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine = gin.New()

	r.engine.Use(middleware.Recovery(r.logger))
	r.engine.Use(middleware.Logger(r.logger))
	if len(r.allowedOrigins) > 0 {
		r.engine.Use(middleware.SetupCORS(r.allowedOrigins))
	}
	r.engine.Use(middleware.ErrorHandler(r.logger))

	api := r.engine.Group(r.basePath + "/api/v1")
	{
		api.POST("/convert", r.uploadHandler.HandleUpload)
		api.GET("/jobs/:jobId/status", r.statusHandler.HandleStatus)
		api.GET("/jobs/:jobId/download", r.downloadHandler.HandleDownload)
		api.GET("/health", r.healthHandler.HandleHealth)
	}

	if r.staticFiles != nil {
		r.setupSPA()
	}

	return r.engine
}

func (r *Router) setupSPA() {
	subFS, err := fs.Sub(r.staticFiles, "dist")
	if err != nil {
		r.logger.WithError(err).Fatal("frontend-verzeichnis nicht gefunden")
	}

	indexHTML, err := fs.ReadFile(subFS, "index.html")
	if err != nil {
		r.logger.WithError(err).Fatal("index.html nicht gefunden")
	}

	fileServer := http.FileServer(http.FS(subFS))

	// Redirect / → /pptx2mp4/
	if r.basePath != "" {
		r.engine.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, r.basePath+"/")
		})
	}

	r.engine.NoRoute(func(c *gin.Context) {
		urlPath := c.Request.URL.Path

		if strings.HasPrefix(urlPath, r.basePath+"/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint nicht gefunden"})
			return
		}

		// Pfade außerhalb des base path → 404
		if r.basePath != "" && !strings.HasPrefix(urlPath, r.basePath) {
			c.JSON(http.StatusNotFound, gin.H{"error": "nicht gefunden"})
			return
		}

		// Base-Path-Prefix entfernen, um den Dateipfad zu ermitteln
		filePath := strings.TrimPrefix(urlPath, r.basePath)
		filePath = strings.TrimPrefix(filePath, "/")

		// Existierende statische Datei ausliefern (außer index.html wegen Redirect-Loop)
		if filePath != "" && filePath != "index.html" {
			if f, err := subFS.Open(filePath); err == nil {
				stat, statErr := f.Stat()
				f.Close()
				if statErr == nil && !stat.IsDir() {
					c.Request.URL.Path = "/" + filePath
					fileServer.ServeHTTP(c.Writer, c.Request)
					return
				}
			}
		}

		// SPA-Fallback: index.html direkt schreiben
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})
}

func (r *Router) Run(addr string) error {
	r.logger.WithField("address", addr).Info("starte HTTP Server")
	return r.engine.Run(addr)
}

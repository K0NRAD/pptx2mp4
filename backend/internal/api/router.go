package api

import (
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
}

func NewRouter(
	uploadHandler *handlers.UploadHandler,
	statusHandler *handlers.StatusHandler,
	downloadHandler *handlers.DownloadHandler,
	healthHandler *handlers.HealthHandler,
	logger *logrus.Logger,
	allowedOrigins []string,
) *Router {
	return &Router{
		uploadHandler:   uploadHandler,
		statusHandler:   statusHandler,
		downloadHandler: downloadHandler,
		healthHandler:   healthHandler,
		logger:          logger,
		allowedOrigins:  allowedOrigins,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine = gin.New()

	r.engine.Use(middleware.Recovery(r.logger))
	r.engine.Use(middleware.Logger(r.logger))
	r.engine.Use(middleware.SetupCORS(r.allowedOrigins))
	r.engine.Use(middleware.ErrorHandler(r.logger))

	api := r.engine.Group("/api/v1")
	{
		api.POST("/convert", r.uploadHandler.HandleUpload)
		api.GET("/jobs/:jobId/status", r.statusHandler.HandleStatus)
		api.GET("/jobs/:jobId/download", r.downloadHandler.HandleDownload)
		api.GET("/health", r.healthHandler.HandleHealth)
	}

	return r.engine
}

func (r *Router) Run(addr string) error {
	r.logger.WithField("address", addr).Info("starte HTTP Server")
	return r.engine.Run(addr)
}

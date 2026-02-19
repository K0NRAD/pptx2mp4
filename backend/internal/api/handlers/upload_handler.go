package handlers

import (
	"net/http"
	"pptx2mp4/backend/internal/domain"
	"pptx2mp4/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UploadHandler struct {
	fileService service.FileService
	jobService  service.JobService
	logger      *logrus.Logger
}

func NewUploadHandler(
	fileService service.FileService,
	jobService service.JobService,
	logger *logrus.Logger,
) *UploadHandler {
	return &UploadHandler{
		fileService: fileService,
		jobService:  jobService,
		logger:      logger,
	}
}

type ConvertRequest struct {
	FPS                int     `form:"fps" binding:"required,min=1,max=60"`
	Resolution         int     `form:"resolution" binding:"required,oneof=720 1080 1440 2160"`
	Duration           int     `form:"duration" binding:"required,min=1,max=60"`
	TransitionDuration float64 `form:"transitionDuration" binding:"min=0,max=3"`
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBind(&req); err != nil {
		h.logger.WithError(err).Error("ungültige Request-Parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validierungsfehler",
			"message": err.Error(),
		})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.logger.WithError(err).Error("keine Datei im Request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datei fehlt",
			"message": "Bitte laden Sie eine PPTX-Datei hoch",
		})
		return
	}

	if err := h.fileService.ValidateUpload(fileHeader); err != nil {
		h.logger.WithError(err).Warn("ungültige Datei hochgeladen")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ungültige Datei",
			"message": err.Error(),
		})
		return
	}

	config, err := domain.NewConversionConfig(req.FPS, req.Resolution, req.Duration, req.TransitionDuration)
	if err != nil {
		h.logger.WithError(err).Error("ungültige Konfiguration")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ungültige Konfiguration",
			"message": err.Error(),
		})
		return
	}

	job := domain.NewJob(fileHeader.Filename, config)

	if _, err := h.fileService.SaveUpload(job.ID, fileHeader); err != nil {
		h.logger.WithError(err).Error("fehler beim Speichern der Datei")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Speicherfehler",
			"message": "Datei konnte nicht gespeichert werden",
		})
		return
	}

	if err := h.jobService.CreateJob(job); err != nil {
		h.logger.WithError(err).Error("fehler beim Erstellen des Jobs")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Job-Erstellungsfehler",
			"message": "Job konnte nicht erstellt werden",
		})
		return
	}

	go func() {
		if err := h.jobService.ProcessJob(job.ID); err != nil {
			h.logger.WithError(err).WithField("jobID", job.ID).Error("Job-Verarbeitung fehlgeschlagen")
		}
	}()

	h.logger.WithFields(logrus.Fields{
		"jobID":    job.ID,
		"filename": fileHeader.Filename,
		"fps":      req.FPS,
		"resolution": req.Resolution,
		"duration": req.Duration,
	}).Info("Job erfolgreich erstellt")

	c.JSON(http.StatusAccepted, gin.H{
		"jobId":  job.ID,
		"status": job.Status,
	})
}

func parseIntParam(c *gin.Context, key string, defaultValue int) int {
	value := c.PostForm(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

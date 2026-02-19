package handlers

import (
	"net/http"
	"pptx2mp4/backend/internal/domain"
	"pptx2mp4/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DownloadHandler struct {
	jobService  service.JobService
	fileService service.FileService
	logger      *logrus.Logger
}

func NewDownloadHandler(
	jobService service.JobService,
	fileService service.FileService,
	logger *logrus.Logger,
) *DownloadHandler {
	return &DownloadHandler{
		jobService:  jobService,
		fileService: fileService,
		logger:      logger,
	}
}

func (h *DownloadHandler) HandleDownload(c *gin.Context) {
	jobID := c.Param("jobId")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ungültige Request",
			"message": "Job-ID fehlt",
		})
		return
	}

	job, err := h.jobService.GetJob(jobID)
	if err != nil {
		if err == domain.ErrJobNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Job nicht gefunden",
				"message": "Der angeforderte Job existiert nicht",
			})
			return
		}

		h.logger.WithError(err).Error("fehler beim Abrufen des Jobs")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Serverfehler",
			"message": "Job konnte nicht abgerufen werden",
		})
		return
	}

	if !job.IsCompleted() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Job nicht abgeschlossen",
			"message": "Die Konvertierung ist noch nicht abgeschlossen",
			"status":  job.Status,
		})
		return
	}

	outputFile, err := h.fileService.GetOutputFile(jobID)
	if err != nil {
		if err == domain.ErrFileNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Datei nicht gefunden",
				"message": "Die konvertierte Datei existiert nicht",
			})
			return
		}

		h.logger.WithError(err).Error("fehler beim Abrufen der Ausgabedatei")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Serverfehler",
			"message": "Datei konnte nicht abgerufen werden",
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"jobID":      jobID,
		"outputFile": outputFile,
	}).Info("starte Download")

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=output.mp4")
	c.Header("Content-Type", "video/mp4")
	c.File(outputFile)
}

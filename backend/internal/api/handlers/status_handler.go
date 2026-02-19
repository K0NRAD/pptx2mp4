package handlers

import (
	"net/http"
	"pptx2mp4/backend/internal/domain"
	"pptx2mp4/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StatusHandler struct {
	jobService service.JobService
	logger     *logrus.Logger
}

func NewStatusHandler(jobService service.JobService, logger *logrus.Logger) *StatusHandler {
	return &StatusHandler{
		jobService: jobService,
		logger:     logger,
	}
}

func (h *StatusHandler) HandleStatus(c *gin.Context) {
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

	response := gin.H{
		"jobId":    job.ID,
		"status":   job.Status,
		"progress": job.Progress,
	}

	if job.Error != "" {
		response["error"] = job.Error
	}

	c.JSON(http.StatusOK, response)
}

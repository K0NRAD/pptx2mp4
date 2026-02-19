package handlers

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HealthHandler struct {
	logger *logrus.Logger
}

func NewHealthHandler(logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

func (h *HealthHandler) HandleHealth(c *gin.Context) {
	libreOfficeAvailable := checkCommand("soffice", "--version")
	ffmpegAvailable := checkCommand("ffmpeg", "-version")
	popplerAvailable := checkCommand("pdftoppm", "-v")

	healthy := libreOfficeAvailable && ffmpegAvailable && popplerAvailable

	status := "ok"
	if !healthy {
		status = "degraded"
	}

	response := gin.H{
		"status":     status,
		"libreoffice": libreOfficeAvailable,
		"ffmpeg":     ffmpegAvailable,
		"poppler":    popplerAvailable,
	}

	statusCode := http.StatusOK
	if !healthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

func checkCommand(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	return err == nil || cmd.ProcessState != nil
}

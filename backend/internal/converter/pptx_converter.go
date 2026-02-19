package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"pptx2mp4/backend/internal/domain"

	"github.com/sirupsen/logrus"
)

type PPTXConverter interface {
	ConvertToPDF(inputPath, outputDir string) (string, error)
}

type LibreOfficeConverter struct {
	logger *logrus.Logger
}

func NewLibreOfficeConverter(logger *logrus.Logger) *LibreOfficeConverter {
	return &LibreOfficeConverter{
		logger: logger,
	}
}

func (c *LibreOfficeConverter) ConvertToPDF(inputPath, outputDir string) (string, error) {
	c.logger.WithFields(logrus.Fields{
		"input":     inputPath,
		"outputDir": outputDir,
	}).Info("starte PPTX zu PDF Konvertierung")

	cmd := exec.Command(
		"soffice",
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outputDir,
		inputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.WithError(err).WithField("output", string(output)).Error("PPTX zu PDF Konvertierung fehlgeschlagen")
		return "", fmt.Errorf("%w: %s", domain.ErrPPTXConversion, string(output))
	}

	outputFile := filepath.Join(outputDir, "input.pdf")

	c.logger.WithField("output", outputFile).Info("PPTX zu PDF Konvertierung erfolgreich")
	return outputFile, nil
}

func (c *LibreOfficeConverter) IsAvailable() bool {
	cmd := exec.Command("soffice", "--version")
	err := cmd.Run()
	return err == nil
}

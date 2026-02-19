package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pptx2mp4/backend/internal/domain"

	"github.com/sirupsen/logrus"
)

type PDFToImagesConverter interface {
	ConvertToImages(pdfPath, outputDir string, resolution int) ([]string, error)
}

type PopplerConverter struct {
	logger *logrus.Logger
}

func NewPopplerConverter(logger *logrus.Logger) *PopplerConverter {
	return &PopplerConverter{
		logger: logger,
	}
}

func (c *PopplerConverter) ConvertToImages(pdfPath, outputDir string, resolution int) ([]string, error) {
	c.logger.WithFields(logrus.Fields{
		"pdf":        pdfPath,
		"outputDir":  outputDir,
		"resolution": resolution,
	}).Info("starte PDF zu Bilder Konvertierung")

	outputPrefix := filepath.Join(outputDir, "slide")

	cmd := exec.Command(
		"pdftoppm",
		"-png",
		"-r", fmt.Sprintf("%d", resolution),
		pdfPath,
		outputPrefix,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.WithError(err).WithField("output", string(output)).Error("PDF zu Bilder Konvertierung fehlgeschlagen")
		return nil, fmt.Errorf("%w: %s", domain.ErrPDFConversion, string(output))
	}

	images, err := filepath.Glob(filepath.Join(outputDir, "slide-*.png"))
	if err != nil {
		c.logger.WithError(err).Error("fehler beim Suchen der generierten Bilder")
		return nil, fmt.Errorf("fehler beim Suchen der generierten Bilder: %w", err)
	}

	if len(images) == 0 {
		c.logger.Error("keine Bilder generiert")
		return nil, fmt.Errorf("%w: keine Bilder generiert", domain.ErrPDFConversion)
	}

	c.logger.WithField("imageCount", len(images)).Info("PDF zu Bilder Konvertierung erfolgreich")
	return images, nil
}

func (c *PopplerConverter) IsAvailable() bool {
	cmd := exec.Command("pdftoppm", "-v")
	err := cmd.Run()
	return err == nil || cmd.ProcessState != nil
}

func (c *PopplerConverter) GetSlideCount(pdfPath string) (int, error) {
	cmd := exec.Command("pdfinfo", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("fehler beim Abrufen der PDF-Informationen: %w", err)
	}

	var pages int
	_, err = fmt.Sscanf(string(output), "Pages: %d", &pages)
	if err != nil {
		return 0, fmt.Errorf("fehler beim Parsen der Seitenzahl: %w", err)
	}

	return pages, nil
}

func (c *PopplerConverter) CleanupImages(images []string) error {
	for _, img := range images {
		if err := os.Remove(img); err != nil {
			c.logger.WithError(err).WithField("image", img).Warn("fehler beim Löschen des Bildes")
		}
	}
	return nil
}

package service

import (
	"fmt"
	"path/filepath"
	"pptx2mp4/backend/internal/converter"
	"pptx2mp4/backend/internal/domain"
	"pptx2mp4/backend/internal/repository"

	"github.com/sirupsen/logrus"
)

type ConversionService interface {
	Convert(job *domain.Job) error
}

type ConversionServiceImpl struct {
	fileRepo        repository.FileRepository
	pptxConverter   converter.PPTXConverter
	pdfConverter    converter.PDFToImagesConverter
	videoEncoder    converter.VideoEncoder
	logger          *logrus.Logger
}

func NewConversionService(
	fileRepo repository.FileRepository,
	pptxConverter converter.PPTXConverter,
	pdfConverter converter.PDFToImagesConverter,
	videoEncoder converter.VideoEncoder,
	logger *logrus.Logger,
) *ConversionServiceImpl {
	return &ConversionServiceImpl{
		fileRepo:      fileRepo,
		pptxConverter: pptxConverter,
		pdfConverter:  pdfConverter,
		videoEncoder:  videoEncoder,
		logger:        logger,
	}
}

func (s *ConversionServiceImpl) Convert(job *domain.Job) error {
	s.logger.WithField("jobID", job.ID).Info("starte Konvertierungs-Pipeline")

	if err := s.fileRepo.EnsureDirectories(job.ID); err != nil {
		return fmt.Errorf("fehler beim Erstellen der Verzeichnisse: %w", err)
	}

	uploadPath := filepath.Join(s.fileRepo.GetUploadPath(job.ID), "input.pptx")
	tempPath := s.fileRepo.GetTempPath(job.ID)
	outputPath := s.fileRepo.GetOutputFilePath(job.ID)

	s.logger.WithFields(logrus.Fields{
		"jobID":      job.ID,
		"uploadPath": uploadPath,
		"tempPath":   tempPath,
		"outputPath": outputPath,
	}).Debug("Pfade konfiguriert")

	job.UpdateProgress(10)

	s.logger.WithField("jobID", job.ID).Info("schritt 1: PPTX zu PDF")
	pdfPath, err := s.pptxConverter.ConvertToPDF(uploadPath, tempPath)
	if err != nil {
		return fmt.Errorf("PPTX zu PDF Konvertierung fehlgeschlagen: %w", err)
	}
	job.UpdateProgress(40)

	s.logger.WithField("jobID", job.ID).Info("schritt 2: PDF zu Bilder")
	images, err := s.pdfConverter.ConvertToImages(pdfPath, tempPath, job.Config.Resolution)
	if err != nil {
		return fmt.Errorf("PDF zu Bilder Konvertierung fehlgeschlagen: %w", err)
	}
	job.UpdateProgress(70)

	s.logger.WithFields(logrus.Fields{
		"jobID":      job.ID,
		"imageCount": len(images),
	}).Info("schritt 3: Bilder zu Video")
	if err := s.videoEncoder.EncodeToMP4(tempPath, outputPath, job.Config); err != nil {
		return fmt.Errorf("video-encoding fehlgeschlagen: %w", err)
	}
	job.UpdateProgress(90)

	job.SetOutputFile(outputPath)

	s.logger.WithFields(logrus.Fields{
		"jobID":      job.ID,
		"outputPath": outputPath,
	}).Info("konvertierung erfolgreich abgeschlossen")

	return nil
}

func (s *ConversionServiceImpl) ValidateDependencies() error {
	if libreOffice, ok := s.pptxConverter.(*converter.LibreOfficeConverter); ok {
		if !libreOffice.IsAvailable() {
			return fmt.Errorf("LibreOffice ist nicht verfügbar")
		}
	}

	if poppler, ok := s.pdfConverter.(*converter.PopplerConverter); ok {
		if !poppler.IsAvailable() {
			return fmt.Errorf("Poppler (pdftoppm) ist nicht verfügbar")
		}
	}

	if ffmpeg, ok := s.videoEncoder.(*converter.FFmpegEncoder); ok {
		if !ffmpeg.IsAvailable() {
			return fmt.Errorf("FFmpeg ist nicht verfügbar")
		}
	}

	return nil
}

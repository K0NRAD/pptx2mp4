package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"pptx2mp4/backend/internal/domain"
	"pptx2mp4/backend/internal/repository"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	MaxFileSize   = 100 * 1024 * 1024 // 100MB
	AllowedMIME   = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	AllowedExtension = ".pptx"
)

type FileService interface {
	ValidateUpload(fileHeader *multipart.FileHeader) error
	SaveUpload(jobID string, fileHeader *multipart.FileHeader) (string, error)
	GetOutputFile(jobID string) (string, error)
	SanitizeFilename(filename string) string
}

type FileServiceImpl struct {
	fileRepo repository.FileRepository
	logger   *logrus.Logger
}

func NewFileService(fileRepo repository.FileRepository, logger *logrus.Logger) *FileServiceImpl {
	return &FileServiceImpl{
		fileRepo: fileRepo,
		logger:   logger,
	}
}

func (s *FileServiceImpl) ValidateUpload(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > MaxFileSize {
		return domain.ErrFileTooLarge
	}

	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), AllowedExtension) {
		return domain.ErrInvalidExtension
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("fehler beim Öffnen der Datei: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("fehler beim Lesen der Datei: %w", err)
	}

	mimeType := http.DetectContentType(buffer)
	if !strings.Contains(mimeType, "application") && !strings.Contains(mimeType, "zip") {
		return domain.ErrInvalidMimeType
	}

	return nil
}

func (s *FileServiceImpl) SaveUpload(jobID string, fileHeader *multipart.FileHeader) (string, error) {
	s.logger.WithFields(logrus.Fields{
		"jobID":    jobID,
		"filename": fileHeader.Filename,
		"size":     fileHeader.Size,
	}).Info("speichere Upload")

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("fehler beim Öffnen der Datei: %w", err)
	}
	defer file.Close()

	sanitizedFilename := s.SanitizeFilename(fileHeader.Filename)
	filePath, err := s.fileRepo.SaveUpload(jobID, file, sanitizedFilename)
	if err != nil {
		s.logger.WithError(err).Error("fehler beim Speichern der Datei")
		return "", err
	}

	s.logger.WithField("filePath", filePath).Info("Upload erfolgreich gespeichert")
	return filePath, nil
}

func (s *FileServiceImpl) GetOutputFile(jobID string) (string, error) {
	outputPath := s.fileRepo.GetOutputFilePath(jobID)

	if !s.fileRepo.FileExists(outputPath) {
		return "", domain.ErrFileNotFound
	}

	return outputPath, nil
}

func (s *FileServiceImpl) SanitizeFilename(filename string) string {
	filename = filepath.Base(filename)

	reg := regexp.MustCompile(`[^a-zA-Z0-9.\-_]`)
	sanitized := reg.ReplaceAllString(filename, "_")

	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized
}

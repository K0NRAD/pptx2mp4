package repository

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"pptx2mp4/backend/internal/domain"
)

type FileRepository interface {
	SaveUpload(jobID string, file multipart.File, filename string) (string, error)
	GetUploadPath(jobID string) string
	GetTempPath(jobID string) string
	GetOutputPath(jobID string) string
	GetOutputFilePath(jobID string) string
	FileExists(path string) bool
	EnsureDirectories(jobID string) error
	CleanupJob(jobID string) error
}

type FileSystemRepository struct {
	basePath string
}

func NewFileSystemRepository(basePath string) *FileSystemRepository {
	return &FileSystemRepository{
		basePath: basePath,
	}
}

func (r *FileSystemRepository) SaveUpload(jobID string, file multipart.File, filename string) (string, error) {
	uploadDir := r.GetUploadPath(jobID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("fehler beim Erstellen des Upload-Verzeichnisses: %w", err)
	}

	destPath := filepath.Join(uploadDir, "input.pptx")
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("fehler beim Erstellen der Zieldatei: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("fehler beim Kopieren der Datei: %w", err)
	}

	return destPath, nil
}

func (r *FileSystemRepository) GetUploadPath(jobID string) string {
	return filepath.Join(r.basePath, "uploads", jobID)
}

func (r *FileSystemRepository) GetTempPath(jobID string) string {
	return filepath.Join(r.basePath, "temp", jobID)
}

func (r *FileSystemRepository) GetOutputPath(jobID string) string {
	return filepath.Join(r.basePath, "output", jobID)
}

func (r *FileSystemRepository) GetOutputFilePath(jobID string) string {
	return filepath.Join(r.GetOutputPath(jobID), "output.mp4")
}

func (r *FileSystemRepository) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (r *FileSystemRepository) EnsureDirectories(jobID string) error {
	dirs := []string{
		r.GetUploadPath(jobID),
		r.GetTempPath(jobID),
		r.GetOutputPath(jobID),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("fehler beim Erstellen des Verzeichnisses %s: %w", dir, err)
		}
	}

	return nil
}

func (r *FileSystemRepository) CleanupJob(jobID string) error {
	paths := []string{
		r.GetUploadPath(jobID),
		r.GetTempPath(jobID),
		r.GetOutputPath(jobID),
	}

	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("fehler beim Löschen des Verzeichnisses %s: %w", path, err)
		}
	}

	return nil
}

func (r *FileSystemRepository) ValidateStoragePath() error {
	if r.basePath == "" {
		return domain.ErrStoragePathInvalid
	}

	info, err := os.Stat(r.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(r.basePath, 0755); err != nil {
				return fmt.Errorf("storage-pfad kann nicht erstellt werden: %w", err)
			}
			return nil
		}
		return fmt.Errorf("storage-pfad konnte nicht überprüft werden: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("storage-pfad ist kein Verzeichnis: %s", r.basePath)
	}

	return nil
}

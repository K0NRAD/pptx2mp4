package main

import (
	"fmt"
	"os"
	"pptx2mp4/backend/internal/api"
	"pptx2mp4/backend/internal/api/handlers"
	"pptx2mp4/backend/internal/config"
	"pptx2mp4/backend/internal/converter"
	"pptx2mp4/backend/internal/repository"
	"pptx2mp4/backend/internal/service"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := setupLogger()

	logger.Info("starte PPTX to MP4 Converter Server")

	cfg := config.LoadConfig()

	if err := validateConfig(cfg); err != nil {
		logger.WithError(err).Fatal("ungültige Konfiguration")
	}

	logger.WithFields(logrus.Fields{
		"port":        cfg.Port,
		"storagePath": cfg.StoragePath,
		"logLevel":    cfg.LogLevel,
	}).Info("konfiguration geladen")

	jobRepo := repository.NewInMemoryJobRepository()
	logger.Info("job-repository initialisiert")

	fileRepo := repository.NewFileSystemRepository(cfg.StoragePath)
	if err := fileRepo.ValidateStoragePath(); err != nil {
		logger.WithError(err).Fatal("ungültiger Storage-Pfad")
	}
	logger.Info("file-repository initialisiert")

	pptxConverter := converter.NewLibreOfficeConverter(logger)
	pdfConverter := converter.NewPopplerConverter(logger)
	videoEncoder := converter.NewFFmpegEncoder(logger)
	logger.Info("converter initialisiert")

	conversionService := service.NewConversionService(
		fileRepo,
		pptxConverter,
		pdfConverter,
		videoEncoder,
		logger,
	)

	if err := conversionService.ValidateDependencies(); err != nil {
		logger.WithError(err).Fatal("externe Abhängigkeiten nicht verfügbar")
	}
	logger.Info("externe Abhängigkeiten validiert")

	fileService := service.NewFileService(fileRepo, logger)
	jobService := service.NewJobService(jobRepo, conversionService, logger)
	logger.Info("services initialisiert")

	uploadHandler := handlers.NewUploadHandler(fileService, jobService, logger)
	statusHandler := handlers.NewStatusHandler(jobService, logger)
	downloadHandler := handlers.NewDownloadHandler(jobService, fileService, logger)
	healthHandler := handlers.NewHealthHandler(logger)
	logger.Info("handlers initialisiert")

	router := api.NewRouter(
		uploadHandler,
		statusHandler,
		downloadHandler,
		healthHandler,
		logger,
		cfg.AllowedOrigins,
	)

	engine := router.Setup()
	logger.Info("router konfiguriert")

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.WithField("address", addr).Info("server bereit")

	if err := engine.Run(addr); err != nil {
		logger.WithError(err).Fatal("server konnte nicht gestartet werden")
	}
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return logger
}

func validateConfig(cfg *config.Config) error {
	if cfg.Port == "" {
		return fmt.Errorf("port darf nicht leer sein")
	}

	if cfg.StoragePath == "" {
		return fmt.Errorf("storage-pfad darf nicht leer sein")
	}

	if cfg.MaxFileSize <= 0 {
		return fmt.Errorf("max-file-size muss größer als 0 sein")
	}

	return nil
}

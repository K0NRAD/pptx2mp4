package domain

import "errors"

var (
	ErrJobNotFound        = errors.New("job nicht gefunden")
	ErrInvalidFile        = errors.New("ungültiges Dateiformat")
	ErrFileTooLarge       = errors.New("datei zu groß")
	ErrInvalidMimeType    = errors.New("ungültiger MIME-Type")
	ErrInvalidExtension   = errors.New("ungültige Dateierweiterung")
	ErrConversionFailed   = errors.New("konvertierung fehlgeschlagen")
	ErrPPTXConversion     = errors.New("PPTX zu PDF Konvertierung fehlgeschlagen")
	ErrPDFConversion      = errors.New("PDF zu Bilder Konvertierung fehlgeschlagen")
	ErrVideoEncoding      = errors.New("video-encoding fehlgeschlagen")
	ErrInvalidConfig      = errors.New("ungültige Konfiguration")
	ErrJobAlreadyExists   = errors.New("job existiert bereits")
	ErrFileNotFound       = errors.New("datei nicht gefunden")
	ErrInvalidJobStatus   = errors.New("ungültiger Job-Status")
	ErrStoragePathInvalid = errors.New("ungültiger Speicherpfad")
)

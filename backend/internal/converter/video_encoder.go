package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"pptx2mp4/backend/internal/domain"

	"github.com/sirupsen/logrus"
)

type VideoEncoder interface {
	EncodeToMP4(imagesDir, outputPath string, config *domain.ConversionConfig) error
}

type FFmpegEncoder struct {
	logger *logrus.Logger
}

func NewFFmpegEncoder(logger *logrus.Logger) *FFmpegEncoder {
	return &FFmpegEncoder{
		logger: logger,
	}
}

func (e *FFmpegEncoder) EncodeToMP4(imagesDir, outputPath string, config *domain.ConversionConfig) error {
	e.logger.WithFields(logrus.Fields{
		"imagesDir":  imagesDir,
		"outputPath": outputPath,
		"fps":        config.FPS,
		"duration":   config.Duration,
	}).Info("starte Video-Encoding")

	framerate := 1.0 / float64(config.Duration)

	pattern := filepath.Join(imagesDir, "slide-*.png")

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-framerate", fmt.Sprintf("%.4f", framerate),
		"-pattern_type", "glob",
		"-i", pattern,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-r", fmt.Sprintf("%d", config.FPS),
		"-vf", fmt.Sprintf("scale=-2:%d", config.Resolution),
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		e.logger.WithError(err).WithField("output", string(output)).Error("video-encoding fehlgeschlagen")
		return fmt.Errorf("%w: %s", domain.ErrVideoEncoding, string(output))
	}

	e.logger.WithField("output", outputPath).Info("video-encoding erfolgreich")
	return nil
}

func (e *FFmpegEncoder) IsAvailable() bool {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	return err == nil
}

func (e *FFmpegEncoder) GetVideoInfo(videoPath string) (map[string]interface{}, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		videoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abrufen der Video-Informationen: %w", err)
	}

	info := make(map[string]interface{})
	info["raw"] = string(output)

	return info, nil
}

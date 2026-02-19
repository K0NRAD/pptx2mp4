package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"pptx2mp4/backend/internal/domain"
	"sort"
	"strings"

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
		"imagesDir":          imagesDir,
		"outputPath":         outputPath,
		"fps":                config.FPS,
		"duration":           config.Duration,
		"transitionDuration": config.TransitionDuration,
	}).Info("starte Video-Encoding")

	images, err := filepath.Glob(filepath.Join(imagesDir, "slide-*.png"))
	if err != nil || len(images) == 0 {
		return fmt.Errorf("%w: keine Slide-Bilder gefunden in %s", domain.ErrVideoEncoding, imagesDir)
	}

	sort.Slice(images, func(i, j int) bool {
		var ni, nj int
		fmt.Sscanf(filepath.Base(images[i]), "slide-%d.png", &ni)
		fmt.Sscanf(filepath.Base(images[j]), "slide-%d.png", &nj)
		return ni < nj
	})

	N := len(images)
	D := float64(config.Duration)
	T := config.TransitionDuration

	args := []string{"-y"}
	for _, img := range images {
		args = append(args, "-loop", "1", "-t", fmt.Sprintf("%.4f", D), "-i", img)
	}

	var filterParts []string
	for i := range images {
		filterParts = append(filterParts,
			fmt.Sprintf("[%d:v]scale=-2:%d,fps=%d,format=yuv420p[v%d]", i, config.Resolution, config.FPS, i))
	}

	lastLabel := "v0"
	if N > 1 && T > 0 {
		for i := 1; i < N; i++ {
			offset := float64(i) * (D - T)
			next := fmt.Sprintf("x%d", i)
			filterParts = append(filterParts,
				fmt.Sprintf("[%s][v%d]xfade=transition=fade:duration=%.4f:offset=%.4f[%s]",
					lastLabel, i, T, offset, next))
			lastLabel = next
		}
	} else if N > 1 {
		var concatInputs string
		for i := range images {
			concatInputs += fmt.Sprintf("[v%d]", i)
		}
		filterParts = append(filterParts,
			fmt.Sprintf("%sconcat=n=%d:v=1:a=0[out]", concatInputs, N))
		lastLabel = "out"
	}

	args = append(args, "-filter_complex", strings.Join(filterParts, ";"))
	args = append(args, "-map", fmt.Sprintf("[%s]", lastLabel))
	args = append(args, "-c:v", "libx264", "-pix_fmt", "yuv420p", outputPath)

	cmd := exec.Command("ffmpeg", args...)
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

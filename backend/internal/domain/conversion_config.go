package domain

type ConversionConfig struct {
	FPS                int     `json:"fps" binding:"required,min=1,max=60"`
	Resolution         int     `json:"resolution" binding:"required,oneof=720 1080 1440 2160"`
	Duration           int     `json:"duration" binding:"required,min=1,max=60"`
	TransitionDuration float64 `json:"transitionDuration"`
}

func NewConversionConfig(fps, resolution, duration int, transitionDuration float64) (*ConversionConfig, error) {
	config := &ConversionConfig{
		FPS:                fps,
		Resolution:         resolution,
		Duration:           duration,
		TransitionDuration: transitionDuration,
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *ConversionConfig) Validate() error {
	if c.FPS < 1 || c.FPS > 60 {
		return ErrInvalidConfig
	}

	if c.Resolution != 720 && c.Resolution != 1080 && c.Resolution != 1440 && c.Resolution != 2160 {
		return ErrInvalidConfig
	}

	if c.Duration < 1 || c.Duration > 60 {
		return ErrInvalidConfig
	}

	if c.TransitionDuration < 0 || c.TransitionDuration >= float64(c.Duration) {
		return ErrInvalidConfig
	}

	return nil
}

func DefaultConfig() *ConversionConfig {
	return &ConversionConfig{
		FPS:                24,
		Resolution:         1080,
		Duration:           5,
		TransitionDuration: 1.0,
	}
}

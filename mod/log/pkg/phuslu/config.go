package phuslu

import (
	"time"
)

// DefaultConfig returns a new Config with default values.
func DefaultConfig() *Config {
	return &Config{
		TimeFormat:     time.Kitchen,
		ColorOutput:    true,
		QuoteString:    true,
		EndWithMessage: true,
	}
}

// Config defines configuration for the logger.
type Config struct {
	TimeFormat     string
	ColorOutput    bool
	QuoteString    bool
	EndWithMessage bool
}

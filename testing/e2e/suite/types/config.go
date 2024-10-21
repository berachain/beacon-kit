package types

import (
	"errors"
	"net/url"
	"time"
)

type Config struct {
	ApiURL      string
	HttpTimeout time.Duration
	MaxRetries  int
}

func (c *Config) Validate() error {
	if _, err := url.Parse(c.ApiURL); err != nil {
		return err
	}

	if c.HttpTimeout <= 0 {
		return errors.New("http timeout must be greater than 0")
	}
	if c.MaxRetries < 0 {
		return errors.New("max retries must be greater than or equal to 0")
	}

	return nil
}

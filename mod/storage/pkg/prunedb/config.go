package prunedb

import "time"

const (
	// defaultPruneInterval is the default interval at which the pruner will run.
	defaultPruneInterval = 10 * time.Second
)

// Config is the configuration for the pruner.
type Config struct {
	// PruneInterval is the interval at which the pruner will run.
	PruneInterval time.Duration `mapstructure:"prune-interval"`
}

// DefaultConfig returns the default configuration for the pruner.
func DefaultConfig() Config {
	return Config{
		PruneInterval: defaultPruneInterval,
	}
}

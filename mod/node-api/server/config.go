package server

// Config is the configuration for the node API server.
type Config struct {
	// Enabled is the flag to enable the node API server.
	Enabled bool `mapstructure:"enabled"`
	// Address is the address to bind the node API server to.
	Address string `mapstructure:"address"`
}

// DefaultConfig returns the default configuration for the node API server.
func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Address: "0.0.0.0:8080",
	}
}

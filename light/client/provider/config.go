package provider

type Config struct {
	ChainID string

	HttpEndpoint string
	WSEndpoint   string
}

func NewConfig(chainID, httpEndpoint, wsEndpoint string) *Config {
	return &Config{
		ChainID:      chainID,
		HttpEndpoint: httpEndpoint,
		WSEndpoint:   wsEndpoint,
	}
}

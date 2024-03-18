package config

const (
	// defaultRPCPort is the default value for the RPC port.
	defaultRPCPort = 4000

	// defaultRPCHost is the default value for the RPC host.
	defaultRPCHost = "127.0.0.1"
)

type RPC struct {
	// Enabled determines if the RPC service is enabled.
	Enabled bool
	// Host is the host of the RPC service.
	Host string
	// Port is the port of the RPC service.
	Port int
	// CertFlag is the flag for the certificate file.
	CertFlag string
	// KeyFlag is the flag for the key file.
	KeyFlag string
}

func DefaultRPCConfig() RPC {
	return RPC{
		Enabled: true,
		Host:    defaultRPCHost,
		Port:    defaultRPCPort,
	}
}

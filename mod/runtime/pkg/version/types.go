package version

// TelemetrySink is an interface for sending telemetry data.
type TelemetrySink interface {
	// IncrementCounter increments a counter metric identified by the provided
	// keys.
	IncrementCounter(key string, args ...string)
}

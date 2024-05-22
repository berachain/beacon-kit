package metrics

import "github.com/cosmos/cosmos-sdk/telemetry"

type TelemetrySink struct{}

// NewTelemetrySink creates a new TelemetrySink.
func NewTelemetrySink() TelemetrySink {
	return TelemetrySink{}
}

// IncrementCounter increments a counter metric identified by the provided
// keys.
func (TelemetrySink) IncrementCounter(keys ...string) {
	telemetry.IncrCounter(1, keys...)
}

// SetGauge sets a gauge metric to the specified value, identified by the
// provided keys.
func (TelemetrySink) SetGauge(value int64, keys ...string) {
	telemetry.SetGauge(float32(value), keys...)
}

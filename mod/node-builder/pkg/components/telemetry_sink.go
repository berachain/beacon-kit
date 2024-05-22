package components

import "github.com/berachain/beacon-kit/mod/node-builder/pkg/components/metrics"

// ProvideTelemetrySink is a function that provides a TelemetrySink.
func ProvideTelemetrySink() *metrics.TelemetrySink {
	return &metrics.TelemetrySink{}
}

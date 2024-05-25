package blockchain

import (
	"time"
)

// chainMetrics is a struct that contains metrics for the chain.
type chainMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newChainMetrics creates a new chainMetrics.
func newChainMetrics(
	sink TelemetrySink,
) *chainMetrics {
	return &chainMetrics{
		sink: sink,
	}
}

// measureStateTransitionDuration measures the time to process
// the state transition for a block.
func (cm *chainMetrics) measureStateTransitionDuration(start time.Time) {
	cm.sink.MeasureSince(
		"beacon_kit.beacon.blockchain.state_transition_duration", start,
	)
}

// measureBlobProcessingDuration measures the time to process
// the blobs for a block.
func (cm *chainMetrics) measureBlobProcessingDuration(start time.Time) {
	cm.sink.MeasureSince(
		"beacon_kit.beacon.blockchain.blob_processing_duration", start,
	)
}

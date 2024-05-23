package blob

import (
	"time"
)

// processorMetrics is a struct that contains metrics for the processor.
type processorMetrics struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
}

// newProcessorMetrics creates a new processorMetrics.
func newProcessorMetrics(
	sink TelemetrySink,
) *processorMetrics {
	return &processorMetrics{
		sink: sink,
	}
}

// MeasureProcessBlobDuration measures the duration of the blob processing.
func (pm *processorMetrics) measureProcessBlobsDuration(startTime time.Time) {
	// TODO: Add NumBlobs
	pm.sink.MeasureSince(
		"beacon_kit.da.blob.processor.process_blob_duration",
		startTime,
	)
}

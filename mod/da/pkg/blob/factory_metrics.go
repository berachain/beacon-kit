package blob

import (
	"time"
)

// factoryMetrics is a struct that contains metrics for the factory.
type factoryMetrics struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
}

// newFactoryMetrics creates a new factoryMetrics.
func newFactoryMetrics(
	sink TelemetrySink,
) *factoryMetrics {
	return &factoryMetrics{
		sink: sink,
	}
}

// measureBuildSidecarDuration measures the duration of the build sidecar.
func (fm *factoryMetrics) measureBuildSidecarDuration(startTime time.Time) {
	// TODO: Add Labels.
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_sidecar_duration",
		startTime,
	)
}

// measureBuildKZGInclusionProofDuration measures the duration of the build KZG inclusion proof.
func (fm *factoryMetrics) measureBuildKZGInclusionProofDuration(startTime time.Time) {
	// TODO: Add Labels.
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_kzg_inclusion_proof_duration",
		startTime,
	)
}

// measureBuildBlockBodyProofDuration measures the duration of the build block body proof.
func (fm *factoryMetrics) measureBuildBlockBodyProofDuration(startTime time.Time) {
	// TODO: Add Labels.
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_block_body_proof_duration",
		startTime,
	)
}

// measureBuildCommitmentProofDuration measures the duration of the build commitment proof.
func (fm *factoryMetrics) measureBuildCommitmentProofDuration(startTime time.Time) {
	// TODO: Add Labels.
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_commitment_proof_duration",
		startTime,
	)
}

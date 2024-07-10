package engineprimitives

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ExecutionPayloadEnvelope is an envelope for execution payloads.
type ExecutionPayloadEnvelope[
	T any,
	BlobsBundleT any,
	ExecutionPayloadT any,
] interface {
	constraints.Nillable
	New(
		blobsBundle BlobsBundleT,
		payload ExecutionPayloadT,
	) T
	// GetExecutionPayload retrieves the associated execution payload.
	GetExecutionPayload() ExecutionPayloadT
	// GetValue returns the Wei value of the block in the execution payload.
	GetValue() math.Wei
	// GetBlobsBundle fetches the associated BlobsBundleV1 if available.
	GetBlobsBundle() BlobsBundleT
	// ShouldOverrideBuilder indicates if the builder should be overridden.
	ShouldOverrideBuilder() bool
}

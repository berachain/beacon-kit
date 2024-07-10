package store

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[
	BeaconBlockBodyT any,
	BlobSidecarsT any,
] interface {
	// IsDataAvailable checks if the data is available for the given slot.
	IsDataAvailable(
		ctx context.Context,
		slot math.Slot,
		body BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(slot math.Slot, sidecars BlobSidecarsT) error
}

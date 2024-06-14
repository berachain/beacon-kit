package middleware

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// ABCIMiddleware is a middleware between ABCI and Beacon logic.
type ABCIMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT interface {
		types.RawBeaconBlock[BeaconBlockBodyT]
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
		Empty(uint32) BeaconBlockT
	},
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconStateT BeaconState,
	BlobSidecarsT ssz.Marshallable,
	StorageBackendT any,
] struct {
	FinalizeBlock *FinalizeBlockMiddleware[
		BeaconBlockT, BeaconStateT, BlobSidecarsT,
	]

	Validator *ValidatorMiddleware[
		AvailabilityStoreT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconStateT,
		BlobSidecarsT,
		StorageBackendT,
	]
}

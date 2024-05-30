package beacondb

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// SSZMarshallable is an interface that combines the ssz.Marshaler and
// ssz.Unmarshaler interfaces.
type SSZMarshallable interface {
	// MarshalSSZTo marshals the object into the provided byte slice and returns
	// it along with any error.
	MarshalSSZTo([]byte) ([]byte, error)
	// MarshalSSZ marshals the object into a new byte slice and returns it along
	// with any error.
	MarshalSSZ() ([]byte, error)
	// UnmarshalSSZ unmarshals the object from the provided byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalSSZ([]byte) error
	// SizeSSZ returns the size in bytes that the object would take when
	// marshaled.
	SizeSSZ() int
}

// Validator represents an interface for a validator in the beacon chain.
type Validator interface {
	// SSZMarshallable is embedded to ensure that the validator can be marshaled
	// and unmarshaled using SSZ (Simple Serialize) encoding.
	SSZMarshallable
	// GetPubkey returns the BLS public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetEffectiveBalance returns the effective balance of the validator in Gwei.
	GetEffectiveBalance() math.Gwei
	// IsActive checks if the validator is active at the given epoch.
	IsActive(epoch math.Epoch) bool
}

package payload

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// AttributesFactory is a factory for creating payload attributes.
type AttributesFactory[
	BeaconStateT any,
	PayloadAttributesT engineprimitives.PayloadAttributes[
		PayloadAttributesT, WithdrawalT,
	],
	WithdrawalT any,
] interface {
	// BuildPayloadAttributes builds a new instance of PayloadAttributes.
	BuildPayloadAttributes(
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		prevHeadRoot [32]byte,
	) (PayloadAttributesT, error)
}

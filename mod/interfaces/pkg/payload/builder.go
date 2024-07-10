package payload

import (
	"context"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	engineprimitives "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Builder is used to build payloads on the execution client.
type Builder[
	// BeaconStateT BeaconState[ExecutionPayloadHeaderT, WithdrawalT],
	BeaconStateT any,
	BlobsBundleT any,
	ExecutionPayloadT types.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	PayloadAttributesT engineprimitives.PayloadAttributes[
		PayloadAttributesT, WithdrawalT,
	],
	PayloadIDT ~[8]byte,
	WithdrawalT any,
] interface {
	Enabled() bool
	// RetrievePayload retrieves the payload for the given slot.
	RetrievePayload(
		ctx context.Context,
		slot math.Slot,
		parentBlockRoot common.Root,
	) (ExecutionPayloadEnvelopeT, error)
	// RequestPayloadAsync requests a payload for the given slot and
	// returns the payload ID without blocking.
	RequestPayloadAsync(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot common.Root,
		headEth1BlockHash gethprimitives.ExecutionHash,
		finalEth1BlockHash gethprimitives.ExecutionHash,
	) (*PayloadIDT, error)
	// RequestPayloadSync requests a payload for the given slot and
	// blocks until the payload is delivered.
	RequestPayloadSync(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		timestamp uint64,
		parentBlockRoot common.Root,
		headEth1BlockHash gethprimitives.ExecutionHash,
		finalEth1BlockHash gethprimitives.ExecutionHash,
	) (ExecutionPayloadEnvelopeT, error)
	// SendForceHeadFCU sends a force head FCU to the execution client.
	SendForceHeadFCU(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
	) error
}

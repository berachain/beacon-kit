package execution

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	engineprimitivesI "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
)

// Engine is Beacon-Kit's interface for the `ExecutionEngine`
// from the Ethereum 2.0 Specification.
type Engine[
	BlobsBundleT any,
	ExecutionPayloadT types.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadEnvelopeT engineprimitivesI.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	PayloadAttributesT engineprimitivesI.PayloadAttributes[
		PayloadAttributesT, WithdrawalT,
	],
	PayloadIDT ~[8]byte,
	WithdrawalT engineprimitivesI.Withdrawal[WithdrawalT],
] interface {
	// GetPayload returns the payload and blobs bundle for the given slot.
	GetPayload(
		ctx context.Context,
		req *engineprimitives.GetPayloadRequest[PayloadIDT],
	) (ExecutionPayloadEnvelopeT, error)
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *engineprimitives.ForkchoiceUpdateRequest[PayloadAttributesT],
	) (*PayloadIDT, *gethprimitives.ExecutionHash, error)
	// VerifyAndNotifyNewPayload verifies the new payload and notifies the
	// execution client.
	VerifyAndNotifyNewPayload(
		ctx context.Context,
		req *engineprimitives.NewPayloadRequest[
			ExecutionPayloadT, WithdrawalT,
		],
	) error
}

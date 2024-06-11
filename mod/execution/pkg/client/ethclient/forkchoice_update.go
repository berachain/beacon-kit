package ethclient

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// ForkchoiceUpdated is a helper function to call the appropriate version of the
func (s *Eth1Client[ExecutionPayloadT]) ForkchoiceUpdated(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs engineprimitives.PayloadAttributer,
	forkVersion uint32,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	switch forkVersion {
	case version.Deneb:
		return s.ForkchoiceUpdatedV3(ctx, state, attrs)
	default:
		return nil, ErrInvalidVersion
	}
}

// ForkchoiceUpdatedV3 calls the engine_forkchoiceUpdatedV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) ForkchoiceUpdatedV3(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs engineprimitives.PayloadAttributer,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	return s.forkchoiceUpdated(ctx, ForkchoiceUpdatedMethodV3, state, attrs)
}

// forkchoiceUpdateCall is a helper function to call to any version
// of the forkchoiceUpdates method.
func (s *Eth1Client[ExecutionPayloadT]) forkchoiceUpdated(
	ctx context.Context,
	method string,
	state *engineprimitives.ForkchoiceStateV1,
	attrs any,
) (*engineprimitives.ForkchoiceResponseV1, error) {
	result := &engineprimitives.ForkchoiceResponseV1{}

	if err := s.Client.Client().CallContext(
		ctx, result, method, state, attrs,
	); err != nil {
		return nil, err
	}

	if (result.PayloadStatus == engineprimitives.PayloadStatusV1{}) {
		return nil, ErrNilResponse
	}

	return result, nil
}

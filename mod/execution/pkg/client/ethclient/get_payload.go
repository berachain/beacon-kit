package ethclient

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// GetPayload is a helper function to call the appropriate version of the
// engine_getPayload method.
func (s *Eth1Client[ExecutionPayloadT]) GetPayload(
	ctx context.Context,
	forkVersion uint32,
	payloadID engineprimitives.PayloadID,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	switch forkVersion {
	case version.Deneb:
		return s.GetPayloadV3(ctx, payloadID)
	default:
		return nil, ErrInvalidVersion
	}
}

// GetPayloadV3 calls the engine_getPayloadV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) GetPayloadV3(
	ctx context.Context, payloadID engineprimitives.PayloadID,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	var t ExecutionPayloadT
	result := &engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadT,
		*engineprimitives.BlobsBundleV1[
			eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
		],
	]{
		ExecutionPayload: t.Empty(version.Deneb),
	}

	if err := s.Client.Client().CallContext(
		ctx, result, GetPayloadMethodV3, payloadID,
	); err != nil {
		return nil, err
	}
	return result, nil
}

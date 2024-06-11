package ethclient

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// NewPayload calls the engine_newPayloadV3 method via JSON-RPC.
func (s *Eth1Client[ExecutionPayloadT]) NewPayload(
	ctx context.Context,
	payload ExecutionPayloadT,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *primitives.Root,
) (*engineprimitives.PayloadStatusV1, error) {
	switch payload.Version() {
	case version.Deneb:
		return s.NewPayloadV3(
			ctx, payload, versionedHashes, parentBlockRoot,
		)
	default:
		return nil, ErrInvalidVersion
	}
}

// newPayload is used to call the underlying JSON-RPC method for newPayload.
func (s *Eth1Client[ExecutionPayloadT]) NewPayloadV3(
	ctx context.Context,
	payload ExecutionPayloadT,
	versionedHashes []common.ExecutionHash,
	parentBlockRoot *primitives.Root,
) (*engineprimitives.PayloadStatusV1, error) {
	result := &engineprimitives.PayloadStatusV1{}
	if err := s.Client.Client().CallContext(
		ctx, result, NewPayloadMethodV3, payload, versionedHashes,
		(*common.ExecutionHash)(parentBlockRoot),
	); err != nil {
		return nil, err
	}
	return result, nil
}

package runtime

import (
	"context"

	ctypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/cosmos/gogoproto/proto"
)

type Middleware interface {
	InitGenesis(
		ctx context.Context,
		bz []byte,
	) (transition.ValidatorUpdates, error)

	PrepareProposal(
		ctx context.Context,
		slotData *types.SlotData[
			*ctypes.AttestationData,
			*ctypes.SlashingInfo],
	) ([]byte, []byte, error)

	ProcessProposal(
		ctx context.Context,
		req proto.Message,
	) (proto.Message, error)

	EndBlock(
		ctx context.Context, req proto.Message,
	) (transition.ValidatorUpdates, error)
}

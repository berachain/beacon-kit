package preblock

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	v1 "github.com/itsdevbear/bolaris/types/v1"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

type BeaconKeeper interface {
	ForkChoiceStore(ctx context.Context) v1.ForkChoiceStore
}

// BeaconPreBlockHandler is responsible for aggregating oracle data from each
// validator and writing the oracle data into the store before any transactions
// are executed/finalized for a given block.
type BeaconPreBlockHandler struct {
	logger log.Logger

	// keeper is the keeper for the oracle module. This is utilized to write
	// oracle data to state.
	keeper BeaconKeeper
}

// NewBeaconPreBlockHandler returns a new BeaconPreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewBeaconPreBlockHandler(
	logger log.Logger,
	beaconKeeper BeaconKeeper,
) *BeaconPreBlockHandler {
	return &BeaconPreBlockHandler{
		logger: logger,
		keeper: beaconKeeper,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *BeaconPreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		h.logger.Info(
			"executing the pre-finalize block hook",
			"height", req.Height,
		)

		beaconBlockData := req.Txs[0] // todo modularize.
		payload := new(enginev1.ExecutionPayloadCapellaWithValue)
		payload.Payload = new(enginev1.ExecutionPayloadCapella)
		if err := payload.Payload.UnmarshalSSZ(beaconBlockData); err != nil {
			h.logger.Error("payload in beacon block could not be unmarshalled", "err", err)
			return nil, err
		}
		// todo handle hardforks without needing codechange.
		data, err := blocks.WrappedExecutionPayloadCapella(
			payload.Payload, blocks.PayloadValueToGwei(payload.Value),
		)
		if err != nil {
			h.logger.Error("failed to wrap payload", "err", err)
			return nil, err
		}

		fmt.Println("Beacon Root is Processing", "execution_block_hash", common.BytesToHash(data.BlockHash()))

		// Finalize the block that is being proposed.
		store := h.keeper.ForkChoiceStore(ctx)
		store.SetFinalizedBlockHash([32]byte(data.BlockHash()))
		store.SetSafeBlockHash([32]byte(data.BlockHash()))
		store.SetLastValidHead([32]byte(data.BlockHash()))
		return &sdk.ResponsePreBlock{}, nil
	}
}

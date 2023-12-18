package proposal

import (
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/cosmos/runtime/miner"
)

// TODO: Need to have the wait for syncing phase at the start to allow the Execution Client
// to sync up and the consensus client shouldn't join the validator set yet.
const PayloadPosition = 0

type Handler struct {
	prepareProposal sdk.PrepareProposalHandler
	processProposal sdk.ProcessProposalHandler
	miner           *miner.Miner
}

func NewHandler(
	miner *miner.Miner,
	prepareProposal sdk.PrepareProposalHandler,
	processProposal sdk.ProcessProposalHandler,
) *Handler {
	return &Handler{
		miner:           miner,
		prepareProposal: prepareProposal,
		processProposal: processProposal,
	}
}

func (h *Handler) PrepareProposalHandler(
	ctx sdk.Context, req *abci.RequestPrepareProposal,
) (*abci.ResponsePrepareProposal, error) {
	logger := ctx.Logger().With("module", "prepare-proposal")
	// Build the block on the execution layer.
	// TODO: manage the different type of engine API errors.
	payload, err := h.miner.BuildBlockV2(ctx)
	if err != nil {
		logger.Error("failed to build block", "err", err)
		return nil, err
	}

	// Run the remainder of the prepare proposal handler.
	resp, err := h.prepareProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	// Marshal the payload.
	bz, err := payload.MarshalSSZ()
	if err != nil {
		return nil, err
	}

	// Inject the payload into the proposal.
	resp.Txs = append([][]byte{bz}, resp.Txs...)
	return resp, nil
}

func (h *Handler) ProcessProposalHandler(
	ctx sdk.Context, req *abci.RequestProcessProposal,
) (*abci.ResponseProcessProposal, error) {
	logger := ctx.Logger().With("module", "process-proposal")

	// Extract the marshalled payload from the proposal
	bz := req.Txs[PayloadPosition]
	req.Txs = req.Txs[1:]

	if bz == nil {
		logger.Error("payload missing from proposal")
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
	}
	payload := new(enginev1.ExecutionPayloadCapellaWithValue)
	payload.Payload = new(enginev1.ExecutionPayloadCapella)
	if err := payload.Payload.UnmarshalSSZ(bz); err != nil {
		logger.Error("failed to unmarshal payload", "err", err)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
	}
	// todo handle hardforks without needing codechange.
	data, err := blocks.WrappedExecutionPayloadCapella(
		payload.Payload, blocks.PayloadValueToGwei(payload.Value),
	)

	if err != nil {
		logger.Error("failed to wrap payload", "err", err)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
	}

	if err = h.miner.ValidateBlock(ctx, data); err != nil {
		logger.Error("failed to validate block", "err", err)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
	}

	return h.processProposal(ctx, req)
}

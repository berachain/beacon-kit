package proposal

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/itsdevbear/bolaris/cosmos/runtime/miner"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// TODO: Need to have the wait for syncing phase at the start to allow the Execution Client
// to sync up and the consensus client shouldn't join the validator set yet.

const PayloadPosition = 0

type ProposalHandler2 struct {
	miner *miner.Miner
}

func NewProposalHandler2(miner *miner.Miner) *ProposalHandler2 {
	return &ProposalHandler2{miner: miner}
}

func (h *ProposalHandler2) PrepareProposalHandler(
	ctx sdk.Context, req *abci.RequestPrepareProposal,
) (*abci.ResponsePrepareProposal, error) {
	var resp abci.ResponsePrepareProposal
	logger := ctx.Logger().With("module", "prepare-proposal")

	// Build the block on the execution layer.
	payload, err := h.miner.BuildBlockV2(ctx)
	// TODO: manage the different type of engine API errors.
	if err != nil {
		logger.Error("failed to build block", "err", err)
		return nil, err
	}

	bz, err := payload.MarshalSSZ()
	if err != nil {
		return nil, err
	}

	// Inject the payload into the proposal.
	resp.Txs = append([][]byte{bz}, resp.Txs...)
	return &resp, nil
}

func (h *ProposalHandler2) ProcessProposalHandler(
	ctx sdk.Context, req *abci.RequestProcessProposal,
) (*abci.ResponseProcessProposal, error) {
	logger := ctx.Logger().With("module", "process-proposal")

	// Extract the marshalled payload from the proposal
	bz := req.Txs[PayloadPosition]
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

	if err := h.miner.ValidateBlock(ctx, data); err != nil {
		logger.Error("failed to validate block", "err", err)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
	}

	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
}

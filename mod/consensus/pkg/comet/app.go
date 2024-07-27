package cometbft

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	abci "github.com/cometbft/cometbft/abci/types"
)

var _ abci.Application = (*Application[engine.Client])(nil)

type Application[NodeT engine.Client] struct {
	// TODO: remove this once we implement the whole interface
	// actually most of these methods are noops, could be useful
	// to keep it instead of implementing the whole interface
	abci.BaseApplication

	// Logger
	logger log.Logger[any]

	// TODO: remove this once we implement the rpc engine
	node NodeT

	// CometBFT Params
	lastFinalizedHeight int64
}

func NewApplication[NodeT engine.Client](
	logger log.Logger[any],
	node NodeT,
) *Application[NodeT] {
	return &Application[NodeT]{
		BaseApplication: abci.BaseApplication{},
		logger:          logger,
		node:            node,
	}
}

func (app *Application[NodeT]) InitChain(
	ctx context.Context,
	req *abci.InitChainRequest,
) (*abci.InitChainResponse, error) {
	app.logger.Info(
		"Initializing chain",
		"chain id", req.ChainId,
		"height", req.InitialHeight,
	)
	valUpdates, stateHash, err := app.node.InitChain(ctx, req.AppStateBytes)
	if err != nil {
		return nil, err
	}

	return &abci.InitChainResponse{
		ConsensusParams: req.ConsensusParams, // TODO: do we need to override this??? how can we abstract this???
		Validators:      convertValidatorUpdates(valUpdates),
		AppHash:         stateHash,
	}, nil
}

func (app *Application[NodeT]) PrepareProposal(
	ctx context.Context,
	req *abci.PrepareProposalRequest,
) (*abci.PrepareProposalResponse, error) {
	app.logger.Info("PrepareProposal", "req", req)
	txs, err := app.node.PrepareProposal(ctx, req)
	if err != nil {
		return nil, err
	}
	return &abci.PrepareProposalResponse{
		Txs: txs,
	}, nil
}

func (app *Application[NodeT]) ProcessProposal(
	ctx context.Context,
	req *abci.ProcessProposalRequest,
) (*abci.ProcessProposalResponse, error) {
	app.logger.Info("ProcessProposal", "req", req)
	var err error
	status := abci.PROCESS_PROPOSAL_STATUS_ACCEPT
	if err = app.node.ProcessProposal(ctx, req); err != nil {
		status = abci.PROCESS_PROPOSAL_STATUS_REJECT
	}
	return &abci.ProcessProposalResponse{
		Status: status,
	}, err
}

func (app *Application[NodeT]) FinalizeBlock(
	ctx context.Context,
	req *abci.FinalizeBlockRequest,
) (*abci.FinalizeBlockResponse, error) {
	app.logger.Info("FinalizeBlock", "req", req)
	resp, err := app.node.FinalizeBlock(ctx, req)
	if err != nil {
		return nil, err
	}
	app.lastFinalizedHeight = req.Height
	return resp, nil
}

func (app *Application[NodeT]) Commit(
	ctx context.Context,
	req *abci.CommitRequest,
) (*abci.CommitResponse, error) {
	app.logger.Info("Commit", "req", req)
	return &abci.CommitResponse{
		RetainHeight: 0, // TODO: implement this
	}, nil
}

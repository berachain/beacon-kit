package engine

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	abci "github.com/cometbft/cometbft/abci/types"
)

// TODO: decouple from the proto types
type Client interface {
	// InitChain is called when the blockchain is first started
	// It returns the validator set and the app hash.
	InitChain(
		ctx context.Context,
		genesisBz []byte,
	) (transition.ValidatorUpdates, []byte, error) // Initialize blockchain w validators/other info from CometBFT

	// PrepareProposal is called when a proposal is made.
	// It returns the txs to be executed in the proposal.
	PrepareProposal(
		ctx context.Context,
		req *abci.PrepareProposalRequest,
	) ([][]byte, error)

	// ProcessProposal is called when a proposal is processed.
	// It returns an error if the proposal is invalid.
	ProcessProposal(
		ctx context.Context,
		req *abci.ProcessProposalRequest,
	) error

	// Deliver the decided block with its txs to the Application
	FinalizeBlock(
		ctx context.Context,
		req *abci.FinalizeBlockRequest,
	) (*abci.FinalizeBlockResponse, error)

	// TODO: snapshot methods
}

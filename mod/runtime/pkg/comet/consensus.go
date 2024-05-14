package comet

import (
	"context"

	math "github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	cmttypes "github.com/cometbft/cometbft/types"
)

type ChainSpec interface {
	// GetCometBFTConfigForSlot returns the CometBFT configuration for the given slot.
	GetCometBFTConfigForSlot(math.Slot) *cmttypes.ConsensusParams
}

// ConsensusParamsStore is a store for consensus parameters.
type ConsensusParamsStore struct {
	cs ChainSpec
}

// NewConsensusParamsStore creates a new ConsensusParamsStore.
func NewConsensusParamsStore(cs ChainSpec) *ConsensusParamsStore {
	return &ConsensusParamsStore{
		cs: cs,
	}
}

// Get retrieves the consensus parameters from the store.
// It returns the consensus parameters and an error, if any.
func (s *ConsensusParamsStore) Get(ctx context.Context) (cmtproto.ConsensusParams, error) {
	return s.cs.GetCometBFTConfigForSlot(0).ToProto(), nil
}

// Has checks if the consensus parameters exist in the store.
// It returns a boolean indicating the presence of the parameters and an error, if any.
func (s *ConsensusParamsStore) Has(ctx context.Context) (bool, error) {
	return true, nil
}

// Set stores the given consensus parameters in the store.
// It returns an error, if any.
func (s *ConsensusParamsStore) Set(ctx context.Context, cp cmtproto.ConsensusParams) error {
	return nil
}

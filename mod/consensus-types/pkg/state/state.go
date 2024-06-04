package state

import (
	deneb "github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconState is the interface for the beacon state.
type BeaconState struct {
	// TODO: decouple from deneb.BeaconState
	*deneb.BeaconState
}

// New creates a new BeaconState.
func (st *BeaconState) New(
	// TODO: handle diffferent versions.
	_ uint32, /*version*/
	genesisValidatorsRoot primitives.Root,
	slot math.Slot,
	fork *types.Fork,
	latestBlockHeader *types.BeaconBlockHeader,
	blockRoots []primitives.Root,
	stateRoots []primitives.Root,
	eth1Data *types.Eth1Data,
	eth1DepositIndex uint64,
	latestExecutionPayloadHeader *types.ExecutionPayloadHeader,
	validators []*types.Validator,
	balances []uint64,
	randaoMixes []primitives.Bytes32,
	nextWithdrawalIndex uint64,
	nextWithdrawalValidatorIndex math.ValidatorIndex,
	slashings []uint64,
	totalSlashing math.Gwei,
) *BeaconState {
	return &BeaconState{
		BeaconState: &deneb.BeaconState{
			Slot:                  slot,
			GenesisValidatorsRoot: genesisValidatorsRoot,
			Fork:                  fork,
			LatestBlockHeader:     latestBlockHeader,
			BlockRoots:            blockRoots,
			StateRoots:            stateRoots,
			//nolint:lll
			LatestExecutionPayloadHeader: latestExecutionPayloadHeader.
				ExecutionPayloadHeader.(*types.ExecutionPayloadHeaderDeneb),
			Eth1Data:                     eth1Data,
			Eth1DepositIndex:             eth1DepositIndex,
			Validators:                   validators,
			Balances:                     balances,
			RandaoMixes:                  randaoMixes,
			NextWithdrawalIndex:          nextWithdrawalIndex,
			NextWithdrawalValidatorIndex: nextWithdrawalValidatorIndex,
			Slashings:                    slashings,
			TotalSlashing:                totalSlashing,
		},
	}
}

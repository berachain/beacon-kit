package state_transition

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// InitializePreminedBeaconStateFromEth1 initializes the beacon state.
//
//nolint:gocognit,funlen // todo fix.
func (sp *StateProcessor[
	_, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT, _, DepositT,
	_, ExecutionPayloadHeaderT, ForkT, _, _, ValidatorT, _, _, _,
]) InitializePreminedBeaconStateFromEth1(
	st BeaconStateT,
	deposits []DepositT,
	executionPayloadHeader ExecutionPayloadHeaderT,
	genesisVersion common.Version,
) (transition.ValidatorUpdates, error) {
	var (
		blkHeader BeaconBlockHeaderT
		blkBody   BeaconBlockBodyT
		fork      ForkT
	)
	fork = fork.New(
		genesisVersion,
		genesisVersion,
		math.U64(constants.GenesisEpoch),
	)

	if err := st.SetSlot(0); err != nil {
		return nil, err
	}

	if err := st.SetFork(fork); err != nil {
		return nil, err
	}

	if err := st.SetEth1DepositIndex(0); err != nil {
		return nil, err
	}

	// TODO: we need to handle common.Version vs
	// uint32 better.
	bodyRoot := blkBody.Empty(
		version.ToUint32(genesisVersion)).HashTreeRoot()
	if err := st.SetLatestBlockHeader(blkHeader.New(
		0, 0, common.Root{}, common.Root{}, bodyRoot,
	)); err != nil {
		return nil, err
	}

	for i := range sp.cs.EpochsPerHistoricalVector() {
		if err := st.UpdateRandaoMixAtIndex(
			i,
			common.Bytes32(executionPayloadHeader.GetBlockHash()),
		); err != nil {
			return nil, err
		}
	}

	for _, deposit := range deposits {
		if err := sp.processDeposit(st, deposit); err != nil {
			return nil, err
		}
	}

	// TODO: process activations.
	validators, err := st.GetValidators()
	if err != nil {
		return nil, err
	}

	if err = st.SetGenesisValidatorsRoot(validators.HashTreeRoot()); err != nil {
		return nil, err
	}

	if err = st.SetLatestExecutionPayloadHeader(
		executionPayloadHeader,
	); err != nil {
		return nil, err
	}

	// Setup a bunch of 0s to prime the DB.
	for i := range sp.cs.HistoricalRootsLimit() {
		//#nosec:G701 // won't overflow in practice.
		if err = st.UpdateBlockRootAtIndex(i, common.Root{}); err != nil {
			return nil, err
		}
		if err = st.UpdateStateRootAtIndex(i, common.Root{}); err != nil {
			return nil, err
		}
	}

	if err = st.SetNextWithdrawalIndex(0); err != nil {
		return nil, err
	}

	if err = st.SetNextWithdrawalValidatorIndex(
		0,
	); err != nil {
		return nil, err
	}

	if err = st.SetTotalSlashing(0); err != nil {
		return nil, err
	}

	var updates transition.ValidatorUpdates
	updates, err = sp.processSyncCommitteeUpdates(st)
	if err != nil {
		return nil, err
	}
	st.Save()
	return updates, nil
}

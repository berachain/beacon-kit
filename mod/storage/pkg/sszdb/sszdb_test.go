package sszdb_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb"
	"github.com/stretchr/testify/require"
)

func testBeaconState() (*deneb.BeaconState, error) {
	bz, err := os.ReadFile("./testdata/beacon.ssz")
	if err != nil {
		return nil, err
	}
	state := &deneb.BeaconState{}
	err = state.UnmarshalSSZ(bz)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func TestDB_Metadata(t *testing.T) {
	beacon, err := testBeaconState()
	require.NoError(t, err)

	beacon.GenesisValidatorsRoot = [32]byte{7, 7, 7, 7}
	beacon.Slot = 777
	beacon.Fork = &types.Fork{
		PreviousVersion: [4]byte{1, 2, 3, 4},
		CurrentVersion:  [4]byte{5, 6, 7, 8},
		Epoch:           123,
	}
	beacon.LatestBlockHeader = &types.BeaconBlockHeader{
		BeaconBlockHeaderBase: types.BeaconBlockHeaderBase{
			Slot:            777,
			ProposerIndex:   123,
			ParentBlockRoot: [32]byte{1, 2, 3, 4},
			StateRoot:       [32]byte{5, 6, 7, 8},
		},
		BodyRoot: [32]byte{9, 10, 11, 12},
	}
	beacon.BlockRoots = []common.Root{
		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16},
	}
	beacon.Validators = []*types.Validator{
		{
			Pubkey:                     [48]byte{1, 2, 3, 4},
			WithdrawalCredentials:      [32]byte{5, 6, 7, 8},
			EffectiveBalance:           123,
			Slashed:                    true,
			ActivationEligibilityEpoch: 123,
			ActivationEpoch:            123,
			ExitEpoch:                  123,
			WithdrawableEpoch:          123,
		},
		{
			Pubkey:                     [48]byte{9, 10, 11, 12},
			WithdrawalCredentials:      [32]byte{13, 14, 15, 16},
			EffectiveBalance:           456,
			Slashed:                    false,
			ActivationEligibilityEpoch: 456,
			ActivationEpoch:            456,
			ExitEpoch:                  456,
			WithdrawableEpoch:          456,
		},
	}

	dir := t.TempDir() + "/sszdb.db"
	db, err := sszdb.New(sszdb.Config{Path: dir})
	require.NoError(t, err)

	err = db.SaveMonolith(beacon)
	require.NoError(t, err)

	schemaDb, err := sszdb.NewSchemaDb(db, beacon)
	require.NoError(t, err)

	bz, err := schemaDb.GetGenesisValidatorsRoot()
	require.NoError(t, err)
	require.True(t, bytes.Equal(bz[:], beacon.GenesisValidatorsRoot[:]))

	slot, err := schemaDb.GetSlot()
	require.NoError(t, err)
	require.Equal(t, beacon.Slot, slot)

	fork, err := schemaDb.GetFork()
	require.NoError(t, err)
	require.Equal(t, beacon.Fork, fork)

	latestHeader, err := schemaDb.GetLatestBlockHeader()
	require.NoError(t, err)
	require.Equal(t, beacon.LatestBlockHeader, latestHeader)

	roots, err := schemaDb.GetBlockRoots()
	require.NoError(t, err)
	require.Equal(t, len(beacon.BlockRoots), len(roots))
	for i, r := range roots {
		require.Equal(t, beacon.BlockRoots[i], r)
	}

	val0, err := schemaDb.GetValidatorAtIndex(0)
	require.NoError(t, err)
	require.Equal(t, beacon.Validators[0], val0)

	vals, err := schemaDb.GetValidators()
	require.NoError(t, err)
	require.Equal(t, len(beacon.Validators), len(vals))
	for i, v := range vals {
		require.Equal(t, beacon.Validators[i], v)
	}
}

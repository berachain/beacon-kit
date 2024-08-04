package sszdb_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestDB_Metadata(t *testing.T) {
	beacon := &types.BeaconState[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.BeaconBlockHeader,
		types.Eth1Data,
		types.ExecutionPayloadHeader,
		types.Fork,
		types.Validator,
		*schema.Codec,
	]{}
	beacon.GenesisValidatorsRoot = [32]byte{7, 7, 7, 7}
	beacon.Slot = 777
	beacon.Fork = &types.Fork{
		PreviousVersion: [4]byte{1, 2, 3, 4},
		CurrentVersion:  [4]byte{5, 6, 7, 8},
		Epoch:           123,
	}
	beacon.LatestBlockHeader = &types.BeaconBlockHeader{
		Slot:            777,
		ProposerIndex:   123,
		ParentBlockRoot: [32]byte{1, 2, 3, 4},
		StateRoot:       [32]byte{5, 6, 7, 8},
		BodyRoot:        [32]byte{9, 10, 11, 12},
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
	beacon.LatestExecutionPayloadHeader = &types.ExecutionPayloadHeader{
		StateRoot:     [32]byte{1, 2, 3, 4},
		ReceiptsRoot:  [32]byte{5, 6, 7, 8},
		Random:        [32]byte{13, 14, 15, 16},
		LogsBloom:     [256]byte{17, 18, 19, 20},
		Number:        123,
		GasLimit:      456,
		GasUsed:       789,
		Timestamp:     101112,
		ExtraData:     []byte{29, 30, 31, 32, 35},
		BaseFeePerGas: uint256.MustFromDecimal("123456"),
	}

	dir := t.TempDir() + "/sszdb.db"
	db, err := sszdb.NewBackend(sszdb.BackendConfig{Path: dir})
	require.NoError(t, err)

	ctx := context.TODO()
	err = db.SaveMonolith(beacon)
	require.NoError(t, err)
	err = db.Commit(ctx)
	require.NoError(t, err)

	beaconDB, err := sszdb.NewSchemaDB(db, beacon)
	require.NoError(t, err)

	/*
		bz, err := beaconDB.GetGenesisValidatorsRoot(ctx)
		require.NoError(t, err)
		require.True(t, bytes.Equal(bz[:], beacon.GenesisValidatorsRoot[:]))

		slot, err := beaconDB.GetSlot(ctx)
		require.NoError(t, err)
		require.Equal(t, beacon.Slot, slot)

		fork, err := beaconDB.GetFork(ctx)
		require.NoError(t, err)
		require.Equal(t, beacon.Fork, fork)

		latestHeader, err := beaconDB.GetLatestBlockHeader(ctx)
		require.NoError(t, err)
		require.Equal(t, beacon.LatestBlockHeader, latestHeader)

		roots, err := beaconDB.GetBlockRoots(ctx)
		require.NoError(t, err)
		require.Equal(t, len(beacon.BlockRoots), len(roots))
		for i, r := range roots {
			require.Equal(t, beacon.BlockRoots[i], r)
		}

		val0, err := beaconDB.GetValidatorAtIndex(ctx, 0)
		require.NoError(t, err)
		require.Equal(t, beacon.Validators[0], val0)

		vals, err := beaconDB.GetValidators(ctx)
		require.NoError(t, err)
		require.Equal(t, len(beacon.Validators), len(vals))
		for i, v := range vals {
			require.Equal(t, beacon.Validators[i], v)
		}
	*/

	headerSSZBytes, err := beacon.LatestExecutionPayloadHeader.MarshalSSZ()
	require.NoError(t, err)
	headerBz, err := beaconDB.GetLatestExecutionPayloadHeader(ctx)
	require.True(t, bytes.Equal(headerSSZBytes, headerBz))
	require.NoError(t, err)
	var header types.ExecutionPayloadHeader
	require.NoError(t, header.UnmarshalSSZ(headerBz))
	require.Equal(t,
		beacon.LatestExecutionPayloadHeader,
		header,
	)
}

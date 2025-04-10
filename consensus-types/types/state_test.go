// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// generateValidBeaconState generates a valid beacon state for the types.
func generateValidBeaconState(forkVersion common.Version) *types.BeaconState {
	beaconState := &types.BeaconState{
		Versionable:           types.NewVersionable(forkVersion),
		GenesisValidatorsRoot: common.Root{0x01, 0x02, 0x03},
		Slot:                  1234,
		BlockRoots: []common.Root{
			{0x04, 0x05, 0x06},
			{0x07, 0x08, 0x09},
		},
		StateRoots: []common.Root{
			{0x0a, 0x0b, 0x0c},
			{0x0d, 0x0e, 0x0f},
		},
		Fork: &types.Fork{
			PreviousVersion: [4]byte{0x01, 0x00, 0x00, 0x00},
			CurrentVersion:  [4]byte{0x02, 0x00, 0x00, 0x00},
			Epoch:           5678,
		},
		Validators: []*types.Validator{
			{
				Pubkey:                     [48]byte{0x01},
				WithdrawalCredentials:      [32]byte{0x02},
				EffectiveBalance:           32000000000,
				Slashed:                    false,
				ActivationEligibilityEpoch: 1,
				ActivationEpoch:            2,
				ExitEpoch:                  18446744073709551615,
				WithdrawableEpoch:          18446744073709551615,
			},
			{
				Pubkey:                     [48]byte{0x03},
				WithdrawalCredentials:      [32]byte{0x04},
				EffectiveBalance:           31000000000,
				Slashed:                    true,
				ActivationEligibilityEpoch: 3,
				ActivationEpoch:            4,
				ExitEpoch:                  5,
				WithdrawableEpoch:          6,
			},
		},
		Balances:                     []uint64{32000000000, 31000000000},
		RandaoMixes:                  generateRandomBytes32(65536),
		Slashings:                    []math.Gwei{1000000000, 2000000000},
		NextWithdrawalIndex:          7,
		NextWithdrawalValidatorIndex: 8,
		TotalSlashing:                3000000000,
		LatestExecutionPayloadHeader: &types.ExecutionPayloadHeader{
			ParentHash:       [32]byte{0x16, 0x17, 0x18},
			FeeRecipient:     [20]byte{0x19, 0x1a, 0x1b},
			StateRoot:        [32]byte{0x1c, 0x1d, 0x1e},
			ReceiptsRoot:     [32]byte{0x1f, 0x20, 0x21},
			LogsBloom:        [256]byte{0x22},
			Random:           [32]byte{0x23, 0x24, 0x25},
			Number:           9876,
			GasLimit:         30000000,
			GasUsed:          25000000,
			Timestamp:        1625097600,
			ExtraData:        []byte{0x26, 0x27, 0x28},
			BaseFeePerGas:    math.NewU256(3906250),
			BlockHash:        [32]byte{0x2c, 0x2d, 0x2e},
			TransactionsRoot: [32]byte{0x2f, 0x30, 0x31},
			WithdrawalsRoot:  [32]byte{0x32, 0x33, 0x34},
		},
		LatestBlockHeader: &types.BeaconBlockHeader{
			Slot:            5678,
			ProposerIndex:   123,
			ParentBlockRoot: [32]byte{0x35, 0x36, 0x37},
			StateRoot:       [32]byte{0x38, 0x39, 0x3a},
			BodyRoot:        [32]byte{0x3b, 0x3c, 0x3d},
		},
		Eth1Data: &types.Eth1Data{
			DepositRoot:  [32]byte{0x3e, 0x3f, 0x40},
			DepositCount: 1000,
			BlockHash:    [32]byte{0x41, 0x42, 0x43},
		},
		Eth1DepositIndex: 100,
	}

	if version.EqualsOrIsAfter(beaconState.GetForkVersion(), version.Electra()) {
		beaconState.PendingPartialWithdrawals = []*types.PendingPartialWithdrawal{
			{
				ValidatorIndex:    123,
				Amount:            32000000000,
				WithdrawableEpoch: 100,
			},
			{
				ValidatorIndex:    124,
				Amount:            100,
				WithdrawableEpoch: 1,
			},
		}
	}
	return beaconState
}

func generateRandomBytes32(count int) []common.Bytes32 {
	result := make([]common.Bytes32, count)
	for i := range result {
		var randomBytes [32]byte
		for j := range randomBytes {
			randomBytes[j] = byte((i + j) % 256)
		}
		result[i] = randomBytes
	}
	return result
}

func TestBeaconStateMarshalUnmarshalSSZ(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		genState := generateValidBeaconState(v)

		data, fastSSZMarshalErr := genState.MarshalSSZ()
		require.NoError(t, fastSSZMarshalErr)
		require.NotNil(t, data)

		newState := types.NewEmptyBeaconStateWithVersion(v)
		err := newState.UnmarshalSSZ(data)
		require.NoError(t, err)

		require.EqualValues(t, genState, newState)

		// Check if the state size is greater than 0
		require.Positive(t, karalabessz.Size(genState))
	})
}

func TestHashTreeRoot(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		state := generateValidBeaconState(v)
		require.NotPanics(t, func() {
			state.HashTreeRoot()
		})
	})
}

func TestGetTree(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		state := generateValidBeaconState(v)
		tree, err := state.GetTree()
		require.NoError(t, err)
		require.NotNil(t, tree)
	})
}

func TestBeaconState_UnmarshalSSZ_Error(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		state := types.NewEmptyBeaconStateWithVersion(v)
		err := state.UnmarshalSSZ([]byte{0x01, 0x02, 0x03}) // Invalid data
		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	})
}

func TestBeaconState_HashTreeRoot(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		state := generateValidBeaconState(v)

		// Get the HashTreeRoot
		root := state.HashTreeRoot()

		// Get the HashConcurrent
		concurrentRoot := common.Root(karalabessz.HashSequential(state))

		// Compare the results
		require.Equal(
			t,
			root,
			concurrentRoot,
			"HashTreeRoot and HashSequential should produce the same result",
		)
	})
}

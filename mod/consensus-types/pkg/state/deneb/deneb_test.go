// SPDX-License-Identifier: MIT
//
// # Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//

package deneb_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// generateValidBeaconState generates a valid beacon state for the Deneb.
func generateValidBeaconState() *deneb.BeaconState {
	return &deneb.BeaconState{
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
		Slashings:                    []uint64{1000000000, 2000000000},
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
			BaseFeePerGas:    [32]byte{0x29, 0x2a, 0x2b},
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
}

func generateRandomBytes32(count int) []common.Bytes32 {
	result := make([]common.Bytes32, count)
	for i := 0; i < count; i++ {
		var randomBytes [32]byte
		for j := 0; j < 32; j++ {
			randomBytes[j] = byte(i + j)
		}
		result[i] = randomBytes
	}
	return result
}

func TestBeaconStateMarshalUnmarshalSSZ(t *testing.T) {
	state := generateValidBeaconState()

	data, fastSSZMarshalErr := state.MarshalSSZ()
	require.NoError(t, fastSSZMarshalErr)
	require.NotNil(t, data)

	newState := &deneb.BeaconState{}
	err := newState.UnmarshalSSZ(data)
	require.NoError(t, err)

	require.EqualValues(t, state, newState)

	// Check if the state size is greater than 0
	require.Positive(t, state.SizeSSZ(false))
}

func TestHashTreeRoot(t *testing.T) {
	state := generateValidBeaconState()
	_, err := state.HashTreeRoot()
	require.NoError(t, err)
}

func TestGetTree(t *testing.T) {
	state := generateValidBeaconState()
	tree, err := state.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestBeaconState_UnmarshalSSZ_Error(t *testing.T) {
	state := &deneb.BeaconState{}
	err := state.UnmarshalSSZ([]byte{0x01, 0x02, 0x03}) // Invalid data
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestBeaconState_MarshalSSZTo(t *testing.T) {
	state := generateValidBeaconState()
	data, err := state.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var buf []byte
	buf, err = state.MarshalSSZTo(buf)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestBeaconState_HashTreeRoot(t *testing.T) {
	state := generateValidBeaconState()

	// Get the HashTreeRoot
	root, err := state.HashTreeRoot()
	require.NoError(t, err)

	// Get the HashConcurrent
	concurrentRoot := karalabessz.HashSequential(state)

	// Compare the results
	require.Equal(
		t,
		root,
		concurrentRoot,
		"HashTreeRoot and HashSequential should produce the same result",
	)
}

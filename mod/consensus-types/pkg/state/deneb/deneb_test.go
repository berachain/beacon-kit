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
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

// generateValidBeaconState generates a valid beacon state for the Deneb.
func generateValidBeaconState() *deneb.BeaconState {
	var byteArray [256]byte
	return &deneb.BeaconState{
		BlockRoots:  []common.Root{},
		StateRoots:  []common.Root{},
		Fork:        &types.Fork{},
		Validators:  []*types.Validator{},
		Balances:    []uint64{},
		RandaoMixes: []common.Bytes32{},
		Slashings:   []uint64{},
		LatestExecutionPayloadHeader: &types.ExecutionPayloadHeaderDeneb{
			LogsBloom: byteArray[:],
			ExtraData: []byte{},
		},
	}
}

func TestBeaconStateMarshalUnmarshalSSZ(t *testing.T) {
	state := generateValidBeaconState()

	data, fastSSZMarshalErr := state.MarshalSSZ()
	require.NoError(t, fastSSZMarshalErr)
	require.NotNil(t, data)

	newState := &deneb.BeaconState{}
	err := newState.UnmarshalSSZ(data)
	require.NoError(t, err)

	require.Equal(t, state, newState)

	// Check if the state size is greater than 0
	require.Positive(t, state.SizeSSZ())
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
	require.ErrorIs(t, err, ssz.ErrSize)
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

func TestBeaconState_MarshalSSZFields(t *testing.T) {
	state := generateValidBeaconState()

	// Test BlockRoots field
	state.BlockRoots = make([]common.Root, 8193) // Exceeding the limit
	_, err := state.MarshalSSZ()
	require.Error(t, err)
	state.BlockRoots = make([]common.Root, 8192) // Within the limit
	_, err = state.MarshalSSZ()
	require.NoError(t, err)

	// Test StateRoots field
	state.StateRoots = make([]common.Root, 8193) // Exceeding the limit
	_, err = state.MarshalSSZ()
	require.Error(t, err)
	state.StateRoots = make([]common.Root, 8192) // Within the limit
	_, err = state.MarshalSSZ()
	require.NoError(t, err)

	// Test LatestExecutionPayloadHeader field
	state.LatestExecutionPayloadHeader = &types.ExecutionPayloadHeaderDeneb{
		LogsBloom: make([]byte, 256), // Initialize LogsBloom with 256 bytes
	}
	_, err = state.MarshalSSZ()
	require.NoError(t, err)

	// Test RandaoMixes field
	state.RandaoMixes = make([]common.Bytes32, 65537) // Exceeding the limit
	_, err = state.MarshalSSZ()
	require.Error(t, err)
	state.RandaoMixes = make([]common.Bytes32, 65536) // Within the limit
	_, err = state.MarshalSSZ()
	require.NoError(t, err)
}

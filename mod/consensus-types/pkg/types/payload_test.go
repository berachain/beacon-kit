// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package types_test

import (
	"encoding/json"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func generateExecutableDataDeneb() *types.ExecutableDataDeneb {
	return &types.ExecutableDataDeneb{
		ParentHash:    common.ExecutionHash{},
		FeeRecipient:  common.ExecutionAddress{},
		StateRoot:     bytes.B32{},
		ReceiptsRoot:  bytes.B32{},
		LogsBloom:     make([]byte, 256),
		Random:        bytes.B32{},
		Number:        math.U64(0),
		GasLimit:      math.U64(0),
		GasUsed:       math.U64(0),
		Timestamp:     math.U64(0),
		ExtraData:     []byte{},
		BaseFeePerGas: math.Wei{},
		BlockHash:     common.ExecutionHash{},
		Transactions:  [][]byte{},
		Withdrawals:   []*engineprimitives.Withdrawal{},
		BlobGasUsed:   math.U64(0),
		ExcessBlobGas: math.U64(0),
	}
}
func TestExecutableDataDeneb_Serialization(t *testing.T) {
	original := generateExecutableDataDeneb()

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.ExecutableDataDeneb
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)

	require.Equal(t, original, &unmarshalled)
}

func TestExecutableDataDeneb_SizeSSZ(t *testing.T) {
	payload := generateExecutableDataDeneb()
	size := payload.SizeSSZ()
	require.Equal(t, 528, size)
}

func TestExecutableDataDeneb_HashTreeRoot(t *testing.T) {
	payload := generateExecutableDataDeneb()
	_, err := payload.HashTreeRoot()
	require.NoError(t, err)
}

func TestExecutableDataDeneb_GetTree(t *testing.T) {
	payload := generateExecutableDataDeneb()
	tree, err := payload.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestExecutableDataDeneb_Getters(t *testing.T) {
	payload := generateExecutableDataDeneb()

	require.Equal(t, common.ExecutionHash{}, payload.GetParentHash())
	require.Equal(t, common.ExecutionAddress{}, payload.GetFeeRecipient())
	require.Equal(t, bytes.B32{}, payload.GetStateRoot())
	require.Equal(t, bytes.B32{}, payload.GetReceiptsRoot())
	require.Equal(t, make([]byte, 256), payload.GetLogsBloom())
	require.Equal(t, bytes.B32{}, payload.GetPrevRandao())
	require.Equal(t, math.U64(0), payload.GetNumber())
	require.Equal(t, math.U64(0), payload.GetGasLimit())
	require.Equal(t, math.U64(0), payload.GetGasUsed())
	require.Equal(t, math.U64(0), payload.GetTimestamp())
	require.Equal(t, []byte{}, payload.GetExtraData())
	require.Equal(t, math.Wei{}, payload.GetBaseFeePerGas())
	require.Equal(t, common.ExecutionHash{}, payload.GetBlockHash())
	require.Equal(t, [][]byte{}, payload.GetTransactions())
	require.Equal(t, []*engineprimitives.Withdrawal{}, payload.GetWithdrawals())
	require.Equal(t, math.U64(0), payload.GetBlobGasUsed())
	require.Equal(t, math.U64(0), payload.GetExcessBlobGas())
}

func TestExecutableDataDeneb_MarshalJSON(t *testing.T) {
	payload := generateExecutableDataDeneb()

	data, err := payload.MarshalJSON()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.ExecutableDataDeneb
	err = unmarshalled.UnmarshalJSON(data)
	require.NoError(t, err)
	require.Equal(t, payload, &unmarshalled)
}

func TestExecutableDataDeneb_UnmarshalJSON_Error(t *testing.T) {
	malformedJSON := `{"invalidField": "invalidValue"}`

	var payload types.ExecutableDataDeneb
	err := json.Unmarshal([]byte(malformedJSON), &payload)

	require.Error(t, err)
	require.ErrorContains(
		t,
		err,
		"missing required field 'parentHash' for ExecutableDataDeneb",
	)
}

// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package deneb_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/holiman/uint256"
	"github.com/itsdevbear/bolaris/crypto/sha256"
	byteslib "github.com/itsdevbear/bolaris/lib/bytes"
	"github.com/itsdevbear/bolaris/types/engine/v1/deneb"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/stretchr/testify/require"
)

func Test_DenebToHeader(t *testing.T) {
	p := deneb.NewWrappedExecutionPayloadDeneb(&enginev1.ExecutionPayloadDeneb{
		Transactions: make([][]byte, 0),
		Withdrawals:  make([]*enginev1.Withdrawal, 0),
	}, uint256.NewInt(1))

	txRoot := sha256.HashRootAndMixinLengthAsBzSlice(p.GetTransactions())

	wdRoot := sha256.HashRootAndMixinLengthAsSlice(p.GetWithdrawals())

	parentHash := make([]byte, 32)
	_, err := rand.Read(parentHash)
	require.NoError(t, err)
	p.ParentHash = parentHash

	feeReceipt := make([]byte, 20)
	_, err = rand.Read(feeReceipt)
	require.NoError(t, err)
	p.FeeRecipient = feeReceipt

	stateRoot := make([]byte, 32)
	_, err = rand.Read(stateRoot)
	require.NoError(t, err)
	p.StateRoot = stateRoot

	receiptsRoot := make([]byte, 32)
	_, err = rand.Read(receiptsRoot)
	require.NoError(t, err)
	p.ReceiptsRoot = receiptsRoot

	logsBloom := make([]byte, 256)
	_, err = rand.Read(logsBloom)
	require.NoError(t, err)
	p.LogsBloom = logsBloom

	prevRandao := make([]byte, 32)
	_, err = rand.Read(prevRandao)
	require.NoError(t, err)
	p.PrevRandao = prevRandao

	extraData := make([]byte, 32)
	_, err = rand.Read(extraData)
	require.NoError(t, err)
	p.ExtraData = extraData

	baseFeePerGas := make([]byte, 256)
	_, err = rand.Read(baseFeePerGas)
	require.NoError(t, err)
	p.BaseFeePerGas = baseFeePerGas

	blockHash := make([]byte, 32)
	_, err = rand.Read(blockHash)
	require.NoError(t, err)
	p.BlockHash = blockHash

	blobGasUsed := uint64(1)
	p.BlobGasUsed = blobGasUsed

	excessBlobGas := uint64(1)
	p.ExcessBlobGas = excessBlobGas

	p.BlockNumber = 1
	p.GasUsed = 1
	p.GasLimit = 1
	p.Timestamp = 1

	x, err := p.ToHeader()
	require.NoError(t, err)
	h, ok := x.(*deneb.WrappedExecutionPayloadHeaderDeneb)
	require.True(t, ok)
	require.True(t, bytes.Equal(h.ParentHash, p.ParentHash))
	require.True(t, bytes.Equal(h.ParentHash, byteslib.SafeCopy(h.ParentHash)))
	require.True(t, bytes.Equal(h.FeeRecipient, byteslib.SafeCopy(h.FeeRecipient)))
	require.True(t, bytes.Equal(h.StateRoot, byteslib.SafeCopy(h.StateRoot)))
	require.True(t, bytes.Equal(h.ReceiptsRoot, byteslib.SafeCopy(h.ReceiptsRoot)))
	require.True(t, bytes.Equal(h.LogsBloom, byteslib.SafeCopy(h.LogsBloom)))
	require.True(t, bytes.Equal(h.PrevRandao, byteslib.SafeCopy(h.PrevRandao)))
	require.True(t, bytes.Equal(h.ExtraData, byteslib.SafeCopy(h.ExtraData)))
	require.True(t, bytes.Equal(h.BaseFeePerGas, byteslib.SafeCopy(h.BaseFeePerGas)))
	require.True(t, bytes.Equal(h.BlockHash, byteslib.SafeCopy(h.BlockHash)))
	require.Equal(t, uint64(1), h.BlockNumber)
	require.Equal(t, uint64(1), h.GasUsed)
	require.Equal(t, uint64(1), h.GasLimit)
	require.Equal(t, uint64(1), h.Timestamp)
	require.True(t, bytes.Equal(h.TransactionsRoot, txRoot))
	require.True(t, bytes.Equal(h.WithdrawalsRoot, wdRoot))
	require.Equal(t, uint64(1), h.BlobGasUsed)
	require.Equal(t, uint64(1), h.ExcessBlobGas)
	require.Equal(t, uint64(1), h.GetValue().Uint64())
}

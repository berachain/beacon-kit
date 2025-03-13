//go:build quick

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

package compare_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"unsafe"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	bytesprimitives "github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	zcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	zdeneb "github.com/protolambda/zrnt/eth2/beacon/deneb"
	ztree "github.com/protolambda/ztyp/tree"
	zview "github.com/protolambda/ztyp/view"
)

// --- Helper Struct ---
// execPayloadExported mirrors the exported fields of ExecutionPayload.
// Note that forkVersion is unexported and omitted.
type execPayloadExported struct {
	// ParentHash is the hash of the parent block.
	ParentHash common.ExecutionHash `json:"parentHash"`
	// FeeRecipient is the address of the fee recipient.
	FeeRecipient common.ExecutionAddress `json:"feeRecipient"`
	// StateRoot is the root of the state trie.
	StateRoot common.Bytes32 `json:"stateRoot"`
	// ReceiptsRoot is the root of the receipts trie.
	ReceiptsRoot common.Bytes32 `json:"receiptsRoot"`
	// LogsBloom is the bloom filter for the logs.
	LogsBloom bytesprimitives.B256 `json:"logsBloom"`
	// Random is the prevRandao value.
	Random common.Bytes32 `json:"prevRandao"`
	// Number is the block number.
	Number math.U64 `json:"blockNumber"`
	// GasLimit is the gas limit for the block.
	GasLimit math.U64 `json:"gasLimit"`
	// GasUsed is the amount of gas used in the block.
	GasUsed math.U64 `json:"gasUsed"`
	// Timestamp is the timestamp of the block.
	Timestamp math.U64 `json:"timestamp"`
	// ExtraData is the extra data of the block.
	ExtraData bytesprimitives.Bytes `json:"extraData"`
	// BaseFeePerGas is the base fee per gas.
	BaseFeePerGas *math.U256 `json:"baseFeePerGas"`
	// BlockHash is the hash of the block.
	BlockHash common.ExecutionHash `json:"blockHash"`
	// Transactions is the list of transactions in the block.
	Transactions engineprimitives.Transactions `json:"transactions"`
	// Withdrawals is the list of withdrawals in the block.
	Withdrawals []*engineprimitives.Withdrawal `json:"withdrawals"`
	// BlobGasUsed is the amount of blob gas used in the block.
	BlobGasUsed math.U64 `json:"blobGasUsed"`
	// ExcessBlobGas is the amount of excess blob gas in the block.
	ExcessBlobGas math.U64 `json:"excessBlobGas"`
}

// --- Local Alias Type ---
// TestExecPayload is our alias for ExecutionPayload.
type TestExecPayload ctypes.ExecutionPayload

// generateWithdrawals generates a slice of *engineprimitives.Withdrawal
// with a random length up to maxLen.
func generateWithdrawals(r *rand.Rand, maxLen int) []*engineprimitives.Withdrawal {
	n := r.Intn(maxLen + 1) // length between 0 and maxLen
	withdrawals := make([]*engineprimitives.Withdrawal, n)
	// For each element, use quick.Value to generate a withdrawal.
	withdrawalType := reflect.TypeOf(engineprimitives.Withdrawal{})
	for i := 0; i < n; i++ {
		v, ok := quick.Value(withdrawalType, r)
		if !ok {
			panic("failed to generate withdrawal")
		}
		w := v.Interface().(engineprimitives.Withdrawal)
		withdrawals[i] = &w
	}
	fmt.Println(withdrawals)
	return withdrawals
}

// Generate implements quick.Generator for *TestExecPayload.
func (p *TestExecPayload) Generate(r *rand.Rand, size int) reflect.Value {
	// Step 1: Generate a value for the helper struct, which contains only exported fields.
	var exp execPayloadExported
	v, ok := quick.Value(reflect.TypeOf(exp), r)
	if !ok {
		panic("failed to generate execPayloadExported")
	}
	exp = v.Interface().(execPayloadExported)

	// Step 2: Copy exported fields from exp into our alias.
	var tep TestExecPayload
	tep.ParentHash = exp.ParentHash
	tep.FeeRecipient = exp.FeeRecipient
	tep.StateRoot = exp.StateRoot
	tep.ReceiptsRoot = exp.ReceiptsRoot
	tep.LogsBloom = exp.LogsBloom
	tep.Random = exp.Random
	tep.Number = exp.Number
	tep.GasLimit = exp.GasLimit
	tep.GasUsed = exp.GasUsed
	tep.Timestamp = exp.Timestamp
	tep.ExtraData = exp.ExtraData
	tep.BaseFeePerGas = exp.BaseFeePerGas
	tep.BlockHash = exp.BlockHash
	tep.Transactions = exp.Transactions
	tep.Withdrawals = exp.Withdrawals
	tep.BlobGasUsed = exp.BlobGasUsed
	tep.ExcessBlobGas = exp.ExcessBlobGas

	// Step 3: Ensure that slices are non-nil. Default withdrawals generation only generates a maximum length of 1.
	const maxWithdrawalLen = 0
	tep.Withdrawals = generateWithdrawals(r, maxWithdrawalLen)

	if tep.Transactions == nil {
		tep.Transactions = engineprimitives.Transactions{}
	}

	if len(tep.Withdrawals) > 2 {
		fmt.Println(tep.Withdrawals)
	}

	// Step 4: Set the unexported forkVersion via the setter.
	// Convert our alias pointer to the production *ctypes.ExecutionPayload.
	orig := (*ctypes.ExecutionPayload)(&tep)
	supported := version.GetSupportedVersions()
	orig.SetForkVersion(supported[r.Intn(len(supported))])

	// Return a reflect.Value representing a pointer to our alias.
	return reflect.ValueOf(&tep)
}

func TestExecutionPayloadHashTreeRootZrnt(t *testing.T) {
	t.Parallel()
	f := func(testPayload *TestExecPayload, logsBloom [256]byte) bool {
		// Convert the generated value back to the production type.
		payload := (*ctypes.ExecutionPayload)(testPayload)
		//fmt.Printf("%+v\n", payload)
		payload.LogsBloom = logsBloom
		payload.BaseFeePerGas = math.NewU256(123)
		typeRoot := payload.HashTreeRoot()

		baseFeePerGas := zview.Uint256View{}
		baseFeePerGas.SetFromBig(payload.BaseFeePerGas.ToBig())
		zpayload := zdeneb.ExecutionPayload{
			ParentHash:    ztree.Root(payload.ParentHash),
			FeeRecipient:  zcommon.Eth1Address(payload.FeeRecipient),
			StateRoot:     ztree.Root(payload.StateRoot),
			ReceiptsRoot:  ztree.Root(payload.ReceiptsRoot),
			LogsBloom:     zcommon.LogsBloom(payload.LogsBloom),
			PrevRandao:    ztree.Root(payload.Random),
			BlockNumber:   zview.Uint64View(payload.Number),
			GasLimit:      zview.Uint64View(payload.GasLimit),
			GasUsed:       zview.Uint64View(payload.GasUsed),
			Timestamp:     zcommon.Timestamp(payload.Timestamp),
			ExtraData:     []byte(payload.ExtraData),
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     ztree.Root(payload.BlockHash),
			Transactions: *(*zcommon.PayloadTransactions)(
				unsafe.Pointer(&payload.Transactions)),
			Withdrawals:   *(*zcommon.Withdrawals)(unsafe.Pointer(&payload.Withdrawals)),
			BlobGasUsed:   zview.Uint64View(payload.BlobGasUsed.Unwrap()),
			ExcessBlobGas: zview.Uint64View(payload.ExcessBlobGas.Unwrap()),
		}

		zRoot := zpayload.HashTreeRoot(spec, hFn)
		containerRoot := payload.HashTreeRoot()

		return bytes.Equal(typeRoot[:], containerRoot[:]) &&
			bytes.Equal(typeRoot[:], zRoot[:])
	}
	if err := quick.Check(f, &c); err != nil {
		t.Error(err)
	}
}

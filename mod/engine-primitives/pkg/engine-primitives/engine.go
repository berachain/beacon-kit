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

//nolint:gochecknoglobals // alias.
package engineprimitives

import (
	"fmt"
	"math/big"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

// There are some types we can borrow from geth.
type (
	ClientVersionV1 = engine.ClientVersionV1
	ExecutableData  = engine.ExecutableData
)

var (
	ExecutableDataToBlock = engine.ExecutableDataToBlock
)

// ExecutableDataToBlock constructs a block from executable data.
// It verifies that the following fields:
//
//		len(extraData) <= 32
//		uncleHash = emptyUncleHash
//		difficulty = 0
//	 	if versionedHashes != nil, versionedHashes match to blob transactions
//
// and that the blockhash of the constructed block matches the parameters. Nil
// Withdrawals value will propagate through the returned block. Empty
// Withdrawals value must be passed via non-nil, length 0 value in params.
func ExecutableDataToBlock2(txs []*types.Transaction, params ExecutableData, versionedHashes []common.ExecutionHash, beaconRoot *common.ExecutionHash) (*types.Block, error) {
	if len(params.ExtraData) > 32 {
		return nil, fmt.Errorf("invalid extradata length: %v", len(params.ExtraData))
	}
	if len(params.LogsBloom) != 256 {
		return nil, fmt.Errorf("invalid logsBloom length: %v", len(params.LogsBloom))
	}
	// Check that baseFeePerGas is not negative or too big
	if params.BaseFeePerGas != nil && (params.BaseFeePerGas.Sign() == -1 || params.BaseFeePerGas.BitLen() > 256) {
		return nil, fmt.Errorf("invalid baseFeePerGas: %v", params.BaseFeePerGas)
	}
	var blobHashes []common.ExecutionHash
	for _, tx := range txs {
		blobHashes = append(blobHashes, tx.BlobHashes()...)
	}
	if len(blobHashes) != len(versionedHashes) {
		return nil, fmt.Errorf("invalid number of versionedHashes: %v blobHashes: %v", versionedHashes, blobHashes)
	}
	for i := 0; i < len(blobHashes); i++ {
		if blobHashes[i] != versionedHashes[i] {
			return nil, fmt.Errorf("invalid versionedHash at %v: %v blobHashes: %v", i, versionedHashes, blobHashes)
		}
	}
	// Only set withdrawalsRoot if it is non-nil. This allows CLs to use
	// ExecutableData before withdrawals are enabled by marshaling
	// Withdrawals as the json null value.
	var withdrawalsRoot *common.ExecutionHash
	if params.Withdrawals != nil {
		h := types.DeriveSha(types.Withdrawals(params.Withdrawals), trie.NewStackTrie(nil))
		withdrawalsRoot = &h
	}
	header := &types.Header{
		ParentHash:       params.ParentHash,
		UncleHash:        types.EmptyUncleHash,
		Coinbase:         params.FeeRecipient,
		Root:             params.StateRoot,
		TxHash:           types.DeriveSha(types.Transactions(txs), trie.NewStackTrie(nil)),
		ReceiptHash:      params.ReceiptsRoot,
		Bloom:            types.BytesToBloom(params.LogsBloom),
		Difficulty:       big.NewInt(0),
		Number:           new(big.Int).SetUint64(params.Number),
		GasLimit:         params.GasLimit,
		GasUsed:          params.GasUsed,
		Time:             params.Timestamp,
		BaseFee:          params.BaseFeePerGas,
		Extra:            params.ExtraData,
		MixDigest:        params.Random,
		WithdrawalsHash:  withdrawalsRoot,
		ExcessBlobGas:    params.ExcessBlobGas,
		BlobGasUsed:      params.BlobGasUsed,
		ParentBeaconRoot: beaconRoot,
	}
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: txs, Uncles: nil, Withdrawals: params.Withdrawals})
	if block.Hash() != params.BlockHash {
		return nil, fmt.Errorf("blockhash mismatch, want %x, got %x", params.BlockHash, block.Hash())
	}
	return block, nil
}

type PayloadStatusStr = string

var (
	// PayloadStatusValid is the status of a valid payload.
	PayloadStatusValid PayloadStatusStr = "VALID"
	// PayloadStatusInvalid is the status of an invalid payload.
	PayloadStatusInvalid PayloadStatusStr = "INVALID"
	// PayloadStatusSyncing is the status returned when the EL is syncing.
	PayloadStatusSyncing PayloadStatusStr = "SYNCING"
	// PayloadStatusAccepted is the status returned when the EL has accepted the
	// payload.
	PayloadStatusAccepted PayloadStatusStr = "ACCEPTED"
)

// ForkchoiceResponseV1 as per the EngineAPI Specification:
// https://github.com/ethereum/execution-apis/blob/main/src/engine/paris.md#response-2
//
//nolint:lll // link.
type ForkchoiceResponseV1 struct {
	// PayloadStatus is the payload status.
	PayloadStatus PayloadStatusV1 `json:"payloadStatus"`
	// PayloadID isthe identifier of the payload build process, it
	// can also be `nil`.
	PayloadID *PayloadID `json:"payloadId"`
}

// ForkchoicStateV1 as per the EngineAPI Specification:
// https://github.com/ethereum/execution-apis/blob/main/src/engine/paris.md#forkchoicestatev1
//
//nolint:lll // link.
type ForkchoiceStateV1 struct {
	// HeadBlockHash is the desired block hash of the head of the canonical
	// chain.
	HeadBlockHash common.ExecutionHash `json:"headBlockHash"`
	// SafeBlockHash is  the "safe" block hash of the canonical chain under
	// certain
	// synchrony and honesty assumptions. This value MUST be either equal to
	// or an ancestor of `HeadBlockHash`.
	SafeBlockHash common.ExecutionHash `json:"safeBlockHash"`
	// FinalizedBlockHash is the desired block hash of the most recent finalized
	// block
	FinalizedBlockHash common.ExecutionHash `json:"finalizedBlockHash"`
}

// PayloadStatusV1 represents the status of a payload as per the EngineAPI
// Specification. For more details, see:
// https://github.com/ethereum/execution-apis/blob/main/src/engine/paris.md#payloadstatusv1
//
//nolint:lll // link.
type PayloadStatusV1 struct {
	// Status string of the payload.
	Status string `json:"status"`
	// LatestValidHash is the hash of the most recent valid block
	// in the branch defined by payload and its ancestors
	LatestValidHash *common.ExecutionHash `json:"latestValidHash"`
	// ValidationError is a message providing additional details on
	// the validation error if the payload is classified as
	// INVALID or INVALID_BLOCK_HASH
	ValidationError *string `json:"validationError"`
}

// PayloadID is an identifier for the payload build process.
type PayloadID = bytes.B8

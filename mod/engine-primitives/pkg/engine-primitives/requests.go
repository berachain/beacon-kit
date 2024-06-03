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

package engineprimitives

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

// NewPayloadRequest as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#modified-newpayloadrequest
//
//nolint:lll
type NewPayloadRequest[
	ExecutionPayloadT interface {
		Empty(uint32) ExecutionPayloadT
		Version() uint32
		ExecutionPayload[WithdrawalT]
	},
	WithdrawalT interface {
		GetIndex() math.U64
		GetAmount() math.U64
		GetAddress() common.ExecutionAddress
		GetValidatorIndex() math.U64
	},
] struct {
	// ExecutionPayload is the payload to the execution client.
	ExecutionPayload ExecutionPayloadT
	// VersionedHashes is the versioned hashes of the execution payload.
	VersionedHashes []common.ExecutionHash
	// ParentBeaconBlockRoot is the root of the parent beacon block.
	ParentBeaconBlockRoot *primitives.Root
	// Optimistic is a flag that indicates if the payload should be
	// optimistically deemed valid. This is useful during syncing.
	Optimistic bool
}

// BuildNewPayloadRequest builds a new payload request.
func BuildNewPayloadRequest[
	ExecutionPayloadT interface {
		Empty(uint32) ExecutionPayloadT
		Version() uint32
		ExecutionPayload[WithdrawalT]
	},
	WithdrawalT interface {
		GetIndex() math.U64
		GetAmount() math.U64
		GetAddress() common.ExecutionAddress
		GetValidatorIndex() math.U64
	},
](
	executionPayload ExecutionPayloadT,
	versionedHashes []common.ExecutionHash,
	parentBeaconBlockRoot *primitives.Root,
	optimistic bool,
) *NewPayloadRequest[ExecutionPayloadT, WithdrawalT] {
	return &NewPayloadRequest[ExecutionPayloadT, WithdrawalT]{
		ExecutionPayload:      executionPayload,
		VersionedHashes:       versionedHashes,
		ParentBeaconBlockRoot: parentBeaconBlockRoot,
		Optimistic:            optimistic,
	}
}

// HasHValidVersionAndBlockHashes checks if the version and block hashes are
// valid.
// As per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_block_hash
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_versioned_hashes
//
//nolint:lll
func (n *NewPayloadRequest[ExecutionPayloadT, WithdrawalT]) HasValidVersionedAndBlockHashes() error {
	var (
		gethWithdrawals []*types.Withdrawal
		withdrawalsHash *common.ExecutionHash
		blobHashes      = make([]common.ExecutionHash, 0)
		payload         = n.ExecutionPayload
		txs             = make(
			[]*types.Transaction,
			len(payload.GetTransactions()),
		)
	)

	// Extracts and validates the blob hashes from the transactions in the
	// execution payload.
	for i, encTx := range payload.GetTransactions() {
		var tx types.Transaction
		if err := tx.UnmarshalBinary(encTx); err != nil {
			return errors.Wrapf(err, "invalid transaction %d", i)
		}
		blobHashes = append(blobHashes, tx.BlobHashes()...)
		txs[i] = &tx
	}

	// Check if the number of blob hashes matches the number of versioned
	// hashes.
	if len(blobHashes) != len(n.VersionedHashes) {
		return errors.Wrapf(
			ErrMismatchedNumVersionedHashes,
			"expected %d, got %d",
			len(n.VersionedHashes),
			len(blobHashes),
		)
	}

	// Validate each blob hash against the corresponding versioned hash.
	for i, blobHash := range blobHashes {
		if blobHash != n.VersionedHashes[i] {
			return errors.Wrapf(
				ErrInvalidVersionedHash,
				"index %d: expected %v, got %v",
				i,
				n.VersionedHashes[i],
				blobHash,
			)
		}
	}

	// Construct the withdrawals and withdrawals hash.
	if payload.GetWithdrawals() != nil {
		gethWithdrawals = make(
			[]*types.Withdrawal,
			len(payload.GetWithdrawals()),
		)
		for i, wd := range payload.GetWithdrawals() {
			gethWithdrawals[i] = &types.Withdrawal{
				Index:     wd.GetIndex().Unwrap(),
				Amount:    wd.GetAmount().Unwrap(),
				Address:   wd.GetAddress(),
				Validator: wd.GetValidatorIndex().Unwrap(),
			}
		}
		h := types.DeriveSha(
			types.Withdrawals(gethWithdrawals),
			trie.NewStackTrie(nil),
		)
		withdrawalsHash = &h
	}

	// Verify that the payload is telling the truth about it's block hash.
	if block := types.NewBlockWithHeader(
		&types.Header{
			ParentHash:       payload.GetParentHash(),
			UncleHash:        types.EmptyUncleHash,
			Coinbase:         payload.GetFeeRecipient(),
			Root:             common.ExecutionHash(payload.GetStateRoot()),
			TxHash:           types.DeriveSha(types.Transactions(txs), trie.NewStackTrie(nil)),
			ReceiptHash:      common.ExecutionHash(payload.GetReceiptsRoot()),
			Bloom:            types.BytesToBloom(payload.GetLogsBloom()),
			Difficulty:       big.NewInt(0),
			Number:           new(big.Int).SetUint64(payload.GetNumber().Unwrap()),
			GasLimit:         payload.GetGasLimit().Unwrap(),
			GasUsed:          payload.GetGasUsed().Unwrap(),
			Time:             payload.GetTimestamp().Unwrap(),
			BaseFee:          payload.GetBaseFeePerGas().UnwrapBig(),
			Extra:            payload.GetExtraData(),
			MixDigest:        common.ExecutionHash(payload.GetPrevRandao()),
			WithdrawalsHash:  withdrawalsHash,
			ExcessBlobGas:    payload.GetExcessBlobGas().UnwrapPtr(),
			BlobGasUsed:      payload.GetBlobGasUsed().UnwrapPtr(),
			ParentBeaconRoot: (*common.ExecutionHash)(n.ParentBeaconBlockRoot),
		},
	).WithBody(types.Body{
		Transactions: txs, Uncles: nil, Withdrawals: gethWithdrawals,
	}); block.Hash() != payload.GetBlockHash() {
		return errors.Wrapf(ErrPayloadBlockHashMismatch,
			"%x, got %x",
			payload.GetBlockHash(), block.Hash(),
		)
	}
	return nil
}

// ForkchoiceUpdateRequest.
type ForkchoiceUpdateRequest struct {
	// State is the forkchoice state.
	State *ForkchoiceStateV1
	// PayloadAttributes is the payload attributer.
	PayloadAttributes PayloadAttributer
	// ForkVersion is the fork version that we
	// are going to be submitting for.
	ForkVersion uint32
}

// BuildForkchoiceUpdateRequest builds a forkchoice update request.
func BuildForkchoiceUpdateRequest(
	state *ForkchoiceStateV1,
	payloadAttributes PayloadAttributer,
	forkVersion uint32,
) *ForkchoiceUpdateRequest {
	return &ForkchoiceUpdateRequest{
		State:             state,
		PayloadAttributes: payloadAttributes,
		ForkVersion:       forkVersion,
	}
}

// GetPayloadRequest represents a request to get a payload.
type GetPayloadRequest struct {
	// PayloadID is the payload ID.
	PayloadID PayloadID
	// ForkVersion is the fork version that we are
	// currently on.
	ForkVersion uint32
}

// BuildGetPayloadRequest builds a get payload request.
func BuildGetPayloadRequest(
	payloadID PayloadID,
	forkVersion uint32,
) *GetPayloadRequest {
	return &GetPayloadRequest{
		PayloadID:   payloadID,
		ForkVersion: forkVersion,
	}
}

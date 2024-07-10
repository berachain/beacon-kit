// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package engineprimitives

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Newa as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#modified-newpayloadrequest
//
//nolint:lll
type NewPayloadRequest[
	ExecutionPayloadT interface {
		constraints.ForkTyped[ExecutionPayloadT]
		GetPrevRandao() common.Bytes32
		GetBlockHash() gethprimitives.ExecutionHash
		GetParentHash() gethprimitives.ExecutionHash
		GetNumber() math.U64
		GetGasLimit() math.U64
		GetGasUsed() math.U64
		GetTimestamp() math.U64
		GetExtraData() []byte
		GetBaseFeePerGas() math.Wei
		GetFeeRecipient() gethprimitives.ExecutionAddress
		GetStateRoot() common.Bytes32
		GetReceiptsRoot() common.Bytes32
		GetLogsBloom() []byte
		GetBlobGasUsed() math.U64
		GetExcessBlobGas() math.U64
		GetWithdrawals() []WithdrawalT
		GetTransactions() [][]byte
	},
	WithdrawalT interface {
		GetIndex() math.U64
		GetAmount() math.U64
		GetAddress() gethprimitives.ExecutionAddress
		GetValidatorIndex() math.U64
	},
] struct {
	// ExecutionPayload is the payload to the execution client.
	ExecutionPayload ExecutionPayloadT
	// VersionedHashes is the versioned hashes of the execution payload.
	VersionedHashes []gethprimitives.ExecutionHash
	// ParentBeaconBlockRoot is the root of the parent beacon block.
	ParentBeaconBlockRoot *common.Root
	// Optimistic is a flag that indicates if the payload should be
	// optimistically deemed valid. This is useful during syncing.
	Optimistic bool
}

// BuildNewPayloadRequest builds a new payload request.
func BuildNewPayloadRequest[
	ExecutionPayloadT interface {
		constraints.ForkTyped[ExecutionPayloadT]
		GetPrevRandao() common.Bytes32
		GetBlockHash() gethprimitives.ExecutionHash
		GetParentHash() gethprimitives.ExecutionHash
		GetNumber() math.U64
		GetGasLimit() math.U64
		GetGasUsed() math.U64
		GetTimestamp() math.U64
		GetExtraData() []byte
		GetBaseFeePerGas() math.Wei
		GetFeeRecipient() gethprimitives.ExecutionAddress
		GetStateRoot() common.Bytes32
		GetReceiptsRoot() common.Bytes32
		GetLogsBloom() []byte
		GetBlobGasUsed() math.U64
		GetExcessBlobGas() math.U64
		GetWithdrawals() []WithdrawalT
		GetTransactions() [][]byte
	},
	WithdrawalT interface {
		GetIndex() math.U64
		GetAmount() math.U64
		GetAddress() gethprimitives.ExecutionAddress
		GetValidatorIndex() math.U64
	},
](
	executionPayload ExecutionPayloadT,
	versionedHashes []gethprimitives.ExecutionHash,
	parentBeaconBlockRoot *common.Root,
	optimistic bool,
) *NewPayloadRequest[ExecutionPayloadT, WithdrawalT] {
	return &NewPayloadRequest[ExecutionPayloadT, WithdrawalT]{
		ExecutionPayload:      executionPayload,
		VersionedHashes:       versionedHashes,
		ParentBeaconBlockRoot: parentBeaconBlockRoot,
		Optimistic:            optimistic,
	}
}

// HasValidVersionedAndBlockHashes checks if the version and block hashes are
// valid.
// As per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_block_hash
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_versioned_hashes
//
//nolint:lll
func (n *NewPayloadRequest[ExecutionPayloadT, WithdrawalT]) HasValidVersionedAndBlockHashes() error {
	var (
		gethWithdrawals []*gethprimitives.Withdrawal
		withdrawalsHash *gethprimitives.ExecutionHash
		blobHashes      = make([]gethprimitives.ExecutionHash, 0)
		payload         = n.ExecutionPayload
		txs             = make(
			[]*gethprimitives.Transaction,
			len(payload.GetTransactions()),
		)
	)

	// Extracts and validates the blob hashes from the transactions in the
	// execution payload.
	for i, encTx := range payload.GetTransactions() {
		var tx gethprimitives.Transaction
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
			[]*gethprimitives.Withdrawal,
			len(payload.GetWithdrawals()),
		)
		for i, wd := range payload.GetWithdrawals() {
			gethWithdrawals[i] = &gethprimitives.Withdrawal{
				Index:     wd.GetIndex().Unwrap(),
				Amount:    wd.GetAmount().Unwrap(),
				Address:   wd.GetAddress(),
				Validator: wd.GetValidatorIndex().Unwrap(),
			}
		}
		h := gethprimitives.DeriveSha(
			gethprimitives.Withdrawals(gethWithdrawals),
			gethprimitives.NewStackTrie(nil),
		)
		withdrawalsHash = &h
	}

	// Verify that the payload is telling the truth about it's block hash.
	if block := gethprimitives.NewBlockWithHeader(
		&gethprimitives.Header{
			ParentHash:       payload.GetParentHash(),
			UncleHash:        gethprimitives.EmptyUncleHash,
			Coinbase:         payload.GetFeeRecipient(),
			Root:             gethprimitives.ExecutionHash(payload.GetStateRoot()),
			TxHash:           gethprimitives.DeriveSha(gethprimitives.Transactions(txs), gethprimitives.NewStackTrie(nil)),
			ReceiptHash:      gethprimitives.ExecutionHash(payload.GetReceiptsRoot()),
			Bloom:            gethprimitives.BytesToBloom(payload.GetLogsBloom()),
			Difficulty:       big.NewInt(0),
			Number:           new(big.Int).SetUint64(payload.GetNumber().Unwrap()),
			GasLimit:         payload.GetGasLimit().Unwrap(),
			GasUsed:          payload.GetGasUsed().Unwrap(),
			Time:             payload.GetTimestamp().Unwrap(),
			BaseFee:          payload.GetBaseFeePerGas().UnwrapBig(),
			Extra:            payload.GetExtraData(),
			MixDigest:        gethprimitives.ExecutionHash(payload.GetPrevRandao()),
			WithdrawalsHash:  withdrawalsHash,
			ExcessBlobGas:    payload.GetExcessBlobGas().UnwrapPtr(),
			BlobGasUsed:      payload.GetBlobGasUsed().UnwrapPtr(),
			ParentBeaconRoot: (*gethprimitives.ExecutionHash)(n.ParentBeaconBlockRoot),
		},
	).WithBody(gethprimitives.Body{
		Transactions: txs, Uncles: nil, Withdrawals: gethWithdrawals,
	}); block.Hash() != payload.GetBlockHash() {
		return errors.Wrapf(ErrPayloadBlockHashMismatch,
			"%x, got %x",
			payload.GetBlockHash(), block.Hash(),
		)
	}
	return nil
}

type ForkchoiceUpdateRequest[PayloadAttributesT any] struct {
	// State is the forkchoice state.
	State *ForkchoiceStateV1
	// PayloadAttributes is the payload attributer.
	PayloadAttributes PayloadAttributesT
	// ForkVersion is the fork version that we
	// are going to be submitting for.
	ForkVersion uint32
}

// BuildForkchoiceUpdateRequest builds a forkchoice update request.
func BuildForkchoiceUpdateRequest[
	PayloadAttributesT any,
](
	state *ForkchoiceStateV1,
	payloadAttributes PayloadAttributesT,
	forkVersion uint32,
) *ForkchoiceUpdateRequest[PayloadAttributesT] {
	return &ForkchoiceUpdateRequest[PayloadAttributesT]{
		State:             state,
		PayloadAttributes: payloadAttributes,
		ForkVersion:       forkVersion,
	}
}

// BuildForkchoiceUpdateRequestNoAttrs builds a forkchoice update request
// without
// any attributes.
func BuildForkchoiceUpdateRequestNoAttrs[
	PayloadAttributesT any,
](
	state *ForkchoiceStateV1,
	forkVersion uint32,
) *ForkchoiceUpdateRequest[PayloadAttributesT] {
	return &ForkchoiceUpdateRequest[PayloadAttributesT]{
		State:       state,
		ForkVersion: forkVersion,
	}
}

// GetPayloadRequest represents a request to get a payload.
type GetPayloadRequest[PayloadIDT ~[8]byte] struct {
	// PayloadID is the payload ID.
	PayloadID PayloadIDT
	// ForkVersion is the fork version that we are
	// currently on.
	ForkVersion uint32
}

// BuildGetPayloadRequest builds a get payload request.
func BuildGetPayloadRequest[PayloadIDT ~[8]byte](
	payloadID PayloadIDT,
	forkVersion uint32,
) *GetPayloadRequest[PayloadIDT] {
	return &GetPayloadRequest[PayloadIDT]{
		PayloadID:   payloadID,
		ForkVersion: forkVersion,
	}
}

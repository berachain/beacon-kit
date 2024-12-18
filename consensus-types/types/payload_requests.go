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

package types

import (
	"math/big"
	"unsafe"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
)

// NewPayloadRequest as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#modified-newpayloadrequest
type NewPayloadRequest struct {
	// ExecutionPayload is the payload to the execution client.
	ExecutionPayload *ExecutionPayload
	// VersionedHashes is the versioned hashes of the execution payload.
	VersionedHashes []common.ExecutionHash
	// ParentBeaconBlockRoot is the root of the parent beacon block.
	ParentBeaconBlockRoot *common.Root
	// Optimistic is a flag that indicates if the payload should be
	// optimistically deemed valid. This is useful during syncing.
	Optimistic bool
}

// BuildNewPayloadRequest builds a new payload request.
func BuildNewPayloadRequest(
	executionPayload *ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBeaconBlockRoot *common.Root,
	optimistic bool,
) *NewPayloadRequest {
	return &NewPayloadRequest{
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
func (n *NewPayloadRequest) HasValidVersionedAndBlockHashes() error {
	var (
		blobHashes = make([]gethprimitives.ExecutionHash, 0)
		payload    = n.ExecutionPayload
		txs        = make(
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
			engineprimitives.ErrMismatchedNumVersionedHashes,
			"expected %d, got %d",
			len(n.VersionedHashes),
			len(blobHashes),
		)
	}

	// Validate each blob hash against the corresponding versioned hash.
	for i, blobHash := range blobHashes {
		if common.ExecutionHash(blobHash) != n.VersionedHashes[i] {
			return errors.Wrapf(
				engineprimitives.ErrInvalidVersionedHash,
				"index %d: expected %v, got %v",
				i,
				n.VersionedHashes[i],
				blobHash,
			)
		}
	}

	wds := payload.GetWithdrawals()
	withdrawalsHash := gethprimitives.DeriveSha(
		wds,
		gethprimitives.NewStackTrie(nil),
	)

	// Verify that the payload is telling the truth about it's block hash.
	//#nosec:G103 // its okay.
	if block := gethprimitives.NewBlockWithHeader(
		&gethprimitives.Header{
			ParentHash:       gethprimitives.ExecutionHash(payload.GetParentHash()),
			UncleHash:        gethprimitives.EmptyUncleHash,
			Coinbase:         gethprimitives.ExecutionAddress(payload.GetFeeRecipient()),
			Root:             gethprimitives.ExecutionHash(payload.GetStateRoot()),
			TxHash:           gethprimitives.DeriveSha(gethprimitives.Transactions(txs), gethprimitives.NewStackTrie(nil)),
			ReceiptHash:      gethprimitives.ExecutionHash(payload.GetReceiptsRoot()),
			Bloom:            gethprimitives.LogsBloom(payload.GetLogsBloom()),
			Difficulty:       big.NewInt(0),
			Number:           new(big.Int).SetUint64(payload.GetNumber().Unwrap()),
			GasLimit:         payload.GetGasLimit().Unwrap(),
			GasUsed:          payload.GetGasUsed().Unwrap(),
			Time:             payload.GetTimestamp().Unwrap(),
			BaseFee:          payload.GetBaseFeePerGas().ToBig(),
			Extra:            payload.GetExtraData(),
			MixDigest:        gethprimitives.ExecutionHash(payload.GetPrevRandao()),
			WithdrawalsHash:  &withdrawalsHash,
			ExcessBlobGas:    payload.GetExcessBlobGas().UnwrapPtr(),
			BlobGasUsed:      payload.GetBlobGasUsed().UnwrapPtr(),
			ParentBeaconRoot: (*gethprimitives.ExecutionHash)(n.ParentBeaconBlockRoot),
		},
	).WithBody(gethprimitives.Body{
		Transactions: txs, Uncles: nil, Withdrawals: *(*gethprimitives.Withdrawals)(unsafe.Pointer(&wds)),
	}); common.ExecutionHash(block.Hash()) != payload.GetBlockHash() {
		return errors.Wrapf(engineprimitives.ErrPayloadBlockHashMismatch,
			"%x, got %x",
			payload.GetBlockHash(), block.Hash(),
		)
	}
	return nil
}

type ForkchoiceUpdateRequest struct {
	// State is the forkchoice state.
	State *engineprimitives.ForkchoiceStateV1
	// PayloadAttributes is the payload attributer.
	PayloadAttributes *engineprimitives.PayloadAttributes
	// ForkVersion is the fork version that we
	// are going to be submitting for.
	ForkVersion uint32
}

// BuildForkchoiceUpdateRequest builds a forkchoice update request.
func BuildForkchoiceUpdateRequest(
	state *engineprimitives.ForkchoiceStateV1,
	payloadAttributes *engineprimitives.PayloadAttributes,
	forkVersion uint32,
) *ForkchoiceUpdateRequest {
	return &ForkchoiceUpdateRequest{
		State:             state,
		PayloadAttributes: payloadAttributes,
		ForkVersion:       forkVersion,
	}
}

// BuildForkchoiceUpdateRequestNoAttrs builds a forkchoice update request
// without
// any attributes.
func BuildForkchoiceUpdateRequestNoAttrs(
	state *engineprimitives.ForkchoiceStateV1,
	forkVersion uint32,
) *ForkchoiceUpdateRequest {
	return &ForkchoiceUpdateRequest{
		State:       state,
		ForkVersion: forkVersion,
	}
}

// GetPayloadRequest represents a request to get a payload.
type GetPayloadRequest struct {
	// PayloadID is the payload ID.
	PayloadID engineprimitives.PayloadID
	// ForkVersion is the fork version that we are
	// currently on.
	ForkVersion uint32
}

// BuildGetPayloadRequest builds a get payload request.
func BuildGetPayloadRequest(
	payloadID engineprimitives.PayloadID,
	forkVersion uint32,
) *GetPayloadRequest {
	return &GetPayloadRequest{
		PayloadID:   payloadID,
		ForkVersion: forkVersion,
	}
}

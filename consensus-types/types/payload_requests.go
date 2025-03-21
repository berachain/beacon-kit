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

package types

import (
	"math/big"
	"unsafe"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/version"
)

// NewPayloadRequest as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/beacon-chain.md#modified-newpayloadrequest
type NewPayloadRequest interface {
	constraints.Versionable
	HasValidVersionedAndBlockHashes() error
	GetExecutionPayload() *ExecutionPayload
	GetVersionedHashes() []common.ExecutionHash
	GetParentBeaconBlockRoot() *common.Root
	GetExecutionRequests() ([]EncodedExecutionRequest, error)
}

type newPayloadRequest struct {
	constraints.Versionable
	// executionPayload is the payload to the execution client.
	executionPayload *ExecutionPayload
	// versionedHashes is the versioned hashes of the execution payload.
	versionedHashes []common.ExecutionHash
	// parentBeaconBlockRoot is the root of the parent beacon block.
	parentBeaconBlockRoot *common.Root
	// ExecutionRequests is introduced in Pectra. It is only non-nil after Pectra.
	executionRequests []EncodedExecutionRequest
}

// BuildNewPayloadRequest builds a new payload request.
func BuildNewPayloadRequest(
	executionPayload *ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBeaconBlockRoot *common.Root,
) NewPayloadRequest {
	return &newPayloadRequest{
		Versionable:           NewVersionable(executionPayload.GetForkVersion()),
		executionPayload:      executionPayload,
		versionedHashes:       versionedHashes,
		parentBeaconBlockRoot: parentBeaconBlockRoot,
	}
}

// BuildNewPayloadRequestWithExecutionRequests builds a new payload post-electra
func BuildNewPayloadRequestWithExecutionRequests(
	executionPayload *ExecutionPayload,
	versionedHashes []common.ExecutionHash,
	parentBeaconBlockRoot *common.Root,
	executionRequests []EncodedExecutionRequest,
) NewPayloadRequest {
	return &newPayloadRequest{
		Versionable:           NewVersionable(executionPayload.GetForkVersion()),
		executionPayload:      executionPayload,
		versionedHashes:       versionedHashes,
		parentBeaconBlockRoot: parentBeaconBlockRoot,
		executionRequests:     executionRequests,
	}
}

func (n *newPayloadRequest) GetExecutionPayload() *ExecutionPayload {
	return n.executionPayload
}

func (n *newPayloadRequest) GetVersionedHashes() []common.ExecutionHash {
	return n.versionedHashes
}

func (n *newPayloadRequest) GetParentBeaconBlockRoot() *common.Root {
	return n.parentBeaconBlockRoot
}

func (n *newPayloadRequest) GetExecutionRequests() ([]EncodedExecutionRequest, error) {
	if version.EqualOrAfter(n.GetForkVersion(), version.Electra()) {
		return nil, ErrForkVersionNotSupported
	}
	if n.executionRequests == nil {
		return nil, errors.Wrap(ErrNilValue, "executionRequests cannot be nil")
	}
	return n.executionRequests, nil
}

// HasValidVersionedAndBlockHashes checks if the version and block hashes are
// valid.
// As per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_block_hash
// https://github.com/ethereum/consensus-specs/blob/v1.4.0-beta.2/specs/deneb/beacon-chain.md#is_valid_versioned_hashes
func (n *newPayloadRequest) HasValidVersionedAndBlockHashes() error {
	block, blobHashes, err := MakeEthBlock(n.executionPayload, n.parentBeaconBlockRoot)
	if err != nil {
		return err
	}

	// Validate the blob hashes from the transactions in the execution payload.
	// Check if the number of blob hashes matches the number of versioned hashes.
	if len(blobHashes) != len(n.versionedHashes) {
		return errors.Wrapf(
			engineprimitives.ErrMismatchedNumVersionedHashes,
			"expected %d, got %d",
			len(blobHashes),
			len(n.versionedHashes),
		)
	}

	// Validate each blob hash against the corresponding versioned hash.
	for i, blobHash := range blobHashes {
		if common.ExecutionHash(blobHash) != n.versionedHashes[i] {
			return errors.Wrapf(
				engineprimitives.ErrInvalidVersionedHash,
				"index %d: expected %v, got %v",
				i,
				blobHash,
				n.versionedHashes[i],
			)
		}
	}

	// Verify that the payload is telling the truth about its block hash.
	if common.ExecutionHash(block.Hash()) != n.executionPayload.GetBlockHash() {
		return errors.Wrapf(engineprimitives.ErrPayloadBlockHashMismatch,
			"expected %x, got %x",
			block.Hash(), n.executionPayload.GetBlockHash(),
		)
	}
	return nil
}

// MakeEthBlock builds an Ethereum block out of given payload and parent block root.
// It also returns blobHashes out of payload to ease up checks.
// Use MakeEthBlockWithExecutionRequests after Pectra.
func MakeEthBlock(
	payload *ExecutionPayload,
	parentBeaconBlockRoot *common.Root,
) (*gethprimitives.Block,
	[]gethprimitives.ExecutionHash,
	error) {
	return makeEthBlock(payload, parentBeaconBlockRoot, nil)
}

// MakeEthBlockWithExecutionRequests is MakeEthBlock with support for executionRequests which is needed post-pectra.
func MakeEthBlockWithExecutionRequests(
	payload *ExecutionPayload,
	parentBeaconBlockRoot *common.Root,
	executionRequests []EncodedExecutionRequest,
) (*gethprimitives.Block,
	[]gethprimitives.ExecutionHash,
	error) {
	return makeEthBlock(payload, parentBeaconBlockRoot, executionRequests)
}

// makeEthBlock builds an Ethereum block out of given payload and parent block root.
// It also returns blobHashes out of payload to ease up checks.
func makeEthBlock(
	payload *ExecutionPayload,
	parentBeaconBlockRoot *common.Root,
	executionRequests []EncodedExecutionRequest,
) (
	*gethprimitives.Block,
	[]gethprimitives.ExecutionHash,
	error,
) {
	var (
		txs        = make([]*gethprimitives.Transaction, 0, len(payload.GetTransactions()))
		blobHashes = make([]gethprimitives.ExecutionHash, 0)
	)

	for i, encTx := range payload.GetTransactions() {
		var tx gethprimitives.Transaction
		if err := tx.UnmarshalBinary(encTx); err != nil {
			return nil, nil, errors.Wrapf(err, "invalid transaction %d", i)
		}
		txs = append(txs, &tx)
		blobHashes = append(blobHashes, tx.BlobHashes()...)
	}

	wds := payload.GetWithdrawals()
	withdrawalsHash := gethprimitives.DeriveSha(wds, gethprimitives.NewStackTrie(nil))

	blkHeader := &gethprimitives.Header{
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
		ParentBeaconRoot: (*gethprimitives.ExecutionHash)(parentBeaconBlockRoot),
	}

	if version.EqualOrAfter(payload.GetForkVersion(), version.Electra()) {
		if executionRequests == nil {
			return nil, nil, ErrNilValue
		}
		result := make([][]byte, len(executionRequests))
		for i, req := range executionRequests {
			result[i] = req // conversion from ExecutionRequest to []byte
		}
		reqHash := gethprimitives.CalcRequestsHash(result)
		blkHeader.RequestsHash = &reqHash
	}

	block := gethprimitives.NewBlockWithHeader(blkHeader).WithBody(
		gethprimitives.Body{
			Transactions: txs,
			Uncles:       nil,
			Withdrawals:  *(*gethprimitives.Withdrawals)(unsafe.Pointer(&wds)), //#nosec:G103 // its okay.
		},
	)
	return block, blobHashes, nil
}

type ForkchoiceUpdateRequest struct {
	// State is the forkchoice state.
	State *engineprimitives.ForkchoiceStateV1
	// PayloadAttributes is the payload attributer.
	PayloadAttributes *engineprimitives.PayloadAttributes
	// ForkVersion is the fork version that we
	// are going to be submitting for.
	ForkVersion common.Version
}

// BuildForkchoiceUpdateRequest builds a forkchoice update request.
func BuildForkchoiceUpdateRequest(
	state *engineprimitives.ForkchoiceStateV1,
	payloadAttributes *engineprimitives.PayloadAttributes,
	forkVersion common.Version,
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
	forkVersion common.Version,
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
	ForkVersion common.Version
}

// BuildGetPayloadRequest builds a get payload request.
func BuildGetPayloadRequest(
	payloadID engineprimitives.PayloadID,
	forkVersion common.Version,
) *GetPayloadRequest {
	return &GetPayloadRequest{
		PayloadID:   payloadID,
		ForkVersion: forkVersion,
	}
}

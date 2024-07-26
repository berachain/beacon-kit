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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"encoding/json"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// ExecutionPayloadHeader is the execution header payload of Deneb.
//
//nolint:lll
type ExecutionPayloadHeader struct {
	// TODO: Enable
	//
	// // Metadata
	// //
	// // version is the fork version of the execution payload header.
	// version uint32

	// Contents
	//
	// ParentHash is the hash of the parent block.
	ParentHash gethprimitives.ExecutionHash `json:"parentHash"       gencodec:"required"`
	// FeeRecipient is the address of the fee recipient.
	FeeRecipient gethprimitives.ExecutionAddress `json:"feeRecipient"     gencodec:"required"`
	// StateRoot is the root of the state trie.
	StateRoot common.Bytes32 `json:"stateRoot"        gencodec:"required"`
	// ReceiptsRoot is the root of the receipts trie.
	ReceiptsRoot common.Bytes32 `json:"receiptsRoot"     gencodec:"required"`
	// LogsBloom is the bloom filter for the logs.
	LogsBloom bytes.B256 `json:"logsBloom"        gencodec:"required"`
	// Random is the prevRandao value.
	Random common.Bytes32 `json:"prevRandao"       gencodec:"required"`
	// Number is the block number.
	Number math.U64 `json:"blockNumber"      gencodec:"required"`
	// GasLimit is the gas limit for the block.
	GasLimit math.U64 `json:"gasLimit"         gencodec:"required"`
	// GasUsed is the amount of gas used in the block.
	GasUsed math.U64 `json:"gasUsed"          gencodec:"required"`
	// Timestamp is the timestamp of the block.
	Timestamp math.U64 `json:"timestamp"        gencodec:"required"`
	// ExtraData is the extra data of the block.
	ExtraData bytes.Bytes `json:"extraData"        gencodec:"required"`
	// BaseFeePerGas is the base fee per gas.
	BaseFeePerGas math.Wei `json:"baseFeePerGas"    gencodec:"required"`
	// BlockHash is the hash of the block.
	BlockHash gethprimitives.ExecutionHash `json:"blockHash"        gencodec:"required"`
	// TransactionsRoot is the root of the transaction trie.
	TransactionsRoot common.Root `json:"transactionsRoot" gencodec:"required"`
	// WithdrawalsRoot is the root of the withdrawals trie.
	WithdrawalsRoot common.Root `json:"withdrawalsRoot"`
	// BlobGasUsed is the amount of blob gas used in the block.
	BlobGasUsed math.U64 `json:"blobGasUsed"`
	// ExcessBlobGas is the amount of excess blob gas in the block.
	ExcessBlobGas math.U64 `json:"excessBlobGas"`
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (h *ExecutionPayloadHeader) Empty(
	forkVersion uint32,
) *ExecutionPayloadHeader {
	// TODO: figure out how to use.
	_ = forkVersion
	return &ExecutionPayloadHeader{}
}

// NewFromSSZ returns a new ExecutionPayloadHeader from the given SSZ bytes.
func (h *ExecutionPayloadHeader) NewFromSSZ(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	h = h.Empty(forkVersion)
	return h, h.UnmarshalSSZ(bz)
}

// NewFromJSON returns a new ExecutionPayloadHeader from the given JSON bytes.
func (h *ExecutionPayloadHeader) NewFromJSON(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	h = h.Empty(forkVersion)
	return h, json.Unmarshal(bz, h)
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns either the static size of the object if fixed == true, or
// the total size otherwise.
func (h *ExecutionPayloadHeader) SizeSSZ(fixed bool) uint32 {
	//nolint:mnd // todo fix.
	var size = uint32(584)
	if fixed {
		return size
	}
	size += ssz.SizeDynamicBytes(h.ExtraData)

	return size
}

// DefineSSZ defines how an object is encoded/decoded.
func (h *ExecutionPayloadHeader) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineStaticBytes(
		codec,
		&h.ParentHash,
	) // Field  ( 0) -       ParentHash -  32 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.FeeRecipient,
	) // Field  ( 1) -     FeeRecipient -  20 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.StateRoot,
	) // Field  ( 2) -        StateRoot -  32 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.ReceiptsRoot,
	) // Field  ( 3) -     ReceiptsRoot -  32 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.LogsBloom,
	) // Field  ( 4) -        LogsBloom - 256 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.Random,
	) // Field  ( 5) -       PrevRandao -  32 bytes
	ssz.DefineUint64(
		codec,
		&h.Number,
	) // Field  ( 6) -      BlockNumber -   8 bytes
	ssz.DefineUint64(
		codec,
		&h.GasLimit,
	) // Field  ( 7) -         GasLimit -   8 bytes
	ssz.DefineUint64(
		codec,
		&h.GasUsed,
	) // Field  ( 8) -          GasUsed -   8 bytes
	ssz.DefineUint64(
		codec,
		&h.Timestamp,
	) // Field  ( 9) -        Timestamp -   8 bytes
	ssz.DefineDynamicBytesOffset(
		codec,
		(*[]byte)(&h.ExtraData),
		//nolint:mnd // todo fix.
		32,
	) // Offset (10) -        ExtraData -   4 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.BaseFeePerGas,
	) // Field  (11) -    BaseFeePerGas -  32 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.BlockHash,
	) // Field  (12) -        BlockHash -  32 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.TransactionsRoot,
	) // Field  (13) - TransactionsRoot -  32 bytes
	ssz.DefineStaticBytes(
		codec,
		&h.WithdrawalsRoot,
	) // Field  (14) -   WithdrawalRoot -  32 bytes
	ssz.DefineUint64(
		codec,
		&h.BlobGasUsed,
	) // Field  (15) -      BlobGasUsed -   8 bytes
	ssz.DefineUint64(
		codec,
		&h.ExcessBlobGas,
	) // Field  (16) -    ExcessBlobGas -   8 bytes

	// Define the dynamic data (fields)
	ssz.DefineDynamicBytesContent(
		codec,
		(*[]byte)(&h.ExtraData),
		//nolint:mnd // todo fix.
		32,
	) // Field  (10) -        ExtraData - ? bytes
}

// MarshalSSZ serializes the ExecutionPayloadHeader object into a slice of
// bytes.
func (h *ExecutionPayloadHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, h.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, h)
}

// UnmarshalSSZ unmarshals the ExecutionPayloadHeaderDeneb object from a source
// array.
func (h *ExecutionPayloadHeader) UnmarshalSSZ(bz []byte) error {
	// TODO: THIS IS A HACK TO GET BEACONSTATE TO WORK
	// h.version = version.Deneb
	return ssz.DecodeFromBytes(bz, h)
}

// HashTreeRootSSZ returns the hash tree root of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(h), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the ExecutionPayloadHeaderDeneb object to a target
// array.
func (h *ExecutionPayloadHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := h.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the ExecutionPayloadHeaderDeneb object with a
// hasher
//
//nolint:mnd // from fastssz.
func (h *ExecutionPayloadHeader) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'ParentHash'
	hh.PutBytes(h.ParentHash[:])

	// Field (1) 'FeeRecipient'
	hh.PutBytes(h.FeeRecipient[:])

	// Field (2) 'StateRoot'
	hh.PutBytes(h.StateRoot[:])

	// Field (3) 'ReceiptsRoot'
	hh.PutBytes(h.ReceiptsRoot[:])

	// Field (4) 'LogsBloom'
	if size := len(h.LogsBloom); size != 256 {
		return fastssz.ErrBytesLengthFn(
			"ExecutionPayloadHeaderDeneb.LogsBloom",
			size,
			256,
		)
	}
	hh.PutBytes(h.LogsBloom[:])

	// Field (5) 'Random'
	hh.PutBytes(h.Random[:])

	// Field (6) 'Number'
	hh.PutUint64(uint64(h.Number))

	// Field (7) 'GasLimit'
	hh.PutUint64(uint64(h.GasLimit))

	// Field (8) 'GasUsed'
	hh.PutUint64(uint64(h.GasUsed))

	// Field (9) 'Timestamp'
	hh.PutUint64(uint64(h.Timestamp))

	// Field (10) 'ExtraData'
	{
		elemIndx := hh.Index()
		byteLen := uint64(len(h.ExtraData))
		if byteLen > 32 {
			return fastssz.ErrIncorrectListSize
		}
		hh.Append(h.ExtraData)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (32+31)/32)
	}

	// Field (11) 'BaseFeePerGas'
	hh.PutBytes(h.BaseFeePerGas[:])

	// Field (12) 'BlockHash'
	hh.PutBytes(h.BlockHash[:])

	// Field (13) 'TransactionsRoot'
	hh.PutBytes(h.TransactionsRoot[:])

	// Field (14) 'WithdrawalsRoot'
	hh.PutBytes(h.WithdrawalsRoot[:])

	// Field (15) 'BlobGasUsed'
	hh.PutUint64(uint64(h.BlobGasUsed))

	// Field (16) 'ExcessBlobGas'
	hh.PutUint64(uint64(h.ExcessBlobGas))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the ExecutionPayloadHeaderDeneb object.
func (h *ExecutionPayloadHeader) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(h)
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// Version returns the version of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutionPayloadHeader is nil.
func (h *ExecutionPayloadHeader) IsNil() bool {
	return h == nil
}

// GetParentHash returns the parent hash of the ExecutionPayloadHeader.
func (
	h *ExecutionPayloadHeader,
) GetParentHash() gethprimitives.ExecutionHash {
	return h.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the
// ExecutionPayloadHeader.
//
//nolint:lll // long variable names.
func (h *ExecutionPayloadHeader) GetFeeRecipient() gethprimitives.ExecutionAddress {
	return h.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetStateRoot() common.Bytes32 {
	return h.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetReceiptsRoot() common.Bytes32 {
	return h.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetLogsBloom() []byte {
	return h.LogsBloom[:]
}

// GetPrevRandao returns the previous Randao value of the
// ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetPrevRandao() common.Bytes32 {
	return h.Random
}

// GetNumber returns the block number of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetNumber() math.U64 {
	return h.Number
}

// GetGasLimit returns the gas limit of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetGasLimit() math.U64 {
	return h.GasLimit
}

// GetGasUsed returns the gas used of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetGasUsed() math.U64 {
	return h.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetTimestamp() math.U64 {
	return h.Timestamp
}

// GetExtraData returns the extra data of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetExtraData() []byte {
	return h.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the
// ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetBaseFeePerGas() math.Wei {
	return h.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionPayloadHeader.
func (
	h *ExecutionPayloadHeader,
) GetBlockHash() gethprimitives.ExecutionHash {
	return h.BlockHash
}

// GetTransactionsRoot returns the transactions root of the
// ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetTransactionsRoot() common.Root {
	return h.TransactionsRoot
}

// GetWithdrawalsRoot returns the withdrawals root of the
// ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetWithdrawalsRoot() common.Root {
	return h.WithdrawalsRoot
}

// GetBlobGasUsed returns the blob gas used of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetBlobGasUsed() math.U64 {
	return h.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the
// ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetExcessBlobGas() math.U64 {
	return h.ExcessBlobGas
}

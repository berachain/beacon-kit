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
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/holiman/uint256"
	"github.com/karalabe/ssz"
)

// ExecutionPayloadHeader is the execution header payload of Deneb.
type ExecutionPayloadHeader struct {
	// TODO: Enable once
	// https://github.com/karalabe/ssz/pull/9/files# is merged.
	//
	// // Metadata
	// //
	// // version is the fork version of the execution payload header.
	// version uint32

	// Contents
	//
	// ParentHash is the hash of the parent block.
	ParentHash gethprimitives.ExecutionHash `json:"parentHash"`
	// FeeRecipient is the address of the fee recipient.
	FeeRecipient gethprimitives.ExecutionAddress `json:"feeRecipient"`
	// StateRoot is the root of the state trie.
	StateRoot common.Bytes32 `json:"stateRoot"`
	// ReceiptsRoot is the root of the receipts trie.
	ReceiptsRoot common.Bytes32 `json:"receiptsRoot"`
	// LogsBloom is the bloom filter for the logs.
	LogsBloom bytes.B256 `json:"logsBloom"`
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
	ExtraData bytes.Bytes `json:"extraData"`
	// BaseFeePerGas is the base fee per gas.
	BaseFeePerGas *uint256.Int `json:"baseFeePerGas"`
	// BlockHash is the hash of the block.
	BlockHash gethprimitives.ExecutionHash `json:"blockHash"`
	// TransactionsRoot is the root of the transaction trie.
	TransactionsRoot common.Root `json:"transactionsRoot"`
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
	ssz.DefineStaticBytes(codec, &h.ParentHash)
	ssz.DefineStaticBytes(codec, &h.FeeRecipient)
	ssz.DefineStaticBytes(codec, &h.StateRoot)
	ssz.DefineStaticBytes(codec, &h.ReceiptsRoot)
	ssz.DefineStaticBytes(codec, &h.LogsBloom)
	ssz.DefineStaticBytes(codec, &h.Random)
	ssz.DefineUint64(codec, &h.Number)
	ssz.DefineUint64(codec, &h.GasLimit)
	ssz.DefineUint64(codec, &h.GasUsed)
	ssz.DefineUint64(codec, &h.Timestamp)
	//nolint:mnd // todo fix.
	ssz.DefineDynamicBytesOffset(codec, (*[]byte)(&h.ExtraData), 32)
	ssz.DefineUint256(codec, &h.BaseFeePerGas)
	ssz.DefineStaticBytes(codec, &h.BlockHash)
	ssz.DefineStaticBytes(codec, &h.TransactionsRoot)
	ssz.DefineStaticBytes(codec, &h.WithdrawalsRoot)
	ssz.DefineUint64(codec, &h.BlobGasUsed)
	ssz.DefineUint64(codec, &h.ExcessBlobGas)

	// Define the dynamic data (fields)
	//nolint:mnd // todo fix.
	ssz.DefineDynamicBytesContent(codec, (*[]byte)(&h.ExtraData), 32)
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
	bz, err := h.BaseFeePerGas.MarshalSSZ()
	if err != nil {
		return err
	}
	hh.PutBytes(bz)

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
/*                                    JSON                                    */
/* -------------------------------------------------------------------------- */

// MarshalJSON marshals as JSON.
func (h ExecutionPayloadHeader) MarshalJSON() ([]byte, error) {
	type ExecutionPayloadHeader struct {
		ParentHash       gethprimitives.ExecutionHash    `json:"parentHash"`
		FeeRecipient     gethprimitives.ExecutionAddress `json:"feeRecipient"`
		StateRoot        bytes.B32                       `json:"stateRoot"`
		ReceiptsRoot     bytes.B32                       `json:"receiptsRoot"`
		LogsBloom        bytes.B256                      `json:"logsBloom"`
		Random           bytes.B32                       `json:"prevRandao"`
		Number           math.U64                        `json:"blockNumber"`
		GasLimit         math.U64                        `json:"gasLimit"`
		GasUsed          math.U64                        `json:"gasUsed"`
		Timestamp        math.U64                        `json:"timestamp"`
		ExtraData        bytes.Bytes                     `json:"extraData"`
		BaseFeePerGas    *math.U256                      `json:"baseFeePerGas"`
		BlockHash        common.ExecutionHash            `json:"blockHash"`
		TransactionsRoot bytes.B32                       `json:"transactionsRoot"`
		WithdrawalsRoot  bytes.B32                       `json:"withdrawalsRoot"`
		BlobGasUsed      math.U64                        `json:"blobGasUsed"`
		ExcessBlobGas    math.U64                        `json:"excessBlobGas"`
	}
	var enc ExecutionPayloadHeader
	enc.ParentHash = h.ParentHash
	enc.FeeRecipient = h.FeeRecipient
	enc.StateRoot = h.StateRoot
	enc.ReceiptsRoot = h.ReceiptsRoot
	enc.LogsBloom = h.LogsBloom
	enc.Random = h.Random
	enc.Number = h.Number
	enc.GasLimit = h.GasLimit
	enc.GasUsed = h.GasUsed
	enc.Timestamp = h.Timestamp
	enc.ExtraData = h.ExtraData
	enc.BaseFeePerGas = h.BaseFeePerGas
	enc.BlockHash = h.BlockHash
	enc.TransactionsRoot = h.TransactionsRoot
	enc.WithdrawalsRoot = h.WithdrawalsRoot
	enc.BlobGasUsed = h.BlobGasUsed
	enc.ExcessBlobGas = h.ExcessBlobGas
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
//
//nolint:funlen // codegen.
func (h *ExecutionPayloadHeader) UnmarshalJSON(input []byte) error {
	type ExecutionPayloadHeader struct {
		ParentHash       *gethprimitives.ExecutionHash    `json:"parentHash"`
		FeeRecipient     *gethprimitives.ExecutionAddress `json:"feeRecipient"`
		StateRoot        *bytes.B32                       `json:"stateRoot"`
		ReceiptsRoot     *bytes.B32                       `json:"receiptsRoot"`
		LogsBloom        *bytes.B256                      `json:"logsBloom"`
		Random           *bytes.B32                       `json:"prevRandao"`
		Number           *math.U64                        `json:"blockNumber"`
		GasLimit         *math.U64                        `json:"gasLimit"`
		GasUsed          *math.U64                        `json:"gasUsed"`
		Timestamp        *math.U64                        `json:"timestamp"`
		ExtraData        *bytes.Bytes                     `json:"extraData"`
		BaseFeePerGas    *math.U256                       `json:"baseFeePerGas"`
		BlockHash        *gethprimitives.ExecutionHash    `json:"blockHash"`
		TransactionsRoot *bytes.B32                       `json:"transactionsRoot"`
		WithdrawalsRoot  *bytes.B32                       `json:"withdrawalsRoot"`
		BlobGasUsed      *math.U64                        `json:"blobGasUsed"`
		ExcessBlobGas    *math.U64                        `json:"excessBlobGas"`
	}
	var dec ExecutionPayloadHeader
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentHash == nil {
		return errors.New(
			"missing required field 'parentHash' for ExecutionPayloadHeader",
		)
	}
	h.ParentHash = *dec.ParentHash
	if dec.FeeRecipient == nil {
		return errors.New(
			"missing required field 'feeRecipient' for ExecutionPayloadHeader",
		)
	}
	h.FeeRecipient = *dec.FeeRecipient
	if dec.StateRoot == nil {
		return errors.New(
			"missing required field 'stateRoot' for ExecutionPayloadHeader",
		)
	}
	h.StateRoot = *dec.StateRoot
	if dec.ReceiptsRoot == nil {
		return errors.New(
			"missing required field 'receiptsRoot' for ExecutionPayloadHeader",
		)
	}
	h.ReceiptsRoot = *dec.ReceiptsRoot
	if dec.LogsBloom == nil {
		return errors.New(
			"missing required field 'logsBloom' for ExecutionPayloadHeader",
		)
	}
	h.LogsBloom = *dec.LogsBloom
	if dec.Random == nil {
		return errors.New(
			"missing required field 'prevRandao' for ExecutionPayloadHeader",
		)
	}
	h.Random = *dec.Random
	if dec.Number == nil {
		return errors.New(
			"missing required field 'blockNumber' for ExecutionPayloadHeader",
		)
	}
	h.Number = *dec.Number
	if dec.GasLimit == nil {
		return errors.New(
			"missing required field 'gasLimit' for ExecutionPayloadHeader",
		)
	}
	h.GasLimit = *dec.GasLimit
	if dec.GasUsed == nil {
		return errors.New(
			"missing required field 'gasUsed' for ExecutionPayloadHeader",
		)
	}
	h.GasUsed = *dec.GasUsed
	if dec.Timestamp == nil {
		return errors.New(
			"missing required field 'timestamp' for ExecutionPayloadHeader",
		)
	}
	h.Timestamp = *dec.Timestamp
	if dec.ExtraData == nil {
		return errors.New(
			"missing required field 'extraData' for ExecutionPayloadHeader",
		)
	}

	// TODO: This is required for the API to be symmetric? But it's not really
	// clear if
	// this matters.
	if len(*dec.ExtraData) != 0 {
		h.ExtraData = *dec.ExtraData
	}

	if dec.BaseFeePerGas == nil {
		return errors.New(
			"missing required field 'baseFeePerGas' for ExecutionPayloadHeader",
		)
	}
	h.BaseFeePerGas = dec.BaseFeePerGas
	if dec.BlockHash == nil {
		return errors.New(
			"missing required field 'blockHash' for ExecutionPayloadHeader",
		)
	}
	h.BlockHash = *dec.BlockHash
	if dec.TransactionsRoot == nil {
		return errors.New(
			"missing required field 'transactionsRoot' for ExecutionPayloadHeader",
		)
	}
	h.TransactionsRoot = *dec.TransactionsRoot
	if dec.WithdrawalsRoot != nil {
		h.WithdrawalsRoot = *dec.WithdrawalsRoot
	}
	if dec.BlobGasUsed != nil {
		h.BlobGasUsed = *dec.BlobGasUsed
	}
	if dec.ExcessBlobGas != nil {
		h.ExcessBlobGas = *dec.ExcessBlobGas
	}
	return nil
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
func (h *ExecutionPayloadHeader) GetLogsBloom() bytes.B256 {
	return h.LogsBloom
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
func (h *ExecutionPayloadHeader) GetBaseFeePerGas() *math.U256 {
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

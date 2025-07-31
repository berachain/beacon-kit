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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
	ssz "github.com/ferranbt/fastssz"
)

// ExecutionPayloadHeaderStaticSize is the static size of the ExecutionPayloadHeader.
const ExecutionPayloadHeaderStaticSize uint32 = 584

// TODO: Re-enable interface assertion once constraints are updated
// var (
// 	_ constraints.SSZVersionedMarshallableRootable = (*ExecutionPayloadHeader)(nil)
// )

// ExecutionPayloadHeader represents the payload header of an execution block.
type ExecutionPayloadHeader struct {
	// NOTE: This version is not required but left in for backwards compatibility.
	//
	// A recommended alternative to `GetForkVersion()` on this struct would be to use the chain
	// spec's `ActiveForkVersionForTimestamp()` on the value of `GetTimestamp()`.
	//
	// This version should still be set to the correct value to avoid potential inconsistencies.
	Versionable

	// Contents
	//
	// ParentHash is the hash of the parent block.
	ParentHash common.ExecutionHash `json:"parentHash"`
	// FeeRecipient is the address of the fee recipient.
	FeeRecipient common.ExecutionAddress `json:"feeRecipient"`
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
	BaseFeePerGas *math.U256 `json:"baseFeePerGas"`
	// BlockHash is the hash of the block.
	BlockHash common.ExecutionHash `json:"blockHash"`
	// TransactionsRoot is the root of the transaction trie.
	TransactionsRoot common.Root `json:"transactionsRoot"`
	// WithdrawalsRoot is the root of the withdrawals trie.
	WithdrawalsRoot common.Root `json:"withdrawalsRoot"`
	// BlobGasUsed is the amount of blob gas used in the block.
	BlobGasUsed math.U64 `json:"blobGasUsed"`
	// ExcessBlobGas is the amount of excess blob gas in the block.
	ExcessBlobGas math.U64 `json:"excessBlobGas"`
}

func NewEmptyExecutionPayloadHeaderWithVersion(version common.Version) *ExecutionPayloadHeader {
	return &ExecutionPayloadHeader{
		Versionable:   NewVersionable(version),
		BaseFeePerGas: &math.U256{},
	}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the total size of the object in SSZ encoding.
func (h *ExecutionPayloadHeader) SizeSSZ() int {
	return int(ExecutionPayloadHeaderStaticSize) + len(h.ExtraData)
}

// MarshalSSZ serializes the ExecutionPayloadHeader object into a slice of
// bytes.
func (h *ExecutionPayloadHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, h.SizeSSZ())
	return h.MarshalSSZTo(buf)
}

func (*ExecutionPayloadHeader) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot returns the hash tree root of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)
	if err := h.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()

}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the ExecutionPayloadHeader object to a target array.
func (h *ExecutionPayloadHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Static fields
	dst = append(dst, h.ParentHash[:]...)
	dst = append(dst, h.FeeRecipient[:]...)
	dst = append(dst, h.StateRoot[:]...)
	dst = append(dst, h.ReceiptsRoot[:]...)
	dst = append(dst, h.LogsBloom[:]...)
	dst = append(dst, h.Random[:]...)
	dst = ssz.MarshalUint64(dst, uint64(h.Number))
	dst = ssz.MarshalUint64(dst, uint64(h.GasLimit))
	dst = ssz.MarshalUint64(dst, uint64(h.GasUsed))
	dst = ssz.MarshalUint64(dst, uint64(h.Timestamp))

	// Offset for ExtraData
	offset := uint32(ExecutionPayloadHeaderStaticSize)
	dst = ssz.MarshalUint32(dst, offset)

	// BaseFeePerGas
	bz, err := h.BaseFeePerGas.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)

	// More static fields
	dst = append(dst, h.BlockHash[:]...)
	dst = append(dst, h.TransactionsRoot[:]...)
	dst = append(dst, h.WithdrawalsRoot[:]...)
	dst = ssz.MarshalUint64(dst, uint64(h.BlobGasUsed))
	dst = ssz.MarshalUint64(dst, uint64(h.ExcessBlobGas))

	// Dynamic field: ExtraData
	dst = append(dst, h.ExtraData...)

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the ExecutionPayloadHeader object.
func (h *ExecutionPayloadHeader) UnmarshalSSZ(buf []byte) error {
	if len(buf) < int(ExecutionPayloadHeaderStaticSize) {
		return ssz.ErrSize
	}

	// Static fields
	copy(h.ParentHash[:], buf[0:32])
	copy(h.FeeRecipient[:], buf[32:52])
	copy(h.StateRoot[:], buf[52:84])
	copy(h.ReceiptsRoot[:], buf[84:116])
	copy(h.LogsBloom[:], buf[116:372])
	copy(h.Random[:], buf[372:404])
	h.Number = math.U64(ssz.UnmarshallUint64(buf[404:412]))
	h.GasLimit = math.U64(ssz.UnmarshallUint64(buf[412:420]))
	h.GasUsed = math.U64(ssz.UnmarshallUint64(buf[420:428]))
	h.Timestamp = math.U64(ssz.UnmarshallUint64(buf[428:436]))

	// Read offset for ExtraData
	extraDataOffset := ssz.UnmarshallUint32(buf[436:440])

	// BaseFeePerGas
	if h.BaseFeePerGas == nil {
		h.BaseFeePerGas = &math.U256{}
	}
	if err := h.BaseFeePerGas.UnmarshalSSZ(buf[440:472]); err != nil {
		return err
	}

	// More static fields
	copy(h.BlockHash[:], buf[472:504])
	copy(h.TransactionsRoot[:], buf[504:536])
	copy(h.WithdrawalsRoot[:], buf[536:568])
	h.BlobGasUsed = math.U64(ssz.UnmarshallUint64(buf[568:576]))
	h.ExcessBlobGas = math.U64(ssz.UnmarshallUint64(buf[576:584]))

	// Dynamic field: ExtraData
	if extraDataOffset > uint32(len(buf)) {
		return ssz.ErrInvalidVariableOffset
	}
	h.ExtraData = append([]byte(nil), buf[extraDataOffset:]...)
	if len(h.ExtraData) > 32 {
		return errors.New("extra data too large")
	}

	return nil
}

// HashTreeRootWith ssz hashes the ExecutionPayloadHeaderDeneb object with a
// hasher
//
//nolint:mnd // from ssz.
func (h *ExecutionPayloadHeader) HashTreeRootWith(hh ssz.HashWalker) error {
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
		return ssz.ErrBytesLengthFn(
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
			return ssz.ErrIncorrectListSize
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
func (h *ExecutionPayloadHeader) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(h)
}

/* -------------------------------------------------------------------------- */
/*                                    JSON                                    */
/* -------------------------------------------------------------------------- */

// MarshalJSON marshals as JSON.
func (h ExecutionPayloadHeader) MarshalJSON() ([]byte, error) {
	type ExecutionPayloadHeader struct {
		ParentHash       common.ExecutionHash    `json:"parentHash"`
		FeeRecipient     common.ExecutionAddress `json:"feeRecipient"`
		StateRoot        bytes.B32               `json:"stateRoot"`
		ReceiptsRoot     bytes.B32               `json:"receiptsRoot"`
		LogsBloom        bytes.B256              `json:"logsBloom"`
		Random           bytes.B32               `json:"prevRandao"`
		Number           math.U64                `json:"blockNumber"`
		GasLimit         math.U64                `json:"gasLimit"`
		GasUsed          math.U64                `json:"gasUsed"`
		Timestamp        math.U64                `json:"timestamp"`
		ExtraData        bytes.Bytes             `json:"extraData"`
		BaseFeePerGas    *math.U256              `json:"baseFeePerGas"`
		BlockHash        common.ExecutionHash    `json:"blockHash"`
		TransactionsRoot common.Root             `json:"transactionsRoot"`
		WithdrawalsRoot  common.Root             `json:"withdrawalsRoot"`
		BlobGasUsed      math.U64                `json:"blobGasUsed"`
		ExcessBlobGas    math.U64                `json:"excessBlobGas"`
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
		ParentHash       *common.ExecutionHash    `json:"parentHash"`
		FeeRecipient     *common.ExecutionAddress `json:"feeRecipient"`
		StateRoot        *bytes.B32               `json:"stateRoot"`
		ReceiptsRoot     *bytes.B32               `json:"receiptsRoot"`
		LogsBloom        *bytes.B256              `json:"logsBloom"`
		Random           *bytes.B32               `json:"prevRandao"`
		Number           *math.U64                `json:"blockNumber"`
		GasLimit         *math.U64                `json:"gasLimit"`
		GasUsed          *math.U64                `json:"gasUsed"`
		Timestamp        *math.U64                `json:"timestamp"`
		ExtraData        *bytes.Bytes             `json:"extraData"`
		BaseFeePerGas    *math.U256               `json:"baseFeePerGas"`
		BlockHash        *common.ExecutionHash    `json:"blockHash"`
		TransactionsRoot *common.Root             `json:"transactionsRoot"`
		WithdrawalsRoot  *common.Root             `json:"withdrawalsRoot"`
		BlobGasUsed      *math.U64                `json:"blobGasUsed"`
		ExcessBlobGas    *math.U64                `json:"excessBlobGas"`
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
/*                                   Getters                                  */
/* -------------------------------------------------------------------------- */

// GetParentHash returns the parent hash of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetParentHash() common.ExecutionHash {
	return h.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutionPayloadHeader.
func (h *ExecutionPayloadHeader) GetFeeRecipient() common.ExecutionAddress {
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
func (h *ExecutionPayloadHeader) GetBlockHash() common.ExecutionHash {
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

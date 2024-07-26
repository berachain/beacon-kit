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
//go:generate go run github.com/fjl/gencodec -type ExecutionPayloadHeader -out payload_header.json.go -field-override executionPayloadHeaderMarshaling
//nolint:lll
type ExecutionPayloadHeader struct {
	version          uint32
	ParentHash       gethprimitives.ExecutionHash    `json:"parentHash"       ssz-size:"32"  gencodec:"required"`
	FeeRecipient     gethprimitives.ExecutionAddress `json:"feeRecipient"     ssz-size:"20"  gencodec:"required"`
	StateRoot        common.Bytes32                  `json:"stateRoot"        ssz-size:"32"  gencodec:"required"`
	ReceiptsRoot     common.Bytes32                  `json:"receiptsRoot"     ssz-size:"32"  gencodec:"required"`
	LogsBloom        bytes.B256                      `json:"logsBloom"        ssz-size:"256" gencodec:"required"`
	Random           common.Bytes32                  `json:"prevRandao"       ssz-size:"32"  gencodec:"required"`
	Number           math.U64                        `json:"blockNumber"                     gencodec:"required"`
	GasLimit         math.U64                        `json:"gasLimit"                        gencodec:"required"`
	GasUsed          math.U64                        `json:"gasUsed"                         gencodec:"required"`
	Timestamp        math.U64                        `json:"timestamp"                       gencodec:"required"`
	ExtraData        []byte                          `json:"extraData"                       gencodec:"required" ssz-max:"32"`
	BaseFeePerGas    math.Wei                        `json:"baseFeePerGas"    ssz-size:"32"  gencodec:"required"`
	BlockHash        gethprimitives.ExecutionHash    `json:"blockHash"        ssz-size:"32"  gencodec:"required"`
	TransactionsRoot common.Root                     `json:"transactionsRoot" ssz-size:"32"  gencodec:"required"`
	WithdrawalsRoot  common.Root                     `json:"withdrawalsRoot"  ssz-size:"32"`
	BlobGasUsed      math.U64                        `json:"blobGasUsed"`
	ExcessBlobGas    math.U64                        `json:"excessBlobGas"`
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (e *ExecutionPayloadHeader) Empty(
	forkVersion uint32,
) *ExecutionPayloadHeader {
	return &ExecutionPayloadHeader{
		version: forkVersion,
	}
}

// NewFromSSZ returns a new ExecutionPayloadHeader from the given SSZ bytes.
func (e *ExecutionPayloadHeader) NewFromSSZ(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	e = e.Empty(forkVersion)
	if err := e.UnmarshalSSZ(bz); err != nil {
		return nil, err
	}
	return e, nil
}

// NewFromJSON returns a new ExecutionPayloadHeader from the given JSON bytes.
func (e *ExecutionPayloadHeader) NewFromJSON(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	e = e.Empty(forkVersion)
	if err := e.UnmarshalJSON(bz); err != nil {
		return nil, err
	}
	return e, nil
}

type executionPayloadHeaderMarshaling struct {
	ExtraData bytes.Bytes
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns either the static size of the object if fixed == true, or
// the total size otherwise.
func (e *ExecutionPayloadHeader) SizeSSZ(fixed bool) uint32 {
	var size = uint32(32 + 20 + 32 + 32 + 256 + 32 + 8 + 8 + 8 + 8 + 4 + 32 + 32 + 32 + 32 + 8 + 8)
	if fixed {
		return size
	}
	size += ssz.SizeDynamicBytes(e.ExtraData)

	return size
}

// DefineSSZ defines how an object is encoded/decoded.
func (e *ExecutionPayloadHeader) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineStaticBytes(codec, &e.ParentHash)           // Field  ( 0) -       ParentHash -  32 bytes
	ssz.DefineStaticBytes(codec, &e.FeeRecipient)         // Field  ( 1) -     FeeRecipient -  20 bytes
	ssz.DefineStaticBytes(codec, &e.StateRoot)            // Field  ( 2) -        StateRoot -  32 bytes
	ssz.DefineStaticBytes(codec, &e.ReceiptsRoot)         // Field  ( 3) -     ReceiptsRoot -  32 bytes
	ssz.DefineStaticBytes(codec, &e.LogsBloom)            // Field  ( 4) -        LogsBloom - 256 bytes
	ssz.DefineStaticBytes(codec, &e.Random)               // Field  ( 5) -       PrevRandao -  32 bytes
	ssz.DefineUint64(codec, &e.Number)                    // Field  ( 6) -      BlockNumber -   8 bytes
	ssz.DefineUint64(codec, &e.GasLimit)                  // Field  ( 7) -         GasLimit -   8 bytes
	ssz.DefineUint64(codec, &e.GasUsed)                   // Field  ( 8) -          GasUsed -   8 bytes
	ssz.DefineUint64(codec, &e.Timestamp)                 // Field  ( 9) -        Timestamp -   8 bytes
	ssz.DefineDynamicBytesOffset(codec, &e.ExtraData, 32) // Offset (10) -        ExtraData -   4 bytes
	ssz.DefineStaticBytes(codec, &e.BaseFeePerGas)        // Field  (11) -    BaseFeePerGas -  32 bytes
	ssz.DefineStaticBytes(codec, &e.BlockHash)            // Field  (12) -        BlockHash -  32 bytes
	ssz.DefineStaticBytes(codec, &e.TransactionsRoot)     // Field  (13) - TransactionsRoot -  32 bytes
	ssz.DefineStaticBytes(codec, &e.WithdrawalsRoot)      // Field  (14) -   WithdrawalRoot -  32 bytes
	ssz.DefineUint64(codec, &e.BlobGasUsed)               // Field  (15) -      BlobGasUsed -   8 bytes
	ssz.DefineUint64(codec, &e.ExcessBlobGas)             // Field  (16) -    ExcessBlobGas -   8 bytes

	// Define the dynamic data (fields)
	ssz.DefineDynamicBytesContent(codec, &e.ExtraData, 32) // Field  (10) -        ExtraData - ? bytes
}

// MarshalSSZ serializes the ExecutionPayloadHeader object into a slice of bytes.
func (e *ExecutionPayloadHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, e.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, e)
}

// UnmarshalSSZ unmarshals the ExecutionPayloadHeaderDeneb object from a source array
func (e *ExecutionPayloadHeader) UnmarshalSSZ(bz []byte) error {
	return ssz.DecodeFromBytes(bz, e)
}

// HashTreeRootSSZ returns the hash tree root of the ExecutionPayloadHeader.
func (e *ExecutionPayloadHeader) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(e), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the ExecutionPayloadHeaderDeneb object to a target array
func (e *ExecutionPayloadHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := e.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the ExecutionPayloadHeaderDeneb object with a hasher
func (e *ExecutionPayloadHeader) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'ParentHash'
	hh.PutBytes(e.ParentHash[:])

	// Field (1) 'FeeRecipient'
	hh.PutBytes(e.FeeRecipient[:])

	// Field (2) 'StateRoot'
	hh.PutBytes(e.StateRoot[:])

	// Field (3) 'ReceiptsRoot'
	hh.PutBytes(e.ReceiptsRoot[:])

	// Field (4) 'LogsBloom'
	if size := len(e.LogsBloom); size != 256 {
		return fastssz.ErrBytesLengthFn("ExecutionPayloadHeaderDeneb.LogsBloom", size, 256)
	}
	hh.PutBytes(e.LogsBloom[:])

	// Field (5) 'Random'
	hh.PutBytes(e.Random[:])

	// Field (6) 'Number'
	hh.PutUint64(uint64(e.Number))

	// Field (7) 'GasLimit'
	hh.PutUint64(uint64(e.GasLimit))

	// Field (8) 'GasUsed'
	hh.PutUint64(uint64(e.GasUsed))

	// Field (9) 'Timestamp'
	hh.PutUint64(uint64(e.Timestamp))

	// Field (10) 'ExtraData'
	{
		elemIndx := hh.Index()
		byteLen := uint64(len(e.ExtraData))
		if byteLen > 32 {
			return fastssz.ErrIncorrectListSize
		}
		hh.Append(e.ExtraData)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (32+31)/32)
	}

	// Field (11) 'BaseFeePerGas'
	hh.PutBytes(e.BaseFeePerGas[:])

	// Field (12) 'BlockHash'
	hh.PutBytes(e.BlockHash[:])

	// Field (13) 'TransactionsRoot'
	hh.PutBytes(e.TransactionsRoot[:])

	// Field (14) 'WithdrawalsRoot'
	hh.PutBytes(e.WithdrawalsRoot[:])

	// Field (15) 'BlobGasUsed'
	hh.PutUint64(uint64(e.BlobGasUsed))

	// Field (16) 'ExcessBlobGas'
	hh.PutUint64(uint64(e.ExcessBlobGas))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the ExecutionPayloadHeaderDeneb object
func (e *ExecutionPayloadHeader) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(e)
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// Version returns the version of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutionPayloadHeader is nil.
func (d *ExecutionPayloadHeader) IsNil() bool {
	return d == nil
}

// IsBlinded checks if the ExecutionPayloadHeader is blinded.
func (d *ExecutionPayloadHeader) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutionPayloadHeader.
func (
	d *ExecutionPayloadHeader,
) GetParentHash() gethprimitives.ExecutionHash {
	return d.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the
// ExecutionPayloadHeader.
//
//nolint:lll // long variable names.
func (d *ExecutionPayloadHeader) GetFeeRecipient() gethprimitives.ExecutionAddress {
	return d.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetStateRoot() common.Bytes32 {
	return d.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetReceiptsRoot() common.Bytes32 {
	return d.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetLogsBloom() []byte {
	return d.LogsBloom[:]
}

// GetPrevRandao returns the previous Randao value of the
// ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetPrevRandao() common.Bytes32 {
	return d.Random
}

// GetNumber returns the block number of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetNumber() math.U64 {
	return d.Number
}

// GetGasLimit returns the gas limit of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetGasLimit() math.U64 {
	return d.GasLimit
}

// GetGasUsed returns the gas used of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetGasUsed() math.U64 {
	return d.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetTimestamp() math.U64 {
	return d.Timestamp
}

// GetExtraData returns the extra data of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetExtraData() []byte {
	return d.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the
// ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetBaseFeePerGas() math.Wei {
	return d.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionPayloadHeader.
func (
	d *ExecutionPayloadHeader,
) GetBlockHash() gethprimitives.ExecutionHash {
	return d.BlockHash
}

// GetTransactionsRoot returns the transactions root of the
// ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetTransactionsRoot() common.Root {
	return d.TransactionsRoot
}

// GetWithdrawalsRoot returns the withdrawals root of the
// ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetWithdrawalsRoot() common.Root {
	return d.WithdrawalsRoot
}

// GetBlobGasUsed returns the blob gas used of the ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetBlobGasUsed() math.U64 {
	return d.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the
// ExecutionPayloadHeader.
func (d *ExecutionPayloadHeader) GetExcessBlobGas() math.U64 {
	return d.ExcessBlobGas
}

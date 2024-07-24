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
	"github.com/karalabe/ssz"
)

// ExecutionPayloadHeader represents an execution payload across
// all fork versions.
//
//go:generate go run github.com/fjl/gencodec -type ExecutionPayloadHeader -out payload_header.json.go -field-override executionPayloadHeaderMarshaling
type ExecutionPayloadHeader struct {
	ParentHash       gethprimitives.ExecutionHash    `json:"parentHash"       ssz-size:"32"  genc:"required"`
	FeeRecipient     gethprimitives.ExecutionAddress `json:"feeRecipient"     ssz-size:"20"  genc:"required"`
	StateRoot        common.Bytes32                  `json:"stateRoot"        ssz-size:"32"  genc:"required"`
	ReceiptsRoot     common.Bytes32                  `json:"receiptsRoot"     ssz-size:"32"  genc:"required"`
	LogsBloom        [256]byte                       `json:"logsBloom"        ssz-size:"256" genc:"required"`
	Random           common.Bytes32                  `json:"prevRandao"       ssz-size:"32"  genc:"required"`
	Number           math.U64                        `json:"blockNumber"                     genc:"required"`
	GasLimit         math.U64                        `json:"gasLimit"                        genc:"required"`
	GasUsed          math.U64                        `json:"gasUsed"                         genc:"required"`
	Timestamp        math.U64                        `json:"timestamp"                       genc:"required"`
	ExtraData        []byte                          `json:"extraData"                       genc:"required" ssz-max:"32"`
	BaseFeePerGas    math.Wei                        `json:"baseFeePerGas"    ssz-size:"32"  genc:"required"`
	BlockHash        gethprimitives.ExecutionHash    `json:"blockHash"        ssz-size:"32"  genc:"required"`
	TransactionsRoot common.Root                     `json:"transactionsRoot" ssz-size:"32"  genc:"required"`
	WithdrawalsRoot  common.Root                     `json:"withdrawalsRoot"  ssz-size:"32"`
	BlobGasUsed      math.U64                        `json:"blobGasUsed"`
	ExcessBlobGas    math.U64                        `json:"excessBlobGas"`
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (e *ExecutionPayloadHeader) Empty(
	forkVersion uint32,
) *ExecutionPayloadHeader {
	return &ExecutionPayloadHeader{}
}

// NewFromSSZ returns a new ExecutionPayloadHeader from the given SSZ bytes.
func (e *ExecutionPayloadHeader) NewFromSSZ(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	e = e.Empty(forkVersion)
	return e, ssz.DecodeFromBytes(bz, e)
}

// NewFromJSON returns a new ExecutionPayloadHeader from the given JSON bytes.
func (e *ExecutionPayloadHeader) NewFromJSON(
	bz []byte, forkVersion uint32,
) (*ExecutionPayloadHeader, error) {
	e = e.Empty(forkVersion)
	return e, e.UnmarshalJSON(bz)
}

type executionPayloadHeaderMarshaling struct {
	ExtraData bytes.Bytes
	LogsBloom bytes.Bytes
}

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

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns either the static size of the object if fixed == true, or
// the total size otherwise.
func (obj *ExecutionPayloadHeader) SizeSSZ(fixed bool) uint32 {
	var size = uint32(32 + 20 + 32 + 32 + 256 + 32 + 8 + 8 + 8 + 8 + 4 + 32 + 32 + 32)
	if fixed {
		return size
	}
	size += ssz.SizeDynamicBytes(obj.ExtraData)

	return size
}

// DefineSSZ defines how an object is encoded/decoded.
func (obj *ExecutionPayloadHeader) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineStaticBytes(codec, &obj.ParentHash)           // Field  ( 0) -       ParentHash -  32 bytes
	ssz.DefineStaticBytes(codec, &obj.FeeRecipient)         // Field  ( 1) -     FeeRecipient -  20 bytes
	ssz.DefineStaticBytes(codec, &obj.StateRoot)            // Field  ( 2) -        StateRoot -  32 bytes
	ssz.DefineStaticBytes(codec, &obj.ReceiptsRoot)         // Field  ( 3) -     ReceiptsRoot -  32 bytes
	ssz.DefineStaticBytes(codec, &obj.LogsBloom)            // Field  ( 4) -        LogsBloom - 256 bytes
	ssz.DefineStaticBytes(codec, &obj.Random)               // Field  ( 5) -       PrevRandao -  32 bytes
	ssz.DefineUint64(codec, &obj.Number)                    // Field  ( 6) -      BlockNumber -   8 bytes
	ssz.DefineUint64(codec, &obj.GasLimit)                  // Field  ( 7) -         GasLimit -   8 bytes
	ssz.DefineUint64(codec, &obj.GasUsed)                   // Field  ( 8) -          GasUsed -   8 bytes
	ssz.DefineUint64(codec, &obj.Timestamp)                 // Field  ( 9) -        Timestamp -   8 bytes
	ssz.DefineDynamicBytesOffset(codec, &obj.ExtraData, 32) // Offset (10) -        ExtraData -   4 bytes
	ssz.DefineStaticBytes(codec, &obj.BaseFeePerGas)        // Field  (11) -    BaseFeePerGas -  32 bytes
	ssz.DefineStaticBytes(codec, &obj.BlockHash)            // Field  (12) -        BlockHash -  32 bytes
	ssz.DefineStaticBytes(codec, &obj.TransactionsRoot)     // Field  (13) - TransactionsRoot -  32 bytes

	// Define the dynamic data (fields)
	ssz.DefineDynamicBytesContent(codec, &obj.ExtraData, 32) // Field  (10) -        ExtraData - ? bytes
}

// MarshalSSZTo marshals the object into a preallocated buffer.
func (obj *ExecutionPayloadHeader) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, obj)
}

// MarshalSSZ marshals the object into a new buffer.
func (obj *ExecutionPayloadHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, obj.SizeSSZ(false))
	return obj.MarshalSSZTo(buf)
}

// UnmarshalSSZ unmarshals the object from the provided buffer.
func (obj *ExecutionPayloadHeader) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, obj)
}

// HashTreeRoot returns the hash tree root of the object.
func (obj *ExecutionPayloadHeader) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(obj), nil
}

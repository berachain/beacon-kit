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
	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// ExecutionPayloadStaticSize is the static size of the ExecutionPayload.
const ExecutionPayloadStaticSize uint32 = 528

// ExecutionPayload represents the payload of an execution block.
type ExecutionPayload struct {
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
	BaseFeePerGas *math.U256 `json:"baseFeePerGas"`
	// BlockHash is the hash of the block.
	BlockHash gethprimitives.ExecutionHash `json:"blockHash"`
	// Transactions is the list of transactions in the block.
	Transactions engineprimitives.Transactions `json:"transactions"`
	// Withdrawals is the list of withdrawals in the block.
	Withdrawals []*engineprimitives.Withdrawal `json:"withdrawals"`
	// BlobGasUsed is the amount of blob gas used in the block.
	BlobGasUsed math.U64 `json:"blobGasUsed"`
	// ExcessBlobGas is the amount of excess blob gas in the block.
	ExcessBlobGas math.U64 `json:"excessBlobGas"`
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns either the static size of the object if fixed == true, or
// the total size otherwise.
func (p *ExecutionPayload) SizeSSZ(fixed bool) uint32 {
	var size = ExecutionPayloadStaticSize
	if fixed {
		return size
	}
	size += ssz.SizeDynamicBytes(p.ExtraData)
	size += ssz.SizeSliceOfDynamicBytes(p.Transactions)
	size += ssz.SizeSliceOfStaticObjects(p.Withdrawals)
	return size
}

// DefineSSZ defines how an object is encoded/decoded.
//
//nolint:mnd // todo fix.
func (p *ExecutionPayload) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineStaticBytes(codec, &p.ParentHash)
	ssz.DefineStaticBytes(codec, &p.FeeRecipient)
	ssz.DefineStaticBytes(codec, &p.StateRoot)
	ssz.DefineStaticBytes(codec, &p.ReceiptsRoot)
	ssz.DefineStaticBytes(codec, &p.LogsBloom)
	ssz.DefineStaticBytes(codec, &p.Random)
	ssz.DefineUint64(codec, &p.Number)
	ssz.DefineUint64(codec, &p.GasLimit)
	ssz.DefineUint64(codec, &p.GasUsed)
	ssz.DefineUint64(codec, &p.Timestamp)
	ssz.DefineDynamicBytesOffset(codec, (*[]byte)(&p.ExtraData), 32)
	ssz.DefineUint256(codec, &p.BaseFeePerGas)
	ssz.DefineStaticBytes(codec, &p.BlockHash)
	ssz.DefineSliceOfDynamicBytesOffset(
		codec,
		(*[][]byte)(&p.Transactions),
		constants.MaxTxsPerPayload,
		constants.MaxBytesPerTx,
	)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &p.Withdrawals, 16)
	ssz.DefineUint64(codec, &p.BlobGasUsed)
	ssz.DefineUint64(codec, &p.ExcessBlobGas)

	// Define the dynamic data (fields)
	ssz.DefineDynamicBytesContent(codec, (*[]byte)(&p.ExtraData), 32)
	ssz.DefineSliceOfDynamicBytesContent(
		codec,
		(*[][]byte)(&p.Transactions),
		constants.MaxTxsPerPayload,
		constants.MaxBytesPerTx,
	)
	ssz.DefineSliceOfStaticObjectsContent(codec, &p.Withdrawals, 16)
}

// MarshalSSZ serializes the ExecutionPayload object into a slice of bytes.
func (p *ExecutionPayload) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, p.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, p)
}

// UnmarshalSSZ unmarshals the ExecutionPayload object from a source array.
func (p *ExecutionPayload) UnmarshalSSZ(bz []byte) error {
	return ssz.DecodeFromBytes(bz, p)
}

// HashTreeRoot returns the hash tree root of the ExecutionPayload.
func (p *ExecutionPayload) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(p), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo serializes the ExecutionPayload object into a writer.
func (p *ExecutionPayload) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := p.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the ExecutionPayload object with a hasher.
//
//nolint:mnd // will be deprecated eventually.
func (p *ExecutionPayload) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'ParentHash'
	hh.PutBytes(p.ParentHash[:])

	// Field (1) 'FeeRecipient'
	hh.PutBytes(p.FeeRecipient[:])

	// Field (2) 'StateRoot'
	hh.PutBytes(p.StateRoot[:])

	// Field (3) 'ReceiptsRoot'
	hh.PutBytes(p.ReceiptsRoot[:])

	// Field (4) 'LogsBloom'
	hh.PutBytes(p.LogsBloom[:])

	// Field (5) 'Random'
	hh.PutBytes(p.Random[:])

	// Field (6) 'Number'
	hh.PutUint64(uint64(p.Number))

	// Field (7) 'GasLimit'
	hh.PutUint64(uint64(p.GasLimit))

	// Field (8) 'GasUsed'
	hh.PutUint64(uint64(p.GasUsed))

	// Field (9) 'Timestamp'
	hh.PutUint64(uint64(p.Timestamp))

	// Field (10) 'ExtraData'
	{
		elemIndx := hh.Index()
		byteLen := uint64(len(p.ExtraData))
		if byteLen > 32 {
			return fastssz.ErrIncorrectListSize
		}
		hh.Append(p.ExtraData)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (32+31)/32)
	}

	// Field (11) 'BaseFeePerGas'
	bz, err := p.BaseFeePerGas.MarshalSSZ()
	if err != nil {
		return err
	}
	hh.PutBytes(bz)

	// Field (12) 'BlockHash'
	hh.PutBytes(p.BlockHash[:])

	// Field (13) 'Transactions'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Transactions))
		if num > 1048576 {
			return fastssz.ErrIncorrectListSize
		}
		for _, elem := range p.Transactions {
			{
				elemIndx := hh.Index()
				byteLen := uint64(len(elem))
				if byteLen > 1073741824 {
					return fastssz.ErrIncorrectListSize
				}
				hh.AppendBytes32(elem)
				hh.MerkleizeWithMixin(elemIndx, byteLen, (1073741824+31)/32)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 1048576)
	}

	// Field (14) 'Withdrawals'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Withdrawals))
		if num > 16 {
			return fastssz.ErrIncorrectListSize
		}
		for _, elem := range p.Withdrawals {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
	}

	// Field (15) 'BlobGasUsed'
	hh.PutUint64(uint64(p.BlobGasUsed))

	// Field (16) 'ExcessBlobGas'
	hh.PutUint64(uint64(p.ExcessBlobGas))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the ExecutionPayload object.
func (p *ExecutionPayload) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(p)
}

/* -------------------------------------------------------------------------- */
/*                                    JSON                                    */
/* -------------------------------------------------------------------------- */

// MarshalJSON marshals as JSON.
func (p *ExecutionPayload) MarshalJSON() ([]byte, error) {
	type ExecutionPayload struct {
		ParentHash    gethprimitives.ExecutionHash    `json:"parentHash"`
		FeeRecipient  gethprimitives.ExecutionAddress `json:"feeRecipient"`
		StateRoot     bytes.B32                       `json:"stateRoot"`
		ReceiptsRoot  bytes.B32                       `json:"receiptsRoot"`
		LogsBloom     bytes.B256                      `json:"logsBloom"`
		Random        bytes.B32                       `json:"prevRandao"`
		Number        math.U64                        `json:"blockNumber"`
		GasLimit      math.U64                        `json:"gasLimit"`
		GasUsed       math.U64                        `json:"gasUsed"`
		Timestamp     math.U64                        `json:"timestamp"`
		ExtraData     bytes.Bytes                     `json:"extraData"`
		BaseFeePerGas *math.U256Hex                   `json:"baseFeePerGas"`
		BlockHash     gethprimitives.ExecutionHash    `json:"blockHash"`
		Transactions  []bytes.Bytes                   `json:"transactions"`
		Withdrawals   []*engineprimitives.Withdrawal  `json:"withdrawals"`
		BlobGasUsed   math.U64                        `json:"blobGasUsed"`
		ExcessBlobGas math.U64                        `json:"excessBlobGas"`
	}
	var enc ExecutionPayload
	enc.ParentHash = p.ParentHash
	enc.FeeRecipient = p.FeeRecipient
	enc.StateRoot = p.StateRoot
	enc.ReceiptsRoot = p.ReceiptsRoot
	enc.LogsBloom = p.LogsBloom
	enc.Random = p.Random
	enc.Number = p.Number
	enc.GasLimit = p.GasLimit
	enc.GasUsed = p.GasUsed
	enc.Timestamp = p.Timestamp
	enc.ExtraData = p.ExtraData
	enc.BaseFeePerGas = (*math.U256Hex)(p.BaseFeePerGas)
	enc.BlockHash = p.BlockHash
	enc.Transactions = make([]bytes.Bytes, len(p.Transactions))
	for k, v := range p.Transactions {
		enc.Transactions[k] = v
	}
	enc.Withdrawals = p.Withdrawals
	enc.BlobGasUsed = p.BlobGasUsed
	enc.ExcessBlobGas = p.ExcessBlobGas
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
//
//nolint:funlen // todo fix.
func (p *ExecutionPayload) UnmarshalJSON(input []byte) error {
	type ExecutionPayload struct {
		ParentHash    *gethprimitives.ExecutionHash    `json:"parentHash"`
		FeeRecipient  *gethprimitives.ExecutionAddress `json:"feeRecipient"`
		StateRoot     *bytes.B32                       `json:"stateRoot"`
		ReceiptsRoot  *bytes.B32                       `json:"receiptsRoot"`
		LogsBloom     *bytes.B256                      `json:"logsBloom"`
		Random        *bytes.B32                       `json:"prevRandao"`
		Number        *math.U64                        `json:"blockNumber"`
		GasLimit      *math.U64                        `json:"gasLimit"`
		GasUsed       *math.U64                        `json:"gasUsed"`
		Timestamp     *math.U64                        `json:"timestamp"`
		ExtraData     *bytes.Bytes                     `json:"extraData"`
		BaseFeePerGas *math.U256Hex                    `json:"baseFeePerGas"`
		BlockHash     *gethprimitives.ExecutionHash    `json:"blockHash"`
		Transactions  []bytes.Bytes                    `json:"transactions"`
		Withdrawals   []*engineprimitives.Withdrawal   `json:"withdrawals"`
		BlobGasUsed   *math.U64                        `json:"blobGasUsed"`
		ExcessBlobGas *math.U64                        `json:"excessBlobGas"`
	}
	var dec ExecutionPayload
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentHash == nil {
		return errors.New(
			"missing required field 'parentHash' for ExecutionPayload",
		)
	}
	p.ParentHash = *dec.ParentHash
	if dec.FeeRecipient == nil {
		return errors.New(
			"missing required field 'feeRecipient' for ExecutionPayload",
		)
	}
	p.FeeRecipient = *dec.FeeRecipient
	if dec.StateRoot == nil {
		return errors.New(
			"missing required field 'stateRoot' for ExecutionPayload",
		)
	}
	p.StateRoot = *dec.StateRoot
	if dec.ReceiptsRoot == nil {
		return errors.New(
			"missing required field 'receiptsRoot' for ExecutionPayload",
		)
	}
	p.ReceiptsRoot = *dec.ReceiptsRoot
	if dec.LogsBloom == nil {
		return errors.New(
			"missing required field 'logsBloom' for ExecutionPayload",
		)
	}
	p.LogsBloom = *dec.LogsBloom
	if dec.Random == nil {
		return errors.New(
			"missing required field 'prevRandao' for ExecutionPayload",
		)
	}
	p.Random = *dec.Random
	if dec.Number == nil {
		return errors.New(
			"missing required field 'blockNumber' for ExecutionPayload",
		)
	}
	p.Number = *dec.Number
	if dec.GasLimit == nil {
		return errors.New(
			"missing required field 'gasLimit' for ExecutionPayload",
		)
	}
	p.GasLimit = *dec.GasLimit
	if dec.GasUsed == nil {
		return errors.New(
			"missing required field 'gasUsed' for ExecutionPayload",
		)
	}
	p.GasUsed = *dec.GasUsed
	if dec.Timestamp == nil {
		return errors.New(
			"missing required field 'timestamp' for ExecutionPayload",
		)
	}
	p.Timestamp = *dec.Timestamp
	if dec.ExtraData == nil {
		return errors.New(
			"missing required field 'extraData' for ExecutionPayload",
		)
	}
	p.ExtraData = *dec.ExtraData
	if dec.BaseFeePerGas == nil {
		return errors.New(
			"missing required field 'baseFeePerGas' for ExecutionPayload",
		)
	}
	p.BaseFeePerGas = (*math.U256)(dec.BaseFeePerGas)
	if dec.BlockHash == nil {
		return errors.New(
			"missing required field 'blockHash' for ExecutionPayload",
		)
	}
	p.BlockHash = *dec.BlockHash
	if dec.Transactions == nil {
		return errors.New(
			"missing required field 'transactions' for ExecutionPayload",
		)
	}
	p.Transactions = make([][]byte, len(dec.Transactions))
	for k, v := range dec.Transactions {
		p.Transactions[k] = v
	}
	if dec.Withdrawals != nil {
		p.Withdrawals = dec.Withdrawals
	}
	if dec.BlobGasUsed != nil {
		p.BlobGasUsed = *dec.BlobGasUsed
	}
	if dec.ExcessBlobGas != nil {
		p.ExcessBlobGas = *dec.ExcessBlobGas
	}
	return nil
}

// Empty returns an empty ExecutionPayload for the given fork version.
func (p *ExecutionPayload) Empty(_ uint32) *ExecutionPayload {
	return &ExecutionPayload{}
}

// Version returns the version of the ExecutionPayload.
func (p *ExecutionPayload) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the ExecutionPayload is nil.
func (p *ExecutionPayload) IsNil() bool {
	return p == nil
}

// IsBlinded checks if the ExecutionPayload is blinded.
func (p *ExecutionPayload) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutionPayload.
func (p *ExecutionPayload) GetParentHash() gethprimitives.ExecutionHash {
	return p.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutionPayload.
func (
	p *ExecutionPayload,
) GetFeeRecipient() gethprimitives.ExecutionAddress {
	return p.FeeRecipient
}

// GetStateRoot returns the state root of the ExecutionPayload.
func (p *ExecutionPayload) GetStateRoot() common.Bytes32 {
	return p.StateRoot
}

// GetReceiptsRoot returns the receipts root of the ExecutionPayload.
func (p *ExecutionPayload) GetReceiptsRoot() common.Bytes32 {
	return p.ReceiptsRoot
}

// GetLogsBloom returns the logs bloom of the ExecutionPayload.
func (p *ExecutionPayload) GetLogsBloom() bytes.B256 {
	return p.LogsBloom
}

// GetPrevRandao returns the previous Randao value of the ExecutionPayload.
func (p *ExecutionPayload) GetPrevRandao() common.Bytes32 {
	return p.Random
}

// GetNumber returns the block number of the ExecutionPayload.
func (p *ExecutionPayload) GetNumber() math.U64 {
	return p.Number
}

// GetGasLimit returns the gas limit of the ExecutionPayload.
func (p *ExecutionPayload) GetGasLimit() math.U64 {
	return p.GasLimit
}

// GetGasUsed returns the gas used of the ExecutionPayload.
func (p *ExecutionPayload) GetGasUsed() math.U64 {
	return p.GasUsed
}

// GetTimestamp returns the timestamp of the ExecutionPayload.
func (p *ExecutionPayload) GetTimestamp() math.U64 {
	return p.Timestamp
}

// GetExtraData returns the extra data of the ExecutionPayload.
func (p *ExecutionPayload) GetExtraData() []byte {
	return p.ExtraData
}

// GetBaseFeePerGas returns the base fee per gas of the ExecutionPayload.
func (p *ExecutionPayload) GetBaseFeePerGas() *math.U256 {
	return p.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionPayload.
func (p *ExecutionPayload) GetBlockHash() gethprimitives.ExecutionHash {
	return p.BlockHash
}

// GetTransactions returns the transactions of the ExecutionPayload.
func (p *ExecutionPayload) GetTransactions() engineprimitives.Transactions {
	return p.Transactions
}

// GetWithdrawals returns the withdrawals of the ExecutionPayload.
func (p *ExecutionPayload) GetWithdrawals() []*engineprimitives.Withdrawal {
	return p.Withdrawals
}

// GetBlobGasUsed returns the blob gas used of the ExecutionPayload.
func (p *ExecutionPayload) GetBlobGasUsed() math.U64 {
	return p.BlobGasUsed
}

// GetExcessBlobGas returns the excess blob gas of the ExecutionPayload.
func (p *ExecutionPayload) GetExcessBlobGas() math.U64 {
	return p.ExcessBlobGas
}

// ToHeader converts the ExecutionPayload to an ExecutionPayloadHeader.
func (p *ExecutionPayload) ToHeader(
	_ uint64,
	eth1ChainID uint64,
) (*ExecutionPayloadHeader, error) {
	var txsRoot common.Root

	// TODO: This is live on bArtio with a bug and needs to be hardforked
	// off of. This is a temporary solution to avoid breaking changes.
	if eth1ChainID == spec.TestnetEth1ChainID {
		txsRoot = engineprimitives.BartioTransactions(
			p.GetTransactions(),
		).HashTreeRoot()
	} else {
		txsRoot = p.GetTransactions().HashTreeRoot()
	}

	switch p.Version() {
	case version.Deneb, version.DenebPlus:
		return &ExecutionPayloadHeader{
			ParentHash:       p.GetParentHash(),
			FeeRecipient:     p.GetFeeRecipient(),
			StateRoot:        p.GetStateRoot(),
			ReceiptsRoot:     p.GetReceiptsRoot(),
			LogsBloom:        p.GetLogsBloom(),
			Random:           p.GetPrevRandao(),
			Number:           p.GetNumber(),
			GasLimit:         p.GetGasLimit(),
			GasUsed:          p.GetGasUsed(),
			Timestamp:        p.GetTimestamp(),
			ExtraData:        p.GetExtraData(),
			BaseFeePerGas:    p.GetBaseFeePerGas(),
			BlockHash:        p.GetBlockHash(),
			TransactionsRoot: txsRoot,
			WithdrawalsRoot: engineprimitives.Withdrawals(p.GetWithdrawals()).
				HashTreeRoot(),
			BlobGasUsed:   p.GetBlobGasUsed(),
			ExcessBlobGas: p.GetExcessBlobGas(),
		}, nil
	default:
		return nil, errors.New("unknown fork version")
	}
}

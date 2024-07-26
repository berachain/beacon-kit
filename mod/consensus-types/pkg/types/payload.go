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
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	essz "github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
	"golang.org/x/sync/errgroup"
)

// ExecutionPayloadStaticSize is the static size of the ExecutionPayload.
const ExecutionPayloadStaticSize uint32 = 528

// ExecutionPayload is the execution payload for Deneb.
//
//nolint:lll
type ExecutionPayload struct {
	// ParentHash is the hash of the parent block.
	ParentHash gethprimitives.ExecutionHash `json:"parentHash"    gencodec:"required"`
	// FeeRecipient is the address of the fee recipient.
	FeeRecipient gethprimitives.ExecutionAddress `json:"feeRecipient"  gencodec:"required"`
	// StateRoot is the root of the state trie.
	StateRoot common.Bytes32 `json:"stateRoot"     gencodec:"required"`
	// ReceiptsRoot is the root of the receipts trie.
	ReceiptsRoot common.Bytes32 `json:"receiptsRoot"  gencodec:"required"`
	// LogsBloom is the bloom filter for the logs.
	LogsBloom bytes.B256 `json:"logsBloom"     gencodec:"required"`
	// Random is the prevRandao value.
	Random common.Bytes32 `json:"prevRandao"    gencodec:"required"`
	// Number is the block number.
	Number math.U64 `json:"blockNumber"   gencodec:"required"`
	// GasLimit is the gas limit for the block.
	GasLimit math.U64 `json:"gasLimit"      gencodec:"required"`
	// GasUsed is the amount of gas used in the block.
	GasUsed math.U64 `json:"gasUsed"       gencodec:"required"`
	// Timestamp is the timestamp of the block.
	Timestamp math.U64 `json:"timestamp"     gencodec:"required"`
	// ExtraData is the extra data of the block.
	ExtraData bytes.Bytes `json:"extraData"     gencodec:"required"`
	// BaseFeePerGas is the base fee per gas.
	BaseFeePerGas math.Wei `json:"baseFeePerGas" gencodec:"required"`
	// BlockHash is the hash of the block.
	BlockHash gethprimitives.ExecutionHash `json:"blockHash"     gencodec:"required"`
	// Transactions is the list of transactions in the block.
	Transactions [][]byte `json:"transactions"  gencodec:"required"`
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
	ssz.DefineStaticBytes(codec, &p.BaseFeePerGas)
	ssz.DefineStaticBytes(codec, &p.BlockHash)
	ssz.DefineSliceOfDynamicBytesOffset(
		codec,
		&p.Transactions,
		1048576,
		1073741824,
	)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &p.Withdrawals, 16)
	ssz.DefineUint64(codec, &p.BlobGasUsed)
	ssz.DefineUint64(codec, &p.ExcessBlobGas)

	// Define the dynamic data (fields)
	ssz.DefineDynamicBytesContent(codec, (*[]byte)(&p.ExtraData), 32)
	ssz.DefineSliceOfDynamicBytesContent(
		codec,
		&p.Transactions,
		1048576,
		1073741824,
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
	hh.PutBytes(p.BaseFeePerGas[:])

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
			if err := elem.HashTreeRootWith(hh); err != nil {
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
func (p *ExecutionPayload) GetLogsBloom() []byte {
	return p.LogsBloom[:]
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
func (p *ExecutionPayload) GetBaseFeePerGas() math.Wei {
	return p.BaseFeePerGas
}

// GetBlockHash returns the block hash of the ExecutionPayload.
func (p *ExecutionPayload) GetBlockHash() gethprimitives.ExecutionHash {
	return p.BlockHash
}

// GetTransactions returns the transactions of the ExecutionPayload.
func (p *ExecutionPayload) GetTransactions() [][]byte {
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
	txsMerkleizer *merkle.Merkleizer[[32]byte, common.Root],
	maxWithdrawalsPerPayload uint64,
) (*ExecutionPayloadHeader, error) {
	// Get the merkle roots of transactions and withdrawals in parallel.
	var (
		g, _            = errgroup.WithContext(context.Background())
		txsRoot         common.Root
		withdrawalsRoot common.Root
	)

	g.Go(func() error {
		var txsRootErr error
		// TODO: This is live on bArtio with a bug and needs to be hardforked
		// off of.
		txsRoot, txsRootErr = engineprimitives.Transactions(
			p.GetTransactions(),
		).HashTreeRootWith(txsMerkleizer)
		return txsRootErr
	})

	g.Go(func() error {
		var withdrawalsRootErr error
		wds := essz.ListFromElements(
			maxWithdrawalsPerPayload,
			p.GetWithdrawals()...)
		withdrawalsRoot, withdrawalsRootErr = wds.HashTreeRoot()
		return withdrawalsRootErr
	})

	// Wait for the goroutines to finish.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	switch p.Version() {
	case version.Deneb, version.DenebPlus:
		return &ExecutionPayloadHeader{
			ParentHash:       p.GetParentHash(),
			FeeRecipient:     p.GetFeeRecipient(),
			StateRoot:        p.GetStateRoot(),
			ReceiptsRoot:     p.GetReceiptsRoot(),
			LogsBloom:        [256]byte(p.GetLogsBloom()),
			Random:           p.GetPrevRandao(),
			Number:           p.GetNumber(),
			GasLimit:         p.GetGasLimit(),
			GasUsed:          p.GetGasUsed(),
			Timestamp:        p.GetTimestamp(),
			ExtraData:        p.GetExtraData(),
			BaseFeePerGas:    p.GetBaseFeePerGas(),
			BlockHash:        p.GetBlockHash(),
			TransactionsRoot: txsRoot,
			WithdrawalsRoot:  withdrawalsRoot,
			BlobGasUsed:      p.GetBlobGasUsed(),
			ExcessBlobGas:    p.GetExcessBlobGas(),
		}, nil
	default:
		return nil, errors.New("unknown fork version")
	}
}

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
	"fmt"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	ssz "github.com/ferranbt/fastssz"
)

const (
	// ExecutionPayloadStaticSize is the static size of the ExecutionPayload.
	ExecutionPayloadStaticSize uint32 = 528

	// ExtraDataSize is the size of ExtraData in bytes.
	ExtraDataSize = 32
)

var (
	_ constraints.SSZVersionedMarshallableRootable = (*ExecutionPayload)(nil)
)

// ExecutionPayload represents the payload of an execution block.
type ExecutionPayload struct {
	Versionable `json:"-"`

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
	// Transactions is the list of transactions in the block.
	Transactions engineprimitives.Transactions `json:"transactions"`
	// Withdrawals is the list of withdrawals in the block.
	Withdrawals []*engineprimitives.Withdrawal `json:"withdrawals"`
	// BlobGasUsed is the amount of blob gas used in the block.
	BlobGasUsed math.U64 `json:"blobGasUsed"`
	// ExcessBlobGas is the amount of excess blob gas in the block.
	ExcessBlobGas math.U64 `json:"excessBlobGas"`
}

func NewEmptyExecutionPayloadWithVersion(forkVersion common.Version) *ExecutionPayload {
	ep := &ExecutionPayload{
		Versionable:   NewVersionable(forkVersion),
		BaseFeePerGas: &math.U256{},
	}

	// For any fork version Capella onwards, non-nil withdrawals are required.
	if version.EqualsOrIsAfter(forkVersion, version.Capella()) {
		ep.Withdrawals = make([]*engineprimitives.Withdrawal, 0)
	}
	return ep
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the total size of the object in SSZ encoding.
func (p *ExecutionPayload) SizeSSZ() int {
	size := int(ExecutionPayloadStaticSize)
	size += len(p.ExtraData)
	// Transactions offset + each tx offset + tx data
	size += len(p.Transactions) * 4
	for _, tx := range p.Transactions {
		size += len(tx)
	}
	// Withdrawals: each withdrawal is 44 bytes
	if p.Withdrawals != nil {
		size += len(p.Withdrawals) * 44
	}
	return size
}

// MarshalSSZ serializes the ExecutionPayload object into a slice of bytes.
func (p *ExecutionPayload) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, p.SizeSSZ())
	return p.MarshalSSZTo(buf)
}

func (p *ExecutionPayload) ValidateAfterDecodingSSZ() error {
	// For any fork version Capella onwards, non-nil withdrawals are required.
	if p.Withdrawals == nil && version.EqualsOrIsAfter(p.GetForkVersion(), version.Capella()) {
		p.Withdrawals = make([]*engineprimitives.Withdrawal, 0)
	}
	return nil
}

// HashTreeRoot returns the hash tree root of the ExecutionPayload.
func (p *ExecutionPayload) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)
	if err := p.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo serializes the ExecutionPayload object into a writer.
func (p *ExecutionPayload) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Static fields
	dst = append(dst, p.ParentHash[:]...)
	dst = append(dst, p.FeeRecipient[:]...)
	dst = append(dst, p.StateRoot[:]...)
	dst = append(dst, p.ReceiptsRoot[:]...)
	dst = append(dst, p.LogsBloom[:]...)
	dst = append(dst, p.Random[:]...)
	dst = ssz.MarshalUint64(dst, uint64(p.Number))
	dst = ssz.MarshalUint64(dst, uint64(p.GasLimit))
	dst = ssz.MarshalUint64(dst, uint64(p.GasUsed))
	dst = ssz.MarshalUint64(dst, uint64(p.Timestamp))

	// Calculate offsets
	extraDataOffset := uint32(ExecutionPayloadStaticSize)
	transactionsOffset := extraDataOffset + uint32(len(p.ExtraData))
	txDataSize := uint32(len(p.Transactions) * 4)
	for _, tx := range p.Transactions {
		txDataSize += uint32(len(tx))
	}
	withdrawalsOffset := transactionsOffset + txDataSize

	// Write offsets
	dst = ssz.MarshalUint32(dst, extraDataOffset)

	// BaseFeePerGas
	var bz []byte
	if p.BaseFeePerGas == nil {
		// Use zero value for nil BaseFeePerGas
		bz = make([]byte, 32)
	} else {
		var err error
		bz, err = p.BaseFeePerGas.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	}
	dst = append(dst, bz...)

	// Static fields continued
	dst = append(dst, p.BlockHash[:]...)
	dst = ssz.MarshalUint32(dst, transactionsOffset)
	dst = ssz.MarshalUint32(dst, withdrawalsOffset)
	dst = ssz.MarshalUint64(dst, uint64(p.BlobGasUsed))
	dst = ssz.MarshalUint64(dst, uint64(p.ExcessBlobGas))

	// Dynamic fields
	// ExtraData
	dst = append(dst, p.ExtraData...)

	// Transactions
	txOffsets := make([]uint32, len(p.Transactions))
	currOffset := transactionsOffset + uint32(len(p.Transactions)*4)
	for i, tx := range p.Transactions {
		txOffsets[i] = currOffset
		currOffset += uint32(len(tx))
	}
	for _, offset := range txOffsets {
		dst = ssz.MarshalUint32(dst, offset)
	}
	for _, tx := range p.Transactions {
		dst = append(dst, tx...)
	}

	// Withdrawals
	for _, w := range p.Withdrawals {
		if w == nil {
			return nil, errors.New("nil withdrawal in ExecutionPayload")
		}
		var err error
		dst, err = w.MarshalSSZTo(dst)
		if err != nil {
			return nil, err
		}
	}

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the ExecutionPayload object.
func (p *ExecutionPayload) UnmarshalSSZ(buf []byte) error {
	if len(buf) < int(ExecutionPayloadStaticSize) {
		return ssz.ErrSize
	}

	// Static fields
	copy(p.ParentHash[:], buf[0:32])
	copy(p.FeeRecipient[:], buf[32:52])
	copy(p.StateRoot[:], buf[52:84])
	copy(p.ReceiptsRoot[:], buf[84:116])
	copy(p.LogsBloom[:], buf[116:372])
	copy(p.Random[:], buf[372:404])
	p.Number = math.U64(ssz.UnmarshallUint64(buf[404:412]))
	p.GasLimit = math.U64(ssz.UnmarshallUint64(buf[412:420]))
	p.GasUsed = math.U64(ssz.UnmarshallUint64(buf[420:428]))
	p.Timestamp = math.U64(ssz.UnmarshallUint64(buf[428:436]))

	// Read offsets
	extraDataOffset := ssz.UnmarshallUint32(buf[436:440])

	// BaseFeePerGas
	if p.BaseFeePerGas == nil {
		p.BaseFeePerGas = &math.U256{}
	}
	if err := p.BaseFeePerGas.UnmarshalSSZ(buf[440:472]); err != nil {
		return err
	}

	// More static fields
	copy(p.BlockHash[:], buf[472:504])
	transactionsOffset := ssz.UnmarshallUint32(buf[504:508])
	withdrawalsOffset := ssz.UnmarshallUint32(buf[508:512])
	p.BlobGasUsed = math.U64(ssz.UnmarshallUint64(buf[512:520]))
	p.ExcessBlobGas = math.U64(ssz.UnmarshallUint64(buf[520:528]))

	// Dynamic fields
	// ExtraData
	if extraDataOffset > uint32(len(buf)) || transactionsOffset > uint32(len(buf)) || extraDataOffset > transactionsOffset {
		return ssz.ErrInvalidVariableOffset
	}
	p.ExtraData = append([]byte(nil), buf[extraDataOffset:transactionsOffset]...)
	if len(p.ExtraData) > 32 {
		return errors.New("extra data too large")
	}

	// Transactions
	if transactionsOffset > uint32(len(buf)) || withdrawalsOffset > uint32(len(buf)) || transactionsOffset > withdrawalsOffset {
		return ssz.ErrInvalidVariableOffset
	}
	txData := buf[transactionsOffset:withdrawalsOffset]
	if len(txData) > 0 {
		// Read transaction offsets
		if len(txData) < 4 {
			return ssz.ErrSize
		}
		firstTxOffset := ssz.UnmarshallUint32(txData[0:4])
		if firstTxOffset < transactionsOffset || firstTxOffset > uint32(len(buf)) {
			return ssz.ErrInvalidVariableOffset
		}
		numTxs := int((firstTxOffset - transactionsOffset) / 4)
		if numTxs > int(constants.MaxTxsPerPayload) {
			return errors.New("too many transactions")
		}

		txOffsets := make([]uint32, numTxs)
		for i := 0; i < numTxs; i++ {
			txOffsets[i] = ssz.UnmarshallUint32(txData[i*4 : (i+1)*4])
		}

		p.Transactions = make([][]byte, numTxs)
		for i := 0; i < numTxs; i++ {
			start := txOffsets[i]
			end := withdrawalsOffset
			if i+1 < numTxs {
				end = txOffsets[i+1]
			}
			if start > uint32(len(buf)) || end > uint32(len(buf)) || start > end {
				return ssz.ErrInvalidVariableOffset
			}
			p.Transactions[i] = append([]byte(nil), buf[start:end]...)
			if len(p.Transactions[i]) > int(constants.MaxBytesPerTx) {
				return errors.New("transaction too large")
			}
		}
	} else {
		p.Transactions = make([][]byte, 0)
	}

	// Withdrawals
	if withdrawalsOffset > uint32(len(buf)) {
		return ssz.ErrInvalidVariableOffset
	}
	wData := buf[withdrawalsOffset:]
	if len(wData)%44 != 0 {
		return ssz.ErrSize
	}
	numWithdrawals := len(wData) / 44
	if numWithdrawals > 16 {
		return errors.New("too many withdrawals")
	}
	p.Withdrawals = make([]*engineprimitives.Withdrawal, numWithdrawals)
	for i := 0; i < numWithdrawals; i++ {
		p.Withdrawals[i] = &engineprimitives.Withdrawal{}
		if err := p.Withdrawals[i].UnmarshalSSZ(wData[i*44 : (i+1)*44]); err != nil {
			return err
		}
	}

	return p.ValidateAfterDecodingSSZ()
}

// HashTreeRootWith ssz hashes the ExecutionPayload object with a hasher.
//
//nolint:mnd // will be deprecated eventually.
func (p *ExecutionPayload) HashTreeRootWith(hh ssz.HashWalker) error {
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
			return ssz.ErrIncorrectListSize
		}
		hh.Append(p.ExtraData)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (32+31)/32)
	}

	// Field (11) 'BaseFeePerGas'
	if p.BaseFeePerGas == nil {
		// Use zero value for nil BaseFeePerGas
		hh.PutBytes(make([]byte, 32))
	} else {
		bz, err := p.BaseFeePerGas.MarshalSSZ()
		if err != nil {
			return err
		}
		hh.PutBytes(bz)
	}

	// Field (12) 'BlockHash'
	hh.PutBytes(p.BlockHash[:])

	// Field (13) 'Transactions'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Transactions))
		if num > 1048576 {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range p.Transactions {
			{
				elemIndx := hh.Index()
				byteLen := uint64(len(elem))
				if byteLen > 1073741824 {
					return ssz.ErrIncorrectListSize
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
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range p.Withdrawals {
			if elem == nil {
				return errors.New("nil withdrawal in ExecutionPayload")
			}
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
func (p *ExecutionPayload) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(p)
}

/* -------------------------------------------------------------------------- */
/*                                    JSON                                    */
/* -------------------------------------------------------------------------- */

// MarshalJSON marshals as JSON.
func (p ExecutionPayload) MarshalJSON() ([]byte, error) {
	type ExecutionPayload struct {
		ParentHash    common.ExecutionHash           `json:"parentHash"`
		FeeRecipient  common.ExecutionAddress        `json:"feeRecipient"`
		StateRoot     bytes.B32                      `json:"stateRoot"`
		ReceiptsRoot  bytes.B32                      `json:"receiptsRoot"`
		LogsBloom     bytes.B256                     `json:"logsBloom"`
		Random        bytes.B32                      `json:"prevRandao"`
		Number        math.U64                       `json:"blockNumber"`
		GasLimit      math.U64                       `json:"gasLimit"`
		GasUsed       math.U64                       `json:"gasUsed"`
		Timestamp     math.U64                       `json:"timestamp"`
		ExtraData     bytes.Bytes                    `json:"extraData"`
		BaseFeePerGas *math.U256Hex                  `json:"baseFeePerGas"`
		BlockHash     common.ExecutionHash           `json:"blockHash"`
		Transactions  []bytes.Bytes                  `json:"transactions"`
		Withdrawals   []*engineprimitives.Withdrawal `json:"withdrawals"`
		BlobGasUsed   math.U64                       `json:"blobGasUsed"`
		ExcessBlobGas math.U64                       `json:"excessBlobGas"`
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
		ParentHash    *common.ExecutionHash          `json:"parentHash"`
		FeeRecipient  *common.ExecutionAddress       `json:"feeRecipient"`
		StateRoot     *bytes.B32                     `json:"stateRoot"`
		ReceiptsRoot  *bytes.B32                     `json:"receiptsRoot"`
		LogsBloom     *bytes.B256                    `json:"logsBloom"`
		Random        *bytes.B32                     `json:"prevRandao"`
		Number        *math.U64                      `json:"blockNumber"`
		GasLimit      *math.U64                      `json:"gasLimit"`
		GasUsed       *math.U64                      `json:"gasUsed"`
		Timestamp     *math.U64                      `json:"timestamp"`
		ExtraData     *bytes.Bytes                   `json:"extraData"`
		BaseFeePerGas *math.U256Hex                  `json:"baseFeePerGas"`
		BlockHash     *common.ExecutionHash          `json:"blockHash"`
		Transactions  []bytes.Bytes                  `json:"transactions"`
		Withdrawals   []*engineprimitives.Withdrawal `json:"withdrawals"`
		BlobGasUsed   *math.U64                      `json:"blobGasUsed"`
		ExcessBlobGas *math.U64                      `json:"excessBlobGas"`
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

/* -------------------------------------------------------------------------- */
/*                                   Getters                                  */
/* -------------------------------------------------------------------------- */

// IsBlinded checks if the ExecutionPayload is blinded.
func (p *ExecutionPayload) IsBlinded() bool {
	return false
}

// GetParentHash returns the parent hash of the ExecutionPayload.
func (p *ExecutionPayload) GetParentHash() common.ExecutionHash {
	return p.ParentHash
}

// GetFeeRecipient returns the fee recipient address of the ExecutionPayload.
func (p *ExecutionPayload) GetFeeRecipient() common.ExecutionAddress {
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
func (p *ExecutionPayload) GetBlockHash() common.ExecutionHash {
	return p.BlockHash
}

// GetTransactions returns the transactions of the ExecutionPayload.
func (p *ExecutionPayload) GetTransactions() engineprimitives.Transactions {
	return p.Transactions
}

// GetWithdrawals returns the withdrawals of the ExecutionPayload.
func (p *ExecutionPayload) GetWithdrawals() engineprimitives.Withdrawals {
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
func (p *ExecutionPayload) ToHeader() (*ExecutionPayloadHeader, error) {
	switch p.GetForkVersion() {
	case version.Deneb(), version.Deneb1(), version.Electra(), version.Electra1():
		// Compute transactions root
		transactionsRoot, err := computeTransactionsRoot(p.GetTransactions())
		if err != nil {
			return nil, err
		}

		// Compute withdrawals root
		withdrawalsRoot, err := computeWithdrawalsRoot(p.GetWithdrawals())
		if err != nil {
			return nil, err
		}

		return &ExecutionPayloadHeader{
			Versionable:      p.Versionable,
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
			TransactionsRoot: transactionsRoot,
			WithdrawalsRoot:  withdrawalsRoot,
			BlobGasUsed:      p.GetBlobGasUsed(),
			ExcessBlobGas:    p.GetExcessBlobGas(),
		}, nil
	default:
		return nil, errors.New("unknown fork version")
	}
}

// computeTransactionsRoot returns the hash tree root of transactions.
func computeTransactionsRoot(transactions engineprimitives.Transactions) (common.Root, error) {
	if transactions == nil {
		return common.Root{}, nil
	}
	root, err := transactions.HashTreeRoot()
	if err != nil {
		return common.Root{}, fmt.Errorf("failed to compute transactions root: %w", err)
	}
	return common.Root(root), nil
}

// computeWithdrawalsRoot returns the hash tree root of withdrawals.
func computeWithdrawalsRoot(withdrawals engineprimitives.Withdrawals) (common.Root, error) {
	if withdrawals == nil {
		return common.Root{}, nil
	}
	root, err := withdrawals.HashTreeRoot()
	if err != nil {
		return common.Root{}, fmt.Errorf("failed to compute withdrawals root: %w", err)
	}
	return common.Root(root), nil
}

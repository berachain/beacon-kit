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

//go:build test

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// ExecutionPayloadHeaderStaticSizeKaralabe is the static size of the ExecutionPayloadHeader.
const ExecutionPayloadHeaderStaticSizeKaralabe uint32 = 584

// Compile-time assertions to ensure ExecutionPayloadHeaderKaralabe implements necessary interfaces.
var _ ssz.DynamicObject = (*ExecutionPayloadHeaderKaralabe)(nil)

// versionableKaralabeHeader is a helper struct that implements the Versionable interface.
type versionableKaralabeHeader struct {
	forkVersion common.Version
}

// NewVersionableKaralabe creates a new versionable object.
func NewVersionableKaralabeHeader(forkVersion common.Version) constraints.Versionable {
	return &versionableKaralabeHeader{forkVersion: forkVersion}
}

// GetForkVersion returns the fork version of the versionable object.
func (v *versionableKaralabeHeader) GetForkVersion() common.Version {
	return v.forkVersion
}

// ExecutionPayloadHeaderKaralabe represents the payload header of an execution block.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
type ExecutionPayloadHeaderKaralabe struct {
	constraints.Versionable

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

// NewEmptyExecutionPayloadHeaderWithVersionKaralabe creates an empty ExecutionPayloadHeader with version.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func NewEmptyExecutionPayloadHeaderWithVersionKaralabe(forkVersion common.Version) *ExecutionPayloadHeaderKaralabe {
	return &ExecutionPayloadHeaderKaralabe{
		Versionable:   NewVersionableKaralabeHeader(forkVersion),
		BaseFeePerGas: &math.U256{},
		ExtraData:     make([]byte, 0),
	}
}

// SizeSSZ returns either the static size of the object if fixed == true, or
// the total size otherwise.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (h *ExecutionPayloadHeaderKaralabe) SizeSSZ(fixed bool) uint32 {
	size := ExecutionPayloadHeaderStaticSizeKaralabe
	if fixed {
		return size
	}
	size += uint32(len(h.ExtraData))

	return size
}

// DefineSSZ defines how an object is encoded/decoded.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (h *ExecutionPayloadHeaderKaralabe) DefineSSZ(codec *ssz.Codec) {
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
	//nolint:mnd // TODO: get from accessible chainspec field params
	ssz.DefineDynamicBytesOffset(codec, (*[]byte)(&h.ExtraData), 32)
	ssz.DefineUint256(codec, &h.BaseFeePerGas)
	ssz.DefineStaticBytes(codec, &h.BlockHash)
	ssz.DefineStaticBytes(codec, &h.TransactionsRoot)
	ssz.DefineStaticBytes(codec, &h.WithdrawalsRoot)
	ssz.DefineUint64(codec, &h.BlobGasUsed)
	ssz.DefineUint64(codec, &h.ExcessBlobGas)

	// Define the dynamic data (fields)
	//nolint:mnd // TODO: get from accessible chainspec field params
	ssz.DefineDynamicBytesContent(codec, (*[]byte)(&h.ExtraData), 32)
}

// MarshalSSZ serializes the ExecutionPayloadHeader object into a slice of
// bytes.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (h *ExecutionPayloadHeaderKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(h))
	return buf, ssz.EncodeToBytes(buf, h)
}

func (*ExecutionPayloadHeaderKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot returns the hash tree root of the ExecutionPayloadHeader.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (h *ExecutionPayloadHeaderKaralabe) HashTreeRoot() common.Root {
	return ssz.HashConcurrent(h)
}

// UnmarshalSSZ unmarshals the ExecutionPayloadHeader object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (h *ExecutionPayloadHeaderKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, h)
}

// TestExecutionPayloadHeaderCompatibility tests that the current ExecutionPayloadHeader implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestExecutionPayloadHeaderCompatibility(t *testing.T) {
	testCases := []struct {
		name    string
		version common.Version
		setup   func() (*types.ExecutionPayloadHeader, *ExecutionPayloadHeaderKaralabe)
	}{
		{
			name:    "Deneb - zero values",
			version: version.Deneb(),
			setup: func() (*types.ExecutionPayloadHeader, *ExecutionPayloadHeaderKaralabe) {
				current := types.NewEmptyExecutionPayloadHeaderWithVersion(version.Deneb())
				if current.ExtraData == nil {
					current.ExtraData = make([]byte, 0)
				}
				karalabe := NewEmptyExecutionPayloadHeaderWithVersionKaralabe(version.Deneb())
				return current, karalabe
			},
		},
		{
			name:    "Deneb - typical header",
			version: version.Deneb(),
			setup: func() (*types.ExecutionPayloadHeader, *ExecutionPayloadHeaderKaralabe) {
				// Create parent hash
				var parentHash common.ExecutionHash
				for i := range parentHash {
					parentHash[i] = byte(i)
				}

				// Create fee recipient
				var feeRecipient common.ExecutionAddress
				for i := range feeRecipient {
					feeRecipient[i] = byte(i * 2)
				}

				// Create state root
				var stateRoot common.Bytes32
				for i := range stateRoot {
					stateRoot[i] = byte(i + 32)
				}

				// Create receipts root
				var receiptsRoot common.Bytes32
				for i := range receiptsRoot {
					receiptsRoot[i] = byte(i + 64)
				}

				// Create logs bloom
				var logsBloom bytes.B256
				for i := range logsBloom {
					logsBloom[i] = byte(i % 256)
				}

				// Create random
				var random common.Bytes32
				for i := range random {
					random[i] = byte(i + 96)
				}

				// Create block hash
				var blockHash common.ExecutionHash
				for i := range blockHash {
					blockHash[i] = byte(255 - i)
				}

				// Create transactions root
				var transactionsRoot common.Root
				for i := range transactionsRoot {
					transactionsRoot[i] = byte(i + 128)
				}

				// Create withdrawals root
				var withdrawalsRoot common.Root
				for i := range withdrawalsRoot {
					withdrawalsRoot[i] = byte(i + 160)
				}

				// Create base fee
				baseFee := math.NewU256(1000000000) // 1 gwei

				// Create extra data
				extraData := []byte("test extra data")

				current := &types.ExecutionPayloadHeader{
					Versionable:      types.NewVersionable(version.Deneb()),
					ParentHash:       parentHash,
					FeeRecipient:     feeRecipient,
					StateRoot:        stateRoot,
					ReceiptsRoot:     receiptsRoot,
					LogsBloom:        logsBloom,
					Random:           random,
					Number:           12345,
					GasLimit:         30000000,
					GasUsed:          21000000,
					Timestamp:        1700000000,
					ExtraData:        extraData,
					BaseFeePerGas:    baseFee,
					BlockHash:        blockHash,
					TransactionsRoot: transactionsRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      131072,
					ExcessBlobGas:    262144,
				}

				karalabe := &ExecutionPayloadHeaderKaralabe{
					Versionable:      NewVersionableKaralabeHeader(version.Deneb()),
					ParentHash:       parentHash,
					FeeRecipient:     feeRecipient,
					StateRoot:        stateRoot,
					ReceiptsRoot:     receiptsRoot,
					LogsBloom:        logsBloom,
					Random:           random,
					Number:           12345,
					GasLimit:         30000000,
					GasUsed:          21000000,
					Timestamp:        1700000000,
					ExtraData:        extraData,
					BaseFeePerGas:    baseFee,
					BlockHash:        blockHash,
					TransactionsRoot: transactionsRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      131072,
					ExcessBlobGas:    262144,
				}

				return current, karalabe
			},
		},
		{
			name:    "Deneb - maximum extra data",
			version: version.Deneb(),
			setup: func() (*types.ExecutionPayloadHeader, *ExecutionPayloadHeaderKaralabe) {
				// Create maximum extra data (32 bytes)
				extraData := make([]byte, 32)
				for i := range extraData {
					extraData[i] = byte(i)
				}

				current := types.NewEmptyExecutionPayloadHeaderWithVersion(version.Deneb())
				current.ExtraData = extraData

				karalabe := NewEmptyExecutionPayloadHeaderWithVersionKaralabe(version.Deneb())
				karalabe.ExtraData = extraData

				return current, karalabe
			},
		},
		{
			name:    "Capella - with withdrawals root",
			version: version.Capella(),
			setup: func() (*types.ExecutionPayloadHeader, *ExecutionPayloadHeaderKaralabe) {
				// Create withdrawals root
				var withdrawalsRoot common.Root
				for i := range withdrawalsRoot {
					withdrawalsRoot[i] = byte(i * 3)
				}

				current := types.NewEmptyExecutionPayloadHeaderWithVersion(version.Capella())
				current.WithdrawalsRoot = withdrawalsRoot
				// Capella doesn't have blob gas fields, so set them to zero
				current.BlobGasUsed = 0
				current.ExcessBlobGas = 0
				if current.ExtraData == nil {
					current.ExtraData = make([]byte, 0)
				}

				karalabe := NewEmptyExecutionPayloadHeaderWithVersionKaralabe(version.Capella())
				karalabe.WithdrawalsRoot = withdrawalsRoot
				karalabe.BlobGasUsed = 0
				karalabe.ExcessBlobGas = 0

				return current, karalabe
			},
		},
		{
			name:    "Deneb - all maximum values",
			version: version.Deneb(),
			setup: func() (*types.ExecutionPayloadHeader, *ExecutionPayloadHeaderKaralabe) {
				// Create all fields with maximum values
				var parentHash common.ExecutionHash
				var feeRecipient common.ExecutionAddress
				var stateRoot common.Bytes32
				var receiptsRoot common.Bytes32
				var logsBloom bytes.B256
				var random common.Bytes32
				var blockHash common.ExecutionHash
				var transactionsRoot common.Root
				var withdrawalsRoot common.Root

				// Fill with max values
				for i := range parentHash {
					parentHash[i] = 0xFF
				}
				for i := range feeRecipient {
					feeRecipient[i] = 0xFF
				}
				for i := range stateRoot {
					stateRoot[i] = 0xFF
				}
				for i := range receiptsRoot {
					receiptsRoot[i] = 0xFF
				}
				for i := range logsBloom {
					logsBloom[i] = 0xFF
				}
				for i := range random {
					random[i] = 0xFF
				}
				for i := range blockHash {
					blockHash[i] = 0xFF
				}
				for i := range transactionsRoot {
					transactionsRoot[i] = 0xFF
				}
				for i := range withdrawalsRoot {
					withdrawalsRoot[i] = 0xFF
				}

				// Max base fee
				baseFee := &math.U256{}
				baseFee.SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

				current := &types.ExecutionPayloadHeader{
					Versionable:      types.NewVersionable(version.Deneb()),
					ParentHash:       parentHash,
					FeeRecipient:     feeRecipient,
					StateRoot:        stateRoot,
					ReceiptsRoot:     receiptsRoot,
					LogsBloom:        logsBloom,
					Random:           random,
					Number:           ^math.U64(0),
					GasLimit:         ^math.U64(0),
					GasUsed:          ^math.U64(0),
					Timestamp:        ^math.U64(0),
					ExtraData:        []byte{},
					BaseFeePerGas:    baseFee,
					BlockHash:        blockHash,
					TransactionsRoot: transactionsRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      ^math.U64(0),
					ExcessBlobGas:    ^math.U64(0),
				}

				karalabe := &ExecutionPayloadHeaderKaralabe{
					Versionable:      NewVersionableKaralabeHeader(version.Deneb()),
					ParentHash:       parentHash,
					FeeRecipient:     feeRecipient,
					StateRoot:        stateRoot,
					ReceiptsRoot:     receiptsRoot,
					LogsBloom:        logsBloom,
					Random:           random,
					Number:           ^math.U64(0),
					GasLimit:         ^math.U64(0),
					GasUsed:          ^math.U64(0),
					Timestamp:        ^math.U64(0),
					ExtraData:        []byte{},
					BaseFeePerGas:    baseFee,
					BlockHash:        blockHash,
					TransactionsRoot: transactionsRoot,
					WithdrawalsRoot:  withdrawalsRoot,
					BlobGasUsed:      ^math.U64(0),
					ExcessBlobGas:    ^math.U64(0),
				}

				return current, karalabe
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current, karalabe := tc.setup()

			// Test Marshal
			currentBytes, err1 := current.MarshalSSZ()
			require.NoError(t, err1, "current MarshalSSZ should not error")

			karalableBytes, err2 := karalabe.MarshalSSZ()
			require.NoError(t, err2, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalableBytes, currentBytes, "marshaled bytes should be identical")

			// Test Size
			require.Equal(t, karalabe.SizeSSZ(false), uint32(current.SizeSSZ()), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := types.NewEmptyExecutionPayloadHeaderWithVersion(tc.version)
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")

			// Normalize ExtraData field for comparison (nil vs empty slice)
			if newCurrent.ExtraData == nil && len(current.ExtraData) == 0 {
				newCurrent.ExtraData = []byte{}
			}
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := NewEmptyExecutionPayloadHeaderWithVersionKaralabe(tc.version)
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")

			// Normalize ExtraData field for comparison (nil vs empty slice)
			if newKaralabe.ExtraData == nil && len(karalabe.ExtraData) == 0 {
				newKaralabe.ExtraData = []byte{}
			}
			require.Equal(t, karalabe, newKaralabe, "unmarshaled karalabe should match original")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabelRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabelRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

// TestExecutionPayloadHeaderCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestExecutionPayloadHeaderCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid header data
		var parentHash common.ExecutionHash
		var feeRecipient common.ExecutionAddress
		var stateRoot common.Bytes32
		var receiptsRoot common.Bytes32
		var logsBloom bytes.B256
		var random common.Bytes32
		var blockHash common.ExecutionHash
		var transactionsRoot common.Root
		var withdrawalsRoot common.Root

		// Use deterministic "random" data based on iteration
		for j := range parentHash {
			parentHash[j] = byte((i + j) % 256)
		}
		for j := range feeRecipient {
			feeRecipient[j] = byte((i*2 + j) % 256)
		}
		for j := range stateRoot {
			stateRoot[j] = byte((i*3 + j) % 256)
		}
		for j := range receiptsRoot {
			receiptsRoot[j] = byte((i*4 + j) % 256)
		}
		for j := range logsBloom {
			logsBloom[j] = byte((i*5 + j) % 256)
		}
		for j := range random {
			random[j] = byte((i*6 + j) % 256)
		}
		for j := range blockHash {
			blockHash[j] = byte((i*7 + j) % 256)
		}
		for j := range transactionsRoot {
			transactionsRoot[j] = byte((i*8 + j) % 256)
		}
		for j := range withdrawalsRoot {
			withdrawalsRoot[j] = byte((i*9 + j) % 256)
		}

		// Variable length extra data
		extraDataLen := i % 33 // 0 to 32 bytes
		extraData := make([]byte, extraDataLen)
		for j := range extraData {
			extraData[j] = byte((i + j) % 256)
		}

		baseFee := math.NewU256(uint64(i) * 1000000)

		current := &types.ExecutionPayloadHeader{
			Versionable:      types.NewVersionable(version.Deneb()),
			ParentHash:       parentHash,
			FeeRecipient:     feeRecipient,
			StateRoot:        stateRoot,
			ReceiptsRoot:     receiptsRoot,
			LogsBloom:        logsBloom,
			Random:           random,
			Number:           math.U64(i * 12345),
			GasLimit:         math.U64(30000000 + i),
			GasUsed:          math.U64(21000000 + i),
			Timestamp:        math.U64(1700000000 + i),
			ExtraData:        extraData,
			BaseFeePerGas:    baseFee,
			BlockHash:        blockHash,
			TransactionsRoot: transactionsRoot,
			WithdrawalsRoot:  withdrawalsRoot,
			BlobGasUsed:      math.U64(i * 131072),
			ExcessBlobGas:    math.U64(i * 262144),
		}

		karalabe := &ExecutionPayloadHeaderKaralabe{
			Versionable:      NewVersionableKaralabeHeader(version.Deneb()),
			ParentHash:       parentHash,
			FeeRecipient:     feeRecipient,
			StateRoot:        stateRoot,
			ReceiptsRoot:     receiptsRoot,
			LogsBloom:        logsBloom,
			Random:           random,
			Number:           math.U64(i * 12345),
			GasLimit:         math.U64(30000000 + i),
			GasUsed:          math.U64(21000000 + i),
			Timestamp:        math.U64(1700000000 + i),
			ExtraData:        extraData,
			BaseFeePerGas:    baseFee,
			BlockHash:        blockHash,
			TransactionsRoot: transactionsRoot,
			WithdrawalsRoot:  withdrawalsRoot,
			BlobGasUsed:      math.U64(i * 131072),
			ExcessBlobGas:    math.U64(i * 262144),
		}

		// Compare marshaling
		currentBytes, err1 := current.MarshalSSZ()
		require.NoError(t, err1)
		karalableBytes, err2 := karalabe.MarshalSSZ()
		require.NoError(t, err2)
		require.Equal(t, karalableBytes, currentBytes, "fuzzing iteration %d: marshaled bytes should be identical", i)

		// Compare roots
		currentRoot, err := current.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, [32]byte(karalabe.HashTreeRoot()), currentRoot, "fuzzing iteration %d: roots should be identical", i)
	}
}

// TestExecutionPayloadHeaderCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestExecutionPayloadHeaderCompatibilityRoundTrip(t *testing.T) {
	// Create a header with specific values
	var parentHash common.ExecutionHash
	for i := range parentHash {
		parentHash[i] = byte(i + 10)
	}

	var feeRecipient common.ExecutionAddress
	for i := range feeRecipient {
		feeRecipient[i] = byte(i + 20)
	}

	var withdrawalsRoot common.Root
	for i := range withdrawalsRoot {
		withdrawalsRoot[i] = byte(i + 30)
	}

	original := &types.ExecutionPayloadHeader{
		Versionable:      types.NewVersionable(version.Deneb()),
		ParentHash:       parentHash,
		FeeRecipient:     feeRecipient,
		StateRoot:        common.Bytes32{40, 41, 42, 43, 44, 45},
		ReceiptsRoot:     common.Bytes32{50, 51, 52, 53, 54, 55},
		LogsBloom:        bytes.B256{60, 61, 62, 63, 64, 65},
		Random:           common.Bytes32{70, 71, 72, 73, 74, 75},
		Number:           999999,
		GasLimit:         30000000,
		GasUsed:          21000000,
		Timestamp:        1700000000,
		ExtraData:        []byte("round trip test"),
		BaseFeePerGas:    math.NewU256(1234567890),
		BlockHash:        common.ExecutionHash{80, 81, 82, 83, 84, 85},
		TransactionsRoot: common.Root{90, 91, 92, 93, 94, 95},
		WithdrawalsRoot:  withdrawalsRoot,
		BlobGasUsed:      131072,
		ExcessBlobGas:    262144,
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := NewEmptyExecutionPayloadHeaderWithVersionKaralabe(version.Deneb())
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := types.NewEmptyExecutionPayloadHeaderWithVersion(version.Deneb())
	err = roundTrip.UnmarshalSSZ(karalableBytes)
	require.NoError(t, err)

	// Verify round trip preserved all data
	require.Equal(t, original, roundTrip, "round trip should preserve all data")

	// Verify both serializations are identical
	require.Equal(t, currentBytes, karalableBytes, "both serializations should be identical")

	// Verify hash roots match throughout
	originalRoot, err := original.HashTreeRoot()
	require.NoError(t, err)
	roundTripRoot, err := roundTrip.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, originalRoot, [32]byte(karalabe.HashTreeRoot()), "hash roots should match")
	require.Equal(t, originalRoot, roundTripRoot, "hash roots should match after round trip")
}

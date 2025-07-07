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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// Eth1DataSizeKaralabe is the size of the Eth1Data object in bytes.
// 32 bytes for DepositRoot + 8 bytes for DepositCount + 32 bytes for BlockHash.
const Eth1DataSizeKaralabe = 72

// Compile-time assertions to ensure Eth1DataKaralabe implements necessary interfaces.
var _ ssz.StaticObject = (*Eth1DataKaralabe)(nil)

// Eth1DataKaralabe is an exact copy of Eth1Data from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// This type uses karalabe/ssz for SSZ operations to ensure compatibility testing.
type Eth1DataKaralabe struct {
	// DepositRoot is the root of the deposit tree.
	DepositRoot common.Root `json:"depositRoot"`
	// DepositCount is the number of deposits in the deposit tree.
	DepositCount math.U64 `json:"depositCount"`
	// BlockHash is the hash of the block corresponding to the Eth1Data.
	BlockHash common.ExecutionHash `json:"blockHash"`
}

// SizeSSZ returns the size of the Eth1Data object in SSZ encoding.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (*Eth1DataKaralabe) SizeSSZ() uint32 {
	return Eth1DataSizeKaralabe
}

// DefineSSZ defines the SSZ encoding for the Eth1Data object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (e *Eth1DataKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &e.DepositRoot)
	ssz.DefineUint64(codec, &e.DepositCount)
	ssz.DefineStaticBytes(codec, &e.BlockHash)
}

// HashTreeRoot computes the SSZ hash tree root of the Eth1Data object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (e *Eth1DataKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(e)
}

// MarshalSSZ marshals the Eth1Data object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (e *Eth1DataKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(e))
	return buf, ssz.EncodeToBytes(buf, e)
}

func (*Eth1DataKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// UnmarshalSSZ unmarshals the Eth1Data object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (e *Eth1DataKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, e)
}

// TestEth1DataCompatibility tests that the current Eth1Data implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestEth1DataCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.Eth1Data, *Eth1DataKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.Eth1Data, *Eth1DataKaralabe) {
				return &types.Eth1Data{}, &Eth1DataKaralabe{}
			},
		},
		{
			name: "typical eth1 data",
			setup: func() (*types.Eth1Data, *Eth1DataKaralabe) {
				depositRoot := common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
				depositCount := math.U64(1000)
				blockHash := common.ExecutionHash{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

				current := &types.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				}
				karalabe := &Eth1DataKaralabe{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.Eth1Data, *Eth1DataKaralabe) {
				var depositRoot common.Root
				var blockHash common.ExecutionHash
				
				// Fill with max values
				for i := range depositRoot {
					depositRoot[i] = 0xFF
				}
				for i := range blockHash {
					blockHash[i] = 0xFF
				}
				
				depositCount := math.U64(^uint64(0)) // max uint64

				current := &types.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				}
				karalabe := &Eth1DataKaralabe{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				}
				return current, karalabe
			},
		},
		{
			name: "specific pattern",
			setup: func() (*types.Eth1Data, *Eth1DataKaralabe) {
				// Create a specific pattern for testing
				var depositRoot common.Root
				var blockHash common.ExecutionHash
				
				// Alternating pattern
				for i := range depositRoot {
					if i%2 == 0 {
						depositRoot[i] = 0xAA
					} else {
						depositRoot[i] = 0x55
					}
				}
				
				// Sequential pattern
				for i := range blockHash {
					blockHash[i] = byte(i)
				}
				
				depositCount := math.U64(123456789)

				current := &types.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				}
				karalabe := &Eth1DataKaralabe{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				}
				return current, karalabe
			},
		},
		{
			name: "genesis eth1 data",
			setup: func() (*types.Eth1Data, *Eth1DataKaralabe) {
				// Simulate genesis eth1 data
				depositRoot := common.Root{} // zero root
				depositCount := math.U64(0)  // no deposits
				blockHash := common.ExecutionHash{} // zero hash

				current := &types.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				}
				karalabe := &Eth1DataKaralabe{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
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
			require.Equal(t, int(Eth1DataSizeKaralabe), current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(Eth1DataSizeKaralabe), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.Eth1Data{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &Eth1DataKaralabe{}
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")
			require.Equal(t, karalabe, newKaralabe, "unmarshaled karalabe should match original")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabelRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabelRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

// TestEth1DataCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestEth1DataCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid eth1 data
		var depositRoot common.Root
		var blockHash common.ExecutionHash
		
		// Use deterministic "random" data based on iteration
		for j := range depositRoot {
			depositRoot[j] = byte((i + j) % 256)
		}
		for j := range blockHash {
			blockHash[j] = byte((i * 2 + j) % 256)
		}
		
		depositCount := math.U64(uint64(i) * 999)

		current := &types.Eth1Data{
			DepositRoot:  depositRoot,
			DepositCount: depositCount,
			BlockHash:    blockHash,
		}
		karalabe := &Eth1DataKaralabe{
			DepositRoot:  depositRoot,
			DepositCount: depositCount,
			BlockHash:    blockHash,
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

// TestEth1DataCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestEth1DataCompatibilityInvalidData(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "insufficient data - missing block hash",
			data: make([]byte, 40), // Only deposit root (32) + count (8), missing block hash
		},
		{
			name: "insufficient data - partial block hash", 
			data: make([]byte, 60), // deposit root (32) + count (8) + partial block hash (20)
		},
		{
			name: "excess data",
			data: make([]byte, 100), // more than required 72 bytes
		},
		{
			name: "exact size but all zeros",
			data: make([]byte, 72), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.Eth1Data{}
			currentErr := current.UnmarshalSSZ(tc.data)
			
			// Test unmarshal with karalabe implementation
			karalabe := &Eth1DataKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)
			
			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				require.Equal(t, current.DepositRoot, karalabe.DepositRoot, "deposit roots should match")
				require.Equal(t, current.DepositCount, karalabe.DepositCount, "deposit counts should match")
				require.Equal(t, current.BlockHash, karalabe.BlockHash, "block hashes should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestEth1DataCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestEth1DataCompatibilityRoundTrip(t *testing.T) {
	// Create eth1 data with specific values
	original := &types.Eth1Data{
		DepositRoot:  common.Root{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210, 220, 230, 240, 250, 1, 2, 3, 4, 5, 6, 7},
		DepositCount: math.U64(987654321),
		BlockHash:    common.ExecutionHash{7, 6, 5, 4, 3, 2, 1, 250, 240, 230, 220, 210, 200, 190, 180, 170, 160, 150, 140, 130, 120, 110, 100, 90, 80, 70, 60, 50, 40, 30, 20, 10},
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &Eth1DataKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.Eth1Data{}
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

// TestEth1DataCompatibilityEndianness verifies that deposit count is encoded in little-endian
func TestEth1DataCompatibilityEndianness(t *testing.T) {
	// Create eth1 data with a specific deposit count that shows endianness
	depositCount := math.U64(0x0102030405060708) // This will show byte order clearly
	
	current := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: depositCount,
		BlockHash:    common.ExecutionHash{},
	}
	
	karalabe := &Eth1DataKaralabe{
		DepositRoot:  common.Root{},
		DepositCount: depositCount, 
		BlockHash:    common.ExecutionHash{},
	}

	// Marshal both
	currentBytes, err := current.MarshalSSZ()
	require.NoError(t, err)
	
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)
	
	// Verify they're identical
	require.Equal(t, karalableBytes, currentBytes, "endianness encoding should be identical")
	
	// Verify the deposit count is encoded in little-endian at offset 32
	// Offset 32 = after DepositRoot (32 bytes)
	expectedLE := []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}
	require.Equal(t, expectedLE, currentBytes[32:40], "deposit count should be little-endian in current")
	require.Equal(t, expectedLE, karalableBytes[32:40], "deposit count should be little-endian in karalabe")
}
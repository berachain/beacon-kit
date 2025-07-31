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
	karalabe "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// sszPendingPartialWithdrawalSizeKaralabe defines the total SSZ serialized size for
// PendingPartialWithdrawal. The fields are assumed to be encoded as follows:
// - ValidatorIndex: 8 bytes (uint64)
// - Amount:         8 bytes (math.Gwei)
// - WithdrawableEpoch: 8 bytes (uint64)
// Total = 8 + 8 + 8 = 24 bytes.
const sszPendingPartialWithdrawalSizeKaralabe = 24

// Compile-time check to ensure PendingPartialWithdrawalKaralabe implements the necessary interfaces.
var _ karalabe.StaticObject = (*PendingPartialWithdrawalKaralabe)(nil)

// PendingPartialWithdrawalKaralabe is an exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// PendingPartialWithdrawal reflects the following spec:
//
//	class PendingPartialWithdrawal(Container):
//	    validator_index: ValidatorIndex
//	    amount: Gwei
//	    withdrawable_epoch: Epoch
type PendingPartialWithdrawalKaralabe struct {
	ValidatorIndex    math.ValidatorIndex
	Amount            math.Gwei
	WithdrawableEpoch math.Epoch
}

// ValidateAfterDecodingSSZ validates the PendingPartialWithdrawal object
// after decoding from SSZ. Customize further validation as needed.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (p *PendingPartialWithdrawalKaralabe) ValidateAfterDecodingSSZ() error {
	return nil
}

// DefineSSZ registers the SSZ encoding for each field in PendingPartialWithdrawal.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (p *PendingPartialWithdrawalKaralabe) DefineSSZ(codec *karalabe.Codec) {
	karalabe.DefineUint64(codec, &p.ValidatorIndex)
	karalabe.DefineUint64(codec, &p.Amount)
	karalabe.DefineUint64(codec, &p.WithdrawableEpoch)
}

// SizeSSZ returns the fixed size of the SSZ serialization for PendingPartialWithdrawal.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (p *PendingPartialWithdrawalKaralabe) SizeSSZ() uint32 {
	return sszPendingPartialWithdrawalSizeKaralabe
}

// MarshalSSZ returns the SSZ encoding of the PendingPartialWithdrawal.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (p *PendingPartialWithdrawalKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabe.Size(p))
	return buf, karalabe.EncodeToBytes(buf, p)
}

// HashTreeRoot computes and returns the hash tree root for the PendingPartialWithdrawal.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (p *PendingPartialWithdrawalKaralabe) HashTreeRoot() common.Root {
	return karalabe.HashSequential(p)
}

// UnmarshalSSZ unmarshals the PendingPartialWithdrawal object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use karalabe.DecodeFromBytes
func (p *PendingPartialWithdrawalKaralabe) UnmarshalSSZ(buf []byte) error {
	return karalabe.DecodeFromBytes(buf, p)
}

// TestPendingPartialWithdrawalCompatibility tests that the current PendingPartialWithdrawal implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestPendingPartialWithdrawalCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.PendingPartialWithdrawal, *PendingPartialWithdrawalKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.PendingPartialWithdrawal, *PendingPartialWithdrawalKaralabe) {
				return &types.PendingPartialWithdrawal{}, &PendingPartialWithdrawalKaralabe{}
			},
		},
		{
			name: "typical withdrawal",
			setup: func() (*types.PendingPartialWithdrawal, *PendingPartialWithdrawalKaralabe) {
				validatorIndex := math.ValidatorIndex(12345)
				amount := math.Gwei(1000000000) // 1 ETH
				withdrawableEpoch := math.Epoch(100)

				current := &types.PendingPartialWithdrawal{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
				}
				karalabe := &PendingPartialWithdrawalKaralabe{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
				}
				return current, karalabe
			},
		},
		{
			name: "large withdrawal",
			setup: func() (*types.PendingPartialWithdrawal, *PendingPartialWithdrawalKaralabe) {
				validatorIndex := math.ValidatorIndex(999999)
				amount := math.Gwei(32000000000) // 32 ETH
				withdrawableEpoch := math.Epoch(50000)

				current := &types.PendingPartialWithdrawal{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
				}
				karalabe := &PendingPartialWithdrawalKaralabe{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.PendingPartialWithdrawal, *PendingPartialWithdrawalKaralabe) {
				validatorIndex := math.ValidatorIndex(^uint64(0))
				amount := math.Gwei(^uint64(0))
				withdrawableEpoch := math.Epoch(^uint64(0))

				current := &types.PendingPartialWithdrawal{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
				}
				karalabe := &PendingPartialWithdrawalKaralabe{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
				}
				return current, karalabe
			},
		},
		{
			name: "minimum non-zero values",
			setup: func() (*types.PendingPartialWithdrawal, *PendingPartialWithdrawalKaralabe) {
				validatorIndex := math.ValidatorIndex(1)
				amount := math.Gwei(1) // 1 gwei
				withdrawableEpoch := math.Epoch(1)

				current := &types.PendingPartialWithdrawal{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
				}
				karalabe := &PendingPartialWithdrawalKaralabe{
					ValidatorIndex:    validatorIndex,
					Amount:            amount,
					WithdrawableEpoch: withdrawableEpoch,
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
			require.Equal(t, int(sszPendingPartialWithdrawalSizeKaralabe), current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(sszPendingPartialWithdrawalSizeKaralabe), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.PendingPartialWithdrawal{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &PendingPartialWithdrawalKaralabe{}
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

// TestPendingPartialWithdrawalCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestPendingPartialWithdrawalCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid withdrawal data
		validatorIndex := math.ValidatorIndex(uint64(i) * 137)
		amount := math.Gwei(uint64(i) * 1000000000)
		withdrawableEpoch := math.Epoch(uint64(i) * 32)

		current := &types.PendingPartialWithdrawal{
			ValidatorIndex:    validatorIndex,
			Amount:            amount,
			WithdrawableEpoch: withdrawableEpoch,
		}
		karalabe := &PendingPartialWithdrawalKaralabe{
			ValidatorIndex:    validatorIndex,
			Amount:            amount,
			WithdrawableEpoch: withdrawableEpoch,
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

// TestPendingPartialWithdrawalCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestPendingPartialWithdrawalCompatibilityInvalidData(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "insufficient data",
			data: make([]byte, 10), // less than required 24 bytes
		},
		{
			name: "excess data",
			data: make([]byte, 50), // more than required 24 bytes
		},
		{
			name: "exact size but all zeros",
			data: make([]byte, 24), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.PendingPartialWithdrawal{}
			currentErr := current.UnmarshalSSZ(tc.data)

			// Test unmarshal with karalabe implementation
			karalabe := &PendingPartialWithdrawalKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)

			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				require.Equal(t, current.ValidatorIndex, karalabe.ValidatorIndex, "validator indices should match")
				require.Equal(t, current.Amount, karalabe.Amount, "amounts should match")
				require.Equal(t, current.WithdrawableEpoch, karalabe.WithdrawableEpoch, "withdrawable epochs should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestPendingPartialWithdrawalCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestPendingPartialWithdrawalCompatibilityRoundTrip(t *testing.T) {
	// Create a withdrawal with specific values
	original := &types.PendingPartialWithdrawal{
		ValidatorIndex:    math.ValidatorIndex(54321),
		Amount:            math.Gwei(5000000000), // 5 ETH
		WithdrawableEpoch: math.Epoch(2048),
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &PendingPartialWithdrawalKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.PendingPartialWithdrawal{}
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

// TestPendingPartialWithdrawalCompatibilityEndianness verifies that all fields are encoded in little-endian
func TestPendingPartialWithdrawalCompatibilityEndianness(t *testing.T) {
	// Create withdrawal with specific values that show endianness
	current := &types.PendingPartialWithdrawal{
		ValidatorIndex:    math.ValidatorIndex(0x0102030405060708),
		Amount:            math.Gwei(0x090A0B0C0D0E0F10),
		WithdrawableEpoch: math.Epoch(0x1112131415161718),
	}

	karalabe := &PendingPartialWithdrawalKaralabe{
		ValidatorIndex:    math.ValidatorIndex(0x0102030405060708),
		Amount:            math.Gwei(0x090A0B0C0D0E0F10),
		WithdrawableEpoch: math.Epoch(0x1112131415161718),
	}

	// Marshal both
	currentBytes, err := current.MarshalSSZ()
	require.NoError(t, err)

	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Verify they're identical
	require.Equal(t, karalableBytes, currentBytes, "endianness encoding should be identical")

	// Verify the fields are encoded in little-endian at correct offsets
	// ValidatorIndex at offset 0
	expectedValidatorIndex := []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}
	require.Equal(t, expectedValidatorIndex, currentBytes[0:8], "validator index should be little-endian")

	// Amount at offset 8
	expectedAmount := []byte{0x10, 0x0F, 0x0E, 0x0D, 0x0C, 0x0B, 0x0A, 0x09}
	require.Equal(t, expectedAmount, currentBytes[8:16], "amount should be little-endian")

	// WithdrawableEpoch at offset 16
	expectedEpoch := []byte{0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11}
	require.Equal(t, expectedEpoch, currentBytes[16:24], "withdrawable epoch should be little-endian")
}

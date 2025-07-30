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
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// ValidatorSizeKaralabe is the size of the Validator struct in bytes.
const ValidatorSizeKaralabe = 121

// Compile-time assertions to ensure ValidatorKaralabe implements necessary interfaces.
var _ ssz.StaticObject = (*ValidatorKaralabe)(nil)

// ValidatorKaralabe is an exact copy of Validator from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// Validator as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#validator
type ValidatorKaralabe struct {
	// Pubkey is the validator's 48-byte BLS public key.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// WithdrawalCredentials are an address that controls the validator.
	WithdrawalCredentials types.WithdrawalCredentials `json:"withdrawalCredentials"`
	// EffectiveBalance is the validator's current effective balance in gwei.
	EffectiveBalance math.Gwei `json:"effectiveBalance"`
	// Slashed indicates whether the validator has been slashed.
	Slashed bool `json:"slashed"`
	// ActivationEligibilityEpoch is the epoch in which the validator became
	// eligible for activation.
	ActivationEligibilityEpoch math.Epoch `json:"activationEligibilityEpoch"`
	// ActivationEpoch is the epoch in which the validator activated.
	ActivationEpoch math.Epoch `json:"activationEpoch"`
	// ExitEpoch is the epoch in which the validator exited.
	ExitEpoch math.Epoch `json:"exitEpoch"`
	// WithdrawableEpoch is the epoch in which the validator can withdraw.
	WithdrawableEpoch math.Epoch `json:"withdrawableEpoch"`
}

// SizeSSZ returns the size of the Validator object in SSZ encoding.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (*ValidatorKaralabe) SizeSSZ() uint32 {
	return ValidatorSizeKaralabe
}

// DefineSSZ defines the SSZ encoding for the Validator object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (v *ValidatorKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &v.Pubkey)
	ssz.DefineStaticBytes(codec, &v.WithdrawalCredentials)
	ssz.DefineUint64(codec, &v.EffectiveBalance)
	ssz.DefineBool(codec, &v.Slashed)
	ssz.DefineUint64(codec, &v.ActivationEligibilityEpoch)
	ssz.DefineUint64(codec, &v.ActivationEpoch)
	ssz.DefineUint64(codec, &v.ExitEpoch)
	ssz.DefineUint64(codec, &v.WithdrawableEpoch)
}

// HashTreeRoot computes the SSZ hash tree root of the Validator object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (v *ValidatorKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(v)
}

// MarshalSSZ marshals the Validator object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (v *ValidatorKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(v))
	return buf, ssz.EncodeToBytes(buf, v)
}

func (*ValidatorKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// UnmarshalSSZ unmarshals the Validator object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (v *ValidatorKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, v)
}

// TestValidatorCompatibility tests that the current Validator implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestValidatorCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.Validator, *ValidatorKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.Validator, *ValidatorKaralabe) {
				return &types.Validator{}, &ValidatorKaralabe{}
			},
		},
		{
			name: "typical active validator",
			setup: func() (*types.Validator, *ValidatorKaralabe) {
				pubkey := crypto.BLSPubkey{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				creds := types.WithdrawalCredentials{0x01} // ETH1_ADDRESS_WITHDRAWAL_PREFIX
				for i := 1; i < 12; i++ {
					creds[i] = 0x00 // padding
				}
				// Set example address bytes
				for i := 12; i < 32; i++ {
					creds[i] = byte(i)
				}
				effectiveBalance := math.Gwei(32000000000) // 32 ETH
				slashed := false
				activationEligibilityEpoch := math.Epoch(100)
				activationEpoch := math.Epoch(105)
				exitEpoch := math.Epoch(constants.FarFutureEpoch)
				withdrawableEpoch := math.Epoch(constants.FarFutureEpoch)

				current := &types.Validator{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
				}
				karalabe := &ValidatorKaralabe{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
				}
				return current, karalabe
			},
		},
		{
			name: "slashed validator",
			setup: func() (*types.Validator, *ValidatorKaralabe) {
				pubkey := crypto.BLSPubkey{100, 101, 102, 103, 104, 105}
				creds := types.WithdrawalCredentials{0x00} // BLS_WITHDRAWAL_PREFIX
				// Fill rest with example BLS pubkey hash
				for i := 1; i < 32; i++ {
					creds[i] = byte(i * 2)
				}
				effectiveBalance := math.Gwei(31000000000) // 31 ETH (reduced due to slashing)
				slashed := true
				activationEligibilityEpoch := math.Epoch(50)
				activationEpoch := math.Epoch(55)
				exitEpoch := math.Epoch(200)               // forced exit due to slashing
				withdrawableEpoch := math.Epoch(200 + 256) // slashed validators have longer withdrawal delay

				current := &types.Validator{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
				}
				karalabe := &ValidatorKaralabe{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
				}
				return current, karalabe
			},
		},
		{
			name: "exited validator",
			setup: func() (*types.Validator, *ValidatorKaralabe) {
				pubkey := crypto.BLSPubkey{200, 201, 202, 203, 204, 205}
				creds := types.WithdrawalCredentials{0x01}
				for i := 1; i < 32; i++ {
					creds[i] = byte(255 - i) // different pattern
				}
				effectiveBalance := math.Gwei(16000000000) // 16 ETH (partial withdrawal)
				slashed := false
				activationEligibilityEpoch := math.Epoch(10)
				activationEpoch := math.Epoch(15)
				exitEpoch := math.Epoch(1000)
				withdrawableEpoch := math.Epoch(1000 + 256)

				current := &types.Validator{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
				}
				karalabe := &ValidatorKaralabe{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.Validator, *ValidatorKaralabe) {
				var pubkey crypto.BLSPubkey
				var creds types.WithdrawalCredentials

				// Fill with max values
				for i := range pubkey {
					pubkey[i] = 0xFF
				}
				for i := range creds {
					creds[i] = 0xFF
				}

				effectiveBalance := math.Gwei(^uint64(0)) // max uint64
				slashed := true
				activationEligibilityEpoch := math.Epoch(^uint64(0)) // max uint64
				activationEpoch := math.Epoch(^uint64(0))            // max uint64
				exitEpoch := math.Epoch(^uint64(0))                  // max uint64
				withdrawableEpoch := math.Epoch(^uint64(0))          // max uint64

				current := &types.Validator{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
				}
				karalabe := &ValidatorKaralabe{
					Pubkey:                     pubkey,
					WithdrawalCredentials:      creds,
					EffectiveBalance:           effectiveBalance,
					Slashed:                    slashed,
					ActivationEligibilityEpoch: activationEligibilityEpoch,
					ActivationEpoch:            activationEpoch,
					ExitEpoch:                  exitEpoch,
					WithdrawableEpoch:          withdrawableEpoch,
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
			require.Equal(t, int(ValidatorSizeKaralabe), current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(ValidatorSizeKaralabe), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.Validator{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &ValidatorKaralabe{}
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

// TestValidatorCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestValidatorCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid validator data
		var pubkey crypto.BLSPubkey
		var creds types.WithdrawalCredentials

		// Use deterministic "random" data based on iteration
		for j := range pubkey {
			pubkey[j] = byte((i + j) % 256)
		}
		for j := range creds {
			creds[j] = byte((i*2 + j) % 256)
		}

		effectiveBalance := math.Gwei(uint64(i) * 1000000000)
		slashed := i%2 == 0
		activationEligibilityEpoch := math.Epoch(uint64(i))
		activationEpoch := math.Epoch(uint64(i + 5))
		exitEpoch := math.Epoch(uint64(i + 1000))
		withdrawableEpoch := math.Epoch(uint64(i + 1256))

		current := &types.Validator{
			Pubkey:                     pubkey,
			WithdrawalCredentials:      creds,
			EffectiveBalance:           effectiveBalance,
			Slashed:                    slashed,
			ActivationEligibilityEpoch: activationEligibilityEpoch,
			ActivationEpoch:            activationEpoch,
			ExitEpoch:                  exitEpoch,
			WithdrawableEpoch:          withdrawableEpoch,
		}
		karalabe := &ValidatorKaralabe{
			Pubkey:                     pubkey,
			WithdrawalCredentials:      creds,
			EffectiveBalance:           effectiveBalance,
			Slashed:                    slashed,
			ActivationEligibilityEpoch: activationEligibilityEpoch,
			ActivationEpoch:            activationEpoch,
			ExitEpoch:                  exitEpoch,
			WithdrawableEpoch:          withdrawableEpoch,
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

// TestValidatorCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestValidatorCompatibilityInvalidData(t *testing.T) {
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
			data: make([]byte, 100), // less than required 121 bytes
		},
		{
			name: "excess data",
			data: make([]byte, 200), // more than required 121 bytes
		},
		{
			name: "exact size but all zeros",
			data: make([]byte, 121), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.Validator{}
			currentErr := current.UnmarshalSSZ(tc.data)

			// Test unmarshal with karalabe implementation
			karalabe := &ValidatorKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)

			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				require.Equal(t, current.Pubkey, karalabe.Pubkey, "pubkeys should match")
				require.Equal(t, current.WithdrawalCredentials, karalabe.WithdrawalCredentials, "withdrawal credentials should match")
				require.Equal(t, current.EffectiveBalance, karalabe.EffectiveBalance, "effective balances should match")
				require.Equal(t, current.Slashed, karalabe.Slashed, "slashed status should match")
				require.Equal(t, current.ActivationEligibilityEpoch, karalabe.ActivationEligibilityEpoch, "activation eligibility epochs should match")
				require.Equal(t, current.ActivationEpoch, karalabe.ActivationEpoch, "activation epochs should match")
				require.Equal(t, current.ExitEpoch, karalabe.ExitEpoch, "exit epochs should match")
				require.Equal(t, current.WithdrawableEpoch, karalabe.WithdrawableEpoch, "withdrawable epochs should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestValidatorCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestValidatorCompatibilityRoundTrip(t *testing.T) {
	// Create a validator with specific values
	original := &types.Validator{
		Pubkey:                     crypto.BLSPubkey{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		WithdrawalCredentials:      types.WithdrawalCredentials{0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
		EffectiveBalance:           math.Gwei(32000000000),
		Slashed:                    false,
		ActivationEligibilityEpoch: math.Epoch(1000),
		ActivationEpoch:            math.Epoch(1005),
		ExitEpoch:                  math.Epoch(constants.FarFutureEpoch),
		WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &ValidatorKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.Validator{}
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

// TestValidatorCompatibilityBooleanEncoding verifies that boolean slashed field is encoded correctly
func TestValidatorCompatibilityBooleanEncoding(t *testing.T) {
	// Test false encoding
	notSlashed := &types.Validator{
		Slashed: false,
		// Other fields can be zero
	}
	notSlashedKaralabe := &ValidatorKaralabe{
		Slashed: false,
	}

	notSlashedBytes, err := notSlashed.MarshalSSZ()
	require.NoError(t, err)
	notSlashedKaralableBytes, err := notSlashedKaralabe.MarshalSSZ()
	require.NoError(t, err)

	// Verify they're identical
	require.Equal(t, notSlashedKaralableBytes, notSlashedBytes, "boolean false encoding should be identical")

	// Verify the boolean is encoded as 0x00 at offset 88
	// Offset 88 = after Pubkey (48) + WithdrawalCredentials (32) + EffectiveBalance (8)
	require.Equal(t, byte(0x00), notSlashedBytes[88], "slashed=false should be encoded as 0x00")

	// Test true encoding
	slashed := &types.Validator{
		Slashed: true,
		// Other fields can be zero
	}
	slashedKaralabe := &ValidatorKaralabe{
		Slashed: true,
	}

	slashedBytes, err := slashed.MarshalSSZ()
	require.NoError(t, err)
	slashedKaralableBytes, err := slashedKaralabe.MarshalSSZ()
	require.NoError(t, err)

	// Verify they're identical
	require.Equal(t, slashedKaralableBytes, slashedBytes, "boolean true encoding should be identical")

	// Verify the boolean is encoded as 0x01 at offset 88
	require.Equal(t, byte(0x01), slashedBytes[88], "slashed=true should be encoded as 0x01")
}

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
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	karalabe "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// Compile-time assertions to ensure ExecutionRequestsKaralabe implements necessary interfaces.
var _ karalabe.DynamicObject = (*ExecutionRequestsKaralabe)(nil)

// ExecutionRequestsKaralabe is a karalabe/ssz version of ExecutionRequests for compatibility testing.
type ExecutionRequestsKaralabe struct {
	Deposits       []*DepositKaralabe
	Withdrawals    []*WithdrawalRequestKaralabe
	Consolidations []*ConsolidationRequestKaralabe
}

// WithdrawalRequestKaralabe is a karalabe/ssz version of WithdrawalRequest.
type WithdrawalRequestKaralabe struct {
	SourceAddress   common.ExecutionAddress
	ValidatorPubKey crypto.BLSPubkey
	Amount          math.Gwei
}

// DefineSSZ defines the SSZ encoding for WithdrawalRequestKaralabe.
func (w *WithdrawalRequestKaralabe) DefineSSZ(c *karalabe.Codec) {
	karalabe.DefineStaticBytes(c, &w.SourceAddress)
	karalabe.DefineStaticBytes(c, &w.ValidatorPubKey)
	karalabe.DefineUint64(c, &w.Amount)
}

// MarshalSSZ marshals the WithdrawalRequestKaralabe to SSZ format.
func (w *WithdrawalRequestKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 76) // 20 + 48 + 8
	return buf, karalabe.EncodeToBytes(buf, w)
}

// UnmarshalSSZ unmarshals the WithdrawalRequestKaralabe from SSZ format.
func (w *WithdrawalRequestKaralabe) UnmarshalSSZ(buf []byte) error {
	return karalabe.DecodeFromBytes(buf, w)
}

// SizeSSZ returns the size of the WithdrawalRequestKaralabe in SSZ format.
func (w *WithdrawalRequestKaralabe) SizeSSZ() uint32 {
	return 76
}

// HashTreeRoot computes the hash tree root of the WithdrawalRequestKaralabe.
func (w *WithdrawalRequestKaralabe) HashTreeRoot() common.Root {
	return karalabe.HashSequential(w)
}

func (w *WithdrawalRequestKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// SizeSSZ returns the SSZ encoded size in bytes for the ExecutionRequestsKaralabe.
func (e *ExecutionRequestsKaralabe) SizeSSZ(fixed bool) uint32 {
	size := uint32(12) // 3 fields * 4 bytes offset each

	if fixed {
		return size
	}

	// Add dynamic sizes
	size += uint32(len(e.Deposits) * 192)       // Each deposit is 192 bytes
	size += uint32(len(e.Withdrawals) * 76)     // Each withdrawal is 76 bytes
	size += uint32(len(e.Consolidations) * 116) // Each consolidation is 116 bytes

	return size
}

// DefineSSZ defines the SSZ encoding for the ExecutionRequestsKaralabe object.
func (e *ExecutionRequestsKaralabe) DefineSSZ(codec *karalabe.Codec) {
	// Define offsets for dynamic fields
	karalabe.DefineSliceOfStaticObjectsOffset(codec, &e.Deposits, 8192)    // MaxDepositRequestsPerPayload
	karalabe.DefineSliceOfStaticObjectsOffset(codec, &e.Withdrawals, 16)   // MaxWithdrawalRequestsPerPayload
	karalabe.DefineSliceOfStaticObjectsOffset(codec, &e.Consolidations, 2) // MaxConsolidationRequestsPerPayload

	// Define content for dynamic fields
	karalabe.DefineSliceOfStaticObjectsContent(codec, &e.Deposits, 8192)
	karalabe.DefineSliceOfStaticObjectsContent(codec, &e.Withdrawals, 16)
	karalabe.DefineSliceOfStaticObjectsContent(codec, &e.Consolidations, 2)
}

// MarshalSSZ marshals the ExecutionRequestsKaralabe to SSZ format.
func (e *ExecutionRequestsKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, e.SizeSSZ(false))
	return buf, karalabe.EncodeToBytes(buf, e)
}

// UnmarshalSSZ unmarshals the ExecutionRequestsKaralabe from SSZ format.
func (e *ExecutionRequestsKaralabe) UnmarshalSSZ(buf []byte) error {
	return karalabe.DecodeFromBytes(buf, e)
}

// HashTreeRoot computes the Merkleization of the ExecutionRequestsKaralabe.
func (e *ExecutionRequestsKaralabe) HashTreeRoot() common.Root {
	return karalabe.HashConcurrent(e)
}

func (e *ExecutionRequestsKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// TestExecutionRequestsCompatibility tests that the current ExecutionRequests implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestExecutionRequestsCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.ExecutionRequests, *ExecutionRequestsKaralabe)
	}{
		{
			name: "all empty slices",
			setup: func() (*types.ExecutionRequests, *ExecutionRequestsKaralabe) {
				return &types.ExecutionRequests{
						Deposits:       []*types.DepositRequest{},
						Withdrawals:    []*types.WithdrawalRequest{},
						Consolidations: []*types.ConsolidationRequest{},
					}, &ExecutionRequestsKaralabe{
						Deposits:       []*DepositKaralabe{},
						Withdrawals:    []*WithdrawalRequestKaralabe{},
						Consolidations: []*ConsolidationRequestKaralabe{},
					}
			},
		},
		{
			name: "single deposit only",
			setup: func() (*types.ExecutionRequests, *ExecutionRequestsKaralabe) {
				pubkey := crypto.BLSPubkey{1, 2, 3, 4, 5, 6, 7, 8}
				creds := types.WithdrawalCredentials{0x01}
				for i := 1; i < 32; i++ {
					creds[i] = byte(i)
				}

				currentDep := &types.Deposit{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      32000000000,
					Signature:   crypto.BLSSignature{9, 10, 11, 12},
					Index:       100,
				}
				karalabeDep := &DepositKaralabe{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      32000000000,
					Signature:   crypto.BLSSignature{9, 10, 11, 12},
					Index:       100,
				}

				return &types.ExecutionRequests{
						Deposits:       []*types.DepositRequest{currentDep},
						Withdrawals:    []*types.WithdrawalRequest{},
						Consolidations: []*types.ConsolidationRequest{},
					}, &ExecutionRequestsKaralabe{
						Deposits:       []*DepositKaralabe{karalabeDep},
						Withdrawals:    []*WithdrawalRequestKaralabe{},
						Consolidations: []*ConsolidationRequestKaralabe{},
					}
			},
		},
		{
			name: "single withdrawal only",
			setup: func() (*types.ExecutionRequests, *ExecutionRequestsKaralabe) {
				addr := common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				pubkey := crypto.BLSPubkey{20, 21, 22, 23, 24}
				amount := math.Gwei(1000000000)

				currentWithdrawal := &types.WithdrawalRequest{
					SourceAddress:   addr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
				}
				karalabeWithdrawal := &WithdrawalRequestKaralabe{
					SourceAddress:   addr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
				}

				return &types.ExecutionRequests{
						Deposits:       []*types.DepositRequest{},
						Withdrawals:    []*types.WithdrawalRequest{currentWithdrawal},
						Consolidations: []*types.ConsolidationRequest{},
					}, &ExecutionRequestsKaralabe{
						Deposits:       []*DepositKaralabe{},
						Withdrawals:    []*WithdrawalRequestKaralabe{karalabeWithdrawal},
						Consolidations: []*ConsolidationRequestKaralabe{},
					}
			},
		},
		{
			name: "single consolidation only",
			setup: func() (*types.ExecutionRequests, *ExecutionRequestsKaralabe) {
				addr := common.ExecutionAddress{30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49}
				srcPubkey := crypto.BLSPubkey{50, 51, 52, 53, 54}
				tgtPubkey := crypto.BLSPubkey{60, 61, 62, 63, 64}

				currentCons := &types.ConsolidationRequest{
					SourceAddress: addr,
					SourcePubKey:  srcPubkey,
					TargetPubKey:  tgtPubkey,
				}
				karalabeCons := &ConsolidationRequestKaralabe{
					SourceAddress: addr,
					SourcePubKey:  srcPubkey,
					TargetPubKey:  tgtPubkey,
				}

				return &types.ExecutionRequests{
						Deposits:       []*types.DepositRequest{},
						Withdrawals:    []*types.WithdrawalRequest{},
						Consolidations: []*types.ConsolidationRequest{currentCons},
					}, &ExecutionRequestsKaralabe{
						Deposits:       []*DepositKaralabe{},
						Withdrawals:    []*WithdrawalRequestKaralabe{},
						Consolidations: []*ConsolidationRequestKaralabe{karalabeCons},
					}
			},
		},
		{
			name: "mixed requests",
			setup: func() (*types.ExecutionRequests, *ExecutionRequestsKaralabe) {
				// Create multiple of each type
				currentDeps := make([]*types.DepositRequest, 2)
				karalabeDeps := make([]*DepositKaralabe, 2)
				for i := 0; i < 2; i++ {
					pubkey := crypto.BLSPubkey{}
					creds := types.WithdrawalCredentials{}
					sig := crypto.BLSSignature{}

					for j := range pubkey {
						pubkey[j] = byte(i*48 + j)
					}
					creds[0] = 0x01
					for j := 1; j < 32; j++ {
						creds[j] = byte(i*32 + j)
					}
					for j := range sig {
						sig[j] = byte(i*96 + j)
					}

					currentDeps[i] = &types.Deposit{
						Pubkey:      pubkey,
						Credentials: creds,
						Amount:      math.Gwei(32000000000 + uint64(i)*1000000000),
						Signature:   sig,
						Index:       uint64(i * 100),
					}
					karalabeDeps[i] = &DepositKaralabe{
						Pubkey:      pubkey,
						Credentials: creds,
						Amount:      math.Gwei(32000000000 + uint64(i)*1000000000),
						Signature:   sig,
						Index:       uint64(i * 100),
					}
				}

				currentWithdrawals := make([]*types.WithdrawalRequest, 3)
				karalabeWithdrawals := make([]*WithdrawalRequestKaralabe, 3)
				for i := 0; i < 3; i++ {
					addr := common.ExecutionAddress{}
					pubkey := crypto.BLSPubkey{}

					for j := range addr {
						addr[j] = byte(i*20 + j)
					}
					for j := range pubkey {
						pubkey[j] = byte(i*48 + j + 100)
					}

					currentWithdrawals[i] = &types.WithdrawalRequest{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          math.Gwei(1000000000 * uint64(i+1)),
					}
					karalabeWithdrawals[i] = &WithdrawalRequestKaralabe{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          math.Gwei(1000000000 * uint64(i+1)),
					}
				}

				currentCons := make([]*types.ConsolidationRequest, 1)
				karalabeCons := make([]*ConsolidationRequestKaralabe, 1)
				addr := common.ExecutionAddress{}
				srcPubkey := crypto.BLSPubkey{}
				tgtPubkey := crypto.BLSPubkey{}
				for j := range addr {
					addr[j] = byte(200 + j)
				}
				for j := range srcPubkey {
					srcPubkey[j] = byte(220 + j)
				}
				for j := range tgtPubkey {
					tgtPubkey[j] = byte(240 + j)
				}
				currentCons[0] = &types.ConsolidationRequest{
					SourceAddress: addr,
					SourcePubKey:  srcPubkey,
					TargetPubKey:  tgtPubkey,
				}
				karalabeCons[0] = &ConsolidationRequestKaralabe{
					SourceAddress: addr,
					SourcePubKey:  srcPubkey,
					TargetPubKey:  tgtPubkey,
				}

				return &types.ExecutionRequests{
						Deposits:       currentDeps,
						Withdrawals:    currentWithdrawals,
						Consolidations: currentCons,
					}, &ExecutionRequestsKaralabe{
						Deposits:       karalabeDeps,
						Withdrawals:    karalabeWithdrawals,
						Consolidations: karalabeCons,
					}
			},
		},
		{
			name: "maximum elements",
			setup: func() (*types.ExecutionRequests, *ExecutionRequestsKaralabe) {
				// Test with maximum allowed elements (using smaller numbers for testing)
				maxDeps := 5 // Using smaller number for test
				maxWithdrawals := 4
				maxCons := 2

				currentDeps := make([]*types.DepositRequest, maxDeps)
				karalabeDeps := make([]*DepositKaralabe, maxDeps)
				for i := 0; i < maxDeps; i++ {
					pubkey := crypto.BLSPubkey{}
					pubkey[0] = byte(i)

					currentDeps[i] = &types.Deposit{
						Pubkey:      pubkey,
						Credentials: types.WithdrawalCredentials{0x01},
						Amount:      32000000000,
						Signature:   crypto.BLSSignature{byte(i)},
						Index:       uint64(i),
					}
					karalabeDeps[i] = &DepositKaralabe{
						Pubkey:      pubkey,
						Credentials: types.WithdrawalCredentials{0x01},
						Amount:      32000000000,
						Signature:   crypto.BLSSignature{byte(i)},
						Index:       uint64(i),
					}
				}

				currentWithdrawals := make([]*types.WithdrawalRequest, maxWithdrawals)
				karalabeWithdrawals := make([]*WithdrawalRequestKaralabe, maxWithdrawals)
				for i := 0; i < maxWithdrawals; i++ {
					addr := common.ExecutionAddress{}
					addr[0] = byte(i)
					pubkey := crypto.BLSPubkey{}
					pubkey[0] = byte(i + 50)

					currentWithdrawals[i] = &types.WithdrawalRequest{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          1000000000,
					}
					karalabeWithdrawals[i] = &WithdrawalRequestKaralabe{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          1000000000,
					}
				}

				currentCons := make([]*types.ConsolidationRequest, maxCons)
				karalabeCons := make([]*ConsolidationRequestKaralabe, maxCons)
				for i := 0; i < maxCons; i++ {
					addr := common.ExecutionAddress{}
					addr[0] = byte(i + 100)
					srcPubkey := crypto.BLSPubkey{}
					srcPubkey[0] = byte(i + 120)
					tgtPubkey := crypto.BLSPubkey{}
					tgtPubkey[0] = byte(i + 140)

					currentCons[i] = &types.ConsolidationRequest{
						SourceAddress: addr,
						SourcePubKey:  srcPubkey,
						TargetPubKey:  tgtPubkey,
					}
					karalabeCons[i] = &ConsolidationRequestKaralabe{
						SourceAddress: addr,
						SourcePubKey:  srcPubkey,
						TargetPubKey:  tgtPubkey,
					}
				}

				return &types.ExecutionRequests{
						Deposits:       currentDeps,
						Withdrawals:    currentWithdrawals,
						Consolidations: currentCons,
					}, &ExecutionRequestsKaralabe{
						Deposits:       karalabeDeps,
						Withdrawals:    karalabeWithdrawals,
						Consolidations: karalabeCons,
					}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current, karalabe := tc.setup()

			// Test Marshal
			currentBytes, err1 := current.MarshalSSZ()
			require.NoError(t, err1, "current MarshalSSZ should not error")

			karalabelBytes, err2 := karalabe.MarshalSSZ()
			require.NoError(t, err2, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalabelBytes, currentBytes, "marshaled bytes should be identical")

			// Test Size
			require.Equal(t, int(karalabe.SizeSSZ(false)), current.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.ExecutionRequests{}
			err := newCurrent.UnmarshalSSZ(karalabelBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, len(current.Deposits), len(newCurrent.Deposits), "deposits lengths should match")
			require.Equal(t, len(current.Withdrawals), len(newCurrent.Withdrawals), "withdrawals lengths should match")
			require.Equal(t, len(current.Consolidations), len(newCurrent.Consolidations), "consolidations lengths should match")

			// Compare each element
			for i := range current.Deposits {
				require.Equal(t, current.Deposits[i], newCurrent.Deposits[i], "deposit at index %d should match", i)
			}
			for i := range current.Withdrawals {
				require.Equal(t, current.Withdrawals[i], newCurrent.Withdrawals[i], "withdrawal at index %d should match", i)
			}
			for i := range current.Consolidations {
				require.Equal(t, current.Consolidations[i], newCurrent.Consolidations[i], "consolidation at index %d should match", i)
			}

			// Test Unmarshal with current marshaled data
			newKaralabe := &ExecutionRequestsKaralabe{}
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")
			require.Equal(t, len(karalabe.Deposits), len(newKaralabe.Deposits), "deposits lengths should match")
			require.Equal(t, len(karalabe.Withdrawals), len(newKaralabe.Withdrawals), "withdrawals lengths should match")
			require.Equal(t, len(karalabe.Consolidations), len(newKaralabe.Consolidations), "consolidations lengths should match")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabelRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabelRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

// TestExecutionRequestsCompatibilityEdgeCases tests edge cases for ExecutionRequests handling
func TestExecutionRequestsCompatibilityEdgeCases(t *testing.T) {
	testCases := []struct {
		name          string
		data          []byte
		expectError   bool
		errorContains string
	}{
		{
			name:        "empty data",
			data:        []byte{},
			expectError: true,
		},
		{
			name:        "insufficient data for offsets",
			data:        []byte{1, 2, 3, 4, 5, 6, 7, 8}, // Only 8 bytes, need 12
			expectError: true,
		},
		{
			name:        "valid empty execution requests",
			data:        []byte{12, 0, 0, 0, 12, 0, 0, 0, 12, 0, 0, 0}, // All offsets point to 12 (after headers)
			expectError: false,
		},
		{
			name:        "invalid offset ordering",
			data:        []byte{20, 0, 0, 0, 16, 0, 0, 0, 12, 0, 0, 0}, // Offsets not in ascending order
			expectError: true,
		},
		{
			name:        "offset beyond data",
			data:        []byte{255, 255, 255, 255, 12, 0, 0, 0, 12, 0, 0, 0}, // First offset way too large
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with current implementation
			current := &types.ExecutionRequests{}
			currentErr := current.UnmarshalSSZ(tc.data)

			// Test with karalabe implementation
			karalabe := &ExecutionRequestsKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)

			// Both should behave the same way
			if tc.expectError {
				require.Error(t, currentErr, "current should error")
				require.Error(t, karalabelErr, "karalabe should error")
			} else {
				require.NoError(t, currentErr, "current should not error")
				require.NoError(t, karalabelErr, "karalabe should not error")
			}
		})
	}
}

// TestExecutionRequestsCompatibilityOrdering tests that ordering is preserved
func TestExecutionRequestsCompatibilityOrdering(t *testing.T) {
	// Create requests with specific ordering
	current := &types.ExecutionRequests{
		Deposits:       make([]*types.DepositRequest, 3),
		Withdrawals:    make([]*types.WithdrawalRequest, 2),
		Consolidations: make([]*types.ConsolidationRequest, 1),
	}
	karalabe := &ExecutionRequestsKaralabe{
		Deposits:       make([]*DepositKaralabe, 3),
		Withdrawals:    make([]*WithdrawalRequestKaralabe, 2),
		Consolidations: make([]*ConsolidationRequestKaralabe, 1),
	}

	// Create deposits with decreasing indices to test ordering
	for i := 0; i < 3; i++ {
		pubkey := crypto.BLSPubkey{}
		pubkey[0] = byte(3 - i) // Reverse order

		current.Deposits[i] = &types.Deposit{
			Pubkey:      pubkey,
			Credentials: types.WithdrawalCredentials{0x01},
			Amount:      32000000000,
			Signature:   crypto.BLSSignature{byte(i)},
			Index:       uint64(1000 - i*100), // Decreasing indices
		}
		karalabe.Deposits[i] = &DepositKaralabe{
			Pubkey:      pubkey,
			Credentials: types.WithdrawalCredentials{0x01},
			Amount:      32000000000,
			Signature:   crypto.BLSSignature{byte(i)},
			Index:       uint64(1000 - i*100),
		}
	}

	// Create withdrawals with specific amounts
	for i := 0; i < 2; i++ {
		addr := common.ExecutionAddress{}
		addr[0] = byte(100 - i*10)

		current.Withdrawals[i] = &types.WithdrawalRequest{
			SourceAddress:   addr,
			ValidatorPubKey: crypto.BLSPubkey{byte(50 - i*10)},
			Amount:          math.Gwei(5000000000 - uint64(i)*1000000000),
		}
		karalabe.Withdrawals[i] = &WithdrawalRequestKaralabe{
			SourceAddress:   addr,
			ValidatorPubKey: crypto.BLSPubkey{byte(50 - i*10)},
			Amount:          math.Gwei(5000000000 - uint64(i)*1000000000),
		}
	}

	// Create single consolidation
	current.Consolidations[0] = &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{200},
		SourcePubKey:  crypto.BLSPubkey{201},
		TargetPubKey:  crypto.BLSPubkey{202},
	}
	karalabe.Consolidations[0] = &ConsolidationRequestKaralabe{
		SourceAddress: common.ExecutionAddress{200},
		SourcePubKey:  crypto.BLSPubkey{201},
		TargetPubKey:  crypto.BLSPubkey{202},
	}

	// Marshal both
	currentBytes, err := current.MarshalSSZ()
	require.NoError(t, err)
	karalabelBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal into new objects
	newCurrent := &types.ExecutionRequests{}
	err = newCurrent.UnmarshalSSZ(karalabelBytes)
	require.NoError(t, err)

	newKaralabe := &ExecutionRequestsKaralabe{}
	err = newKaralabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Verify ordering is preserved for deposits
	for i := 0; i < 3; i++ {
		require.Equal(t, byte(3-i), newCurrent.Deposits[i].Pubkey[0], "current deposit ordering preserved at index %d", i)
		require.Equal(t, uint64(1000-i*100), newCurrent.Deposits[i].Index, "current deposit index preserved at index %d", i)

		require.Equal(t, byte(3-i), newKaralabe.Deposits[i].Pubkey[0], "karalabe deposit ordering preserved at index %d", i)
		require.Equal(t, uint64(1000-i*100), newKaralabe.Deposits[i].Index, "karalabe deposit index preserved at index %d", i)
	}

	// Verify ordering for withdrawals
	for i := 0; i < 2; i++ {
		require.Equal(t, byte(100-i*10), newCurrent.Withdrawals[i].SourceAddress[0], "current withdrawal ordering preserved at index %d", i)
		require.Equal(t, math.Gwei(5000000000-uint64(i)*1000000000), newCurrent.Withdrawals[i].Amount, "current withdrawal amount preserved at index %d", i)

		require.Equal(t, byte(100-i*10), newKaralabe.Withdrawals[i].SourceAddress[0], "karalabe withdrawal ordering preserved at index %d", i)
		require.Equal(t, math.Gwei(5000000000-uint64(i)*1000000000), newKaralabe.Withdrawals[i].Amount, "karalabe withdrawal amount preserved at index %d", i)
	}
}

// TestExecutionRequestsCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestExecutionRequestsCompatibilityRoundTrip(t *testing.T) {
	// Create an ExecutionRequests with all types of requests
	original := &types.ExecutionRequests{
		Deposits: []*types.DepositRequest{
			{
				Pubkey:      crypto.BLSPubkey{1, 2, 3, 4, 5},
				Credentials: types.WithdrawalCredentials{0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
				Amount:      32000000000,
				Signature:   crypto.BLSSignature{40, 41, 42, 43, 44, 45},
				Index:       1337,
			},
		},
		Withdrawals: []*types.WithdrawalRequest{
			{
				SourceAddress:   common.ExecutionAddress{50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69},
				ValidatorPubKey: crypto.BLSPubkey{70, 71, 72, 73, 74},
				Amount:          16000000000,
			},
			{
				SourceAddress:   common.ExecutionAddress{80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99},
				ValidatorPubKey: crypto.BLSPubkey{100, 101, 102, 103, 104},
				Amount:          8000000000,
			},
		},
		Consolidations: []*types.ConsolidationRequest{
			{
				SourceAddress: common.ExecutionAddress{110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129},
				SourcePubKey:  crypto.BLSPubkey{130, 131, 132, 133, 134},
				TargetPubKey:  crypto.BLSPubkey{140, 141, 142, 143, 144},
			},
		},
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &ExecutionRequestsKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalabelBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.ExecutionRequests{}
	err = roundTrip.UnmarshalSSZ(karalabelBytes)
	require.NoError(t, err)

	// Verify round trip preserved all data
	require.Equal(t, len(original.Deposits), len(roundTrip.Deposits), "deposits length should match")
	require.Equal(t, len(original.Withdrawals), len(roundTrip.Withdrawals), "withdrawals length should match")
	require.Equal(t, len(original.Consolidations), len(roundTrip.Consolidations), "consolidations length should match")

	// Verify each element
	for i := range original.Deposits {
		require.Equal(t, original.Deposits[i], roundTrip.Deposits[i], "deposit %d should match after round trip", i)
	}
	for i := range original.Withdrawals {
		require.Equal(t, original.Withdrawals[i], roundTrip.Withdrawals[i], "withdrawal %d should match after round trip", i)
	}
	for i := range original.Consolidations {
		require.Equal(t, original.Consolidations[i], roundTrip.Consolidations[i], "consolidation %d should match after round trip", i)
	}

	// Verify both serializations are identical
	require.Equal(t, currentBytes, karalabelBytes, "both serializations should be identical")

	// Verify hash roots match throughout
	originalRoot, err := original.HashTreeRoot()
	require.NoError(t, err)
	roundTripRoot, err := roundTrip.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, originalRoot, [32]byte(karalabe.HashTreeRoot()), "hash roots should match")
	require.Equal(t, originalRoot, roundTripRoot, "hash roots should match after round trip")
}

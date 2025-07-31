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
	karalabe "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// Compile-time assertions to ensure SigningDataKaralabe implements necessary interfaces.
var _ karalabe.StaticObject = (*SigningDataKaralabe)(nil)

// SigningDataKaralabe as defined in the Ethereum 2.0 specification - exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#signingdata
type SigningDataKaralabe struct {
	// ObjectRoot is the hash tree root of the object being signed.
	ObjectRoot common.Root
	// Domain is the domain the object is being signed in.
	Domain common.Domain
}

// SizeSSZ returns the size of the SigningData object in SSZ encoding.
func (*SigningDataKaralabe) SizeSSZ() uint32 {
	//nolint:mnd // 32*2 = 64.
	return 64
}

// DefineSSZ defines the SSZ encoding for the SigningData object.
func (s *SigningDataKaralabe) DefineSSZ(codec *karalabe.Codec) {
	karalabe.DefineStaticBytes(codec, &s.ObjectRoot)
	karalabe.DefineStaticBytes(codec, &s.Domain)
}

// HashTreeRoot computes the SSZ hash tree root of the SigningData object.
func (s *SigningDataKaralabe) HashTreeRoot() common.Root {
	return karalabe.HashSequential(s)
}

// MarshalSSZ marshals the SigningData object to SSZ format.
func (s *SigningDataKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabe.Size(s))
	return buf, karalabe.EncodeToBytes(buf, s)
}

// UnmarshalSSZ unmarshals the SigningData object from SSZ format.
func (s *SigningDataKaralabe) UnmarshalSSZ(buf []byte) error {
	return karalabe.DecodeFromBytes(buf, s)
}

// TestSigningDataSSZRegression ensures that the SSZ encoding for SigningData
// remains stable and backward compatible.
func TestSigningDataSSZRegression(t *testing.T) {
	testCases := []struct {
		name        string
		signingData *types.SigningData
		expectedSSZ []byte // Pre-computed expected SSZ encoding
	}{
		{
			name: "zero values",
			signingData: &types.SigningData{
				ObjectRoot: common.Root{},
				Domain:     common.Domain{},
			},
			// Expected SSZ: 32 zero bytes (root) + 32 zero bytes (domain) = 64 bytes
			expectedSSZ: make([]byte, 64),
		},
		{
			name: "typical signing data",
			signingData: &types.SigningData{
				ObjectRoot: common.Root{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
				},
				Domain: common.Domain{
					0x07, 0x00, 0x00, 0x00, // Domain type beacon proposer
					0x01, 0x02, 0x03, 0x04, // Fork version
					0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, // Genesis validators root
					0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
					0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8,
				},
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 64)
				// ObjectRoot
				copy(ssz[0:32], []byte{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
				})
				// Domain
				copy(ssz[32:64], []byte{
					0x07, 0x00, 0x00, 0x00, // Domain type beacon proposer
					0x01, 0x02, 0x03, 0x04, // Fork version
					0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, // Genesis validators root
					0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
					0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8,
				})
				return ssz
			}(),
		},
		{
			name: "all max values",
			signingData: &types.SigningData{
				ObjectRoot: common.Root{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				Domain: common.Domain{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 64)
				for i := range ssz {
					ssz[i] = 0xff
				}
				return ssz
			}(),
		},
		{
			name: "specific domain values",
			signingData: &types.SigningData{
				ObjectRoot: common.Root{
					0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
					0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
					0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
					0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
				},
				Domain: common.Domain{
					0x00, 0x00, 0x00, 0x00, // Domain type beacon proposer
					0xaa, 0xbb, 0xcc, 0xdd, // Fork version
					0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, // Genesis validators root
					0x33, 0x33, 0x33, 0x33, 0x44, 0x44, 0x44, 0x44,
					0x55, 0x55, 0x55, 0x55, 0x66, 0x66, 0x66, 0x66,
				},
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 64)
				copy(ssz[0:32], []byte{
					0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
					0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
					0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
					0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
				})
				copy(ssz[32:64], []byte{
					0x00, 0x00, 0x00, 0x00, // Domain type beacon proposer
					0xaa, 0xbb, 0xcc, 0xdd, // Fork version
					0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, // Genesis validators root
					0x33, 0x33, 0x33, 0x33, 0x44, 0x44, 0x44, 0x44,
					0x55, 0x55, 0x55, 0x55, 0x66, 0x66, 0x66, 0x66,
				})
				return ssz
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Marshal
			actualSSZ, err := tc.signingData.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ, "SSZ encoding should match expected")

			// Test Size
			require.Equal(t, 64, tc.signingData.SizeSSZ(), "size should be 64 bytes")

			// Test Unmarshal
			unmarshaled := &types.SigningData{}
			err = unmarshaled.UnmarshalSSZ(tc.expectedSSZ)
			require.NoError(t, err, "UnmarshalSSZ should not error")
			require.Equal(t, tc.signingData, unmarshaled, "unmarshaled object should match original")

			// Test MarshalSSZTo
			buf := make([]byte, 0, tc.signingData.SizeSSZ())
			actualSSZ2, err := tc.signingData.MarshalSSZTo(buf)
			require.NoError(t, err, "MarshalSSZTo should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ2, "MarshalSSZTo should produce same output")

			// Test HashTreeRoot consistency
			root1, err := tc.signingData.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot should not error")
			root2, err := unmarshaled.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot of unmarshaled should not error")
			require.Equal(t, root1, root2, "hash tree roots should match")
		})
	}
}

// TestSigningDataSSZInvalidData tests error handling for invalid SSZ data
func TestSigningDataSSZInvalidData(t *testing.T) {
	testCases := []struct {
		name          string
		data          []byte
		expectedError string
	}{
		{
			name:          "empty data",
			data:          []byte{},
			expectedError: "incorrect size",
		},
		{
			name:          "insufficient data",
			data:          make([]byte, 30), // less than required 64 bytes
			expectedError: "incorrect size",
		},
		{
			name:          "excess data",
			data:          make([]byte, 100), // more than required 64 bytes
			expectedError: "incorrect size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signingData := &types.SigningData{}
			err := signingData.UnmarshalSSZ(tc.data)
			require.Error(t, err, "UnmarshalSSZ should error on invalid data")
			require.Contains(t, err.Error(), tc.expectedError, "error should contain expected message")
		})
	}
}

// TestSigningDataSSZRoundTrip tests round-trip encoding/decoding with various data patterns
func TestSigningDataSSZRoundTrip(t *testing.T) {
	patterns := []struct {
		name  string
		setup func() *types.SigningData
	}{
		{
			name: "all zeros",
			setup: func() *types.SigningData {
				return &types.SigningData{}
			},
		},
		{
			name: "incremental pattern",
			setup: func() *types.SigningData {
				var root common.Root
				for i := range root {
					root[i] = byte(i)
				}
				var domain common.Domain
				for i := range domain {
					domain[i] = byte(i + 32)
				}
				return &types.SigningData{
					ObjectRoot: root,
					Domain:     domain,
				}
			},
		},
		{
			name: "specific values",
			setup: func() *types.SigningData {
				return &types.SigningData{
					ObjectRoot: common.Root{
						0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
						0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00,
						0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
						0x0f, 0x1e, 0x2d, 0x3c, 0x4b, 0x5a, 0x69, 0x78,
					},
					Domain: common.Domain{
						0x07, 0x00, 0x00, 0x00, // Domain type
						0x01, 0x02, 0x03, 0x04, // Fork version
						0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
						0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00,
						0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
					},
				}
			},
		},
	}

	for _, pattern := range patterns {
		t.Run(pattern.name, func(t *testing.T) {
			original := pattern.setup()

			// Marshal
			data, err := original.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")

			// Unmarshal
			decoded := &types.SigningData{}
			err = decoded.UnmarshalSSZ(data)
			require.NoError(t, err, "UnmarshalSSZ should not error")

			// Compare
			require.Equal(t, original, decoded, "round trip should preserve data")

			// Verify hash tree roots match
			root1, err := original.HashTreeRoot()
			require.NoError(t, err)
			root2, err := decoded.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, root1, root2, "hash tree roots should match after round trip")
		})
	}
}

// TestSigningDataComputeSigningRoot tests that ComputeSigningRoot produces expected results
func TestSigningDataComputeSigningRoot(t *testing.T) {
	// Define a test type that implements HashTreeRoot
	type testObject struct {
		value uint64
	}

	// Simple HashTreeRoot implementation for testing
	// This mimics how a real object would compute its hash tree root
	hashTreeRootFunc := func(obj *testObject) ([32]byte, error) {
		var root [32]byte
		// Simple deterministic hash for testing
		for i := 0; i < 8; i++ {
			root[i] = byte(obj.value >> (i * 8))
		}
		return root, nil
	}

	testCases := []struct {
		name   string
		value  uint64
		domain common.Domain
	}{
		{
			name:   "zero value and domain",
			value:  0,
			domain: common.Domain{},
		},
		{
			name:  "specific value and domain",
			value: 12345,
			domain: common.Domain{
				0x07, 0x00, 0x00, 0x00, // Domain type
				0x01, 0x02, 0x03, 0x04, // Fork version
				0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
				0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create object and compute its HTR
			obj := &testObject{value: tc.value}
			htr, err := hashTreeRootFunc(obj)
			require.NoError(t, err)

			// Use ComputeSigningRootFromHTR
			signingRoot := types.ComputeSigningRootFromHTR(htr, tc.domain)

			// Manually create SigningData and compute its root
			signingData := &types.SigningData{
				ObjectRoot: common.Root(htr),
				Domain:     tc.domain,
			}
			expectedRoot, err := signingData.HashTreeRoot()
			require.NoError(t, err)

			// Verify they match
			require.Equal(t, common.Root(expectedRoot), signingRoot, "signing roots should match")
		})
	}
}

// TestSigningDataComputeSigningRootUInt64 tests the uint64 specific signing root function
func TestSigningDataComputeSigningRootUInt64(t *testing.T) {
	testCases := []struct {
		name   string
		value  uint64
		domain common.Domain
	}{
		{
			name:   "zero value",
			value:  0,
			domain: common.Domain{0x01, 0x00, 0x00, 0x00}, // Some domain
		},
		{
			name:   "max uint64",
			value:  ^uint64(0),
			domain: common.Domain{0x02, 0x00, 0x00, 0x00}, // Another domain
		},
		{
			name:  "specific value",
			value: 999999,
			domain: common.Domain{
				0x03, 0x00, 0x00, 0x00, // Domain type
				0xaa, 0xbb, 0xcc, 0xdd, // Fork version
				0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22,
				0x33, 0x33, 0x33, 0x33, 0x44, 0x44, 0x44, 0x44,
				0x55, 0x55, 0x55, 0x55, 0x66, 0x66, 0x66, 0x66,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Compute signing root using the function
			signingRoot := types.ComputeSigningRootUInt64(tc.value, tc.domain)

			// Verify it's deterministic
			signingRoot2 := types.ComputeSigningRootUInt64(tc.value, tc.domain)
			require.Equal(t, signingRoot, signingRoot2, "signing root should be deterministic")

			// Different value should produce different root
			if tc.value != 0 {
				differentRoot := types.ComputeSigningRootUInt64(tc.value-1, tc.domain)
				require.NotEqual(t, signingRoot, differentRoot, "different values should produce different roots")
			}

			// Different domain should produce different root
			var differentDomain common.Domain
			copy(differentDomain[:], tc.domain[:])
			differentDomain[0] ^= 0xff
			differentRoot := types.ComputeSigningRootUInt64(tc.value, differentDomain)
			require.NotEqual(t, signingRoot, differentRoot, "different domains should produce different roots")
		})
	}
}

// TestSigningDataCompatibility tests that current and karalabe implementations produce identical results
func TestSigningDataCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.SigningData, *SigningDataKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.SigningData, *SigningDataKaralabe) {
				return &types.SigningData{},
					&SigningDataKaralabe{}
			},
		},
		{
			name: "typical signing data",
			setup: func() (*types.SigningData, *SigningDataKaralabe) {
				objectRoot := common.Root{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
				}
				domain := common.Domain{
					0x07, 0x00, 0x00, 0x00, // Domain type beacon proposer
					0x01, 0x02, 0x03, 0x04, // Fork version
					0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, // Genesis validators root
					0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
					0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8,
				}

				current := &types.SigningData{
					ObjectRoot: objectRoot,
					Domain:     domain,
				}
				karalabe := &SigningDataKaralabe{
					ObjectRoot: objectRoot,
					Domain:     domain,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.SigningData, *SigningDataKaralabe) {
				var objectRoot common.Root
				var domain common.Domain
				for i := range objectRoot {
					objectRoot[i] = 0xFF
				}
				for i := range domain {
					domain[i] = 0xFF
				}

				current := &types.SigningData{
					ObjectRoot: objectRoot,
					Domain:     domain,
				}
				karalabe := &SigningDataKaralabe{
					ObjectRoot: objectRoot,
					Domain:     domain,
				}
				return current, karalabe
			},
		},
		{
			name: "specific values",
			setup: func() (*types.SigningData, *SigningDataKaralabe) {
				objectRoot := common.Root{
					0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
					0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
					0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
					0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
				}
				domain := common.Domain{
					0x00, 0x00, 0x00, 0x00, // Domain type beacon proposer
					0xaa, 0xbb, 0xcc, 0xdd, // Fork version
					0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, // Genesis validators root
					0x33, 0x33, 0x33, 0x33, 0x44, 0x44, 0x44, 0x44,
					0x55, 0x55, 0x55, 0x55, 0x66, 0x66, 0x66, 0x66,
				}

				current := &types.SigningData{
					ObjectRoot: objectRoot,
					Domain:     domain,
				}
				karalabe := &SigningDataKaralabe{
					ObjectRoot: objectRoot,
					Domain:     domain,
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

			karalabeBytes, err2 := karalabe.MarshalSSZ()
			require.NoError(t, err2, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalabeBytes, currentBytes, "marshaled bytes should be identical")

			// Test Size
			require.Equal(t, 64, current.SizeSSZ(), "current size should be 64")
			require.Equal(t, uint32(64), karalabe.SizeSSZ(), "karalabe size should be 64")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.SigningData{}
			err := newCurrent.UnmarshalSSZ(karalabeBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &SigningDataKaralabe{}
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")
			require.Equal(t, karalabe, newKaralabe, "unmarshaled karalabe should match original")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabeRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabeRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

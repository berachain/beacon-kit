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
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	ssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// SimplifiedBeaconStateKaralabe is a simplified version of BeaconState that focuses on
// the fork-specific SSZ logic. This is extracted from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// in beacon-kit-karalabe at path consensus-types/types/state.go
type SimplifiedBeaconStateKaralabe struct {
	constraints.Versionable `json:"-"`

	// Core fields (simplified to focus on fork logic)
	Slot          math.Slot `json:"slot,omitempty"`
	TotalSlashing math.Gwei `json:"total_slashing,omitempty"`

	// Dynamic field to represent other fields without implementing full SSZ
	DummyDynamicField []uint64 `json:"dummy_dynamic_field,omitempty"`

	// Fork-specific field: only present in Electra and later
	PendingPartialWithdrawals []*PendingPartialWithdrawalKaralabe `json:"pending_partial_withdrawals,omitempty"`
}

// NewEmptySimplifiedBeaconStateWithVersionKaralabe returns a new empty SimplifiedBeaconStateKaralabe with the given fork version.
func NewEmptySimplifiedBeaconStateWithVersionKaralabe(version common.Version) *SimplifiedBeaconStateKaralabe {
	return &SimplifiedBeaconStateKaralabe{
		Versionable: types.NewVersionable(version),
	}
}

// SizeSSZ returns the ssz encoded size in bytes for the SimplifiedBeaconStateKaralabe object.
// Exact copy of fork-specific logic from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (st *SimplifiedBeaconStateKaralabe) SizeSSZ(fixed bool) uint32 {
	// Base size: Slot (8) + TotalSlashing (8) + DummyDynamicField offset (4) = 20
	var size uint32 = 20

	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		// Add 4 for PendingPartialWithdrawals offset after Electra
		size += 4
	}

	if fixed {
		return size
	}

	// Dynamic size fields
	size += ssz.SizeSliceOfUint64s(st.DummyDynamicField)
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		size += ssz.SizeSliceOfStaticObjects(st.PendingPartialWithdrawals)
	}

	return size
}

// DefineSSZ defines the SSZ encoding for the SimplifiedBeaconStateKaralabe object.
// Fork-specific logic extracted from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (st *SimplifiedBeaconStateKaralabe) DefineSSZ(codec *ssz.Codec) {
	// Fixed fields
	ssz.DefineUint64(codec, &st.Slot)
	ssz.DefineUint64(codec, (*uint64)(&st.TotalSlashing))

	// Dynamic field offset
	ssz.DefineSliceOfUint64sOffset(codec, &st.DummyDynamicField, 100)

	// Electra-specific offset
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		ssz.DefineSliceOfStaticObjectsOffset(codec, &st.PendingPartialWithdrawals, constants.PendingPartialWithdrawalsLimit)
	}

	// Dynamic content
	ssz.DefineSliceOfUint64sContent(codec, &st.DummyDynamicField, 100)

	// Electra-specific content
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		ssz.DefineSliceOfStaticObjectsContent(codec, &st.PendingPartialWithdrawals, constants.PendingPartialWithdrawalsLimit)
	}
}

// MarshalSSZ marshals the SimplifiedBeaconStateKaralabe into SSZ format.
func (st *SimplifiedBeaconStateKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(st))
	return buf, ssz.EncodeToBytes(buf, st)
}

// UnmarshalSSZ unmarshals the SimplifiedBeaconStateKaralabe from SSZ format.
func (st *SimplifiedBeaconStateKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, st)
}

// HashTreeRoot computes the Merkleization of the SimplifiedBeaconStateKaralabe.
func (st *SimplifiedBeaconStateKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(st)
}

func (st *SimplifiedBeaconStateKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// SimplifiedBeaconState is a simplified version of the current BeaconState implementation
// that focuses on the fork-specific SSZ logic for testing compatibility
type SimplifiedBeaconState struct {
	types.Versionable `json:"-"`

	// Core fields (simplified to focus on fork logic)
	Slot          math.Slot `json:"slot,omitempty"`
	TotalSlashing math.Gwei `json:"total_slashing,omitempty"`

	// Dynamic field to represent other fields without implementing full SSZ
	DummyDynamicField []uint64 `json:"dummy_dynamic_field,omitempty"`

	// Fork-specific field: only present in Electra and later
	PendingPartialWithdrawals []*types.PendingPartialWithdrawal `json:"pending_partial_withdrawals,omitempty"`
}

// NewEmptySimplifiedBeaconStateWithVersion returns a new empty SimplifiedBeaconState with the given fork version.
func NewEmptySimplifiedBeaconStateWithVersion(version common.Version) *SimplifiedBeaconState {
	return &SimplifiedBeaconState{
		Versionable: types.NewVersionable(version),
	}
}

// SizeSSZ returns the ssz encoded size in bytes for the SimplifiedBeaconState object.
func (st *SimplifiedBeaconState) SizeSSZ() int {
	// Base size: Slot (8) + TotalSlashing (8) + DummyDynamicField offset (4) = 20
	var size = 20

	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		// Add 4 for PendingPartialWithdrawals offset after Electra
		size += 4
	}

	// Dynamic size fields
	size += len(st.DummyDynamicField) * 8
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		size += len(st.PendingPartialWithdrawals) * 24 // Each pending withdrawal is 24 bytes
	}

	return size
}

// MarshalSSZ marshals the SimplifiedBeaconState into SSZ format.
func (st *SimplifiedBeaconState) MarshalSSZ() ([]byte, error) {
	return st.MarshalSSZTo(make([]byte, 0, st.SizeSSZ()))
}

// MarshalSSZTo marshals the SimplifiedBeaconState to a target array.
func (st *SimplifiedBeaconState) MarshalSSZTo(buf []byte) ([]byte, error) {
	// Encode Slot
	buf = ssz.MarshalUint64(buf, uint64(st.Slot))

	// Encode TotalSlashing
	buf = ssz.MarshalUint64(buf, uint64(st.TotalSlashing))

	// Calculate offsets
	offset := uint32(20) // After fixed fields
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		offset += 4
	}

	// Write DummyDynamicField offset
	buf = ssz.MarshalUint32(buf, offset)
	offset += uint32(len(st.DummyDynamicField) * 8)

	// Write PendingPartialWithdrawals offset if Electra
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		buf = ssz.MarshalUint32(buf, offset)
	}

	// Write DummyDynamicField content
	for _, val := range st.DummyDynamicField {
		buf = ssz.MarshalUint64(buf, val)
	}

	// Write PendingPartialWithdrawals content if Electra
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		for _, ppw := range st.PendingPartialWithdrawals {
			ppwBytes, err := ppw.MarshalSSZ()
			if err != nil {
				return nil, err
			}
			buf = append(buf, ppwBytes...)
		}
	}

	return buf, nil
}

// UnmarshalSSZ unmarshals the SimplifiedBeaconState from SSZ format.
func (st *SimplifiedBeaconState) UnmarshalSSZ(buf []byte) error {
	// This is a simplified implementation for testing
	// A full implementation would handle all fields properly

	expectedMinSize := 20
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		expectedMinSize = 24
	}

	if len(buf) < expectedMinSize {
		return ssz.ErrSize
	}

	// Decode Slot
	st.Slot = math.Slot(ssz.UnmarshallUint64(buf[0:8]))

	// Decode TotalSlashing
	st.TotalSlashing = math.Gwei(ssz.UnmarshallUint64(buf[8:16]))

	// Decode DummyDynamicField offset
	dummyOffset := ssz.UnmarshallUint32(buf[16:20])

	// Decode PendingPartialWithdrawals offset if Electra
	var ppwOffset uint32
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		ppwOffset = ssz.UnmarshallUint32(buf[20:24])
	}

	// Decode DummyDynamicField
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		if ppwOffset > dummyOffset {
			numDummy := (ppwOffset - dummyOffset) / 8
			st.DummyDynamicField = make([]uint64, numDummy)
			for i := uint32(0); i < numDummy; i++ {
				start := dummyOffset + i*8
				st.DummyDynamicField[i] = ssz.UnmarshallUint64(buf[start : start+8])
			}
		}

		// Decode PendingPartialWithdrawals
		if ppwOffset < uint32(len(buf)) {
			remaining := buf[ppwOffset:]
			numPPW := len(remaining) / 24
			st.PendingPartialWithdrawals = make([]*types.PendingPartialWithdrawal, numPPW)
			for i := 0; i < numPPW; i++ {
				st.PendingPartialWithdrawals[i] = &types.PendingPartialWithdrawal{}
				if err := st.PendingPartialWithdrawals[i].UnmarshalSSZ(remaining[i*24 : (i+1)*24]); err != nil {
					return err
				}
			}
		}
	} else {
		// Deneb: no PendingPartialWithdrawals
		if dummyOffset < uint32(len(buf)) {
			remaining := buf[dummyOffset:]
			numDummy := len(remaining) / 8
			st.DummyDynamicField = make([]uint64, numDummy)
			for i := 0; i < numDummy; i++ {
				st.DummyDynamicField[i] = ssz.UnmarshallUint64(remaining[i*8 : (i+1)*8])
			}
		}
	}

	return nil
}

// HashTreeRoot computes the hash tree root of the SimplifiedBeaconState.
func (st *SimplifiedBeaconState) HashTreeRoot() ([32]byte, error) {
	// For testing purposes, we'll use a simplified hash
	// The actual implementation would use ssz.HashTreeRoot
	return [32]byte{}, nil
}

// TestBeaconStateCompatibility tests that the fork-specific SSZ encoding logic
// is compatible between karalabe/ssz and fastssz implementations
func TestBeaconStateCompatibility(t *testing.T) {
	testCases := []struct {
		name    string
		version common.Version
		setup   func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe)
	}{
		{
			name:    "Deneb - no PendingPartialWithdrawals",
			version: version.Deneb(),
			setup: func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe) {
				current := NewEmptySimplifiedBeaconStateWithVersion(version.Deneb())
				current.Slot = 12345
				current.TotalSlashing = 1000000
				current.DummyDynamicField = []uint64{1, 2, 3, 4, 5}

				karalabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Deneb())
				karalabe.Slot = 12345
				karalabe.TotalSlashing = 1000000
				karalabe.DummyDynamicField = []uint64{1, 2, 3, 4, 5}

				return current, karalabe
			},
		},
		{
			name:    "Electra - with empty PendingPartialWithdrawals",
			version: version.Electra(),
			setup: func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe) {
				current := NewEmptySimplifiedBeaconStateWithVersion(version.Electra())
				current.Slot = 54321
				current.TotalSlashing = 2000000
				current.DummyDynamicField = []uint64{10, 20, 30}
				current.PendingPartialWithdrawals = []*types.PendingPartialWithdrawal{}

				karalabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Electra())
				karalabe.Slot = 54321
				karalabe.TotalSlashing = 2000000
				karalabe.DummyDynamicField = []uint64{10, 20, 30}
				karalabe.PendingPartialWithdrawals = []*PendingPartialWithdrawalKaralabe{}

				return current, karalabe
			},
		},
		{
			name:    "Electra - with PendingPartialWithdrawals",
			version: version.Electra(),
			setup: func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe) {
				ppw1 := &types.PendingPartialWithdrawal{
					ValidatorIndex:    100,
					Amount:            50000000,
					WithdrawableEpoch: 1000,
				}
				ppw2 := &types.PendingPartialWithdrawal{
					ValidatorIndex:    200,
					Amount:            75000000,
					WithdrawableEpoch: 2000,
				}

				ppw1Karalabe := &PendingPartialWithdrawalKaralabe{
					ValidatorIndex:    100,
					Amount:            50000000,
					WithdrawableEpoch: 1000,
				}
				ppw2Karalabe := &PendingPartialWithdrawalKaralabe{
					ValidatorIndex:    200,
					Amount:            75000000,
					WithdrawableEpoch: 2000,
				}

				current := NewEmptySimplifiedBeaconStateWithVersion(version.Electra())
				current.Slot = 99999
				current.TotalSlashing = 3000000
				current.DummyDynamicField = []uint64{100, 200}
				current.PendingPartialWithdrawals = []*types.PendingPartialWithdrawal{ppw1, ppw2}

				karalabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Electra())
				karalabe.Slot = 99999
				karalabe.TotalSlashing = 3000000
				karalabe.DummyDynamicField = []uint64{100, 200}
				karalabe.PendingPartialWithdrawals = []*PendingPartialWithdrawalKaralabe{ppw1Karalabe, ppw2Karalabe}

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

			karalabelBytes, err2 := karalabe.MarshalSSZ()
			require.NoError(t, err2, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalabelBytes, currentBytes, "marshaled bytes should be identical")

			// Test Size
			karalabelSize := karalabe.SizeSSZ(false)
			require.Equal(t, int(karalabelSize), current.SizeSSZ(), "sizes should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := NewEmptySimplifiedBeaconStateWithVersion(tc.version)
			err := newCurrent.UnmarshalSSZ(karalabelBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current.Slot, newCurrent.Slot, "slot should match after unmarshal")
			require.Equal(t, current.TotalSlashing, newCurrent.TotalSlashing, "total slashing should match after unmarshal")
			require.Equal(t, current.DummyDynamicField, newCurrent.DummyDynamicField, "dummy field should match after unmarshal")
			require.Equal(t, len(current.PendingPartialWithdrawals), len(newCurrent.PendingPartialWithdrawals), "pending withdrawals count should match")

			// Test Unmarshal with current marshaled data
			newKaralabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(tc.version)
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")
			require.Equal(t, karalabe.Slot, newKaralabe.Slot, "slot should match after unmarshal")
			require.Equal(t, karalabe.TotalSlashing, newKaralabe.TotalSlashing, "total slashing should match after unmarshal")
			require.Equal(t, karalabe.DummyDynamicField, newKaralabe.DummyDynamicField, "dummy field should match after unmarshal")
			require.Equal(t, len(karalabe.PendingPartialWithdrawals), len(newKaralabe.PendingPartialWithdrawals), "pending withdrawals count should match")

			// For Electra, verify PendingPartialWithdrawals content
			if tc.version == version.Electra() && len(current.PendingPartialWithdrawals) > 0 {
				for i, ppw := range current.PendingPartialWithdrawals {
					require.Equal(t, ppw.ValidatorIndex, newCurrent.PendingPartialWithdrawals[i].ValidatorIndex, "validator index should match")
					require.Equal(t, ppw.Amount, newCurrent.PendingPartialWithdrawals[i].Amount, "amount should match")
					require.Equal(t, ppw.WithdrawableEpoch, newCurrent.PendingPartialWithdrawals[i].WithdrawableEpoch, "withdrawable epoch should match")
				}
			}
		})
	}
}

// TestBeaconStateCompatibilityForkTransition tests the fork transition behavior
func TestBeaconStateCompatibilityForkTransition(t *testing.T) {
	// Create a Deneb state
	denebCurrent := NewEmptySimplifiedBeaconStateWithVersion(version.Deneb())
	denebCurrent.Slot = 8888
	denebCurrent.TotalSlashing = 500000
	denebCurrent.DummyDynamicField = []uint64{7, 8, 9}

	denebKaralabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Deneb())
	denebKaralabe.Slot = 8888
	denebKaralabe.TotalSlashing = 500000
	denebKaralabe.DummyDynamicField = []uint64{7, 8, 9}

	// Marshal Deneb states
	denebCurrentBytes, err := denebCurrent.MarshalSSZ()
	require.NoError(t, err)

	denebKaralabelBytes, err := denebKaralabe.MarshalSSZ()
	require.NoError(t, err)

	// Verify Deneb encoding is identical
	require.Equal(t, denebKaralabelBytes, denebCurrentBytes, "Deneb marshaled bytes should be identical")

	// Verify Deneb has no PendingPartialWithdrawals offset (20 bytes base + dynamic content)
	require.Equal(t, 20, len(denebCurrentBytes)-len(denebCurrent.DummyDynamicField)*8, "Deneb should have 20 bytes of fixed fields + offsets")

	// Create an Electra state with same base fields
	electraCurrent := NewEmptySimplifiedBeaconStateWithVersion(version.Electra())
	electraCurrent.Slot = 8888
	electraCurrent.TotalSlashing = 500000
	electraCurrent.DummyDynamicField = []uint64{7, 8, 9}
	electraCurrent.PendingPartialWithdrawals = []*types.PendingPartialWithdrawal{} // Empty but present

	electraKaralabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Electra())
	electraKaralabe.Slot = 8888
	electraKaralabe.TotalSlashing = 500000
	electraKaralabe.DummyDynamicField = []uint64{7, 8, 9}
	electraKaralabe.PendingPartialWithdrawals = []*PendingPartialWithdrawalKaralabe{} // Empty but present

	// Marshal Electra states
	electraCurrentBytes, err := electraCurrent.MarshalSSZ()
	require.NoError(t, err)

	electraKaralabelBytes, err := electraKaralabe.MarshalSSZ()
	require.NoError(t, err)

	// Verify Electra encoding is identical
	require.Equal(t, electraKaralabelBytes, electraCurrentBytes, "Electra marshaled bytes should be identical")

	// Verify Electra has PendingPartialWithdrawals offset (24 bytes base + dynamic content)
	require.Equal(t, 24, len(electraCurrentBytes)-len(electraCurrent.DummyDynamicField)*8, "Electra should have 24 bytes of fixed fields + offsets")

	// Verify Deneb and Electra encodings are different due to fork-specific field
	require.NotEqual(t, denebCurrentBytes, electraCurrentBytes, "Deneb and Electra encodings should differ")
}

// TestBeaconStateCompatibilityEdgeCases tests edge cases for SSZ encoding
func TestBeaconStateCompatibilityEdgeCases(t *testing.T) {
	testCases := []struct {
		name    string
		version common.Version
		setup   func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe)
	}{
		{
			name:    "Electra - maximum PendingPartialWithdrawals",
			version: version.Electra(),
			setup: func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe) {
				// Create many pending withdrawals (not actually at limit for test performance)
				var ppws []*types.PendingPartialWithdrawal
				var ppwsKaralabe []*PendingPartialWithdrawalKaralabe
				for i := 0; i < 10; i++ {
					ppws = append(ppws, &types.PendingPartialWithdrawal{
						ValidatorIndex:    math.ValidatorIndex(i),
						Amount:            math.Gwei(i * 1000000),
						WithdrawableEpoch: math.Epoch(i * 100),
					})
					ppwsKaralabe = append(ppwsKaralabe, &PendingPartialWithdrawalKaralabe{
						ValidatorIndex:    math.ValidatorIndex(i),
						Amount:            math.Gwei(i * 1000000),
						WithdrawableEpoch: math.Epoch(i * 100),
					})
				}

				current := NewEmptySimplifiedBeaconStateWithVersion(version.Electra())
				current.Slot = 77777
				current.TotalSlashing = 9999999
				current.DummyDynamicField = []uint64{}
				current.PendingPartialWithdrawals = ppws

				karalabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Electra())
				karalabe.Slot = 77777
				karalabe.TotalSlashing = 9999999
				karalabe.DummyDynamicField = []uint64{}
				karalabe.PendingPartialWithdrawals = ppwsKaralabe

				return current, karalabe
			},
		},
		{
			name:    "All zero values - Deneb",
			version: version.Deneb(),
			setup: func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe) {
				current := NewEmptySimplifiedBeaconStateWithVersion(version.Deneb())
				karalabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Deneb())
				return current, karalabe
			},
		},
		{
			name:    "All zero values - Electra",
			version: version.Electra(),
			setup: func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe) {
				current := NewEmptySimplifiedBeaconStateWithVersion(version.Electra())
				current.PendingPartialWithdrawals = []*types.PendingPartialWithdrawal{}

				karalabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Electra())
				karalabe.PendingPartialWithdrawals = []*PendingPartialWithdrawalKaralabe{}

				return current, karalabe
			},
		},
		{
			name:    "Maximum values",
			version: version.Electra(),
			setup: func() (*SimplifiedBeaconState, *SimplifiedBeaconStateKaralabe) {
				current := NewEmptySimplifiedBeaconStateWithVersion(version.Electra())
				current.Slot = math.Slot(^uint64(0))
				current.TotalSlashing = math.Gwei(^uint64(0))
				current.DummyDynamicField = []uint64{^uint64(0), ^uint64(0)}
				current.PendingPartialWithdrawals = []*types.PendingPartialWithdrawal{
					{
						ValidatorIndex:    math.ValidatorIndex(^uint64(0)),
						Amount:            math.Gwei(^uint64(0)),
						WithdrawableEpoch: math.Epoch(^uint64(0)),
					},
				}

				karalabe := NewEmptySimplifiedBeaconStateWithVersionKaralabe(version.Electra())
				karalabe.Slot = math.Slot(^uint64(0))
				karalabe.TotalSlashing = math.Gwei(^uint64(0))
				karalabe.DummyDynamicField = []uint64{^uint64(0), ^uint64(0)}
				karalabe.PendingPartialWithdrawals = []*PendingPartialWithdrawalKaralabe{
					{
						ValidatorIndex:    math.ValidatorIndex(^uint64(0)),
						Amount:            math.Gwei(^uint64(0)),
						WithdrawableEpoch: math.Epoch(^uint64(0)),
					},
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

			karalabelBytes, err2 := karalabe.MarshalSSZ()
			require.NoError(t, err2, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalabelBytes, currentBytes, "marshaled bytes should be identical")

			// Test round-trip
			newCurrent := NewEmptySimplifiedBeaconStateWithVersion(tc.version)
			err := newCurrent.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal should not error")

			// Re-marshal and verify it's still the same
			newCurrentBytes, err := newCurrent.MarshalSSZ()
			require.NoError(t, err, "re-marshal should not error")
			require.Equal(t, currentBytes, newCurrentBytes, "round-trip should preserve encoding")
		})
	}
}

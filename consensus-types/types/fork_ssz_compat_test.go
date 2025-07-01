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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// ForkSSZ is a reference implementation of Fork that uses only karalabe/ssz
// for testing compatibility between the two SSZ libraries
type ForkSSZ struct {
	PreviousVersion common.Version
	CurrentVersion  common.Version
	Epoch           math.Epoch
}

func (f *ForkSSZ) SizeSSZ(*ssz.Sizer) uint32 {
	return 16 // 4 + 4 + 8
}

func (f *ForkSSZ) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &f.PreviousVersion)
	ssz.DefineStaticBytes(codec, &f.CurrentVersion)
	ssz.DefineUint64(codec, &f.Epoch)
}

func (f *ForkSSZ) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(f))
	return buf, ssz.EncodeToBytes(buf, f)
}

func (f *ForkSSZ) HashTreeRoot() common.Root {
	return ssz.HashSequential(f)
}

// TestForkSSZCompatibility verifies that the fastssz implementation produces
// identical results to the karalabe/ssz implementation
func TestForkSSZCompatibility(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		fork *types.Fork
	}{
		{
			name: "zero values",
			fork: &types.Fork{
				PreviousVersion: common.Version{0, 0, 0, 0},
				CurrentVersion:  common.Version{0, 0, 0, 0},
				Epoch:           math.Epoch(0),
			},
		},
		{
			name: "typical values",
			fork: &types.Fork{
				PreviousVersion: common.Version{1, 2, 3, 4},
				CurrentVersion:  common.Version{5, 6, 7, 8},
				Epoch:           math.Epoch(1000),
			},
		},
		{
			name: "max values",
			fork: &types.Fork{
				PreviousVersion: common.Version{255, 255, 255, 255},
				CurrentVersion:  common.Version{255, 255, 255, 255},
				Epoch:           math.Epoch(^uint64(0)),
			},
		},
		{
			name: "real world example - deneb",
			fork: &types.Fork{
				PreviousVersion: common.Version{0x03, 0x00, 0x00, 0x00},
				CurrentVersion:  common.Version{0x04, 0x00, 0x00, 0x00},
				Epoch:           math.Epoch(269568),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create reference fork using karalabe/ssz
			forkSSZ := &ForkSSZ{
				PreviousVersion: tc.fork.PreviousVersion,
				CurrentVersion:  tc.fork.CurrentVersion,
				Epoch:           tc.fork.Epoch,
			}

			// Test serialization
			sszBytes, err := forkSSZ.MarshalSSZ()
			require.NoError(t, err)

			fastSSZBytes, err := tc.fork.MarshalSSZ()
			require.NoError(t, err)

			require.Equal(t, sszBytes, fastSSZBytes)

			sszRoot := forkSSZ.HashTreeRoot()
			fastSSZRoot := tc.fork.HashTreeRoot()
			require.Equal(t, sszRoot, fastSSZRoot)

			// Also test unmarshaling
			unmarshaledFork := new(types.Fork)
			err = unmarshaledFork.UnmarshalSSZ(sszBytes)
			require.NoError(t, err)
			require.Equal(t, tc.fork, unmarshaledFork, "Unmarshaling SSZ bytes should produce original fork")
		})
	}
}

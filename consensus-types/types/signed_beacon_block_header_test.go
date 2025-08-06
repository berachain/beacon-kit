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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/encoding/sszutil"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

func TestSignedBeaconBlockHeader_Serialization(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{0xde, 0xad, 0xbe, 0xef},
		common.Root{0xca, 0xca, 0xca, 0xfe},
		common.Root{0xde, 0xad, 0xca, 0xfe},
	)
	sig := crypto.BLSSignature{0xde, 0xad, 0xc4, 0xc4}
	orig := &types.SignedBeaconBlockHeader{
		Header:    header,
		Signature: sig,
	}

	data, err := orig.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalled := new(types.SignedBeaconBlockHeader)
	err = sszutil.Unmarshal(data, unmarshalled)
	require.NoError(t, err)
	require.Equal(t, orig, unmarshalled)

	// Test that MarshalSSZ works correctly
	buf, err := orig.MarshalSSZ()
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestSignedBeaconBlockHeader_EmptySerialization(t *testing.T) {
	t.Parallel()
	orig := &types.SignedBeaconBlockHeader{}
	data, err := orig.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalled := new(types.SignedBeaconBlockHeader)
	err = sszutil.Unmarshal(data, unmarshalled)
	require.NoError(t, err)
	require.NotNil(t, unmarshalled)
	require.NotNil(t, unmarshalled.GetHeader())
	require.NotNil(t, unmarshalled.GetSignature())
	require.Equal(t, &types.BeaconBlockHeader{}, unmarshalled.GetHeader())

	// Test that MarshalSSZ works correctly
	buf, err := orig.MarshalSSZ()
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestSignedBeaconBlockHeader_SizeSSZ(t *testing.T) {
	t.Parallel()
	sigHeader := types.NewSignedBeaconBlockHeader(
		types.NewBeaconBlockHeader(
			math.Slot(100),
			math.ValidatorIndex(200),
			common.Root{0xaa},
			common.Root{0xbb},
			common.Root{0xcc},
		),
		crypto.BLSSignature{0xff},
	)

	size := sigHeader.SizeSSZ()
	require.Equal(t, 208, size)
}

func TestSignedBeaconBlockHeader_HashTreeRoot(t *testing.T) {
	t.Parallel()
	sigHeader := types.NewSignedBeaconBlockHeader(
		types.NewBeaconBlockHeader(
			math.Slot(100),
			math.ValidatorIndex(200),
			common.Root{0xaa},
			common.Root{0xbb},
			common.Root{0xcc},
		),
		crypto.BLSSignature{0xff},
	)
	_, _ = sigHeader.HashTreeRoot()
}

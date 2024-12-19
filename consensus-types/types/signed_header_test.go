// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/primitives/math"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func TestSignedBeaconBlockHeader_Serialization(t *testing.T) {
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{0xde, 0xad, 0xbe, 0xef},
		common.Root{0xca, 0xca, 0xca, 0xfe},
		common.Root{0xde, 0xad, 0xca, 0xfe},
	)
	sig := crypto.BLSSignature{0xde, 0xad, 0xc4, 0xc4}
	orig := new(types.SignedBeaconBlockHeader).Empty()
	orig.SetHeader(header)
	orig.SetSignature(sig)

	data, err := orig.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)
	var unmarshalled types.SignedBeaconBlockHeader
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, orig, &unmarshalled)

	buf := make([]byte, karalabessz.Size(orig))
	err = karalabessz.EncodeToBytes(buf, orig)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestSignedBeaconBlockHeader_EmptySerialization(t *testing.T) {
	orig := new(types.SignedBeaconBlockHeader)
	data, err := orig.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.SignedBeaconBlockHeader
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.NotNil(t, unmarshalled.GetHeader())
	require.NotNil(t, unmarshalled.GetSignature())
	require.Equal(t, types.BeaconBlockHeader{}, *unmarshalled.GetHeader())

	buf := make([]byte, karalabessz.Size(orig))
	err = karalabessz.EncodeToBytes(buf, orig)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestSignedBeaconBlockHeader_SizeSSZ(t *testing.T) {
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

	size := karalabessz.Size(sigHeader)
	require.Equal(t, uint32(208), size)
}

func TestSignedBeaconBlockHeader_HashTreeRoot(_ *testing.T) {
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
	_ = sigHeader.HashTreeRoot()
}

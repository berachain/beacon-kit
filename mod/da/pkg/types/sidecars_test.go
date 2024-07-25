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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"testing"

	ctypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types/v2"
	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestEmptySidecarMarshalling(t *testing.T) {
	// Create an empty BlobSidecar
	sidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.BeaconBlockHeader{},
		&eip4844.Blob{},
		eip4844.KZGCommitment{},
		[48]byte{},
		[][32]byte{
			byteslib.ToBytes32([]byte("1")),
			byteslib.ToBytes32([]byte("2")),
			byteslib.ToBytes32([]byte("3")),
			byteslib.ToBytes32([]byte("4")),
			byteslib.ToBytes32([]byte("5")),
			byteslib.ToBytes32([]byte("6")),
			byteslib.ToBytes32([]byte("7")),
			byteslib.ToBytes32([]byte("8")),
		},
	)

	// Marshal the empty sidecar
	marshalled, err := sidecar.MarshalSSZ()
	require.NoError(
		t,
		err,
		"Marshalling empty sidecar should not produce an error",
	)
	require.NotNil(
		t,
		marshalled,
		"Marshalling empty sidecar should produce a result",
	)

	// Unmarshal the empty sidecar
	unmarshalled := &types.BlobSidecar{}
	err = unmarshalled.UnmarshalSSZ(marshalled)
	require.NoError(
		t,
		err,
		"Unmarshalling empty sidecar should not produce an error",
	)

	// Compare the original and unmarshalled empty sidecars
	require.Equal(
		t,
		sidecar,
		unmarshalled,
		"The original and unmarshalled empty sidecars should be equal",
	)
}

func TestValidateBlockRoots(t *testing.T) {
	// Create a sample BlobSidecar with valid roots
	validSidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.BeaconBlockHeader{
			StateRoot: [32]byte{1},
			BodyRoot:  [32]byte{2},
		},

		&eip4844.Blob{},
		[48]byte{},
		[48]byte{},
		[][32]byte{
			byteslib.ToBytes32([]byte("1")),
			byteslib.ToBytes32([]byte("2")),
			byteslib.ToBytes32([]byte("3")),
			byteslib.ToBytes32([]byte("4")),
			byteslib.ToBytes32([]byte("5")),
			byteslib.ToBytes32([]byte("6")),
			byteslib.ToBytes32([]byte("7")),
			byteslib.ToBytes32([]byte("8")),
		},
	)

	// Validate the sidecar with valid roots
	sidecars := types.BlobSidecars{
		Sidecars: []*types.BlobSidecar{validSidecar},
	}
	err := sidecars.ValidateBlockRoots()
	require.NoError(
		t,
		err,
		"Validating sidecar with valid roots should not produce an error",
	)

	// Create a sample BlobSidecar with invalid roots
	differentBlockRootSidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.BeaconBlockHeader{
			StateRoot: [32]byte{1},
			BodyRoot:  [32]byte{3},
		},
		&eip4844.Blob{},
		eip4844.KZGCommitment{},
		eip4844.KZGProof{},
		[][32]byte{
			byteslib.ToBytes32([]byte("1")),
			byteslib.ToBytes32([]byte("2")),
			byteslib.ToBytes32([]byte("3")),
			byteslib.ToBytes32([]byte("4")),
			byteslib.ToBytes32([]byte("5")),
			byteslib.ToBytes32([]byte("6")),
			byteslib.ToBytes32([]byte("7")),
			byteslib.ToBytes32([]byte("8")),
		},
	)

	// Validate the sidecar with invalid roots
	sidecarsInvalid := types.BlobSidecars{
		Sidecars: []*types.BlobSidecar{
			validSidecar,
			differentBlockRootSidecar,
		},
	}
	err = sidecarsInvalid.ValidateBlockRoots()
	require.Error(
		t,
		err,
		"Validating sidecar with invalid roots should produce an error",
	)
}

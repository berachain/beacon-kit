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

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/encoding/sszutil"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/math/log"
	"github.com/berachain/beacon-kit/primitives/version"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

func generateBeaconBlockBody(t *testing.T, v common.Version) types.BeaconBlockBody {
	versionable := types.NewVersionable(v)
	body := types.BeaconBlockBody{
		Versionable:  versionable,
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data:     &types.Eth1Data{},
		Graffiti:     [32]byte{4, 5, 6},
		Deposits:     []*types.Deposit{},
		ExecutionPayload: &types.ExecutionPayload{
			Versionable:   versionable,
			BaseFeePerGas: math.NewU256(0),
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
	body.SetProposerSlashings(types.ProposerSlashings{})
	body.SetAttesterSlashings(types.AttesterSlashings{})
	body.SetAttestations(types.Attestations{})
	body.SetSyncAggregate(&types.SyncAggregate{})
	body.SetVoluntaryExits(types.VoluntaryExits{})
	body.SetBlsToExecutionChanges(types.BlsToExecutionChanges{})
	if version.EqualsOrIsAfter(v, version.Electra()) {
		require.NoError(t, body.SetExecutionRequests(&types.ExecutionRequests{}))
	}
	return body
}

func TestBeaconBlockBodyBase(t *testing.T) {
	t.Parallel()
	body := types.BeaconBlockBody{
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data:     &types.Eth1Data{},
		Graffiti:     [32]byte{4, 5, 6},
		Deposits:     []*types.Deposit{},
	}

	require.Equal(t, bytes.B96{1, 2, 3}, body.GetRandaoReveal())
	require.NotNil(t, body.GetEth1Data())

	newGraffiti := [32]byte{7, 8, 9}
	body.SetGraffiti(newGraffiti)

	require.Equal(t, newGraffiti, [32]byte(body.GetGraffiti()))
	require.NotNil(t, body.GetDeposits())
}

func TestBeaconBlockBody(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		versionable := types.NewVersionable(v)
		body := types.BeaconBlockBody{
			Versionable:  versionable,
			RandaoReveal: [96]byte{1, 2, 3},
			Eth1Data:     &types.Eth1Data{},
			Graffiti:     [32]byte{4, 5, 6},
			Deposits:     []*types.Deposit{},
			ExecutionPayload: &types.ExecutionPayload{
				Versionable:   versionable,
				BaseFeePerGas: math.NewU256(0),
			},
			BlobKzgCommitments: []eip4844.KZGCommitment{},
		}
		require.NotNil(t, body.GetExecutionPayload())
		require.NotNil(t, body.GetBlobKzgCommitments())
		if version.EqualsOrIsAfter(v, version.Electra()) {
			require.Equal(t, types.BodyLengthElectra, body.Length())
		} else {
			require.Equal(t, types.BodyLengthDeneb, body.Length())
		}
	})
}

func TestBeaconBlockBody_SetBlobKzgCommitments(t *testing.T) {
	t.Parallel()
	body := types.BeaconBlockBody{}
	commitments := eip4844.KZGCommitments[common.ExecutionHash]{}
	body.SetBlobKzgCommitments(commitments)

	require.Equal(t, commitments, body.GetBlobKzgCommitments())
}

func TestBeaconBlockBody_SetRandaoReveal(t *testing.T) {
	t.Parallel()
	body := types.BeaconBlockBody{}
	randaoReveal := crypto.BLSSignature{1, 2, 3}
	body.SetRandaoReveal(randaoReveal)

	require.Equal(t, randaoReveal, body.GetRandaoReveal())
}

func TestBeaconBlockBody_SetEth1Data(t *testing.T) {
	t.Parallel()
	body := types.BeaconBlockBody{}
	eth1Data := &types.Eth1Data{}
	body.SetEth1Data(eth1Data)

	require.Equal(t, eth1Data, body.GetEth1Data())
}

func TestBeaconBlockBody_SetDeposits(t *testing.T) {
	t.Parallel()
	body := types.BeaconBlockBody{}
	deposits := types.Deposits{}
	body.SetDeposits(deposits)

	require.Equal(t, deposits, body.GetDeposits())
}

func TestBeaconBlockBody_MarshalSSZ(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		body := types.BeaconBlockBody{
			Versionable:        types.NewVersionable(v),
			RandaoReveal:       [96]byte{1, 2, 3},
			Eth1Data:           &types.Eth1Data{},
			Graffiti:           [32]byte{4, 5, 6},
			Deposits:           []*types.Deposit{},
			ExecutionPayload:   &types.ExecutionPayload{},
			BlobKzgCommitments: []eip4844.KZGCommitment{},
		}
		data, err := body.MarshalSSZ()

		require.NoError(t, err)
		require.NotNil(t, data)
	})
}

// TestBeaconBlockBody_FastSSZ tests the fastssz methods of BeaconBlockBody.
func TestBeaconBlockBody_FastSSZ(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		body := generateBeaconBlockBody(t, v)

		t.Run("MarshalSSZTo", func(t *testing.T) {
			dst := make([]byte, 0)
			result, err := body.MarshalSSZTo(dst)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Greater(t, len(result), 0)
		})

		t.Run("SizeSSZ", func(t *testing.T) {
			size := body.SizeSSZ()
			require.Greater(t, size, 0)
		})

		t.Run("HashTreeRootWith", func(t *testing.T) {
			hh := ssz.NewHasher()
			err := body.HashTreeRootWith(hh)
			require.NoError(t, err)
		})

		t.Run("GetTree", func(t *testing.T) {
			tree, err := body.GetTree()
			require.NoError(t, err)
			require.NotNil(t, tree)
		})
	})
}

// TestBeaconBlockBody_UnmarshalSSZ tests that UnmarshalSSZ returns an error for incorrect buffer size.
func TestBeaconBlockBody_UnmarshalSSZ(t *testing.T) {
	t.Parallel()
	body := &types.BeaconBlockBody{}
	err := body.UnmarshalSSZ([]byte{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "incorrect size")
}

func TestBeaconBlockBody_GetTopLevelRoots(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		body := generateBeaconBlockBody(t, v)
		roots, err := body.GetTopLevelRoots()
		require.NoError(t, err)
		require.NotNil(t, roots)
		if version.EqualsOrIsAfter(v, version.Electra()) {
			require.Equal(t, types.BodyLengthElectra, uint64(len(roots)))
		} else {
			require.Equal(t, types.BodyLengthDeneb, uint64(len(roots)))
		}
	})
}

func TestBeaconBlockBody_Empty(t *testing.T) {
	t.Parallel()
	body := types.BeaconBlockBody{}
	require.NotNil(t, body)
}

// Ensure that the ProposerSlashings field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedProposerSlashingsEnforcement(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		blockBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		// Test that marshaling with non-empty unused fields fails
		unused := common.UnusedType(1)
		blockBody.SetProposerSlashings(types.ProposerSlashings{&unused})
		_, err := blockBody.MarshalSSZ()
		require.Error(t, err)
		require.Contains(t, err.Error(), "ProposerSlashings must be unused")

		// Test that marshaling with empty unused fields succeeds
		blockBody.SetProposerSlashings(types.ProposerSlashings{})
		buf, err := blockBody.MarshalSSZ()
		require.NoError(t, err)

		// Test that unmarshaling enforces unused constraint
		unmarshalledBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		err = sszutil.Unmarshal(buf, unmarshalledBody)
		require.NoError(t, err)
	})
}

// Ensure that the AttesterSlashings field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedAttesterSlashingsEnforcement(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		blockBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		// Test that marshaling with non-empty unused fields fails
		unused := common.UnusedType(1)
		blockBody.SetAttesterSlashings(types.AttesterSlashings{&unused})
		_, err := blockBody.MarshalSSZ()
		require.Error(t, err)
		require.Contains(t, err.Error(), "AttesterSlashings must be unused")

		// Test that marshaling with empty unused fields succeeds
		blockBody.SetAttesterSlashings(types.AttesterSlashings{})
		buf, err := blockBody.MarshalSSZ()
		require.NoError(t, err)

		// Test that unmarshaling enforces unused constraint
		unmarshalledBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		err = sszutil.Unmarshal(buf, unmarshalledBody)
		require.NoError(t, err)
	})
}

// Ensure that the Attestations field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedAttestationsEnforcement(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		blockBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		// Test that marshaling with non-empty unused fields fails
		unused := common.UnusedType(1)
		blockBody.SetAttestations(types.Attestations{&unused})
		_, err := blockBody.MarshalSSZ()
		require.Error(t, err)
		require.Contains(t, err.Error(), "Attestations must be unused")

		// Test that marshaling with empty unused fields succeeds
		blockBody.SetAttestations(types.Attestations{})
		buf, err := blockBody.MarshalSSZ()
		require.NoError(t, err)

		// Test that unmarshaling enforces unused constraint
		unmarshalledBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		err = sszutil.Unmarshal(buf, unmarshalledBody)
		require.NoError(t, err)
	})
}

// Ensure that the VoluntaryExits field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedVoluntaryExitsEnforcement(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		blockBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		// Test that marshaling with non-empty unused fields fails
		unused := common.UnusedType(1)
		blockBody.SetVoluntaryExits(types.VoluntaryExits{&unused})
		_, err := blockBody.MarshalSSZ()
		require.Error(t, err)
		require.Contains(t, err.Error(), "VoluntaryExits must be unused")

		// Test that marshaling with empty unused fields succeeds
		blockBody.SetVoluntaryExits(types.VoluntaryExits{})
		buf, err := blockBody.MarshalSSZ()
		require.NoError(t, err)

		// Test that unmarshaling enforces unused constraint
		unmarshalledBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		err = sszutil.Unmarshal(buf, unmarshalledBody)
		require.NoError(t, err)
	})
}

// Ensure that the BlsToExecutionChanges field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedBlsToExecutionChangesEnforcement(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		blockBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		// Test that marshaling with non-empty unused fields fails
		unused := common.UnusedType(1)
		blockBody.SetBlsToExecutionChanges(types.BlsToExecutionChanges{&unused})
		_, err := blockBody.MarshalSSZ()
		require.Error(t, err)
		require.Contains(t, err.Error(), "BlsToExecutionChanges must be unused")

		// Test that marshaling with empty unused fields succeeds
		blockBody.SetBlsToExecutionChanges(types.BlsToExecutionChanges{})
		buf, err := blockBody.MarshalSSZ()
		require.NoError(t, err)

		// Test that unmarshaling enforces unused constraint
		unmarshalledBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		err = sszutil.Unmarshal(buf, unmarshalledBody)
		require.NoError(t, err)
	})
}

func TestBeaconBlockBody_RoundTrip_HashTreeRoot(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		body := generateBeaconBlockBody(t, v)
		data, err := body.MarshalSSZ()
		require.NoError(t, err)
		require.NotNil(t, data)

		unmarshalledBody := types.NewEmptyBeaconBlockBodyWithVersion(v)
		err = sszutil.Unmarshal(data, unmarshalledBody)
		require.NoError(t, err)
		bodyRoot, err := body.HashTreeRoot()
		require.NoError(t, err)
		unmarshalledRoot, err := unmarshalledBody.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, bodyRoot, unmarshalledRoot)
	})
}

// This test explains the calculation of the KZG commitment' inclusion proof depth.
func Test_KZGCommitmentInclusionProofDepth(t *testing.T) {
	t.Parallel()
	maxUint8 := uint64(^uint8(0))
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)

	// Depth of the partial BeaconBlockBody merkle tree. This is partial
	// because we only include as much as we need to prove the inclusion of
	// the KZG commitments.
	blockBodyMerkleDepth := uint64(log.ILog2Floor(uint64(types.KZGGeneralizedIndex)))
	require.Less(t, blockBodyMerkleDepth, maxUint8)

	// The depth of the merkle tree of the KZG Commitments, including the +1
	// for the length mixin.
	commitmentProofMerkleDepth := uint64(log.ILog2Ceil(cs.MaxBlobCommitmentsPerBlock())) + 1
	require.Less(t, commitmentProofMerkleDepth, maxUint8)

	// InclusionProofDepth is the combined depth of all of these things.
	expectedInclusionProofDepth := blockBodyMerkleDepth + commitmentProofMerkleDepth
	require.Less(t, expectedInclusionProofDepth, maxUint8)

	// Grab the inclusionProofDepth from beacon-kit.
	actualInclusionProofDepth := types.KZGInclusionProofDepth
	require.Less(t, uint64(actualInclusionProofDepth), maxUint8)

	require.Equal(t, uint8(expectedInclusionProofDepth), uint8(actualInclusionProofDepth))
}

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

package deneb_test

import (
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/deneb"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/math/log"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func generateBeaconBlockBody() deneb.BeaconBlockBody {
	body := deneb.BeaconBlockBody{
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data:     &deneb.Eth1Data{},
		Graffiti:     [32]byte{4, 5, 6},
		Deposits:     []*deneb.Deposit{},
		ExecutionPayload: &deneb.ExecutionPayload{
			BaseFeePerGas: math.NewU256(0),
			EpVersion:     version.Deneb1(),
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
	body.SetProposerSlashings(deneb.ProposerSlashings{})
	body.SetAttesterSlashings(deneb.AttesterSlashings{})
	body.SetAttestations(deneb.Attestations{})
	body.SetSyncAggregate(&deneb.SyncAggregate{})
	body.SetVoluntaryExits(deneb.VoluntaryExits{})
	body.SetBlsToExecutionChanges(deneb.BlsToExecutionChanges{})
	return body
}

func TestBeaconBlockBodyBase(t *testing.T) {
	t.Parallel()
	body := deneb.BeaconBlockBody{
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data:     &deneb.Eth1Data{},
		Graffiti:     [32]byte{4, 5, 6},
		Deposits:     []*deneb.Deposit{},
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
	body := deneb.BeaconBlockBody{
		RandaoReveal:       [96]byte{1, 2, 3},
		Eth1Data:           &deneb.Eth1Data{},
		Graffiti:           [32]byte{4, 5, 6},
		Deposits:           []*deneb.Deposit{},
		ExecutionPayload:   (&deneb.ExecutionPayload{}).Empty(version.Deneb1()),
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}

	require.False(t, body.IsNil())
	require.NotNil(t, body.GetExecutionPayload())
	require.NotNil(t, body.GetBlobKzgCommitments())
	require.Equal(t, deneb.BodyLengthDeneb, body.Length())
}

func TestBeaconBlockBody_SetBlobKzgCommitments(t *testing.T) {
	t.Parallel()
	body := deneb.BeaconBlockBody{}
	commitments := eip4844.KZGCommitments[common.ExecutionHash]{}
	body.SetBlobKzgCommitments(commitments)

	require.Equal(t, commitments, body.GetBlobKzgCommitments())
}

func TestBeaconBlockBody_SetRandaoReveal(t *testing.T) {
	t.Parallel()
	body := deneb.BeaconBlockBody{}
	randaoReveal := crypto.BLSSignature{1, 2, 3}
	body.SetRandaoReveal(randaoReveal)

	require.Equal(t, randaoReveal, body.GetRandaoReveal())
}

func TestBeaconBlockBody_SetEth1Data(t *testing.T) {
	t.Parallel()
	body := deneb.BeaconBlockBody{}
	eth1Data := &deneb.Eth1Data{}
	body.SetEth1Data(eth1Data)

	require.Equal(t, eth1Data, body.GetEth1Data())
}

func TestBeaconBlockBody_SetDeposits(t *testing.T) {
	t.Parallel()
	body := deneb.BeaconBlockBody{}
	deposits := deneb.Deposits{}
	body.SetDeposits(deposits)

	require.Equal(t, deposits, body.GetDeposits())
}

func TestBeaconBlockBody_MarshalSSZ(t *testing.T) {
	t.Parallel()
	body := deneb.BeaconBlockBody{
		RandaoReveal:       [96]byte{1, 2, 3},
		Eth1Data:           &deneb.Eth1Data{},
		Graffiti:           [32]byte{4, 5, 6},
		Deposits:           []*deneb.Deposit{},
		ExecutionPayload:   (&deneb.ExecutionPayload{}).Empty(version.Deneb1()),
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
	data, err := body.MarshalSSZ()

	require.NoError(t, err)
	require.NotNil(t, data)
}

func TestBeaconBlockBody_GetTopLevelRoots(t *testing.T) {
	t.Parallel()
	body := generateBeaconBlockBody()
	roots := body.GetTopLevelRoots()
	require.NotNil(t, roots)
}

func TestBeaconBlockBody_Empty(t *testing.T) {
	t.Parallel()
	body := deneb.BeaconBlockBody{}
	require.NotNil(t, body)
}

// Ensure that the ProposerSlashings field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedProposerSlashingsEnforcement(t *testing.T) {
	t.Parallel()
	blockBody := deneb.BeaconBlockBody{}
	unused := deneb.UnusedType(1)
	blockBody.SetProposerSlashings(deneb.ProposerSlashings{&unused})
	_, err := blockBody.MarshalSSZ()
	require.Error(t, err)

	buf := make([]byte, ssz.Size(&blockBody))
	err = ssz.EncodeToBytes(buf, &blockBody)
	require.NoError(t, err)

	unmarshalledBody := &deneb.BeaconBlockBody{}
	err = unmarshalledBody.UnmarshalSSZ(buf)
	require.ErrorContains(t, err, "must be unused")
}

// Ensure that the AttesterSlashings field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedAttesterSlashingsEnforcement(t *testing.T) {
	t.Parallel()
	blockBody := deneb.BeaconBlockBody{}
	unused := deneb.UnusedType(1)
	blockBody.SetAttesterSlashings(deneb.AttesterSlashings{&unused})
	_, err := blockBody.MarshalSSZ()
	require.Error(t, err)

	buf := make([]byte, ssz.Size(&blockBody))
	err = ssz.EncodeToBytes(buf, &blockBody)
	require.NoError(t, err)

	unmarshalledBody := &deneb.BeaconBlockBody{}
	err = unmarshalledBody.UnmarshalSSZ(buf)
	require.ErrorContains(t, err, "must be unused")
}

// Ensure that the Attestations field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedAttestationsEnforcement(t *testing.T) {
	t.Parallel()
	blockBody := deneb.BeaconBlockBody{}
	unused := deneb.UnusedType(1)
	blockBody.SetAttestations(deneb.Attestations{&unused})
	_, err := blockBody.MarshalSSZ()
	require.Error(t, err)

	buf := make([]byte, ssz.Size(&blockBody))
	err = ssz.EncodeToBytes(buf, &blockBody)
	require.NoError(t, err)

	unmarshalledBody := &deneb.BeaconBlockBody{}
	err = unmarshalledBody.UnmarshalSSZ(buf)
	require.ErrorContains(t, err, "must be unused")
}

// Ensure that the VoluntaryExits field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedVoluntaryExitsEnforcement(t *testing.T) {
	t.Parallel()
	blockBody := deneb.BeaconBlockBody{}
	unused := deneb.UnusedType(1)
	blockBody.SetVoluntaryExits(deneb.VoluntaryExits{&unused})
	_, err := blockBody.MarshalSSZ()
	require.Error(t, err)

	buf := make([]byte, ssz.Size(&blockBody))
	err = ssz.EncodeToBytes(buf, &blockBody)
	require.NoError(t, err)

	unmarshalledBody := &deneb.BeaconBlockBody{}
	err = unmarshalledBody.UnmarshalSSZ(buf)
	require.ErrorContains(t, err, "must be unused")
}

// Ensure that the BlsToExecutionChanges field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedBlsToExecutionChangesEnforcement(t *testing.T) {
	t.Parallel()
	blockBody := deneb.BeaconBlockBody{}
	unused := deneb.UnusedType(1)
	blockBody.SetBlsToExecutionChanges(deneb.BlsToExecutionChanges{&unused})
	_, err := blockBody.MarshalSSZ()
	require.Error(t, err)

	buf := make([]byte, ssz.Size(&blockBody))
	err = ssz.EncodeToBytes(buf, &blockBody)
	require.NoError(t, err)

	unmarshalledBody := &deneb.BeaconBlockBody{}
	err = unmarshalledBody.UnmarshalSSZ(buf)
	require.ErrorContains(t, err, "must be unused")
}

func TestBeaconBlockBody_RoundTrip_HashTreeRoot(t *testing.T) {
	t.Parallel()
	body := generateBeaconBlockBody()
	data, err := body.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalledBody := &deneb.BeaconBlockBody{}
	err = unmarshalledBody.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, body.HashTreeRoot(), unmarshalledBody.HashTreeRoot())
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
	blockBodyMerkleDepth := uint64(log.ILog2Floor(uint64(deneb.KZGGeneralizedIndex)))
	require.Less(t, blockBodyMerkleDepth, maxUint8)

	// The depth of the merkle tree of the KZG Commitments, including the +1
	// for the length mixin.
	commitmentProofMerkleDepth := uint64(log.ILog2Ceil(cs.MaxBlobCommitmentsPerBlock())) + 1
	require.Less(t, commitmentProofMerkleDepth, maxUint8)

	// InclusionProofDepth is the combined depth of all of these things.
	expectedInclusionProofDepth := blockBodyMerkleDepth + commitmentProofMerkleDepth
	require.Less(t, expectedInclusionProofDepth, maxUint8)

	// Grab the inclusionProofDepth from beacon-kit.
	actualInclusionProofDepth := deneb.KZGInclusionProofDepth
	require.Less(t, uint64(actualInclusionProofDepth), maxUint8)

	require.Equal(t, uint8(expectedInclusionProofDepth), uint8(actualInclusionProofDepth))
}

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
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/math/log"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func generateBeaconBlockBody(version common.Version) types.BeaconBlockBody {
	body := types.BeaconBlockBody{
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data:     &types.Eth1Data{},
		Graffiti:     [32]byte{4, 5, 6},
		Deposits:     []*types.Deposit{},
		ExecutionPayload: &types.ExecutionPayload{
			BaseFeePerGas: math.NewU256(0),
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
	body.SetForkVersion(version)
	body.GetExecutionPayload().SetForkVersion(version)
	body.SetProposerSlashings(types.ProposerSlashings{})
	body.SetAttesterSlashings(types.AttesterSlashings{})
	body.SetAttestations(types.Attestations{})
	body.SetSyncAggregate(&types.SyncAggregate{})
	body.SetVoluntaryExits(types.VoluntaryExits{})
	body.SetBlsToExecutionChanges(types.BlsToExecutionChanges{})
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
	body := types.BeaconBlockBody{
		RandaoReveal:       [96]byte{1, 2, 3},
		Eth1Data:           &types.Eth1Data{},
		Graffiti:           [32]byte{4, 5, 6},
		Deposits:           []*types.Deposit{},
		ExecutionPayload:   (&types.ExecutionPayload{}).Empty(version.Deneb1()),
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}

	require.False(t, body.IsNil())
	require.NotNil(t, body.GetExecutionPayload())
	require.NotNil(t, body.GetBlobKzgCommitments())
	require.Equal(t, types.BodyLengthDeneb, body.Length())
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
	body := types.BeaconBlockBody{
		RandaoReveal:       [96]byte{1, 2, 3},
		Eth1Data:           &types.Eth1Data{},
		Graffiti:           [32]byte{4, 5, 6},
		Deposits:           []*types.Deposit{},
		ExecutionPayload:   (&types.ExecutionPayload{}).Empty(version.Deneb1()),
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
	data, err := body.MarshalSSZ()

	require.NoError(t, err)
	require.NotNil(t, data)
}

func TestBeaconBlockBody_GetTopLevelRoots(t *testing.T) {
	t.Parallel()
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			body := generateBeaconBlockBody(v)
			roots := body.GetTopLevelRoots()
			require.NotNil(t, roots)
		})
	}
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
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			blockBody := types.BeaconBlockBody{}
			unused := types.UnusedType(1)
			blockBody.SetProposerSlashings(types.ProposerSlashings{&unused})
			_, err := blockBody.MarshalSSZ()
			require.Error(t, err)

			buf := make([]byte, ssz.Size(&blockBody))
			err = ssz.EncodeToBytes(buf, &blockBody)
			require.NoError(t, err)

			unmarshalledBody := &types.BeaconBlockBody{}
			err = unmarshalledBody.UnmarshalSSZ(buf, v)
			require.ErrorContains(t, err, "must be unused")
		})
	}
}

// Ensure that the AttesterSlashings field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedAttesterSlashingsEnforcement(t *testing.T) {
	t.Parallel()
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			blockBody := types.BeaconBlockBody{}
			unused := types.UnusedType(1)
			blockBody.SetAttesterSlashings(types.AttesterSlashings{&unused})
			_, err := blockBody.MarshalSSZ()
			require.Error(t, err)

			buf := make([]byte, ssz.Size(&blockBody))
			err = ssz.EncodeToBytes(buf, &blockBody)
			require.NoError(t, err)

			unmarshalledBody := &types.BeaconBlockBody{}
			err = unmarshalledBody.UnmarshalSSZ(buf, v)
			require.ErrorContains(t, err, "must be unused")
		})
	}
}

// Ensure that the Attestations field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedAttestationsEnforcement(t *testing.T) {
	t.Parallel()
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			blockBody := types.BeaconBlockBody{}
			unused := types.UnusedType(1)
			blockBody.SetAttestations(types.Attestations{&unused})
			_, err := blockBody.MarshalSSZ()
			require.Error(t, err)

			buf := make([]byte, ssz.Size(&blockBody))
			err = ssz.EncodeToBytes(buf, &blockBody)
			require.NoError(t, err)

			unmarshalledBody := &types.BeaconBlockBody{}
			err = unmarshalledBody.UnmarshalSSZ(buf, v)
			require.ErrorContains(t, err, "must be unused")
		})
	}
}

// Ensure that the VoluntaryExits field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedVoluntaryExitsEnforcement(t *testing.T) {
	t.Parallel()
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			blockBody := types.BeaconBlockBody{}
			unused := types.UnusedType(1)
			blockBody.SetVoluntaryExits(types.VoluntaryExits{&unused})
			_, err := blockBody.MarshalSSZ()
			require.Error(t, err)

			buf := make([]byte, ssz.Size(&blockBody))
			err = ssz.EncodeToBytes(buf, &blockBody)
			require.NoError(t, err)

			unmarshalledBody := &types.BeaconBlockBody{}
			err = unmarshalledBody.UnmarshalSSZ(buf, v)
			require.ErrorContains(t, err, "must be unused")
		})
	}
}

// Ensure that the BlsToExecutionChanges field cannot be unmarshaled with data in it,
// enforcing that it's unused.
func TestBeaconBlockBody_UnusedBlsToExecutionChangesEnforcement(t *testing.T) {
	t.Parallel()
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			blockBody := types.BeaconBlockBody{}
			unused := types.UnusedType(1)
			blockBody.SetBlsToExecutionChanges(types.BlsToExecutionChanges{&unused})
			_, err := blockBody.MarshalSSZ()
			require.Error(t, err)

			buf := make([]byte, ssz.Size(&blockBody))
			err = ssz.EncodeToBytes(buf, &blockBody)
			require.NoError(t, err)

			unmarshalledBody := &types.BeaconBlockBody{}
			err = unmarshalledBody.UnmarshalSSZ(buf, v)
			require.ErrorContains(t, err, "must be unused")
		})
	}
}

func TestBeaconBlockBody_RoundTrip_HashTreeRoot(t *testing.T) {
	t.Parallel()
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			body := generateBeaconBlockBody(v)
			data, err := body.MarshalSSZ()
			require.NoError(t, err)
			require.NotNil(t, data)

			unmarshalledBody := &types.BeaconBlockBody{}
			// We must set the version first for correct marshalling
			unmarshalledBody.SetForkVersion(v)
			err = unmarshalledBody.UnmarshalSSZ(data, v)
			require.NoError(t, err)
			require.Equal(t, body.HashTreeRoot(), unmarshalledBody.HashTreeRoot())
		})
	}
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

func TestBeaconBlockBody_ExecutionRequests(t *testing.T) {
	t.Parallel()
	body := generateBeaconBlockBody(version.Electra())
	err := body.SetExecutionRequests(&types.ExecutionRequests{
		Deposits: []*types.DepositRequest{
			{
				Pubkey:                bytes.B48{0, 1, 2},
				WithdrawalCredentials: types.WithdrawalCredentials(common.Bytes32{0, 1, 2}),
				Amount:                math.Gwei(1000),
				Signature:             bytes.B96{0, 1, 2},
				Index:                 69,
			},
			{
				Pubkey:                bytes.B48{0, 3, 4},
				WithdrawalCredentials: types.WithdrawalCredentials(common.Bytes32{0, 1, 2}),
				Amount:                math.Gwei(1000),
				Signature:             bytes.B96{0, 1, 2},
				Index:                 70,
			},
		},
		Withdrawals: []*types.WithdrawalRequest{
			{
				SourceAddress:   common.NewExecutionAddressFromHex("0xFF00000000000000000000000000000000000010"),
				ValidatorPubKey: bytes.B48{0, 1, 2},
				Amount:          math.Gwei(1000),
			},
		},
	})
	require.NoError(t, err)
	originalRequests, err := body.GetExecutionRequests()
	require.NoError(t, err)
	require.NotNil(t, originalRequests)
	data, err := body.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalledBody := &types.BeaconBlockBody{}
	err = unmarshalledBody.UnmarshalSSZ(data, body.GetForkVersion())
	require.NoError(t, err)
	executionRequests, err := unmarshalledBody.GetExecutionRequests()
	require.NoError(t, err)
	require.NotNil(t, executionRequests)
	require.Equal(t, body.HashTreeRoot(), unmarshalledBody.HashTreeRoot())
	require.Equal(t, executionRequests.Deposits[0], originalRequests.Deposits[0])
	require.Equal(t, executionRequests.Deposits[1], originalRequests.Deposits[1])
	require.Equal(t, executionRequests.Withdrawals[0], originalRequests.Withdrawals[0])
}

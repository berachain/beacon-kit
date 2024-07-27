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

package blob_test

// TODO: Create a mock such that core/types doesn't need
// to be imported here.

// MockSpec is a mock implementation of the ChainSpec interface used for
// testing.
type MockSpec struct{}

// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments per
// block.
// This mock implementation always returns 16.
func (m *MockSpec) MaxBlobCommitmentsPerBlock() uint64 {
	return 16
}

// TODO: Re-enable once we can easily decouple from core/types.
// func TestBuildKZGInclusionProof(t *testing.T) {
// 	chainspec := &MockSpec{}
// 	factory := da.NewSidecarFactory[da.BeaconBlockBody](
// 		chainspec,
// 		5,
// 	)
// 	body := mockBody()
// 	// Test for a valid index
//
// 	index := uint64(0)
// 	proof, err := factory.BuildKZGInclusionProof(body, index)
// 	require.NoError(
// 		t,
// 		err,
// 		"Building KZG inclusion proof should not produce an error",
// 	)
// 	require.NotNil(t, proof, "Proof should not be nil")

// 	bodyRoot, err := body.HashTreeRoot()
// 	require.NoError(t, err, "Hashing the body should not produce an error")

// 	// Verify the valid KZG inclusion proof
// 	validProof := merkle.VerifyProof(
// 		bodyRoot,
// 		body.GetBlobKzgCommitments()[index].ToHashChunks()[0],
// 		types.KZGOffset(chainspec.MaxBlobCommitmentsPerBlock())+index,
// 		proof,
// 	)
// 	require.True(t, validProof, "The KZG inclusion proof should be valid")

// 	// Test for an invalid index
// 	invalidIndex := uint64(100) // Assuming this is out of range
// 	_, err = factory.BuildKZGInclusionProof(body, invalidIndex)
// 	require.Error(
// 		t,
// 		err,
// 		"Building KZG inclusion proof with invalid index should produce an error",
// 	)

// 	require.True(t, validProof, "The KZG inclusion proof should be valid")

// 	// Attempt to verify the invalid KZG inclusion proof and expect failure
// 	invalidProof, err := factory.BuildKZGInclusionProof(body, invalidIndex)
// 	require.Error(
// 		t,
// 		err,
// 		"Building KZG inclusion proof should produce an error",
// 	)
// 	validInvalidProof := merkle.VerifyProof(
// 		bodyRoot,
// 		body.GetBlobKzgCommitments()[index].ToHashChunks()[0],
// 		types.KZGOffset(chainspec.MaxBlobCommitmentsPerBlock())+index,
// 		invalidProof,
// 	)
// 	require.False(
// 		t,
// 		validInvalidProof,
// 		"The KZG inclusion proof for an invalid index should be invalid",
// 	)
// }

// func mockBody() da.BeaconBlockBody {
// 	// Create a real ExecutionPayloadDeneb and BeaconBlockBody
// 	executionPayload := &engineprimitives.ExecutionPayload{
// 		ParentHash:    common.HexToHash("0x01"),
// 		FeeRecipient:  common.HexToAddress("0x02"),
// 		StateRoot:     common.HexToHash("0x03"),
// 		ReceiptsRoot:  common.HexToHash("0x04"),
// 		LogsBloom:     bytes.Repeat([]byte("b"), 256),
// 		Random:        common.HexToHash("0x05"),
// 		BaseFeePerGas: math.Wei(bytes.Repeat([]byte("f"), 32)),
// 		BlockHash:     common.HexToHash("0x06"),
// 		Transactions:  [][]byte{[]byte("tx1"), []byte("tx2")},
// 		ExtraData:     []byte("extra"),
// 	}

// 	return &types.BeaconBlockBodyDeneb{
// 		RandaoReveal: [96]byte{0x01},
// 		Eth1Data: &primitives.Eth1Data{
// 			DepositRoot:  common.Root{},
// 			DepositCount: 0,
// 			BlockHash:    common.ZeroHash,
// 		},
// 		ExecutionPayload: executionPayload,
// 		BlobKzgCommitments: kzg.Commitments{
// 			[48]byte(bytes.Repeat([]byte{0x01}, 48)),
// 			[48]byte(bytes.Repeat([]byte{0x10}, 48)),
// 			[48]byte(bytes.Repeat([]byte{0x11}, 48)),
// 		},
// 	}
// }

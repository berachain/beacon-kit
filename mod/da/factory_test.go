// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package da_test

import (
	"bytes"
	"testing"

	// TODO: Create a mock such that core/types doesn't need
	// to be imported here.
	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/da"
	"github.com/berachain/beacon-kit/mod/merkle"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// MockSpec is a mock implementation of the ChainSpec interface used for
// testing.
type MockSpec struct{}

// MaxBlobCommitmentsPerBlock returns the maximum number of blob commitments per
// block.
// This mock implementation always returns 16.
func (m *MockSpec) MaxBlobCommitmentsPerBlock() uint64 {
	return 16
}
func TestBuildKZGInclusionProof(t *testing.T) {
	chainspec := &MockSpec{}
	factory := da.NewSidecarFactory[da.BeaconBlockBody](
		chainspec,
		4,
	)
	body := mockBody()
	// Test for a valid index
	index := uint64(0)
	proof, err := factory.BuildKZGInclusionProof(body, index)
	require.NoError(
		t,
		err,
		"Building KZG inclusion proof should not produce an error",
	)
	require.NotNil(t, proof, "Proof should not be nil")

	bodyRoot, err := body.HashTreeRoot()
	require.NoError(t, err, "Hashing the body should not produce an error")

	// Verify the valid KZG inclusion proof
	validProof := merkle.VerifyProof(
		bodyRoot,
		body.GetBlobKzgCommitments()[index].ToHashChunks()[0],
		types.KZGOffset(chainspec.MaxBlobCommitmentsPerBlock())+index,
		proof,
	)
	require.True(t, validProof, "The KZG inclusion proof should be valid")

	// Test for an invalid index
	invalidIndex := uint64(100) // Assuming this is out of range
	_, err = factory.BuildKZGInclusionProof(body, invalidIndex)
	require.Error(
		t,
		err,
		"Building KZG inclusion proof with invalid index should produce an error",
	)

	require.True(t, validProof, "The KZG inclusion proof should be valid")

	// Attempt to verify the invalid KZG inclusion proof and expect failure
	invalidProof, err := factory.BuildKZGInclusionProof(body, invalidIndex)
	require.Error(
		t,
		err,
		"Building KZG inclusion proof should produce an error",
	)
	validInvalidProof := merkle.VerifyProof(
		bodyRoot,
		body.GetBlobKzgCommitments()[index].ToHashChunks()[0],
		types.KZGOffset(chainspec.MaxBlobCommitmentsPerBlock())+index,
		invalidProof,
	)
	require.False(
		t,
		validInvalidProof,
		"The KZG inclusion proof for an invalid index should be invalid",
	)
}

func mockBody() da.BeaconBlockBody {
	// Create a real ExecutionPayloadDeneb and BeaconBlockBody
	executionPayload := &engineprimitives.ExecutableDataDeneb{
		ParentHash:    common.HexToHash("0x01"),
		FeeRecipient:  common.HexToAddress("0x02"),
		StateRoot:     common.HexToHash("0x03"),
		ReceiptsRoot:  common.HexToHash("0x04"),
		LogsBloom:     bytes.Repeat([]byte("b"), 256),
		Random:        common.HexToHash("0x05"),
		BaseFeePerGas: primitives.Wei(bytes.Repeat([]byte("f"), 32)),
		BlockHash:     common.HexToHash("0x06"),
		Transactions:  [][]byte{[]byte("tx1"), []byte("tx2")},
		ExtraData:     []byte("extra"),
	}

	return &types.BeaconBlockBodyDeneb{
		RandaoReveal:     [96]byte{0x01},
		ExecutionPayload: executionPayload,
		BlobKzgCommitments: kzg.Commitments{
			[48]byte(bytes.Repeat([]byte{0x01}, 48)),
			[48]byte(bytes.Repeat([]byte{0x10}, 48)),
			[48]byte(bytes.Repeat([]byte{0x11}, 48)),
		},
	}
}

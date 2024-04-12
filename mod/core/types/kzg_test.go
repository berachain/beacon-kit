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

package types_test

import (
	"bytes"
	"testing"

	"github.com/berachain/beacon-kit/mod/core/types"
	enginetypes "github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/merkle"
	"github.com/berachain/beacon-kit/mod/merkle/htr"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/gohashtree"
	"github.com/stretchr/testify/require"
)

const LogMaxBlobCommitments = 4

func mockBody() *types.BeaconBlockBodyDeneb {
	// Create a real ExecutionPayloadDeneb and BeaconBlockBody
	executionPayload := &enginetypes.ExecutableDataDeneb{
		ParentHash:    common.HexToHash("0x01"),
		FeeRecipient:  common.HexToAddress("0x02"),
		StateRoot:     common.HexToHash("0x03"),
		ReceiptsRoot:  common.HexToHash("0x04"),
		LogsBloom:     bytes.Repeat([]byte("b"), 256),
		Random:        common.HexToHash("0x05"),
		BaseFeePerGas: bytes.Repeat([]byte("f"), 32),
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

// This test explains the calculation of the
// KZG commitment root's Merkle index
// in the Body's Merkle tree based on the
// index of the KZG commitment list in the Body.
func Test_KZGRootIndex(t *testing.T) {
	// Level of the KZG commitment root's parent.
	kzgParentRootLevel := types.LogBodyLengthDeneb
	// Merkle index of the KZG commitment root's parent.
	// The parent's left child is the KZG commitment root,
	// and its right child is the KZG commitment size.
	kzgParentRootIndex := types.KZGPositionDeneb + (1 << kzgParentRootLevel)
	// The KZG commitment root is the left child of its parent.
	// Its Merkle index is the double of its parent's Merkle index.
	require.Equal(t, uint64(types.KZGMerkleIndex), (2 * kzgParentRootIndex))
}

func Test_BodyProof(t *testing.T) {
	// Create a real ExecutionPayloadDeneb and BeaconBlockBody
	body := mockBody()

	// The body has the commitments.
	commitments := body.GetBlobKzgCommitments()

	// Generate leaves from commitments
	leaves := types.LeavesFromCommitments(commitments)

	// Calculate the depth the given tree will have.
	// depth := uint64(math.Ceil(math.Sqrt(float64(len(commitments)))))
	depth := uint8(LogMaxBlobCommitments)

	// Generate a sparse Merkle tree from the leaves.
	tree, err := merkle.NewTreeFromLeavesWithDepth(leaves, depth)
	require.NoError(t, err, "Failed to generate tree from items")

	// Get the root of the tree.
	root, err := tree.HashTreeRoot()
	require.NoError(t, err, "Failed to generate root hash")

	// Generate a proof for the index.
	var proof [][32]byte
	for index := uint64(0); index < uint64(len(leaves)); index++ {
		proof, err = tree.MerkleProofWithMixin(index)
		require.NoError(t, err, "Failed to generate Merkle proof")
		require.NotNil(t, proof, "Merkle proof should not be nil")
		require.Len(t, proof, int(depth)+1)

		// Verify the Merkle proof
		valid := merkle.VerifyProof(
			root,
			leaves[index],
			index,
			proof,
		)
		require.True(t, valid, "Merkle proof should be valid")

		proof, err = types.BodyProof(commitments, index)
		require.NoError(t, err, "Failed to generate Merkle proof")
		require.NotNil(t, proof, "Merkle proof should not be nil")
		require.Len(t, proof, int(depth)+1)

		valid = merkle.VerifyProof(
			root,
			leaves[index],
			index,
			proof,
		)
		require.True(t, valid, "Merkle proof should be valid")
	}
}

func Test_TopLevelRoots(t *testing.T) {
	body := mockBody()

	// Commitments
	commitments := body.GetBlobKzgCommitments()
	commitmentsRoot, err := htr.ListSSZ[kzg.Commitment](
		commitments, types.MaxBlobCommitmentsPerBlock)
	require.NoError(t, err, "Failed to generate root hash")

	// Body
	bodyMembersRoots, err := types.GetTopLevelRoots(body)
	require.NoError(t, err, "Failed to get top level roots")
	// Add the commitments root to the body members roots.
	// For this test only. We don't need to do this when
	// generating the proof.
	bodyMembersRoots[types.KZGPositionDeneb] = commitmentsRoot
	bodyTree, err := merkle.NewTreeFromLeavesWithDepth(
		bodyMembersRoots,
		types.LogBodyLengthDeneb,
	)
	require.NoError(t, err, "Failed to generate tree from member roots")
	bodySparseRoot, err := bodyTree.HashTreeRoot()
	require.NoError(t, err, "Failed to generate root hash")

	topProof, err := bodyTree.MerkleProofWithMixin(types.KZGPositionDeneb)
	require.NoError(t, err, "Failed to generate Merkle proof")

	// Verify the Merkle proof
	valid := merkle.VerifyProof(
		bodySparseRoot,
		commitmentsRoot,
		types.KZGPositionDeneb,
		topProof,
	)
	require.True(t, valid, "Merkle proof should be valid")
}

func Test_MerkleProofKZGCommitment(t *testing.T) {
	kzgs := kzg.Commitments{
		[48]byte(bytes.Repeat([]byte{0x01}, 48)),
		[48]byte(bytes.Repeat([]byte{0x10}, 48)),
		[48]byte(bytes.Repeat([]byte{0x11}, 48)),
	}
	body := mockBody()

	index := uint64(1)
	proof, err := types.MerkleProofKZGCommitment(body, index)
	require.NoError(t, err)
	require.Len(t,
		proof,
		int(LogMaxBlobCommitments)+1+int(types.LogBodyLengthDeneb))

	chunk := make([][32]byte, 2)
	copy(chunk[0][:], kzgs[index][:])
	copy(chunk[1][:], kzgs[index][32:])
	gohashtree.HashChunks(chunk, chunk)
	root, err := body.HashTreeRoot()
	require.NoError(t, err)

	commitments := body.GetBlobKzgCommitments()
	commitmentsRoot, err := htr.ListSSZ[kzg.Commitment](
		commitments, types.MaxBlobCommitmentsPerBlock)
	require.NoError(t, err, "Failed to generate root hash")

	require.True(t,
		merkle.VerifyProof(
			commitmentsRoot,
			chunk[0],
			index,
			proof[:LogMaxBlobCommitments+1],
		),
	)

	// Body
	bodyMembersRoots, err := types.GetTopLevelRoots(body)
	require.NoError(t, err, "Failed to get top level roots")
	// Add the commitments root to the body members roots.
	// For this test only. We don't need to do this when
	// generating the proof.
	// bodyMembersRoots[types.KZGPositionDeneb] = commitmentsRoot[:]
	bodyTree, err := merkle.NewTreeFromLeavesWithDepth(
		bodyMembersRoots,
		types.LogBodyLengthDeneb,
	)
	require.NoError(t, err, "Failed to generate tree from member roots")
	topProof, err := bodyTree.MerkleProofWithMixin(types.KZGPositionDeneb)
	require.NoError(t, err, "Failed to generate Merkle proof")
	require.Equal(t,
		topProof[:len(topProof)-1],
		proof[LogMaxBlobCommitments+1:],
	)

	require.Len(t,
		proof[LogMaxBlobCommitments+1:],
		int(types.LogBodyLengthDeneb),
	)
	require.True(t,
		merkle.VerifyProof(
			root,
			commitmentsRoot,
			types.KZGPositionDeneb,
			proof[LogMaxBlobCommitments+1:],
		),
	)

	require.True(t,
		merkle.VerifyProof(
			root,
			chunk[0],
			index+types.KZGOffset,
			proof,
		),
	)
}

// func Test_MerkleProofKZGCommitment(t *testing.T) {
// 	body := mockBody()

// 	blk := &types.BeaconBlockDeneb{
// 		Slot:            1,
// 		ProposerIndex:   1,
// 		ParentBlockRoot: common.HexToHash("0x07"),
// 		Body:            body,
// 	}

// 	kzgs := body.GetBlobKzgCommitments()
// 	index := 1
// 	_, err := types.MerkleProofKZGCommitment(
// 		blk,
// 		len(kzgs)+1,
// 	)
// 	require.NotNil(t, err)
// 	proof, err := types.MerkleProofKZGCommitment(blk, index)
// 	require.NoError(t, err)

// 	chunk := make([][32]byte, 2)
// 	copy(chunk[0][:], kzgs[index][:])
// 	copy(chunk[1][:], kzgs[index][32:])
// 	gohashtree.HashChunks(chunk, chunk)
// 	root, err := body.HashTreeRoot()
// 	require.NoError(t, err)
// 	kzgOffset := 54 * 4096
// 	for i := 0; i <= 54; i++ {
// 		ok := tree.VerifyProof(
// 			root[:],
// 			chunk[0][:],
// 			uint64(index+i*4096),
// 			proof,
// 		)
// 		fmt.Printf("i: %d, ok: %v\n", i, ok)
// 	}

// 	require.True(
// 		t,
// 		tree.VerifyProof(
// 			root[:],
// 			chunk[0][:],
// 			uint64(index+kzgOffset),
// 			proof,
// 		),
// 	)
// }

// func Benchmark_MerkleProofKZGCommitment(b *testing.B) {
// 	kzgs := make([][]byte, 3)
// 	kzgs[0] = make([]byte, 48)
// 	_, err := rand.Read(kzgs[0])
// 	require.NoError(b, err)
// 	kzgs[1] = make([]byte, 48)
// 	_, err = rand.Read(kzgs[1])
// 	require.NoError(b, err)
// 	kzgs[2] = make([]byte, 48)
// 	_, err = rand.Read(kzgs[2])
// 	require.NoError(b, err)
// 	pbBody := &ethpb.BeaconBlockBodyDeneb{
// 		SyncAggregate: &ethpb.SyncAggregate{
// 			SyncCommitteeBits:      make([]byte,
// fieldparams.SyncAggregateSyncCommitteeBytesLength),
// 			SyncCommitteeSignature: make([]byte, fieldparams.BLSSignatureLength),
// 		},
// 		ExecutionPayload: &enginev1.ExecutionPayloadDeneb{
// 			ParentHash:    make([]byte, fieldparams.RootLength),
// 			FeeRecipient:  make([]byte, 20),
// 			StateRoot:     make([]byte, fieldparams.RootLength),
// 			ReceiptsRoot:  make([]byte, fieldparams.RootLength),
// 			LogsBloom:     make([]byte, 256),
// 			PrevRandao:    make([]byte, fieldparams.RootLength),
// 			BaseFeePerGas: make([]byte, fieldparams.RootLength),
// 			BlockHash:     make([]byte, fieldparams.RootLength),
// 			Transactions:  make([][]byte, 0),
// 			ExtraData:     make([]byte, 0),
// 		},
// 		Eth1Data: &ethpb.Eth1Data{
// 			DepositRoot: make([]byte, fieldparams.RootLength),
// 			BlockHash:   make([]byte, fieldparams.RootLength),
// 		},
// 		BlobKzgCommitments: kzgs,
// 	}

// 	body, err := NewBeaconBlockBody(pbBody)
// 	require.NoError(b, err)
// 	index := 1
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_, err := MerkleProofKZGCommitment(body, index)
// 		require.NoError(b, err)
// 	}
// }

// func Test_VerifyKZGInclusionProof(t *testing.T) {
// 	kzgs := make([][]byte, 3)
// 	kzgs[0] = make([]byte, 48)
// 	_, err := rand.Read(kzgs[0])
// 	require.NoError(t, err)
// 	kzgs[1] = make([]byte, 48)
// 	_, err = rand.Read(kzgs[1])
// 	require.NoError(t, err)
// 	kzgs[2] = make([]byte, 48)
// 	_, err = rand.Read(kzgs[2])
// 	require.NoError(t, err)
// 	pbBody := &ethpb.BeaconBlockBodyDeneb{
// 		SyncAggregate: &ethpb.SyncAggregate{
// 			SyncCommitteeBits:      make([]byte,
// fieldparams.SyncAggregateSyncCommitteeBytesLength),
// 			SyncCommitteeSignature: make([]byte, fieldparams.BLSSignatureLength),
// 		},
// 		ExecutionPayload: &enginev1.ExecutionPayloadDeneb{
// 			ParentHash:    make([]byte, fieldparams.RootLength),
// 			FeeRecipient:  make([]byte, 20),
// 			StateRoot:     make([]byte, fieldparams.RootLength),
// 			ReceiptsRoot:  make([]byte, fieldparams.RootLength),
// 			LogsBloom:     make([]byte, 256),
// 			PrevRandao:    make([]byte, fieldparams.RootLength),
// 			BaseFeePerGas: make([]byte, fieldparams.RootLength),
// 			BlockHash:     make([]byte, fieldparams.RootLength),
// 			Transactions:  make([][]byte, 0),
// 			ExtraData:     make([]byte, 0),
// 		},
// 		Eth1Data: &ethpb.Eth1Data{
// 			DepositRoot: make([]byte, fieldparams.RootLength),
// 			BlockHash:   make([]byte, fieldparams.RootLength),
// 		},
// 		BlobKzgCommitments: kzgs,
// 	}

// 	body, err := NewBeaconBlockBody(pbBody)
// 	require.NoError(t, err)
// 	root, err := body.HashTreeRoot()
// 	require.NoError(t, err)
// 	index := 1
// 	proof, err := MerkleProofKZGCommitment(body, index)
// 	require.NoError(t, err)

// 	header := &ethpb.BeaconBlockHeader{
// 		BodyRoot:   root[:],
// 		ParentRoot: make([]byte, 32),
// 		StateRoot:  make([]byte, 32),
// 	}
// 	signedHeader := &ethpb.SignedBeaconBlockHeader{
// 		Header: header,
// 	}
// 	sidecar := &ethpb.BlobSidecar{
// 		Index:                    uint64(index),
// 		KzgCommitment:            kzgs[index],
// 		CommitmentInclusionProof: proof,
// 		SignedBlockHeader:        signedHeader,
// 	}
// 	blob, err := NewROBlob(sidecar)
// 	require.NoError(t, err)
// 	require.NoError(t, VerifyKZGInclusionProof(blob))
// 	proof[2] = make([]byte, 32)
// 	require.ErrorIs(t, errInvalidInclusionProof, VerifyKZGInclusionProof(blob))
// }

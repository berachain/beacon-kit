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

	"github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/trie"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/gohashtree"
	"github.com/stretchr/testify/require"
)

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
		BlobKzgCommitments: [][48]byte{
			[48]byte(bytes.Repeat([]byte("1"), 48)),
			[48]byte(bytes.Repeat([]byte("2"), 48)),
			[48]byte(bytes.Repeat([]byte("3"), 48)),
		},
	}
}

func Test_BodyProof(t *testing.T) {
	// Create a real ExecutionPayloadDeneb and BeaconBlockBody
	body := mockBody()

	// The body has the commitments.
	commitments := body.GetBlobKzgCommitments()

	// Generate leaves from commitments
	leaves := types.LeavesFromCommitments(commitments)

	// Calculate the depth the given trie will have.
	// depth := uint64(math.Ceil(math.Sqrt(float64(len(commitments)))))
	depth := types.LogMaxBlobCommitments

	// Generate a sparse Merkle tree from the leaves.
	sparse, err := trie.NewFromItems(leaves, depth)
	require.NoError(t, err, "Failed to generate trie from items")
	require.Equal(t, len(leaves), sparse.NumOfItems())

	// Get the root of the tree.
	root, err := sparse.HashTreeRoot()
	require.NoError(t, err, "Failed to generate root hash")

	// Generate a proof for the index.
	var proof [][]byte
	for index := uint64(0); index < uint64(len(leaves)); index++ {
		proof, err = sparse.MerkleProof(index)
		require.NoError(t, err, "Failed to generate Merkle proof")
		require.NotNil(t, proof, "Merkle proof should not be nil")
		require.Len(t, proof, int(depth)+1)

		// Verify the Merkle proof
		valid := trie.VerifyMerkleProof(
			root[:],
			leaves[index],
			index,
			proof,
		)
		require.True(t, valid, "Merkle proof should be valid")

		proof, err = types.BodyProof(commitments, index)
		require.NoError(t, err, "Failed to generate Merkle proof")
		require.NotNil(t, proof, "Merkle proof should not be nil")
		require.Len(t, proof, int(depth)+1)

		valid = trie.VerifyMerkleProof(
			root[:],
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
	commitmentsRoot, err := types.GetBlobKzgCommitmentsRoot(commitments)
	require.NoError(t, err, "Failed to generate root hash")

	// Body
	bodyMembersRoots, err := types.GetTopLevelRoots(body)
	require.NoError(t, err, "Failed to get top level roots")
	// Add the commitments root to the body members roots.
	// For this test only. We don't need to do this when
	// generating the proof.
	bodyMembersRoots[types.KZGPosition] = commitmentsRoot[:]
	bodySparse, err := trie.NewFromItems(
		bodyMembersRoots,
		types.LogBodyLength,
	)
	require.NoError(t, err, "Failed to generate trie from member roots")
	bodySparseRoot, err := bodySparse.HashTreeRoot()
	require.NoError(t, err, "Failed to generate root hash")
	require.Equal(t, types.BodyLength, bodySparse.NumOfItems())

	topProof, err := bodySparse.MerkleProof(types.KZGPosition)
	require.NoError(t, err, "Failed to generate Merkle proof")

	// Verify the Merkle proof
	valid := trie.VerifyMerkleProof(
		bodySparseRoot[:],
		commitmentsRoot[:],
		uint64(types.KZGPosition),
		topProof,
	)
	require.True(t, valid, "Merkle proof should be valid")
}

func Test_MerkleProofKZGCommitment(t *testing.T) {
	kzgs := [][48]byte{
		[48]byte(bytes.Repeat([]byte("1"), 48)),
		[48]byte(bytes.Repeat([]byte("2"), 48)),
		[48]byte(bytes.Repeat([]byte("3"), 48)),
	}
	body := mockBody()

	blk := &types.BeaconBlockDeneb{
		Slot:          1,
		ProposerIndex: 1,
		ParentBlockRoot: primitives.HashRoot(
			common.HexToHash("0x07").Bytes()),
		Body: body,
	}

	index := uint64(1)
	proof, err := types.MerkleProofKZGCommitment(blk, index)
	require.NoError(t, err)
	require.Len(t,
		proof,
		int(types.LogMaxBlobCommitments)+1+int(types.LogBodyLength))

	chunk := make([][32]byte, 2)
	copy(chunk[0][:], kzgs[index][:])
	copy(chunk[1][:], kzgs[index][32:])
	gohashtree.HashChunks(chunk, chunk)
	root, err := body.HashTreeRoot()
	require.NoError(t, err)

	commitments := body.GetBlobKzgCommitments()
	commitmentsRoot, err := types.GetBlobKzgCommitmentsRoot(commitments)
	require.NoError(t, err, "Failed to generate root hash")

	require.True(t,
		trie.VerifyMerkleProofWithDepth(
			commitmentsRoot[:],
			chunk[0][:],
			index,
			proof[:types.LogMaxBlobCommitments+1],
			types.LogMaxBlobCommitments,
		),
	)

	// Body
	bodyMembersRoots, err := types.GetTopLevelRoots(body)
	require.NoError(t, err, "Failed to get top level roots")
	// Add the commitments root to the body members roots.
	// For this test only. We don't need to do this when
	// generating the proof.
	// bodyMembersRoots[types.KZGPosition] = commitmentsRoot[:]
	bodySparse, err := trie.NewFromItems(
		bodyMembersRoots,
		types.LogBodyLength,
	)
	require.NoError(t, err, "Failed to generate trie from member roots")
	require.Equal(t, types.BodyLength, bodySparse.NumOfItems())
	topProof, err := bodySparse.MerkleProof(types.KZGPosition)
	require.NoError(t, err, "Failed to generate Merkle proof")
	require.Equal(t,
		topProof[:len(topProof)-1],
		proof[types.LogMaxBlobCommitments+1:],
	)

	require.Len(t,
		proof[types.LogMaxBlobCommitments+1:],
		int(types.LogBodyLength),
	)
	require.True(t,
		trie.VerifyMerkleProof(
			root[:],
			commitmentsRoot[:],
			uint64(types.KZGPosition),
			proof[types.LogMaxBlobCommitments+1:],
		),
	)

	require.True(t,
		trie.VerifyMerkleProof(
			root[:],
			chunk[0][:],
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
// 		ok := trie.VerifyMerkleProof(
// 			root[:],
// 			chunk[0][:],
// 			uint64(index+i*4096),
// 			proof,
// 		)
// 		fmt.Printf("i: %d, ok: %v\n", i, ok)
// 	}

// 	require.True(
// 		t,
// 		trie.VerifyMerkleProof(
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

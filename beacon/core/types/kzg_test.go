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
	"math"
	"testing"

	"github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/trie"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_BodyProof(t *testing.T) {
	// Create a real ExecutionPayloadDeneb and BeaconBlockBody
	executionPayload := &enginetypes.ExecutableDataDeneb{
		ParentHash:    common.HexToHash("0x01"),
		FeeRecipient:  common.HexToAddress("0x02"),
		StateRoot:     common.HexToHash("0x03"),
		ReceiptsRoot:  common.HexToHash("0x04"),
		LogsBloom:     []byte("bloom"),
		Random:        common.HexToHash("0x05"),
		BaseFeePerGas: []byte("baseFee"),
		BlockHash:     common.HexToHash("0x06"),
		Transactions:  [][]byte{[]byte("tx1"), []byte("tx2")},
		ExtraData:     []byte("extra"),
	}

	body := &types.BeaconBlockBodyDeneb{
		RandaoReveal:     [96]byte{0x01},
		ExecutionPayload: executionPayload,
		BlobKzgCommitments: [][48]byte{
			[48]byte(bytes.Repeat([]byte("1"), 48)),
		},
	}

	// The body has the commitments.
	commitments := body.GetBlobKzgCommitments()

	// Generate leaves from commitments
	leaves := types.LeavesFromCommitments(commitments)

	// Calculate the depth the given trie will have.
	depth := uint64(math.Ceil(math.Sqrt(float64(len(commitments)))))

	// Generate a sparse Merkle tree from the leaves.
	sparse, err := trie.GenerateTrieFromItems(leaves, depth)
	require.NoError(t, err, "Failed to generate trie from items")

	// Get the root of the tree.
	root, err := sparse.HashTreeRoot()
	require.NoError(t, err, "Failed to generate root hash")

	// Generate a proof for the index.
	index := 0
	proof, err := sparse.MerkleProof(index)
	require.NoError(t, err, "Failed to generate Merkle proof")
	require.NotNil(t, proof, "Merkle proof should not be nil")

	// Verify the Merkle proof
	valid := trie.VerifyMerkleProof(
		root[:],
		leaves[index],
		uint64(index),
		proof,
	)
	require.True(t, valid, "Merkle proof should be valid")
}

func Test_TopLevelRoots(t *testing.T) {
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

	body := &types.BeaconBlockBodyDeneb{
		RandaoReveal:     [96]byte{0x01},
		ExecutionPayload: executionPayload,
		BlobKzgCommitments: [][48]byte{
			[48]byte(bytes.Repeat([]byte("1"), 48)),
		},
	}

	// Commitments
	commitments := body.GetBlobKzgCommitments()
	commitmentsLeaves := types.LeavesFromCommitments(commitments)
	depth := uint64(math.Ceil(math.Sqrt(float64(len(commitments)))))
	commitmentsSparse, err := trie.GenerateTrieFromItems(commitmentsLeaves, depth)
	require.NoError(t, err, "Failed to generate trie from items")
	commitmentsRoot, err := commitmentsSparse.HashTreeRoot()
	require.NoError(t, err, "Failed to generate root hash")

	// Body
	bodyMembersRoots, err := types.GetTopLevelRoots(body)
	require.NoError(t, err, "Failed to get top level roots")
	bodySparse, err := trie.GenerateTrieFromItems(bodyMembersRoots, 3)
	require.NoError(t, err, "Failed to generate trie from member roots")
	bodyRoot, err := bodySparse.HashTreeRoot()
	require.NoError(t, err, "Failed to generate root hash")

	commitmentsIndex := 5
	topProof, err := bodySparse.MerkleProof(commitmentsIndex)
	require.NoError(t, err, "Failed to generate Merkle proof")

	// Verify the Merkle proof
	valid := trie.VerifyMerkleProof(
		bodyRoot[:],
		commitmentsRoot[:],
		uint64(commitmentsIndex),
		topProof,
	)
	require.True(t, valid, "Merkle proof should be valid")
}

// func Test_MerkleProofKZGCommitment(t *testing.T) {
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
// 	// pbBody := &beacontypes.BeaconBlockBodyDeneb{
// 	// 	ExecutionPayload: &enginev1.ExecutionPayloadDeneb{
// 	// 		ParentHash:    make([]byte, fieldparams.RootLength),
// 	// 		FeeRecipient:  make([]byte, 20),
// 	// 		StateRoot:     make([]byte, fieldparams.RootLength),
// 	// 		ReceiptsRoot:  make([]byte, fieldparams.RootLength),
// 	// 		LogsBloom:     make([]byte, 256),
// 	// 		PrevRandao:    make([]byte, fieldparams.RootLength),
// 	// 		BaseFeePerGas: make([]byte, fieldparams.RootLength),
// 	// 		BlockHash:     make([]byte, fieldparams.RootLength),
// 	// 		Transactions:  make([][]byte, 0),
// 	// 		ExtraData:     make([]byte, 0),
// 	// 	},
// 	// 	Eth1Data: &ethpb.Eth1Data{
// 	// 		DepositRoot: make([]byte, fieldparams.RootLength),
// 	// 		BlockHash:   make([]byte, fieldparams.RootLength),
// 	// 	},
// 	// 	BlobKzgCommitments: kzgs,
// 	// }

// 	kzgs48 := make([][48]byte, 3)
// 	for i, kzg := range kzgs {
// 		copy(kzgs48[i][:], kzg)
// 	}
// 	body := &types.BeaconBlockBodyDeneb{
// 		ExecutionPayload: &enginetypes.ExecutableDataDeneb{
// 			ParentHash:    common.Hash{},
// 			FeeRecipient:  common.Address{},
// 			StateRoot:     common.Hash{},
// 			ReceiptsRoot:  common.Hash{},
// 			LogsBloom:     make([]byte, 256),
// 			Random:        common.Hash{},
// 			BaseFeePerGas: make([]byte, 32),
// 			BlockHash:     common.Hash{},
// 			Transactions:  make([][]byte, 0),
// 			ExtraData:     make([]byte, 0),
// 		},
// 		BlobKzgCommitments: kzgs48,
// 	}
// 	blk := &types.BeaconBlockDeneb{
// 		Slot:            0,
// 		ProposerIndex:   0,
// 		ParentBlockRoot: [32]byte{},
// 		Body:            body,
// 	}
// 	require.NoError(t, err)
// 	index := 1
// 	_, err = types.MerkleProofKZGCommitment(blk, 10)
// 	require.NotNil(t, err)
// 	proof, err := types.MerkleProofKZGCommitment(blk, index)
// 	require.NoError(t, err)

// 	chunk := make([][32]byte, 2)
// 	copy(chunk[0][:], kzgs[index])
// 	copy(chunk[1][:], kzgs[index][32:])
// 	gohashtree.HashChunks(chunk, chunk)
// 	root, err := body.HashTreeRoot()
// 	require.NoError(t, err)
// 	kzgOffset := 54 * 4096
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

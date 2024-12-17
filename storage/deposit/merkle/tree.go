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

package merkle

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/crypto/sha256"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/ethereum/go-ethereum/common"
)

// DepositTree is the Merkle tree representation of deposits.
type DepositTree struct {
	tree                    TreeNode
	mixInLength             uint64
	finalizedExecutionBlock executionBlock
	hasher                  merkle.Hasher[[32]byte]
}

type executionBlock struct {
	Hash  [32]byte
	Depth uint64
}

// NewDepositTree creates an empty deposit tree.
func NewDepositTree() *DepositTree {
	var (
		hasher = merkle.NewHasher[[32]byte](sha256.Hash)
		leaves [][32]byte
	)
	merkle := create(hasher, leaves, DepositContractDepth)
	return &DepositTree{
		tree:                    merkle,
		mixInLength:             0,
		finalizedExecutionBlock: executionBlock{},
		hasher:                  hasher,
	}
}

// GetSnapshot returns a deposit tree snapshot.
func (d *DepositTree) GetSnapshot() (DepositTreeSnapshot, error) {
	var finalized [][32]byte
	mixInLength, finalized := d.tree.GetFinalized(finalized)
	return fromTreeParts(finalized, mixInLength, d.finalizedExecutionBlock)
}

// Finalize marks a deposit as finalized.
func (d *DepositTree) Finalize(
	eth1DepositIndex uint64,
	executionHash common.Hash,
	executionNumber uint64,
) error {
	var blockHash [32]byte
	copy(blockHash[:], executionHash[:])
	d.finalizedExecutionBlock = executionBlock{
		Hash:  blockHash,
		Depth: executionNumber,
	}
	mixInLength := eth1DepositIndex + 1
	_, err := d.tree.Finalize(mixInLength, DepositContractDepth)
	if err != nil {
		return err
	}
	return nil
}

// getProof returns the deposit tree proof.
func (d *DepositTree) getProof(index uint64) ([32]byte, [][32]byte, error) {
	if d.mixInLength <= 0 {
		return [32]byte{}, nil, ErrInvalidDepositCount
	}
	if index >= d.mixInLength {
		return [32]byte{}, nil, ErrInvalidIndex
	}
	finalizedDeposits, _ := d.tree.GetFinalized([][32]byte{})
	finalizedIdx := -1
	if finalizedDeposits != 0 {
		fd, err := math.Int(finalizedDeposits)
		if err != nil {
			return [32]byte{}, nil, err
		}
		finalizedIdx = fd - 1
	}
	i, err := math.Int(index)
	if err != nil {
		return [32]byte{}, nil, err
	}
	if finalizedDeposits > 0 && i <= finalizedIdx {
		return [32]byte{}, nil, ErrInvalidIndex
	}
	leaf, proof := generateProof(d.tree, index, DepositContractDepth)

	mixInLength := [32]byte{}
	binary.LittleEndian.PutUint64(mixInLength[:], d.mixInLength)
	proof = append(proof, mixInLength)
	return leaf, proof, nil
}

// getRoot returns the root of the deposit tree.
func (d *DepositTree) getRoot() [32]byte {
	var enc [32]byte
	binary.LittleEndian.PutUint64(enc[:], d.mixInLength)

	root := d.tree.GetRoot()
	return d.hasher.Combi(root, enc)
}

// pushLeaf adds a new leaf to the tree.
func (d *DepositTree) pushLeaf(leaf [32]byte) error {
	var err error
	d.tree, err = d.tree.PushLeaf(leaf, DepositContractDepth)
	if err != nil {
		return err
	}
	d.mixInLength++
	return nil
}

// Insert is defined as part of MerkleTree interface and adds a new leaf to the
// tree.
func (d *DepositTree) Insert(item []byte, _ int) error {
	var leaf [32]byte
	copy(leaf[:], item[:32])
	return d.pushLeaf(leaf)
}

// HashTreeRoot is defined as part of MerkleTree interface and calculates the
// hash tree root.
func (d *DepositTree) HashTreeRoot() ([32]byte, error) {
	root := d.getRoot()
	if root == [32]byte{} {
		return [32]byte{}, errors.New("could not retrieve hash tree root")
	}
	return root, nil
}

// NumOfItems is defined as part of MerkleTree interface and returns the number
// of deposits in the tree.
func (d *DepositTree) NumOfItems() int {
	return int(d.mixInLength)
}

// MerkleProof is defined as part of MerkleTree interface and generates a merkle
// proof.
func (d *DepositTree) MerkleProof(index int) ([][]byte, error) {
	_, proof, err := d.getProof(uint64(index))
	if err != nil {
		return nil, err
	}
	byteSlices := make([][]byte, len(proof))
	for i, p := range proof {
		copied := p
		byteSlices[i] = copied[:]
	}
	return byteSlices, nil
}

// Copy performs a deep copy of the tree.
func (d *DepositTree) Copy() (*DepositTree, error) {
	snapshot, err := d.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return fromSnapshot(d.hasher, snapshot)
}

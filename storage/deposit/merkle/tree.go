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

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto/sha256"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/merkle"
)

// DepositTree is the Merkle tree representation of deposits.
type DepositTree struct {
	tree                    TreeNode
	mixInLength             uint64
	finalizedExecutionBlock executionBlock
	hasher                  merkle.Hasher[common.Root]
}

type executionBlock struct {
	Hash  common.ExecutionHash
	Depth math.U64
}

// NewDepositTree creates an empty deposit tree.
//
// NOTE: Not safe for concurrent use as it uses a single hasher.
func NewDepositTree() *DepositTree {
	var (
		hasher = merkle.NewHasher[common.Root](sha256.Hash)
		leaves []common.Root
	)
	merkle := create(hasher, leaves, constants.DepositContractDepth)
	return &DepositTree{
		tree:                    merkle,
		mixInLength:             0,
		finalizedExecutionBlock: executionBlock{},
		hasher:                  hasher,
	}
}

// GetSnapshot returns a deposit tree snapshot.
func (d *DepositTree) GetSnapshot() DepositTreeSnapshot {
	var finalized []common.Root
	mixInLength, finalized := d.tree.GetFinalized(finalized)
	return fromTreeParts(
		d.hasher,
		finalized,
		mixInLength,
		d.finalizedExecutionBlock,
	)
}

// Finalize marks a deposit as finalized.
func (d *DepositTree) Finalize(
	eth1DepositIndex uint64,
	executionHash common.ExecutionHash,
	executionNumber math.U64,
) error {
	d.finalizedExecutionBlock = executionBlock{
		Hash:  executionHash,
		Depth: executionNumber,
	}
	mixInLength := eth1DepositIndex + 1
	_, err := d.tree.Finalize(mixInLength, constants.DepositContractDepth)
	if err != nil {
		return err
	}
	return nil
}

// getProof returns the deposit tree proof.
func (d *DepositTree) getProof(index uint64) (
	common.Root, [constants.DepositContractDepth + 1]common.Root, error,
) {
	var proof [constants.DepositContractDepth + 1]common.Root

	if d.mixInLength <= 0 {
		return common.Root{}, proof, ErrInvalidDepositCount
	}
	if index >= d.mixInLength {
		return common.Root{}, proof, ErrInvalidIndex
	}

	finalizedDeposits, _ := d.tree.GetFinalized([]common.Root{})
	finalizedIdx := -1
	if finalizedDeposits != 0 {
		fd, err := math.Int(finalizedDeposits)
		if err != nil {
			return common.Root{}, proof, err
		}
		finalizedIdx = fd - 1
	}
	i, err := math.Int(index)
	if err != nil {
		return common.Root{}, proof, err
	}
	if finalizedDeposits > 0 && i <= finalizedIdx {
		return common.Root{}, proof, ErrInvalidIndex
	}

	leaf, proofWithoutMixin := generateProof(d.tree, index, constants.DepositContractDepth)
	copy(proof[:constants.DepositContractDepth], proofWithoutMixin[:])

	mixInLength := common.Root{}
	binary.LittleEndian.PutUint64(mixInLength[:], d.mixInLength)
	proof[constants.DepositContractDepth] = mixInLength
	return leaf, proof, nil
}

// getRoot returns the root of the deposit tree.
func (d *DepositTree) getRoot() common.Root {
	var enc common.Root
	binary.LittleEndian.PutUint64(enc[:], d.mixInLength)

	root := d.tree.GetRoot()
	return d.hasher.Combi(root, enc)
}

// pushLeaf adds a new leaf to the tree.
func (d *DepositTree) pushLeaf(leaf common.Root) error {
	var err error
	d.tree, err = d.tree.PushLeaf(leaf, constants.DepositContractDepth)
	if err != nil {
		return err
	}
	d.mixInLength++
	return nil
}

// Insert is defined as part of MerkleTree interface and adds a new leaf to the tree.
func (d *DepositTree) Insert(item common.Root) error {
	return d.pushLeaf(item)
}

// HashTreeRoot is defined as part of MerkleTree interface and calculates the hash tree root.
func (d *DepositTree) HashTreeRoot() common.Root {
	return d.getRoot()
}

// NumOfItems is defined as part of MerkleTree interface and returns the number of deposits in the tree.
func (d *DepositTree) NumOfItems() uint64 {
	return d.mixInLength
}

// MerkleProof is defined as part of MerkleTree interface and generates a merkle proof.
func (d *DepositTree) MerkleProof(index uint64) (
	[constants.DepositContractDepth + 1]common.Root, error,
) {
	_, proof, err := d.getProof(index)
	return proof, err
}

// Copy performs a deep copy of the tree.
func (d *DepositTree) Copy() (*DepositTree, error) {
	return fromSnapshot(d.GetSnapshot())
}

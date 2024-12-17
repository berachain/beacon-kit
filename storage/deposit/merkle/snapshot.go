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
	"github.com/berachain/beacon-kit/primitives/math/pow"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/berachain/beacon-kit/primitives/merkle/zero"
)

// DepositTreeSnapshot represents the data used to create a deposit tree given a
// snapshot.
type DepositTreeSnapshot struct {
	finalized      [][32]byte
	depositRoot    [32]byte
	depositCount   uint64
	executionBlock executionBlock
	hasher         merkle.Hasher[[32]byte]
}

// CalculateRoot returns the root of a deposit tree snapshot.
func (ds *DepositTreeSnapshot) CalculateRoot() ([32]byte, error) {
	size := ds.depositCount
	index := len(ds.finalized)
	root := zero.Hashes[0]
	for i := range DepositContractDepth {
		if (size & 1) == 1 {
			if index == 0 {
				break
			}
			index--
			root = ds.hasher.Combi(ds.finalized[index], root)
		} else {
			root = ds.hasher.Combi(root, zero.Hashes[i])
		}
		size >>= 1
	}
	return ds.hasher.MixIn(root, ds.depositCount), nil
}

// fromSnapshot returns a deposit tree from a deposit tree snapshot.
func fromSnapshot(
	hasher merkle.Hasher[[32]byte],
	snapshot DepositTreeSnapshot,
) (*DepositTree, error) {
	root, err := snapshot.CalculateRoot()
	if err != nil {
		return nil, err
	}
	if snapshot.depositRoot != root {
		return nil, ErrInvalidSnapshotRoot
	}
	if snapshot.depositCount >= pow.TwoToThePowerOf(DepositContractDepth) {
		return nil, ErrTooManyDeposits
	}
	tree, err := fromSnapshotParts(
		hasher,
		snapshot.finalized,
		snapshot.depositCount,
		DepositContractDepth,
	)
	if err != nil {
		return nil, err
	}
	return &DepositTree{
		tree:                    tree,
		mixInLength:             snapshot.depositCount,
		finalizedExecutionBlock: snapshot.executionBlock,
	}, nil
}

// fromTreeParts constructs the deposit tree from pre-existing data.
func fromTreeParts(
	finalised [][32]byte,
	depositCount uint64,
	executionBlock executionBlock,
) (DepositTreeSnapshot, error) {
	snapshot := DepositTreeSnapshot{
		finalized:      finalised,
		depositRoot:    zero.Hashes[0],
		depositCount:   depositCount,
		executionBlock: executionBlock,
	}
	root, err := snapshot.CalculateRoot()
	if err != nil {
		return snapshot, ErrInvalidSnapshotRoot
	}
	snapshot.depositRoot = root
	return snapshot, nil
}

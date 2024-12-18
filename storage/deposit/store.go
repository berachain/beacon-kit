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

package deposit

import (
	"sync"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/deposit/merkle"
)

// Store is a simple memory store based implementation that
// maintains a merkle tree of the deposits for verification.
type Store struct {
	// tree is the EIP-4881 compliant deposit merkle tree.
	tree *merkle.DepositTree

	// pendingDeposits holds the pending deposits for blocks that have yet to be
	// processed by the CL.
	pendingDeposits []*Block

	// lastUsedIndex is the index of the last deposit included in CL blocks.
	lastUsedIndex uint64

	// mu protects store for concurrent access.
	mu sync.RWMutex
}

// NewStore creates a new deposit store.
func NewStore() *Store {
	res := &Store{
		tree:            merkle.NewDepositTree(),
		pendingDeposits: make([]*Block, 0),
	}
	return res
}

// GetDepositsByIndex returns the first N deposits starting from the given
// index. If N is greater than the number of deposits, it returns up to the
// last deposit available. It also returns the deposit tree root at the end of
// the range.
func (s *Store) GetDepositsByIndex(
	startIndex, numView uint64,
) (ctypes.Deposits, common.Root, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var (
		deposits = ctypes.Deposits{}
		// maxIndex    = startIndex + numView
		depTreeRoot common.Root
	)

	// for i := startIndex; i < maxIndex; i++ {
	// 	deposit, err := s.store.Get(context.TODO(), i)
	// 	if err == nil {
	// 		deposits = append(deposits, deposit)
	// 		continue
	// 	}

	// 	if errors.Is(err, sdkcollections.ErrNotFound) {
	// 		depTreeRoot = s.pendingDepositsToRoots[i-1]
	// 		break
	// 	}

	// 	return nil, common.Root{}, errors.Wrapf(
	// 		err, "failed to get deposit %d, start: %d, end: %d", i, startIndex, maxIndex,
	// 	)
	// }

	// if depTreeRoot == (common.Root{}) {
	// 	depTreeRoot = s.pendingDepositsToRoots[maxIndex-1]
	// }

	return deposits, depTreeRoot, nil
}

// EnqueueDepositDatas pushes multiple deposits to the queue for a given EL block.
//
// TODO: ensure that in-order is maintained. i.e. ignore any deposits we've already seen.
func (s *Store) EnqueueDepositDatas(
	depositDatas []*ctypes.DepositData,
	indexes []uint64,
	executionHash common.ExecutionHash,
	executionNumber math.U64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Build the deposits information for the block while inserting into the deposit tree.
	block := &Block{
		executionHash:   executionHash,
		executionNumber: executionNumber,
		deposits:        make(ctypes.Deposits, len(depositDatas)),
		root:            make([]common.Root, len(depositDatas)),
	}
	for i, depositData := range depositDatas {
		if err := s.tree.Insert(depositData.HashTreeRoot()); err != nil {
			return errors.Wrapf(err,
				"failed to insert deposit %d into merkle tree, execution number: %d",
				indexes[i], executionNumber,
			)
		}

		proof, err := s.tree.MerkleProof(indexes[i])
		if err != nil {
			return errors.Wrapf(err,
				"failed to get merkle proof for deposit %d, execution number: %d",
				indexes[i], executionNumber,
			)
		}
		block.deposits[i] = ctypes.NewDeposit(proof, depositData)
		block.root[i] = s.tree.HashTreeRoot()
	}
	s.pendingDeposits = append(s.pendingDeposits, block)

	// Finalize the block's deposits in the tree. Error returned here means the
	// EIP 4881 merkle library is broken. // TODO: could move this to when we delete block.
	lastDepositIndex := indexes[len(indexes)-1]
	if err := s.tree.Finalize(lastDepositIndex, executionHash, executionNumber); err != nil {
		return errors.Wrapf(
			err, "failed to finalize deposits in tree for block %d, index %d",
			executionNumber, lastDepositIndex,
		)
	}

	return nil
}

// Prune removes the deposits from the given height.
// NO-OP for now since the pruning call is not at the right time.
func (s *Store) Prune(height uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return nil
}

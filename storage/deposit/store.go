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
	"fmt"
	"sync"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/storage/deposit/merkle"
	"github.com/berachain/beacon-kit/storage/pruner"
)

const KeyDepositPrefix = "deposit"

// Store is a simple memory store based implementation that
// maintains a merkle tree of the deposits for verification.
type Store struct {
	// tree is the EIP-4881 compliant deposit merkle tree.
	tree *merkle.DepositTree

	// pendingDepositsToRoots maps the deposit tree root after each deposit.
	pendingDepositsToRoots map[uint64]common.Root

	// mu protects store for concurrent access.
	mu sync.Mutex
}

// NewStore creates a new deposit store.
func NewStore() *Store {
	res := &Store{
		tree:                   merkle.NewDepositTree(),
		pendingDepositsToRoots: make(map[uint64]common.Root),
	}
	return res
}

// GetDepositsByIndex returns the first N deposits starting from the given
// index. If N is greater than the number of deposits, it returns up to the
// last deposit available. It also returns the deposit tree root at the end of
// the range.
//
// TODO: figure out when to finalize. Need to do after proof has been generated.
func (s *Store) GetDepositsByIndex(
	startIndex, numView uint64,
) (ctypes.Deposits, common.Root, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var (
		deposits    = ctypes.Deposits{}
		maxIndex    = startIndex + numView
		depTreeRoot common.Root
	)

	for i := startIndex; i < maxIndex; i++ {
		proof, err := s.tree.MerkleProof(i)
		if err != nil {
			return nil, common.Root{}, errors.Wrapf(
				err, "failed to get merkle proof for deposit %d", i,
			)
		}
		deposits = append(deposits, ctypes.NewDeposit(proof, nil))
		delete(s.pendingDepositsToRoots, i-1)
	}

	if depTreeRoot == (common.Root{}) {
		depTreeRoot = s.pendingDepositsToRoots[maxIndex-1]
		delete(s.pendingDepositsToRoots, maxIndex-1)
	}
	return deposits, depTreeRoot, nil
}

// EnqueueDepositDatas pushes multiple deposits to the queue.
//
// TODO: ensure that in-order is maintained. i.e. ignore any deposits we've already seen.
func (s *Store) EnqueueDepositDatas(deposits []*ctypes.DepositData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, deposit := range deposits {
		if err := s.tree.Insert(deposit.HashTreeRoot()); err != nil {
			return errors.Wrap(err, "failed to insert deposit into merkle tree")
		}

		// s.pendingDepositsToRoots[idx] = s.tree.HashTreeRoot()
	}

	return nil
}

// Prune removes the [start, end) deposits from the store.
func (s *Store) Prune(start, end uint64) error {
	if start > end {
		return fmt.Errorf(
			"DepositStore prune start: %d, end: %d: %w", start, end, pruner.ErrInvalidRange,
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range end {
		delete(s.pendingDepositsToRoots, start+i)
	}

	return nil
}

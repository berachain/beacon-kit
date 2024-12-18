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
	"context"
	"fmt"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/deposit/merkle"
	"github.com/berachain/beacon-kit/storage/encoding"
	"github.com/berachain/beacon-kit/storage/pruner"
)

const KeyDepositPrefix = "deposit"

// Store is a simple KV store based implementation that assumes
// the deposit indexes are tracked outside of the s store.
// It also maintains a merkle tree of the deposits for verification,
// which will remove the need for indexed based tracking.
type Store struct {
	// tree is the EIP-4881 compliant deposit merkle tree.
	tree *merkle.DepositTree

	// store is the KV store that holds the deposits.
	store sdkcollections.Map[uint64, *ctypes.DepositData]

	// mu protects store for concurrent access.
	mu sync.RWMutex
}

// NewStore creates a new deposit store.
func NewStore(kvsp store.KVStoreService) *Store {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	res := &Store{
		tree: merkle.NewDepositTree(),
		store: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte(KeyDepositPrefix)),
			KeyDepositPrefix,
			sdkcollections.Uint64Key,
			encoding.SSZValueCodec[*ctypes.DepositData]{},
		),
	}
	if _, err := schemaBuilder.Build(); err != nil {
		panic(fmt.Errorf("failed building Store schema: %w", err))
	}
	return res
}

// GetDepositsByIndex returns the first N deposits starting from the given
// index. If N is greater than the number of deposits, it returns up to the
// last deposit.
func (s *Store) GetDepositsByIndex(
	startIndex uint64,
	numView uint64,
) (ctypes.Deposits, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var (
		deposits = ctypes.Deposits{}
		endIdx   = startIndex + numView
	)

	for i := startIndex; i < endIdx; i++ {
		deposit, err := s.store.Get(context.TODO(), i)
		switch {
		case err == nil:
			var proof [constants.DepositContractDepth + 1]common.Root
			proof, err = s.tree.MerkleProof(i)
			if err != nil {
				return deposits, errors.Wrapf(err, "failed to get merkle proof for deposit %d", i)
			}
			deposits = append(deposits, ctypes.NewDeposit(proof, deposit))
		case errors.Is(err, sdkcollections.ErrNotFound):
			return deposits, nil
		default:
			return deposits, errors.Wrapf(
				err, "failed to get deposit %d, start: %d, end: %d", i, startIndex, endIdx,
			)
		}
	}

	return deposits, nil
}

// GetDepositsRoot returns the root of the deposit merkle tree. This is the hash tree
// root of the deposit datas.
func (s *Store) GetDepositsRoot() common.Root {
	return s.tree.HashTreeRoot()
}

// GetDepositsCount returns the number of deposits in the store.
func (s *Store) GetDepositsCount() uint64 {
	return s.tree.NumOfItems()
}

// EnqueueDepositDatas pushes multiple deposits to the queue.
func (s *Store) EnqueueDepositDatas(
	deposits []*ctypes.DepositData,
	executionBlockHash common.ExecutionHash,
	executionBlockNumber math.U64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, deposit := range deposits {
		idx := deposit.GetIndex().Unwrap()
		if err := s.store.Set(context.TODO(), idx, deposit); err != nil {
			return errors.Wrapf(err, "failed to set deposit %d in KVStore", idx)
		}

		if err := s.tree.Insert(deposit.HashTreeRoot()); err != nil {
			return errors.Wrapf(err, "failed to insert deposit %d into merkle tree", idx)
		}

		if err := s.tree.Finalize(idx, executionBlockHash, executionBlockNumber); err != nil {
			return errors.Wrapf(err, "failed to finalize deposit %d in merkle tree", idx)
		}
	}

	return nil
}

// Prune removes the [start, end) deposits from the store.
func (s *Store) Prune(start, end uint64) error {
	if start > end {
		return fmt.Errorf(
			"DepositKVStore Prune start: %d, end: %d: %w", start, end, pruner.ErrInvalidRange,
		)
	}

	var ctx = context.TODO()
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range end {
		// This only errors if the key passed in cannot be encoded.
		if err := s.store.Remove(ctx, start+i); err != nil {
			return errors.Wrapf(err, "failed to prune deposit %d", start+i)
		}
	}

	return nil
}

// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"errors"
	"fmt"
	"sync"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	depositstorev1 "github.com/berachain/beacon-kit/storage/deposit/v1"
	dbm "github.com/cosmos/cosmos-db"
)

const (
	unset uint8 = 0
	V1    uint8 = 1
)

type Store interface {
	GetDepositsByIndex(ctx context.Context, startIndex uint64, depRange uint64) (ctypes.Deposits, common.Root, error)
	EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error
	Prune(ctx context.Context, start, end uint64) error
	Close() error
}

type StoreManager interface {
	Store
}

var (
	_ Store        = (*depositstorev1.KVStore)(nil)
	_ StoreManager = (*generalStore)(nil)

	ErrUnknownStoreVersion = errors.New("unknown deposit store version")
)

// We have changed in time the way we stored deposits. generalStore is meant to offer
// a single way to access deposits and to handle the data migration among versions when needed
type generalStore struct {
	// mu protects storex for concurrent access
	mu             sync.RWMutex
	currentVersion uint8
	storeV1        *depositstorev1.KVStore
	logger         log.Logger
}

func NewStore(dbV1 dbm.DB, logger log.Logger) StoreManager {
	storeV1 := depositstorev1.NewStore(dbV1, logger)

	currentVersion := V1
	return &generalStore{
		currentVersion: currentVersion,
		storeV1:        storeV1,
		logger:         logger,
	}
}

func (gs *generalStore) GetDepositsByIndex(
	ctx context.Context,
	startIndex uint64,
	depRange uint64,
) (ctypes.Deposits, common.Root, error) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	switch gs.currentVersion {
	case V1:
		return gs.storeV1.GetDepositsByIndex(ctx, startIndex, depRange)
	default:
		return nil, common.Root{}, fmt.Errorf("%w, version %d", ErrUnknownStoreVersion, gs.currentVersion)
	}
}

func (gs *generalStore) EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	switch gs.currentVersion {
	case V1:
		return gs.storeV1.EnqueueDeposits(ctx, deposits)
	default:
		return fmt.Errorf("%w, version %d", ErrUnknownStoreVersion, gs.currentVersion)
	}
}

func (gs *generalStore) Close() error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	return gs.storeV1.Close()
}

func (gs *generalStore) Prune(ctx context.Context, start, end uint64) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	switch gs.currentVersion {
	case V1:
		return gs.storeV1.Prune(ctx, start, end)
	default:
		return fmt.Errorf("%w, version %d", ErrUnknownStoreVersion, gs.currentVersion)
	}
}

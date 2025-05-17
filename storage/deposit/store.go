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

	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	depositstorev1 "github.com/berachain/beacon-kit/storage/deposit/v1"
	depositstorev2 "github.com/berachain/beacon-kit/storage/deposit/v2"
)

type Store interface {
	GetDepositsByIndex(ctx context.Context, startIndex uint64, depRange uint64) (ctypes.Deposits, common.Root, error)
	EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error
	Prune(ctx context.Context, start, end uint64) error
	Close() error
}

var (
	_ Store = (*depositstorev1.KVStore)(nil)
	_ Store = (*depositstorev2.KVStore)(nil)
)

type CloseFunc func() error

// We have changed in time the way we stored deposits. store is meant to offer
// a single way to access deposits
type generalStore struct {
	cs chain.Spec

	storeV1 *depositstorev1.KVStore

	//nolint:unused // TODO ABENEGIA: wip, the the whole point of generalStore is to switch among store versions
	storeV2 *depositstorev2.KVStore
}

func NewStore(
	cs chain.Spec,
	kvsp store.KVStoreService,
	closeFunc CloseFunc,
	logger log.Logger,
) Store {
	storeV1 := depositstorev1.NewStore(kvsp, depositstorev1.CloseFunc(closeFunc), logger)
	return &generalStore{
		cs:      cs,
		storeV1: storeV1,
	}
}

func (gs *generalStore) GetDepositsByIndex(
	ctx context.Context,
	startIndex uint64,
	depRange uint64,
) (ctypes.Deposits, common.Root, error) {
	// TODO ABENEGIA: add switch at Electra fork
	return gs.storeV1.GetDepositsByIndex(ctx, startIndex, depRange)
}

func (gs *generalStore) EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error {
	// TODO ABENEGIA: add switch at Electra fork
	return gs.storeV1.EnqueueDeposits(ctx, deposits)
}
func (gs *generalStore) Close() error {
	// TODO ABENEGIA: add switch at Electra fork
	return gs.storeV1.Close()
}

func (gs *generalStore) Prune(ctx context.Context, start, end uint64) error {
	// TODO ABENEGIA: add switch at Electra fork
	return gs.storeV1.Prune(ctx, start, end)
}

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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	depositstorev1 "github.com/berachain/beacon-kit/storage/deposit/v1"
	depositstorev2 "github.com/berachain/beacon-kit/storage/deposit/v2"
	dbm "github.com/cosmos/cosmos-db"
)

const (
	unset uint8 = 0
	V1    uint8 = 1
	V2    uint8 = 2
)

type Store interface {
	// TODO ABENEGIA: consider having a GetDepositsByIndex and GetDepositsUpToIndex methods
	GetDepositsByIndex(ctx context.Context, startIndex uint64, depRange uint64) (ctypes.Deposits, common.Root, error)
	EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error
	Prune(ctx context.Context, start, end uint64) error
	Close() error
}

type StoreManager interface {
	Store
	MigrateV1ToV2() error
	SelectVersion(uint8) error
}

var (
	_ Store        = (*depositstorev1.KVStore)(nil)
	_ Store        = (*depositstorev2.KVStore)(nil)
	_ StoreManager = (*generalStore)(nil)

	ErrUnknownStoreVersion = errors.New("unknown deposit store version")
)

// We have changed in time the way we stored deposits. generalStore is meant to offer
// a single way to access deposits and to handle the data migration among versions when needed
type generalStore struct {
	currentVersion uint8
	storeV1        *depositstorev1.KVStore
	storeV2        *depositstorev2.KVStore
}

func NewStore(
	dbV1 dbm.DB,
	dbV2 dbm.DB,

	logger log.Logger,
) StoreManager {
	storeV1 := depositstorev1.NewStore(dbV1, logger)
	storeV2 := depositstorev2.NewStore(dbV2, logger)
	return &generalStore{
		currentVersion: unset,
		storeV1:        storeV1,
		storeV2:        storeV2,
	}
}

func (gs *generalStore) GetDepositsByIndex(
	ctx context.Context,
	startIndex uint64,
	depRange uint64,
) (ctypes.Deposits, common.Root, error) {
	switch gs.currentVersion {
	case V1:
		return gs.storeV1.GetDepositsByIndex(ctx, startIndex, depRange)
	case V2:
		return gs.storeV2.GetDepositsByIndex(ctx, startIndex, depRange)
	default:
		return nil, common.Root{}, fmt.Errorf("%w, version %d", ErrUnknownStoreVersion, gs.currentVersion)
	}
}

func (gs *generalStore) EnqueueDeposits(ctx context.Context, deposits []*ctypes.Deposit) error {
	switch gs.currentVersion {
	case V1:
		return gs.storeV1.EnqueueDeposits(ctx, deposits)
	case V2:
		return gs.storeV2.EnqueueDeposits(ctx, deposits)
	default:
		return fmt.Errorf("%w, version %d", ErrUnknownStoreVersion, gs.currentVersion)
	}
}
func (gs *generalStore) Close() error {
	// TODO ABENEGIA: add switch at Electra fork
	return errors.Join(
		gs.storeV1.Close(),
		gs.storeV2.Close(),
	)
}

func (gs *generalStore) Prune(ctx context.Context, start, end uint64) error {
	switch gs.currentVersion {
	case V1:
		return gs.storeV1.Prune(ctx, start, end)
	case V2:
		return gs.storeV2.Prune(ctx, start, end)
	default:
		return fmt.Errorf("%w, version %d", ErrUnknownStoreVersion, gs.currentVersion)
	}
}

func (gs *generalStore) SelectVersion(v uint8) error {
	if v != V1 && v != V2 {
		return fmt.Errorf("%w, version %d", ErrUnknownStoreVersion, v)
	}
	gs.currentVersion = v
	return nil
}

// Simply copies over storeV1 content to storeV2 content when called
// Relies heavily on the fact that deposits are indexed by Deposit.Index
// which is contiguous and starts from zero
func (gs *generalStore) MigrateV1ToV2() error {
	ctx := context.TODO()
	// Note: under the hood GetDepositsByIndex allocates a slice up to depRange.
	// So we cannot just get all deposits from zero up to math.MaxUint64
	const span = 69
	var startIdx = constants.FirstDepositIndex
	for {
		v1Deposits, _, err := gs.storeV1.GetDepositsByIndex(ctx, startIdx, span)
		if err != nil {
			return fmt.Errorf(
				"failed loading v1 deposits from %d to %d: %w",
				startIdx, startIdx+span, err,
			)
		}
		if err = gs.storeV2.EnqueueDeposits(ctx, v1Deposits); err != nil {
			return fmt.Errorf(
				"failed copying to v2 deposits from %d to %d: %w",
				startIdx, startIdx+span, err,
			)
		}
		if uint64(len(v1Deposits)) < span {
			break // done
		}
		startIdx += span
	}
	return nil
}

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

package deposits

import (
	"context"
	"fmt"
	stdmath "math"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/execution/deposit"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	depositstore "github.com/berachain/beacon-kit/storage/deposit"
)

// If on the first block of Fulu, catchup the previous block's deposits. Between
// Prepare/ProcessProposal and FinalizeBlock, this only needs to be done once.
func CatchupFuluDeposits(
	ctx context.Context,
	depositContract deposit.Contract,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	chainSpec ChainSpec,
	depositStore depositstore.StoreManager,
	logger log.Logger,
) error {
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}
	prevBlockForkVersion := chainSpec.ActiveForkVersionForTimestamp(lph.GetTimestamp())
	isFirstFuluBlock := version.Equals(prevBlockForkVersion, version.Electra1()) &&
		version.Equals(blk.GetForkVersion(), version.Fulu())
	if !isFirstFuluBlock {
		return nil
	}

	// If we already fetched deposits for this block, we don't need to do it again.
	if depositContract.LastBlockNumber() == lph.GetNumber() {
		return nil
	}

	deposits, err := depositContract.ReadDeposits(ctx, lph.GetNumber())
	if err != nil {
		return err
	}
	if len(deposits) == 0 {
		logger.Info("Deposits catchup for Fulu, nothing to fetch")
		return nil
	}

	logger.Info("Found deposits to catchup for Fulu", "num", len(deposits))
	if err = depositStore.EnqueueDeposits(ctx, deposits); err != nil {
		logger.Error("Failed to store catchup deposits for Fulu", "error", err)
		return err
	}

	depositContract.SetLastBlockNumber(lph.GetNumber())

	return nil
}

// FetchPreviousDepositsPreFulu fetches deposits from the EL at the given eth1 follow distance.
func FetchPreviousDepositsPreFulu(
	ctx context.Context,
	depositContract deposit.Contract,
	blk *ctypes.BeaconBlock,
	eth1FollowDistance math.U64,
	depositStore depositstore.StoreManager,
	logger log.Logger,
) {
	// If after Fulu, we don't need to fetch previous deposits since EIP-6110 is used.
	if version.EqualsOrIsAfter(blk.GetForkVersion(), version.Fulu()) {
		return
	}

	// Fetch and store the deposit for the block.
	blockNum := blk.GetBody().GetExecutionPayload().GetNumber()
	if blockNum <= eth1FollowDistance {
		logger.Info(
			"depositFetcher, nothing to fetch",
			"block num", blockNum,
			"eth1FollowDistance", eth1FollowDistance,
		)
		return
	}
	blockToFetch := blockNum - eth1FollowDistance
	deposits, err := depositContract.ReadDeposits(ctx, blockToFetch)
	if err != nil {
		logger.Error("Failed to read deposits", "block", blockNum, "error", err)
		return
	}
	if len(deposits) == 0 {
		logger.Info(
			"depositFetcher, nothing to fetch",
			"block", blockNum, "eth1FollowDistance", eth1FollowDistance,
		)
	} else {
		logger.Info(
			"Found deposits on execution layer", "block", blockNum, "deposits", len(deposits),
		)
	}

	if err = depositStore.EnqueueDeposits(ctx, deposits); err != nil {
		logger.Error("Failed to store deposits", "block", blockNum, "error", err)
	}
}

// SetDepositsOnBlockBody sets the deposits on the block body, used by the block builder.
func SetDepositsOnBlockBody(
	ctx context.Context,
	st *statedb.StateDB,
	body *ctypes.BeaconBlockBody,
	chainSpec ChainSpec,
	depositStore depositstore.StoreManager,
	logger log.Logger,
) error {
	// Dequeue deposits from the state.
	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return fmt.Errorf("failed loading eth1 deposit index: %w", err)
	}

	forkVersion := body.GetForkVersion()
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}
	prevBlockForkVersion := chainSpec.ActiveForkVersionForTimestamp(lph.GetTimestamp())

	// Only if called before Fulu or on the first block of Fulu, do we set deposits on the
	// block body.
	var depositRange uint64
	switch {
	case version.IsBefore(forkVersion, version.Fulu()):
		depositRange = depositIndex + chainSpec.MaxDepositsPerBlock()
	case version.Equals(prevBlockForkVersion, version.Electra1()) &&
		version.Equals(forkVersion, version.Fulu()):
		// For the first block of Fulu catchup deposits, we will include as many as are required to exhaust the
		// queue. Since after this block in Fulu, we no longer use the deposit queue and
		// instead follow EIP-6110 deposit requests.
		depositRange = stdmath.MaxUint64
	default:
		// We don't set deposits on the block body after the first block of Fulu.
		return nil
	}

	// Grab all previous deposits from genesis up to the current index + max deposits per block.
	deposits, localDepositRoot, err := depositStore.GetDepositsByIndex(
		ctx, constants.FirstDepositIndex, depositRange,
	)
	if err != nil {
		return err
	}
	if uint64(len(deposits)) < depositIndex {
		return errors.Wrapf(ErrDepositStoreIncomplete,
			"all historical deposits not available, expected: %d, got: %d",
			depositIndex, len(deposits),
		)
	}

	logger.Info(
		"Building block body with local deposits",
		"start_index", depositIndex, "num_deposits", uint64(len(deposits))-depositIndex,
	)
	body.SetEth1Data(ctypes.NewEth1Data(localDepositRoot))
	body.SetDeposits(deposits[depositIndex:])
	return nil
}

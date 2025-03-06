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

package core

import (
	"context"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
)

func validateGenesisDeposits(
	st *statedb.StateDB, deposits []*ctypes.Deposit, validatorSetCap uint64,
) error {
	eth1DepositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}
	if eth1DepositIndex != constants.FirstDepositIndex {
		return errors.New("Eth1DepositIndex should be 0 at genesis")
	}

	if len(deposits) == 0 {
		// there should be at least a validator in genesis
		return errors.Wrap(ErrDepositsLengthMismatch, "at least one validator should be in genesis")
	}
	for i, deposit := range deposits {
		// deposit indices should be contiguous
		// #nosec G115
		if deposit.GetIndex() != math.U64(i) {
			return errors.Wrapf(ErrDepositIndexOutOfOrder,
				"genesis deposit index: %d, expected index: %d", deposit.GetIndex().Unwrap(), i,
			)
		}
	}

	// BeaconKit enforces a cap on the validator set size.
	// If genesis deposits breaches the cap we return an error.
	//#nosec:G701 // can't overflow.
	if uint64(len(deposits)) > validatorSetCap {
		return errors.Wrapf(
			ErrValSetCapExceeded,
			"validator set cap %d, deposits count %d",
			validatorSetCap, len(deposits),
		)
	}
	return nil
}

func ValidateNonGenesisDeposits(
	ctx context.Context,
	st *statedb.StateDB,
	depositStore *depositdb.KVStore,
	maxDepositsPerBlock uint64,
	blkDeposits []*ctypes.Deposit,
	blkDepositRoot common.Root,
) error {
	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}

	// Grab all previous deposits from genesis up to the current index + max deposits per block.
	localDeposits, err := depositStore.GetDepositsByIndex(
		ctx,
		constants.FirstDepositIndex,
		depositIndex+maxDepositsPerBlock,
	)
	if err != nil {
		return err
	}

	// First verify that the number of block deposits matches the number of local deposits.
	totalBlockDeposits := depositIndex + uint64(len(blkDeposits))
	if uint64(len(localDeposits)) != totalBlockDeposits {
		return errors.Wrapf(ErrDepositsLengthMismatch,
			"block deposit count: %d, expected deposit count: %d",
			totalBlockDeposits, len(localDeposits),
		)
	}

	// Then check that the block's deposits 1) have contiguous indices and 2) match the local
	// view of the block's deposits.
	for i, blkDeposit := range blkDeposits {
		blkDepositIndex := blkDeposit.GetIndex().Unwrap()
		//#nosec:G115 // won't overflow in practice.
		if blkDepositIndex != depositIndex+uint64(i) {
			return errors.Wrapf(ErrDepositIndexOutOfOrder,
				"deposit index: %d, expected index: %d", blkDepositIndex, i,
			)
		}

		if !localDeposits[blkDepositIndex].Equals(blkDeposit) {
			return errors.Wrapf(ErrDepositMismatch,
				"deposit index: %d, expected deposit: %+v, actual deposit: %+v",
				blkDepositIndex, *localDeposits[blkDepositIndex], *blkDeposit,
			)
		}
	}

	// Finally check that the historical deposits root matches locally what's on the beacon block.
	if !localDeposits.HashTreeRoot().Equals(blkDepositRoot) {
		return ErrDepositsRootMismatch
	}

	return nil
}

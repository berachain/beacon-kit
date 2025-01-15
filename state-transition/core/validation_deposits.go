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

package core

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

func (sp *StateProcessor[
	_, _,
]) validateGenesisDeposits(
	st *statedb.StateDB,
	deposits []*ctypes.Deposit,
) error {
	eth1DepositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}
	if eth1DepositIndex != 0 {
		return errors.New("Eth1DepositIndex should be 0 at genesis")
	}

	if len(deposits) == 0 {
		// there should be at least a validator in genesis
		return errors.Wrap(ErrDepositsLengthMismatch, "at least one validator should be in genesis")
	}
	for i, deposit := range deposits {
		// deposit indices should be contiguous
		if deposit.GetIndex() != math.U64(i) {
			return errors.Wrapf(ErrDepositIndexOutOfOrder,
				"genesis deposit index: %d, expected index: %d", deposit.GetIndex().Unwrap(), i,
			)
		}
	}

	// BeaconKit enforces a cap on the validator set size.
	// If genesis deposits breaches the cap we return an error.
	//#nosec:G701 // can't overflow.
	if uint64(len(deposits)) > sp.cs.ValidatorSetCap() {
		return errors.Wrapf(ErrValSetCapExceeded,
			"validator set cap %d, deposits count %d", sp.cs.ValidatorSetCap(), len(deposits),
		)
	}
	return nil
}

func (sp *StateProcessor[
	_, _,
]) validateNonGenesisDeposits(
	st *statedb.StateDB,
	blkDeposits []*ctypes.Deposit,
	blkDepositRoot common.Root,
) error {
	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}
	for i, deposit := range blkDeposits {
		// deposit indices should be contiguous
		if deposit.GetIndex() != math.U64(depositIndex)+math.U64(i) {
			return errors.Wrapf(ErrDepositIndexOutOfOrder,
				"deposit index: %d, expected index: %d", deposit.GetIndex().Unwrap(), i,
			)
		}
	}

	var deposits ctypes.Deposits
	deposits, err = sp.ds.GetDepositsByIndex(0, depositIndex+uint64(len(blkDeposits)))
	if err != nil {
		return err
	}

	if !blkDepositRoot.Equals(deposits.HashTreeRoot()) {
		return ErrDepositsRootMismatch
	}
	return nil
}

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
	"fmt"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT,
	_, _, _, _, _, _, _, _, _, _, _,
]) validateGenesisDeposits(
	st BeaconStateT,
	deposits []DepositT,
) error {
	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// Bartio does not properly validate deposits index
		// We skip checks for backward compatibility
		return nil

	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID:
		// Boonet inherited the bug from Bartio and it may have added some
		// validators before we activate the fork. So we skip validation
		// before fork activation
		return nil

	default:
		if _, err := st.GetEth1DepositIndex(); err == nil {
			// there should not be Eth1DepositIndex stored before
			// genesis first deposit
			return errors.Wrap(
				ErrDepositMismatch,
				"Eth1DepositIndex should be unset at genesis",
			)
		}
		if len(deposits) == 0 {
			// there should be at least a validator in genesis
			return errors.Wrap(
				ErrDepositsLengthMismatch,
				"at least one validator should be in genesis",
			)
		}
		for i, deposit := range deposits {
			// deposit indices should be contiguous
			if deposit.GetIndex() != math.U64(i) {
				return errors.Wrapf(
					ErrDepositIndexOutOfOrder,
					"genesis deposit index: %d, expected index: %d",
					deposit.GetIndex().Unwrap(), i,
				)
			}
		}
		return nil
	}
}

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT,
	_, _, _, _, _, _, _, _, _, _, _,
]) validateNonGenesisDeposits(
	st BeaconStateT,
	deposits []DepositT,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return fmt.Errorf(
			"failed loading slot while processing deposits: %w", err,
		)
	}
	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// Bartio does not properly validate deposits index
		// We skip checks for backward compatibility
		return nil

	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork2Height):
		// Boonet inherited the bug from Bartio and it may have added some
		// validators before we activate the fork. So we skip validation
		// before fork activation
		return nil

	default:
		// Verify that outstanding deposits match those listed by contract
		var depositIndex uint64
		depositIndex, err = st.GetEth1DepositIndex()
		if err != nil {
			return err
		}
		expectedStartIdx := depositIndex + 1

		var localDeposits []DepositT
		localDeposits, err = sp.ds.GetDepositsByIndex(
			expectedStartIdx,
			sp.cs.MaxDepositsPerBlock(),
		)
		if err != nil {
			return err
		}

		sp.logger.Info(
			"processOperations",
			"Expected deposit start index", expectedStartIdx,
			"Expected deposits length", len(localDeposits),
		)

		if len(localDeposits) != len(deposits) {
			return errors.Wrapf(
				ErrDepositsLengthMismatch,
				"local: %d, payload: %d", len(localDeposits), len(deposits),
			)
		}

		for i, sd := range localDeposits {
			// Deposit indices should be contiguous
			//#nosec:G701 // i never negative
			expectedIdx := expectedStartIdx + uint64(i)
			if sd.GetIndex().Unwrap() != expectedIdx {
				return errors.Wrapf(
					ErrDepositIndexOutOfOrder,
					"local deposit index: %d, expected index: %d",
					sd.GetIndex().Unwrap(), expectedIdx,
				)
			}

			if !sd.Equals(deposits[i]) {
				return errors.Wrapf(
					ErrDepositMismatch,
					"local deposit: %+v, payload deposit: %+v",
					sd, deposits[i],
				)
			}
		}

		return nil
	}
}

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
)

func (sp *StateProcessor[
	_, BeaconStateT, _, _,
]) validateNonGenesisDeposits(
	st BeaconStateT,
	deposits []*ctypes.Deposit,
) error {
	// Verify that outstanding deposits match those listed by contract
	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}

	var localDeposits ctypes.Deposits
	localDeposits, _, err = sp.ds.GetDepositsByIndex(
		depositIndex, sp.cs.MaxDepositsPerBlock(),
	)
	if err != nil {
		return err
	}

	sp.logger.Info(
		"Processing deposits in range",
		"expected_start_index", depositIndex, "expected_range_length", len(localDeposits),
	)

	if len(localDeposits) != len(deposits) {
		return errors.Wrapf(
			ErrDepositsLengthMismatch,
			"local: %d, payload: %d", len(localDeposits), len(deposits),
		)
	}

	for i, sd := range localDeposits {
		// DepositData indices should be contiguous
		//#nosec:G701 // i never negative
		expectedIdx := depositIndex + uint64(i)
		if sd.Data.GetIndex().Unwrap() != expectedIdx {
			return errors.Wrapf(
				ErrDepositIndexOutOfOrder,
				"local deposit index: %d, expected index: %d",
				sd.Data.GetIndex().Unwrap(), expectedIdx,
			)
		}

		if !sd.Data.Equals(deposits[i].Data) {
			return errors.Wrapf(
				ErrDepositMismatch,
				"local deposit: %+v, payload deposit: %+v", sd, deposits[i],
			)
		}
	}

	return nil
}

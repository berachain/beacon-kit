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
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// InitiateValidatorExit initiates the exit of the validator with index `idx`.
func (sp *StateProcessor) InitiateValidatorExit(st *statedb.StateDB, idx math.ValidatorIndex) error {
	validator, err := st.ValidatorByIndex(idx)
	if err != nil {
		return err
	}
	// We will use the fork version from the state to determine how to exit the validator.
	fork, err := st.GetFork()
	if err != nil {
		return err
	}
	var withdrawableEpoch, exitEpoch math.Epoch
	if version.EqualsOrIsAfter(fork.CurrentVersion, version.Electra()) {
		// Return if the validator already initiated an exit, making sure to only exit validators once.
		if validator.GetExitEpoch() != constants.FarFutureEpoch {
			return nil
		}
		slot, slotErr := st.GetSlot()
		if slotErr != nil {
			return slotErr
		}
		nextEpoch := sp.cs.SlotToEpoch(slot) + 1
		// We continue to have no cap on validator churn, choosing not to adopt any churn-related electra changes.
		exitEpoch = nextEpoch
		// The withdrawable Epoch is `MinValidatorWithdrawabilityDelay` epoch's after `exitEpoch`.
		withdrawableEpoch = math.Epoch(uint64(exitEpoch) + sp.cs.MinValidatorWithdrawabilityDelay())
	} else {
		// Before Electra, this was the logic for exiting a validator.
		// It would only be triggered if the maximum validator set size was reached.
		// It did not add the `MinValidatorWithdrawabilityDelay`.
		slot, slotErr := st.GetSlot()
		if slotErr != nil {
			return slotErr
		}
		nextEpoch := sp.cs.SlotToEpoch(slot) + 1
		exitEpoch = nextEpoch
		// The withdrawable Epoch is the next epoch after `exitEpoch`.
		withdrawableEpoch = nextEpoch + 1
	}

	// Set validator exit epoch and withdrawable epoch.
	validator.SetExitEpoch(exitEpoch)
	validator.SetWithdrawableEpoch(withdrawableEpoch)
	if updateErr := st.UpdateValidatorAtIndex(idx, validator); updateErr != nil {
		return updateErr
	}
	return nil
}

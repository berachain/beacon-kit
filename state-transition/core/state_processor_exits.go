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
	// Return if the validator already initiated an exit, making sure to only exit validators once.
	validator, err := st.ValidatorByIndex(idx)
	if err != nil {
		return err
	}
	if validator.GetExitEpoch() != constants.FarFutureEpoch {
		return nil
	}

	// We will use the fork version from the state to determine how to exit the validator.
	fork, err := st.GetFork()
	if err != nil {
		return err
	}
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// We still have no cap on validator churn, choosing not to adopt any churn-related
	// Electra changes, so exit epoch is at the next epoch.
	exitEpoch := sp.cs.SlotToEpoch(slot) + 1

	// The withdrawable epoch is `MinValidatorWithdrawabilityDelay` epoch's after `exitEpoch`.
	var withdrawableEpoch math.Epoch
	if version.IsBefore(fork.CurrentVersion, version.Electra()) {
		// Before Electra, this was the logic for exiting a validator: only trigger if the validator
		// set cap was reached. The withdrawable epoch does not include
		// `MinValidatorWithdrawabilityDelay`, but is instead the next epoch after exiting.
		withdrawableEpoch = exitEpoch + 1
	} else {
		// The withdrawable Epoch is `MinValidatorWithdrawabilityDelay` epoch's after `exitEpoch`.
		withdrawableEpoch = exitEpoch + sp.cs.MinValidatorWithdrawabilityDelay()
	}

	// Set validator exit epoch and withdrawable epoch.
	validator.SetExitEpoch(exitEpoch)
	validator.SetWithdrawableEpoch(withdrawableEpoch)
	return st.UpdateValidatorAtIndex(idx, validator)
}

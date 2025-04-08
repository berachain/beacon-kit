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
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// GetActivationExitChurnLimit returns the churn limit for the current epoch dedicated to activations and exits.
func (sp *StateProcessor) GetActivationExitChurnLimit(st *statedb.StateDB) math.Gwei {
	// TODO(pectra): get this value from config or constant
	var maxPerEpochActivationExitChurnLimit math.Gwei
	return min(maxPerEpochActivationExitChurnLimit, st.GetBalanceChurnLimit())
}

func (sp *StateProcessor) ComputeActivationExitEpoch(epoch math.Epoch) math.Epoch {
	// TODO(pectra): get this value from config or constant
	var maxSeedLookahead math.Epoch
	return epoch + 1 + maxSeedLookahead
}

func (sp *StateProcessor) ComputeExitEpochAndUpdateChurn(st *statedb.StateDB, exitBalance math.Gwei) (math.Epoch, error) {
	slot, err := st.GetSlot()
	if err != nil {
		return 0, err
	}
	earliestExitEpoch := max(st.GetEarliestExitEpoch(), sp.ComputeActivationExitEpoch(sp.cs.SlotToEpoch(slot)))
	perEpochChurn := sp.GetActivationExitChurnLimit(st)
	var exitBalanceToConsume math.Gwei
	if st.GetEarliestExitEpoch() < earliestExitEpoch {
		exitBalanceToConsume = perEpochChurn
	} else {
		exitBalanceToConsume = st.GetExitBalanceToConsume()
	}

	if exitBalance > exitBalanceToConsume {
		balanceToProcess := exitBalance - exitBalanceToConsume
		additionalEpochs := (balanceToProcess-1)/perEpochChurn + 1
		earliestExitEpoch += additionalEpochs
		exitBalanceToConsume += additionalEpochs * perEpochChurn
	}

	// Consume the balance and update state variables.
	st.SetExitBalanceToConsume(exitBalanceToConsume - exitBalance)
	st.SetEarliestExitEpoch(earliestExitEpoch)
	return st.GetEarliestExitEpoch(), nil
}

// InitiateValidatorExit initiates the exit of the validator with index `idx`.
func (sp *StateProcessor) InitiateValidatorExit(st *statedb.StateDB, idx math.ValidatorIndex) error {
	validator, err := st.ValidatorByIndex(idx)
	if err != nil {
		return err
	}
	// Return if validator already initiated exit.
	if validator.GetExitEpoch() != math.Epoch(constants.FarFutureEpoch) {
		return nil
	}
	// Compute the exit queue epoch.
	exitQueueEpoch, err := sp.ComputeExitEpochAndUpdateChurn(st, validator.GetEffectiveBalance())
	if err != nil {
		return err
	}

	// Set validator exit epoch and withdrawable epoch.
	validator.SetExitEpoch(exitQueueEpoch)
	// TODO(pectra): Get chainspec value
	minValidatorWithdrawabilityDelay := math.Epoch(0)
	validator.SetWithdrawableEpoch(validator.GetExitEpoch() + minValidatorWithdrawabilityDelay)
	err = st.UpdateValidatorAtIndex(idx, validator)
	if err != nil {
		return err
	}
	return nil
}

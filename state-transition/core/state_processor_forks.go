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
	"fmt"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// PrepareStateForFork prepares the state for the fork version at the given timestamp.
//   - If this function is called for the same version as the state's current version,
//     it will do nothing.
//   - If this function is called for a version before the state's current version,
//     it will return error as this is not allowed.
//   - If this function is called for a version after the state's current version,
//     it will upgrade the state to the new version.
func (sp *StateProcessor) PrepareStateForFork(
	st *statedb.StateDB, timestamp math.U64, slot math.Slot, logUpgrade bool,
) error {
	stateFork, err := st.GetFork()
	if err != nil {
		return err
	}

	// Return early if the given fork version is before or equal to the current state fork version.
	forkVersion := sp.cs.ActiveForkVersionForTimestamp(timestamp)
	if version.IsBefore(forkVersion, stateFork.CurrentVersion) {
		return fmt.Errorf(
			"cannot downgrade state from %s to %s", stateFork.CurrentVersion, forkVersion,
		)
	} else if version.Equals(forkVersion, stateFork.CurrentVersion) {
		return nil
	}

	// Upgrade the state to the new version.
	switch forkVersion {
	case version.Deneb():
		return sp.upgradeToDeneb(stateFork.CurrentVersion, timestamp, slot, logUpgrade)
	case version.Deneb1():
		return sp.upgradeToDeneb1(stateFork.CurrentVersion, timestamp, slot, logUpgrade)
	case version.Electra():
		return sp.upgradeToElectra(st, stateFork, timestamp, slot, logUpgrade)
	default:
		return fmt.Errorf("unsupported fork version: %s", forkVersion)
	}
}

// Do nothing to the state. NOTE: Deneb is the genesis version of Berachain
// mainnet and Bepolia testnet. At genesis, InitializePreminedBeaconStateFromEth1
// should be called, which adequately prepares the BeaconState for Deneb.
func (sp *StateProcessor) upgradeToDeneb(
	previousVersion common.Version, timestamp math.U64, slot math.Slot, logUpgrade bool,
) error {
	// Log the upgrade to Deneb if requested.
	if logUpgrade {
		sp.logger.Info(fmt.Sprintf(`


	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️

	+ ✅  upgraded to deneb (0x04000000) fork! 🎉
	+ 🚝  previous fork: %s (%s)
	+ ⏱️   deneb fork time: %d
	+ 🍴  first slot / timestamp of deneb: %d / %d
	+ ⛓️   current beacon epoch: %d

	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️


`,
			version.Name(previousVersion), previousVersion.String(),
			timestamp.Unwrap(), // TODO: this should be fetched from the chain spec.
			slot.Unwrap(), timestamp.Unwrap(),
			sp.cs.SlotToEpoch(slot).Unwrap(),
		))
	}
	return nil
}

// upgradeToDeneb1 does nothing to the state. NOTE: Deneb1 is the first hard fork of Berachain
// mainnet and Bepolia testnet. In this fork, the Fork struct on BeaconState is NOT updated. In
// future hard forks, the Fork struct should be updated.
func (sp *StateProcessor) upgradeToDeneb1(
	previousVersion common.Version, timestamp math.U64, slot math.Slot, logUpgrade bool,
) error {
	// Log the upgrade to Deneb1 if requested.
	if logUpgrade {
		sp.logger.Info(fmt.Sprintf(`


	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️

	+ ✅  upgraded to deneb1 (0x04010000) fork! 🎉
	+ 🚝  previous fork: %s (%s)
	+ ⏱️   deneb1 fork time: %d
	+ 🍴  first slot / timestamp of deneb1: %d / %d
	+ ⛓️   current beacon epoch: %d

	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️


`,
			version.Name(previousVersion), previousVersion.String(),
			sp.cs.Deneb1ForkTime(),
			slot.Unwrap(), timestamp.Unwrap(),
			sp.cs.SlotToEpoch(slot).Unwrap(),
		))
	}
	return nil
}

// upgradeToElectra upgrades the state to the Electra fork version. It is modified from the ETH 2.0
// spec (https://ethereum.github.io/consensus-specs/specs/electra/fork/#upgrading-the-state)
// to only upgrade the Fork struct in the BeaconState.
func (sp *StateProcessor) upgradeToElectra(
	st *statedb.StateDB, fork *types.Fork, timestamp math.U64, slot math.Slot, logUpgrade bool,
) error {
	// Set the fork on BeaconState.
	fork.PreviousVersion = fork.CurrentVersion
	fork.CurrentVersion = version.Electra()
	fork.Epoch = sp.cs.SlotToEpoch(slot)
	if err := st.SetFork(fork); err != nil {
		return err
	}

	// Initialize the pending partial withdrawals to an empty array.
	if err := st.SetPendingPartialWithdrawals([]*types.PendingPartialWithdrawal{}); err != nil {
		return err
	}

	// Log the upgrade to Electra if requested.
	if logUpgrade {
		sp.logger.Info(fmt.Sprintf(`


	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️

	+ ✅  upgraded to electra (0x05000000) fork! 🎉
	+ 🚝  previous fork: %s (%s)
	+ ⏱️   electra fork time: %d
	+ 🍴  first slot / timestamp of electra: %d / %d
	+ ⛓️   current beacon epoch: %d

	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️


`,
			version.Name(fork.PreviousVersion), fork.PreviousVersion.String(),
			sp.cs.ElectraForkTime(),
			slot.Unwrap(), timestamp.Unwrap(),
			fork.Epoch.Unwrap(),
		))
	}
	return nil
}

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
	"sync"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// ProcessFork prepares the state for the fork version at the given timestamp.
//   - If this function is called for the same version as the state's current version,
//     it will do nothing. Unless it is the genesis slot, in which case we want to
//     prepare the state for the genesis fork version.
//   - If this function is called for a version before the state's current version,
//     it will return error as this is not allowed.
//   - If this function is called for a version after the state's current version,
//     it will upgrade the state to the new version.
//
// NOTE for caller: `ProcessSlots` must be called before this function. If we are
// crossing into a new fork, the first slot of the new fork will be retrieved from
// the state. The state must be prepared for this new slot.
func (sp *StateProcessor) ProcessFork(
	st *statedb.StateDB, timestamp math.U64, logUpgrade bool,
) error {
	stateFork, err := st.GetFork()
	if err != nil {
		return err
	}
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Return early if the given fork version is before or equal to the current state fork version.
	forkVersion := sp.cs.ActiveForkVersionForTimestamp(timestamp)
	if version.IsBefore(forkVersion, stateFork.CurrentVersion) {
		return fmt.Errorf(
			"cannot downgrade state from %s to %s", stateFork.CurrentVersion, forkVersion,
		)
	} else if slot > 0 && version.Equals(forkVersion, stateFork.CurrentVersion) {
		// If we are past genesis and the fork version remains consistent, do nothing.
		return nil
	}

	// If we are at genesis or moving to a new fork version, upgrade the state.
	switch forkVersion {
	case version.Deneb():
		// Do nothing to the state. NOTE: Deneb is the genesis version of Berachain mainnet and
		// Bepolia testnet.

		// Log the upgrade to Deneb if requested.
		if logUpgrade {
			sp.logDenebFork(timestamp)
		}
	case version.Deneb1():
		// Do nothing to the state. NOTE: Deneb1 is the first hard fork of Berachain mainnet and
		// Bepolia testnet. In this fork, the Fork struct on BeaconState is NOT updated. In
		// future hard forks, the Fork struct should be updated.

		// Log the upgrade to Deneb1 if requested.
		if logUpgrade {
			sp.logDeneb1Fork(stateFork.PreviousVersion, timestamp, slot)
		}
	case version.Electra():
		if err = sp.upgradeToElectra(st, stateFork, slot); err != nil {
			return err
		}

		// Log the upgrade to Electra if requested.
		if logUpgrade {
			sp.logElectraFork(stateFork.PreviousVersion, timestamp, slot)
		}
	default:
		panic(fmt.Sprintf("unsupported fork version: %s", forkVersion))
	}

	return nil
}

// logDenebFork logs information about the Deneb fork.
func (sp *StateProcessor) logDenebFork(timestamp math.U64) {
	// Since Deneb is the earliest fork version supported by beacon-kit, if we are
	// entering Deneb it must be at genesis, which means the fork time of Deneb is
	// the timestamp of the genesis block itself.
	denebForkTime := timestamp.Unwrap()

	sp.logger.Info(fmt.Sprintf(`


	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️

	+ ✅  welcome to the deneb (0x04000000) fork! 🎉
	+ ⏱️   deneb fork time: %d
	+ 🍴  first slot / timestamp of deneb: %d / %d
	+ ⛓️   current beacon epoch: %d

	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️


`,
		denebForkTime,
		constants.GenesisSlot.Unwrap(), denebForkTime,
		constants.GenesisEpoch.Unwrap(),
	))
}

// logDeneb1Fork logs information about the Deneb1 fork.
func (sp *StateProcessor) logDeneb1Fork(
	previousVersion common.Version, timestamp math.U64, slot math.Slot,
) {
	// Since state fork is not updating to Deneb1, every block observes Deneb1 as "new fork" during
	// Deneb1. Hence, we must wrap this in a OnceFunc to ensure it is logged only the first time
	// we process a Deneb1 block.
	sync.OnceFunc(func() {
		sp.logger.Info(fmt.Sprintf(`


	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️

	+ ✅  welcome to the deneb1 (0x04010000) fork! 🎉
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
	})()
}

// upgradeToElectra upgrades the state to the Electra fork version. It is modified from the ETH 2.0
// spec (https://ethereum.github.io/consensus-specs/specs/electra/fork/#upgrading-the-state) to:
//   - update the Fork struct in the BeaconState
//   - initialize the pending partial withdrawals to an empty array
func (sp *StateProcessor) upgradeToElectra(
	st *statedb.StateDB, fork *types.Fork, slot math.Slot,
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

	return nil
}

// logElectraFork logs information about the Electra fork.
func (sp *StateProcessor) logElectraFork(
	previousVersion common.Version, timestamp math.U64, slot math.Slot,
) {
	sp.logger.Info(fmt.Sprintf(`


	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️

	+ ✅  welcome to the electra (0x05000000) fork! 🎉
	+ 🚝  previous fork: %s (%s)
	+ ⏱️   electra fork time: %d
	+ 🍴  first slot / timestamp of electra: %d / %d
	+ ⛓️   current beacon epoch: %d

	⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️⏭️


`,
		version.Name(previousVersion), previousVersion.String(),
		sp.cs.ElectraForkTime(),
		slot.Unwrap(), timestamp.Unwrap(),
		sp.cs.SlotToEpoch(slot).Unwrap(),
	))
}

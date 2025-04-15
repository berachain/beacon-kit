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

// prepareStateForFork prepares the state for the given fork version.
//   - If this function is called for the same version as the state's current version,
//     it will do nothing.
//   - If this function is called for a version before the state's current version,
//     it will return error as this is not allowed.
//   - If this function is called for a version after the state's current version,
//     it will upgrade the state to the new version.
func (sp *StateProcessor) prepareStateForFork(
	st *statedb.StateDB, forkVersion common.Version, slot math.Slot,
) error {
	stateFork, err := st.GetFork()
	if err != nil {
		return err
	}

	// Return early if the given fork version is before or equal to the current state fork version.
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
		// Do nothing. NOTE: Deneb is the genesis version of Berachain.
		// At genesis, InitializePreminedBeaconStateFromEth1 should be called,
		// which adequately prepares the BeaconState for Deneb.
	case version.Deneb1():
		// Do nothing. NOTE: Deneb1 is the first hard fork of Berachain.
		// In this fork, the Fork struct on BeaconState is NOT updated.
		// In future hard forks, the Fork struct WILL be updated.
	case version.Electra():
		return sp.upgradeToElectra(st, stateFork, slot)
	default:
		return fmt.Errorf("unsupported fork version: %s", forkVersion)
	}

	return nil
}

// upgradeToElectra upgrades the state to the Electra fork version. It is modified from the ETH 2.0
// spec (https://ethereum.github.io/consensus-specs/specs/electra/fork/#upgrading-the-state)
// to only upgrade the Fork struct in the BeaconState.
func (sp *StateProcessor) upgradeToElectra(
	st *statedb.StateDB, fork *types.Fork, slot math.Slot,
) error {
	fork.PreviousVersion = fork.CurrentVersion
	fork.CurrentVersion = version.Electra()
	fork.Epoch = sp.cs.SlotToEpoch(slot)

	return st.SetFork(fork)
}

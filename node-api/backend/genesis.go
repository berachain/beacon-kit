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

package backend

import (
	"fmt"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// GetGenesis returns the genesis state of the beacon chain.
func (b Backend) GenesisValidatorsRoot(slot math.Slot) (common.Root, error) {
	// needs genesis_time and genesis_fork_version
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return common.Root{}, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}
	return st.GetGenesisValidatorsRoot()
}

// GenesisForkVersion returns the genesis fork version of the beacon chain.
func (b Backend) GenesisForkVersion(genesisSlot math.Slot) (common.Version, error) {
	st, _, err := b.stateFromSlot(genesisSlot)
	if err != nil {
		return common.Version{}, errors.Wrapf(err, "failed to get state from slot %d", genesisSlot)
	}
	fork, err := st.GetFork()
	if err != nil {
		return common.Version{}, errors.Wrapf(err, "failed to get fork from state")
	}
	return fork.CurrentVersion, nil
}

// GenesisTime returns the genesis time of the beacon chain.
func (b Backend) GenesisTime(genesisSlot math.Slot) (math.U64, error) {
	st, _, err := b.stateFromSlot(genesisSlot)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get state from slot %d", genesisSlot)
	}
	genesisTime, err := st.GetGenesisTime()
	fmt.Println("genesisTime in backend", genesisTime)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get genesis time from state")
	}
	return genesisTime, nil
}

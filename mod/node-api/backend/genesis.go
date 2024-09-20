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

package backend

import (
	"fmt"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// GetGenesis returns the genesis state of the beacon chain.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GenesisValidatorsRoot(slot math.Slot) (common.Root, error) {
	// needs genesis_time and gensis_fork_version
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}
	return st.GetGenesisValidatorsRoot()
}

// GetGenesisForkVersion returns the genesis fork version of the beacon chain.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetGenesisForkVersion(genesisSlot math.Slot) (common.Version, error) {
	st, _, err := b.stateFromSlot(genesisSlot)
	if err != nil {
		return common.Version{}, err
	}
	fork, err := st.GetFork()
	if err != nil {
		return common.Version{}, err
	}
	return fork.GetPreviousVersion(), nil
}

// GetGenesisTime returns the genesis time of the beacon chain.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetGenesisTime(slot math.Slot) (uint64, error) {

	//state := b.sb.StateFromContext(ctx)
	//if state == nil {
	//	fmt.Errorf("Failed to retrieve beacon state")
	//	return 0, errors.New("failed to retrieve beacon state")
	//}
	//
	//genesisTime, err := state.GetGenesisTime()
	//if err != nil {
	//	fmt.Errorf("Failed to get genesis time", err)
	//	return 0, err
	//}
	//
	//fmt.Printf("Retrieved genesis time %v", genesisTime)
	//return genesisTime, nil

	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return 0, err
	}

	fmt.Printf("state nidzi %v", st)
	return st.GetGenesisTime()
}

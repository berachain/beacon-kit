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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// GenesisTime returns the genesis time of the beacon chain.
func (b *Backend) GenesisTime() math.U64 {
	return *b.genesisTime.Load()
}

// GenesisState returns the genesis state of the beacon chain.
func (b *Backend) GenesisState() *statedb.StateDB {
	return b.stateFetcher.GetGenesisState()
}

// GenesisBlockHeader returns the genesis block header of the beacon chain.
func (b *Backend) GenesisBlockHeader() *ctypes.BeaconBlockHeader {
	genesisState := b.GenesisState()
	if genesisState == nil {
		return nil
	}

	// For genesis state, the latest block header is the genesis block header
	header, err := genesisState.GetLatestBlockHeader()
	if err != nil {
		return nil
	}
	return header
}

// GenesisForkVersion returns the genesis fork version of the beacon chain.
func (b *Backend) GenesisForkVersion() common.Version {
	// Derive the genesis fork version from the genesis time.
	genesisTime := b.GenesisTime()

	return b.cs.ActiveForkVersionForTimestamp(genesisTime)
}

// GenesisBlockRoot returns the genesis block root of the beacon chain.
func (b *Backend) GenesisBlockRoot() common.Root {
	genesisState := b.GenesisState()
	if genesisState == nil {
		return common.Root{}
	}
	// Return the hash tree root of the genesis header after updating the state root in it.
	// This is similar to how we get the block root.
	genesisHeader, err := genesisState.GetLatestBlockHeader()
	if err != nil {
		return common.Root{}
	}
	genesisHeader.SetStateRoot(genesisState.HashTreeRoot())
	return genesisHeader.HashTreeRoot()
}

// GenesisValidators returns the genesis validators of the beacon chain.
func (b *Backend) GenesisValidators() []*ctypes.Validator {
	genesisState := b.GenesisState()
	if genesisState == nil {
		return nil
	}
	validators, err := genesisState.GetValidators()
	if err != nil {
		return nil
	}
	return validators
}

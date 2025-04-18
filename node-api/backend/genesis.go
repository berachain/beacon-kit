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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// GenesisValidatorsRoot returns the genesis validators root of the beacon chain.
func (b *Backend) GenesisValidatorsRoot() (common.Root, error) {
	// First check if the value is cached.
	root := b.genesisValidatorsRoot.Load()
	if root != nil && *root != (common.Root{}) {
		return *root, nil
	}

	// If not cached, read state from the beacon state at the tip of chain.
	st, _, err := b.StateAtSlot(utils.Head)
	if err != nil {
		return common.Root{}, errors.Wrapf(err, "failed to get state from tip of chain")
	}

	// Get the genesis validators root.
	validatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return common.Root{}, errors.Wrap(err, "failed to get genesis validators root from state")
	}

	// Cache the value for future use.
	b.genesisValidatorsRoot.Store(&validatorsRoot)

	return validatorsRoot, nil
}

// GenesisForkVersion returns the genesis fork version of the beacon chain.
func (b *Backend) GenesisForkVersion() (common.Version, error) {
	return *b.genesisForkVersion.Load(), nil
}

// GenesisTime returns the genesis time of the beacon chain.
func (b *Backend) GenesisTime() (math.U64, error) {
	return *b.genesisTime.Load(), nil
}

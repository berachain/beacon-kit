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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// GetGenesisValidatorsRoot returns the genesis validators root of the beacon chain.
func (b *Backend) GenesisValidatorsRoot() (common.Root, error) {
	// Fast path: read lock for checking cached value
	b.genesisValidatorsRootMu.RLock()
	if b.genesisValidatorsRoot != (common.Root{}) {
		root := b.genesisValidatorsRoot
		b.genesisValidatorsRootMu.RUnlock()
		return root, nil
	}
	b.genesisValidatorsRootMu.RUnlock()

	// Slow path: write lock for initialization
	b.genesisValidatorsRootMu.Lock()
	defer b.genesisValidatorsRootMu.Unlock()

	// Double check after acquiring write lock
	if b.genesisValidatorsRoot != (common.Root{}) {
		return b.genesisValidatorsRoot, nil
	}

	// If not cached, read state from the genesis slot
	st, _, err := b.stateFromSlot(0)
	if err != nil {
		return common.Root{}, errors.Wrapf(err, "failed to get state from tip of chain")
	}
	// Get the genesis validators root
	root, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return common.Root{}, errors.Wrap(err, "failed to get genesis validators root from state")
	}

	// Cache the value for future use
	b.genesisValidatorsRoot = root

	return root, nil
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
	// Get the execution payload header from the beacon state,
	// and return the timestamp as the genesis time.
	execPayloadHeader, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get execution payload header from state")
	}
	return execPayloadHeader.Timestamp, nil
}

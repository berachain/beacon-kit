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

package beacon

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Backend is the interface for backend of the beacon API.
type Backend interface {
	GenesisBackend
	BlockBackend
	RandaoBackend
	StateBackend
	ValidatorBackend
	HistoricalBackend
	// GetSlotByBlockRoot retrieves the slot by a given root from the store.
	GetSlotByBlockRoot(root common.Root) (math.Slot, error)
	// GetSlotByStateRoot retrieves the slot by a given root from the store.
	GetSlotByStateRoot(root common.Root) (math.Slot, error)
}

type GenesisBackend interface {
	GenesisValidatorsRoot(slot math.Slot) (common.Root, error)
}

type HistoricalBackend interface {
	StateRootAtSlot(slot math.Slot) (common.Root, error)
	StateForkAtSlot(slot math.Slot) (*ctypes.Fork, error)
}

type RandaoBackend interface {
	RandaoAtEpoch(slot math.Slot, epoch math.Epoch) (common.Bytes32, error)
}

type BlockBackend interface {
	BlockRootAtSlot(slot math.Slot) (common.Root, error)
	BlockRewardsAtSlot(slot math.Slot) (*types.BlockRewardsData, error)
	BlockHeaderAtSlot(slot math.Slot) (*ctypes.BeaconBlockHeader, error)
}

type StateBackend interface {
	StateRootAtSlot(slot math.Slot) (common.Root, error)
	StateForkAtSlot(slot math.Slot) (*ctypes.Fork, error)
}

type ValidatorBackend interface {
	ValidatorByID(
		slot math.Slot, id string,
	) (*types.ValidatorData, error)
	ValidatorsByIDs(
		slot math.Slot,
		ids []string,
		statuses []string,
	) ([]*types.ValidatorData, error)
	ValidatorBalancesByIDs(
		slot math.Slot,
		ids []string,
	) ([]*types.ValidatorBalanceData, error)
}

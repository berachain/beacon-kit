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
	"github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Backend is the interface for backend of the beacon API.
type Backend[BlockHeaderT, ForkT, ValidatorT any] interface {
	GenesisBackend
	BlockBackend[BlockHeaderT]
	RandaoBackend
	StateBackend[ForkT]
	ValidatorBackend[ValidatorT]
	HistoricalBackend[ForkT]
	GetSlotByRoot(root [32]byte) (math.Slot, error)
}

type GenesisBackend interface {
	GenesisValidatorsRoot(slot math.Slot) (common.Root, error)
}

type HistoricalBackend[ForkT any] interface {
	StateRootAtSlot(slot math.Slot) (common.Root, error)
	StateForkAtSlot(slot math.Slot) (ForkT, error)
}

type RandaoBackend interface {
	RandaoAtEpoch(slot math.Slot, epoch math.Epoch) (common.Bytes32, error)
}

type BlockBackend[BeaconBlockHeaderT any] interface {
	BlockRootAtSlot(slot math.Slot) (common.Root, error)
	BlockRewardsAtSlot(slot math.Slot) (*types.BlockRewardsData, error)
	BlockHeaderAtSlot(slot math.Slot) (BeaconBlockHeaderT, error)
}

type StateBackend[ForkT any] interface {
	StateRootAtSlot(slot math.Slot) (common.Root, error)
	StateForkAtSlot(slot math.Slot) (ForkT, error)
}

type ValidatorBackend[ValidatorT any] interface {
	ValidatorByID(
		slot math.Slot, id string,
	) (*types.ValidatorData[ValidatorT], error)
	ValidatorsByIDs(
		slot math.Slot,
		ids []string,
		statuses []string,
	) ([]*types.ValidatorData[ValidatorT], error)
	ValidatorBalancesByIDs(
		slot math.Slot,
		ids []string,
	) ([]*types.ValidatorBalanceData, error)
}

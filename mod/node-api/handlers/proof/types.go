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

package proof

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlockHeader is the interface for a beacon block header.
type BeaconBlockHeader[BeaconBlockHeaderT any] interface {
	New(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		stateRoot common.Root,
		bodyRoot common.Root,
	) BeaconBlockHeaderT
	HashTreeRoot() ([32]byte, error)
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetStateRoot() common.Root
	SetStateRoot(common.Root)
	GetBodyRoot() common.Root
}

// Backend is the interface for backend of the proof API.
type Backend[
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	ValidatorT any,
] interface {
	BlockBackend[BeaconBlockHeaderT]
	ValidatorBackend[ValidatorT]
}

type BlockBackend[BeaconBlockHeaderT any] interface {
	BlockHeader(slot uint64) (BeaconBlockHeaderT, error)
}

type ValidatorBackend[ValidatorT any] interface {
	AllValidators(slot uint64) ([]ValidatorT, error)
}

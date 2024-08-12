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

package components

import (
	"context"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ConsensusEngineInput is the input for the consensus engine.
type ConsensusEngineInput[
	AttestationDataT any,
	BeaconStateT any,
	MiddlewareT Middleware[AttestationDataT, SlashingInfoT, SlotDataT],
	SlashingInfoT any,
	SlotDataT any,
	StorageBackendT interface {
		StateFromContext(context.Context) BeaconStateT
	},
] struct {
	depinject.In
	ConsensusMiddleware MiddlewareT
	StorageBackend      StorageBackendT
}

// ProvideConsensusEngine is a depinject provider for the consensus engine.
func ProvideConsensusEngine[
	AttestationDataT AttestationData[AttestationDataT],
	BeaconStateT interface {
		// GetValidatorIndexByCometBFTAddress returns the validator index by the
		ValidatorIndexByCometBFTAddress(
			cometBFTAddress []byte,
		) (math.ValidatorIndex, error)
		// HashTreeRoot returns the hash tree root of the beacon state.
		HashTreeRoot() common.Root
	},
	MiddlewareT Middleware[AttestationDataT, SlashingInfoT, SlotDataT],
	SlashingInfoT SlashingInfo[SlashingInfoT],
	SlotDataT SlotData[AttestationDataT, SlashingInfoT, SlotDataT],
	StorageBackendT interface {
		StateFromContext(context.Context) BeaconStateT
	},
	ValidatorUpdateT any,
](
	in ConsensusEngineInput[
		AttestationDataT, BeaconStateT, MiddlewareT,
		SlashingInfoT, SlotDataT, StorageBackendT,
	],
) (*cometbft.ConsensusEngine[
	AttestationDataT, BeaconStateT, MiddlewareT,
	SlashingInfoT, SlotDataT, StorageBackendT, ValidatorUpdateT,
], error) {
	return cometbft.NewConsensusEngine[
		AttestationDataT,
		BeaconStateT,
		MiddlewareT,
		SlashingInfoT,
		SlotDataT,
		StorageBackendT,
		ValidatorUpdateT,
	](
		in.ConsensusMiddleware,
		in.StorageBackend,
	), nil
}

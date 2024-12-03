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

package params

import (
	math "github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	cmttypes "github.com/cometbft/cometbft/types"
)

type ChainSpec interface {
	// GetCometBFTConfigForSlot returns the CometBFT configuration for the given
	// slot.
	GetCometBFTConfigForSlot(math.Slot) any
}

// ConsensusParamsStore is a store for consensus parameters.
type ConsensusParamsStore struct {
	cs ChainSpec
}

// NewConsensusParamsStore creates a new ConsensusParamsStore.
func NewConsensusParamsStore(cs ChainSpec) *ConsensusParamsStore {
	return &ConsensusParamsStore{
		cs: cs,
	}
}

// Get retrieves the consensus parameters from the store.
// It returns the consensus parameters and an error, if any.
func (s *ConsensusParamsStore) Get() *cmtproto.ConsensusParams {
	//nolint:errcheck // TODO (fridrik): Is this safe?
	p := s.cs.
		GetCometBFTConfigForSlot(0).(*cmttypes.ConsensusParams).
		ToProto()
	return &p
}

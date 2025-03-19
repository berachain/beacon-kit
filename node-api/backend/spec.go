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
	"github.com/berachain/beacon-kit/node-api/handlers/config/types"
)

const (
	Phase0InactivityPenaltyQuotient = 67108864
	AltairInactivityPenaltyQuotient = 50331648
)

func (b *Backend) Spec() (*types.SpecData, error) {
	if b.cs == nil {
		return nil, errors.New("chain spec not found")
	}
	return &types.SpecData{
		DepositContractAddress: b.cs.DepositContractAddress(),
		// Network ID is same as eth1 chain ID.
		DepositNetworkID:        b.cs.DepositEth1ChainID(),
		DomainAggregateAndProof: b.cs.DomainTypeAggregateAndProof(),
		// These values are taken from ETH2.0 consensus spec. Currently these are placeholders, will be replaced
		// with the correct values for our versions like Deneb, Deneb1 etc once we implement slashing for inactivity.
		// https://github.com/ethereum/consensus-specs/blob/v1.3.0/specs/phase0/beacon-chain.md#rewards-and-penalties for phase 0.
		InactivityPenaltyQuotient: Phase0InactivityPenaltyQuotient,
		// https://github.com/ethereum/consensus-specs/blob/v1.3.0/specs/altair/beacon-chain.md#updated-penalty-values for altair.
		InactivityPenaltyQuotientAltair: AltairInactivityPenaltyQuotient,
	}, nil
}

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
	"github.com/berachain/beacon-kit/mod/errors"
	configtypes "github.com/berachain/beacon-kit/mod/node-api/handlers/config/types"
)

// GetSpec retrieves the spec from the store.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetSpec() (configtypes.SpecData, error) {
	chainSpec := b.ChainSpec()
	if chainSpec == nil {
		return configtypes.SpecData{}, errors.New("chain spec not found")
	}
	return configtypes.SpecData{
		DepositContractAddress: chainSpec.DepositContractAddress(),
		// TODO: Get the Network ID, introduce in chainSpec
		DepositNetworkID:        chainSpec.DepositEth1ChainID(),
		DomainAggregateAndProof: chainSpec.DomainTypeAggregateAndProof(),
		// TODO: put data for InactivityPenaltyQuotient in chainSpec
		InactivityPenaltyQuotient: chainSpec.InactivityPenaltyQuotient(),
		// TODO: introduce InactivityPenaltyQuotientAltair as
		// it does not exits in chainSpec
		InactivityPenaltyQuotientAltair: chainSpec.InactivityPenaltyQuotient(),
	}, nil
}

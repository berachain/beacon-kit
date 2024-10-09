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

package config

import (
	configtypes "github.com/berachain/beacon-kit/mod/node-api/handlers/config/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
)

// GetSpec returns the spec data.
func (h *Handler[ContextT]) GetSpec(_ ContextT) (any, error) {
	chainSpec, err := h.backend.GetSpec()
	if err != nil {
		return nil, err
	}

	return types.Wrap(
		configtypes.SpecData{
			DepositContractAddress: chainSpec.DepositContractAddress(),
			// TODO: Get the Network ID, introduce in chainSpec
			DepositNetworkID:                chainSpec.DepositEth1ChainID(),
			DomainAggregateAndProof:         chainSpec.DomainTypeAggregateAndProof(),
			InactivityPenaltyQuotient:       chainSpec.InactivityPenaltyQuotient(),
			InactivityPenaltyQuotientAltair: chainSpec.InactivityPenaltyQuotient(),
			// TODO: stubbed InactivityPenaltyQuotientAltair as
			// it doesn't exist in the spec
		}), nil
}

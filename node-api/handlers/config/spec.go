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

package config

import (
	"net/http"

	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/handlers/config/types"
	"github.com/berachain/beacon-kit/primitives/math"
)

// InactivityPenaltyQuotientPlaceholder is a placeholder value for the inactivity penalty quotient.
const InactivityPenaltyQuotientPlaceholder = "0"

// GetSpec returns the spec of the beacon chain.
func (h *Handler) GetSpec(handlers.Context) (any, error) {
	cs, err := h.backend.Spec()
	if err != nil {
		return nil, handlers.NewHTTPError(http.StatusInternalServerError, "failed to get spec: %v", err)
	}
	return types.SpecResponse{Data: types.SpecData{
		DepositContractAddress: cs.DepositContractAddress().String(),

		// Network ID is same as eth1 chain ID.
		DepositNetworkID: math.U64(cs.DepositEth1ChainID()).Base10(),

		DomainAggregateAndProof: cs.DomainTypeAggregateAndProof().String(),

		// Currently these are placeholders, will be replaced with the correct values for our
		// versions like Deneb, Deneb1 etc once we implement slashing for inactivity.
		InactivityPenaltyQuotient:       InactivityPenaltyQuotientPlaceholder,
		InactivityPenaltyQuotientAltair: InactivityPenaltyQuotientPlaceholder,
	}}, nil
}

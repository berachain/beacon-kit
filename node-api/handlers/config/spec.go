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

// These values are taken from ETH2.0 consensus spec.
const (
	//https://github.com/ethereum/consensus-specs/blob/v1.3.0/specs/phase0/beacon-chain.md#rewards-and-penalties
	Phase0InactivityPenaltyQuotient math.U64 = 67108864
	// https://github.com/ethereum/consensus-specs/blob/v1.3.0/specs/altair/beacon-chain.md#updated-penalty-values
	AltairInactivityPenaltyQuotient math.U64 = 50331648
)

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
		InactivityPenaltyQuotient:       Phase0InactivityPenaltyQuotient.Base10(),
		InactivityPenaltyQuotientAltair: AltairInactivityPenaltyQuotient.Base10(),
	}}, nil
}

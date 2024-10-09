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

package types

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

//nolint:lll // tags get long
type SpecData struct {
	DepositContractAddress          common.ExecutionAddress `json:"DEPOSIT_CONTRACT_ADDRESS"`
	DepositNetworkID                uint64                  `json:"DEPOSIT_NETWORK_ID"`
	DomainAggregateAndProof         common.DomainType       `json:"DOMAIN_AGGREGATE_AND_PROOF"`
	InactivityPenaltyQuotient       uint64                  `json:"INACTIVITY_PENALTY_QUOTIENT"`
	InactivityPenaltyQuotientAltair uint64                  `json:"INACTIVITY_PENALTY_QUOTIENT_ALTAIR"`
}

//nolint:lll // tags get long
type specJSON struct {
	DepositContractAddress          string `json:"DEPOSIT_CONTRACT_ADDRESS"`
	DepositNetworkID                string `json:"DEPOSIT_NETWORK_ID"`
	DomainAggregateAndProof         string `json:"DOMAIN_AGGREGATE_AND_PROOF"`
	InactivityPenaltyQuotient       string `json:"INACTIVITY_PENALTY_QUOTIENT"`
	InactivityPenaltyQuotientAltair string `json:"INACTIVITY_PENALTY_QUOTIENT_ALTAIR"`
}

func (sd SpecData) MarshalJSON() ([]byte, error) {
	return json.Marshal(specJSON{
		DepositContractAddress:          "0x" + hex.EncodeToString(sd.DepositContractAddress[:]),
		DepositNetworkID:                strconv.FormatUint(sd.DepositNetworkID, 10),
		DomainAggregateAndProof:         "0x" + hex.EncodeToString(sd.DomainAggregateAndProof[:]),
		InactivityPenaltyQuotient:       strconv.FormatUint(sd.InactivityPenaltyQuotient, 10),
		InactivityPenaltyQuotientAltair: strconv.FormatUint(sd.InactivityPenaltyQuotientAltair, 10),
	})
}

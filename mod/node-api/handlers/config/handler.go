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
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/server/context"
)

type Handler[ContextT context.Context] struct {
	*handlers.BaseHandler[ContextT]
}

func NewHandler[ContextT context.Context]() *Handler[ContextT] {
	h := &Handler[ContextT]{
		BaseHandler: handlers.NewBaseHandler(
			handlers.NewRouteSet[ContextT](""),
		),
	}
	return h
}

type Fork struct {
	PreviousVersion string `json:"previous_version"`
	CurrentVersion  string `json:"current_version"`
	Epoch           string `json:"epoch"`
}

type ForksResponse struct {
	Data []Fork `json:"data"`
}

type SpecData struct {
	SecondsPerSlot                  string `json:"SECONDS_PER_SLOT"`
	DepositContractAddress          string `json:"DEPOSIT_CONTRACT_ADDRESS"`
	DepositNetworkID                string `json:"DEPOSIT_NETWORK_ID"`
	DomainAggregateAndProof         string `json:"DOMAIN_AGGREGATE_AND_PROOF"`
	InactivityPenaltyQuotient       string `json:"INACTIVITY_PENALTY_QUOTIENT"`
	InactivityPenaltyQuotientAltair string `json:"INACTIVITY_PENALTY_QUOTIENT_ALTAIR"`
	BellatrixForkVersion            string `json:"BELLATRIX_FORK_VERSION"`
	DenebForkVersion                string `json:"DENEb_FORK_VERSION"`
	CapellaForkVersion              string `json:"CAPELLA_FORK_VERSION"`
}

type SpecResponse struct {
	Data SpecData `json:"data"`
}

type DepositContractData struct {
	ChainID string `json:"chain_id"`
	Address string `json:"address"`
}

type DepositContractResponse struct {
	Data DepositContractData `json:"data"`
}

func (h *Handler[Context]) GetForkSchedule(c Context) (any, error) {
	response := ForksResponse{
		Data: []Fork{
			{
				PreviousVersion: "0x00000000",
				CurrentVersion:  "0x00000000",
				Epoch:           "1",
			},
		},
	}
	return response, nil
}

func (h *Handler[Context]) GetSpec(c Context) (any, error) {
	response := SpecResponse{
		Data: SpecData{
			SecondsPerSlot:                  "2",
			DepositNetworkID:                "80087",
			DepositContractAddress:          "0x4242424242424242424242424242424242424242",
			InactivityPenaltyQuotient:       "1",
			InactivityPenaltyQuotientAltair: "1",
			DomainAggregateAndProof:         "0x06000000",
			BellatrixForkVersion:            "0x02000000",
			DenebForkVersion:                "0x02000000",
			CapellaForkVersion:              "0x02000000",
		},
	}
	return response, nil
}

func (h *Handler[Context]) GetDepositContract(c Context) (any, error) {
	response := DepositContractResponse{
		Data: DepositContractData{
			ChainID: "80087",
			Address: "0x4242424242424242424242424242424242424242",
		},
	}
	return response, nil
}

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

package cometbft

import (
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ConsensusEngine is used to decouple the Comet consensus engine from
// the Cosmos SDK.
// Right now, it is very coupled to the sdk base app and we will
// eventually fully decouple this.
type ConsensusEngine[
	AttestationDataT AttestationData[AttestationDataT],
	BeaconStateT BeaconState,
	MiddlewareT Middleware[AttestationDataT, SlashingInfoT, SlotDataT],
	SlashingInfoT SlashingInfo[SlashingInfoT],
	SlotDataT SlotData[AttestationDataT, SlashingInfoT, SlotDataT],
	StorageBackendT StorageBackend[BeaconStateT],
	ValidatorUpdateT any,
] struct {
	m  MiddlewareT
	sb StorageBackendT
}

// NewConsensusEngine returns a new consensus middleware.
func NewConsensusEngine[
	AttestationDataT AttestationData[AttestationDataT],
	BeaconStateT BeaconState,
	MiddlewareT Middleware[AttestationDataT, SlashingInfoT, SlotDataT],
	SlashingInfoT SlashingInfo[SlashingInfoT],
	SlotDataT SlotData[AttestationDataT, SlashingInfoT, SlotDataT],
	StorageBackendT StorageBackend[BeaconStateT],
	ValidatorUpdateT any,
](
	m MiddlewareT,
	sb StorageBackendT,
) *ConsensusEngine[
	AttestationDataT, BeaconStateT, MiddlewareT,
	SlashingInfoT, SlotDataT, StorageBackendT, ValidatorUpdateT,
] {
	return &ConsensusEngine[
		AttestationDataT,
		BeaconStateT,
		MiddlewareT,
		SlashingInfoT,
		SlotDataT,
		StorageBackendT,
		ValidatorUpdateT,
	]{
		m:  m,
		sb: sb,
	}
}

// TODO: Decouple Comet Types
func (c *ConsensusEngine[_, _, _, _, _, _, _]) PrepareProposal(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	slotData, err := c.convertPrepareProposalToSlotData(
		ctx,
		req,
	)
	if err != nil {
		return nil, err
	}
	blkBz, sidecarsBz, err := c.m.PrepareProposal(
		ctx,
		slotData,
	)
	if err != nil {
		return nil, err
	}
	return &cmtabci.PrepareProposalResponse{
		Txs: [][]byte{blkBz, sidecarsBz},
	}, nil
}

// TODO: Decouple Comet Types
func (c *ConsensusEngine[_, _, _, _, _, _, ValidatorUpdateT]) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	resp, err := c.m.ProcessProposal(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*cmtabci.ProcessProposalResponse), nil
}

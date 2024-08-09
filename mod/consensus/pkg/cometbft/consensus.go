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
	"context"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sourcegraph/conc/iter"
)

// ConsensusEngine is used to decouple the Comet consensus engine from
// the Cosmos SDK.
// Right now, it is very coupled to the sdk base app and we will
// eventually fully decouple this.
type ConsensusEngine[
	AttestationDataT AttestationData[AttestationDataT],
	BeaconStateT BeaconState,
	SlashingInfoT SlashingInfo[SlashingInfoT],
	SlotDataT SlotData[AttestationDataT, SlashingInfoT, SlotDataT],
	StorageBackendT StorageBackend[BeaconStateT],
	ValidatorUpdateT any,
] struct {
	Middleware[AttestationDataT, SlashingInfoT, SlotDataT]
	sb StorageBackendT
}

// NewConsensusEngine returns a new consensus middleware.
func NewConsensusEngine[
	AttestationDataT AttestationData[AttestationDataT],
	BeaconStateT BeaconState,
	SlashingInfoT SlashingInfo[SlashingInfoT],
	SlotDataT SlotData[AttestationDataT, SlashingInfoT, SlotDataT],
	StorageBackendT StorageBackend[BeaconStateT],
	ValidatorUpdateT any,
](
	m Middleware[AttestationDataT, SlashingInfoT, SlotDataT],
	sb StorageBackendT,
) *ConsensusEngine[
	AttestationDataT,
	BeaconStateT,
	SlashingInfoT,
	SlotDataT,
	StorageBackendT,
	ValidatorUpdateT,
] {
	return &ConsensusEngine[
		AttestationDataT,
		BeaconStateT,
		SlashingInfoT,
		SlotDataT,
		StorageBackendT,
		ValidatorUpdateT,
	]{
		Middleware: m,
		sb:         sb,
	}
}

func (c *ConsensusEngine[_, _, _, _, _, ValidatorUpdateT]) InitGenesis(
	ctx context.Context,
	genesisBz []byte,
) ([]ValidatorUpdateT, error) {
	updates, err := c.Middleware.InitGenesis(ctx, genesisBz)
	if err != nil {
		return nil, err
	}
	// Convert updates into the Cosmos SDK format.
	return iter.MapErr(updates, convertValidatorUpdate[ValidatorUpdateT])
}

// TODO: Decouple Comet Types
func (c *ConsensusEngine[_, _, _, _, _, _]) PrepareProposal(
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
	blkBz, sidecarsBz, err := c.Middleware.PrepareProposal(
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
func (c *ConsensusEngine[_, _, _, _, _, ValidatorUpdateT]) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	resp, err := c.Middleware.ProcessProposal(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*cmtabci.ProcessProposalResponse), nil
}

func (c *ConsensusEngine[_, _, _, _, _, ValidatorUpdateT]) FinalizeBlock(
	ctx context.Context, req *cmtabci.FinalizeBlockRequest,
) ([]ValidatorUpdateT, error) {
	updates, err := c.Middleware.FinalizeBlock(ctx, req)
	if err != nil {
		return nil, err
	}
	return iter.MapErr(updates, convertValidatorUpdate[ValidatorUpdateT])
}

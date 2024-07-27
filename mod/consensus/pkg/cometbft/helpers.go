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
	"sort"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	v1 "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// convertValidatorUpdate abstracts the conversion of a
// transition.ValidatorUpdate to an appmodulev2.ValidatorUpdate.
// TODO: this is so hood, bktypes -> sdktypes -> generic is crazy
// maybe make this some kind of codec/func that can be passed in?
func convertValidatorUpdate[ValidatorUpdateT any](
	u **transition.ValidatorUpdate,
) (ValidatorUpdateT, error) {
	var valUpdate ValidatorUpdateT
	update := *u
	if update == nil {
		return valUpdate, ErrUndefinedValidatorUpdate
	}
	return any(appmodulev2.ValidatorUpdate{
		PubKey:     update.Pubkey[:],
		PubKeyType: crypto.CometBLSType,
		//#nosec:G701 // this is safe.
		Power: int64(update.EffectiveBalance.Unwrap()),
	}).(ValidatorUpdateT), nil
}

// convertPrepareProposalToSlotData converts a prepare proposal request to
// a slot data.
func (c *ConsensusEngine[
	_, _, _, SlotDataT, _, _,
]) convertPrepareProposalToSlotData(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (SlotDataT, error) {
	var t SlotDataT

	// Get the attestation data from the votes.
	attestationData, err := c.attestationsFromVotes(
		ctx,
		req.LocalLastCommit.Votes,
		//#nosec:G701 // safe.
		math.Slot(req.Height),
	)
	if err != nil {
		return t, err
	}

	// Get the slashing info from the misbehaviors.
	slashingInfo, err := c.slashingInfoFromMisbehaviors(
		ctx,
		req.Misbehavior,
	)
	if err != nil {
		return t, err
	}

	// Create the slot data.
	t = t.New(
		math.U64(req.Height),
		attestationData,
		slashingInfo,
	)
	return t, nil
}

// attestationsFromVotes returns a list of attestation data from the votes.
func (c *ConsensusEngine[
	AttestationDataT, _, _, _, _, _,
]) attestationsFromVotes(
	ctx sdk.Context,
	votes []v1.ExtendedVoteInfo,
	slot math.Slot,
) ([]AttestationDataT, error) {
	var err error
	var index math.U64
	attestations := make([]AttestationDataT, len(votes))
	st := c.sb.StateFromContext(ctx)
	root, err := st.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	for i, vote := range votes {
		index, err = st.ValidatorIndexByCometBFTAddress(vote.Validator.Address)
		if err != nil {
			return nil, err
		}

		var t AttestationDataT
		t = t.New(
			slot,
			index,
			root,
		)
		attestations[i] = t
	}

	// Attestations are sorted by index.
	sort.Slice(attestations, func(i, j int) bool {
		return attestations[i].GetIndex() < attestations[j].GetIndex()
	})
	return attestations, nil
}

// slashingInfoFromMisbehaviors returns a list of slashing info from the
// comet misbehaviors.
func (c *ConsensusEngine[
	_, _, SlashingInfoT, _, _, _,
]) slashingInfoFromMisbehaviors(
	ctx sdk.Context,
	misbehaviors []v1.Misbehavior,
) ([]SlashingInfoT, error) {
	var err error
	var index math.U64
	st := c.sb.StateFromContext(ctx)
	slashingInfo := make([]SlashingInfoT, len(misbehaviors))
	for i, misbehavior := range misbehaviors {
		index, err = st.ValidatorIndexByCometBFTAddress(
			misbehavior.Validator.Address,
		)
		if err != nil {
			return nil, err
		}
		var t SlashingInfoT
		t = t.New(
			//#nosec:G701 // safe.
			math.Slot(misbehavior.GetHeight()),
			index,
		)
		slashingInfo[i] = t
	}
	return slashingInfo, nil
}

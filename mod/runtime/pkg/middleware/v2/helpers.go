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

package middleware

import (
	"context"
	"sort"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	abci "github.com/cometbft/cometbft/abci/types"
)

// convertPrepareProposalToSlotData converts a prepare proposal request to
// a slot data.
func (c *ABCIMiddleware[
	_, _, _, _, _, _, _, _, _, SlotDataT, _,
]) convertPrepareProposalToSlotData(
	ctx context.Context,
	req *abci.PrepareProposalRequest,
) (SlotDataT, error) {
	var t SlotDataT

	// Get the attestation data from the votes.
	attestationData, err := c.attestationsFromVotes(
		ctx,
		req.LocalLastCommit.Votes,
		//#nosec:G701 // safe.
		math.U64(req.Height),
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
func (c *ABCIMiddleware[
	AttestationDataT, _, _, _, _, _, _, _, _, _, _,
]) attestationsFromVotes(
	ctx context.Context,
	votes []abci.ExtendedVoteInfo,
	slot math.U64,
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
func (c *ABCIMiddleware[
	_, _, _, _, _, _, _, _, SlashingInfoT, _, _,
]) slashingInfoFromMisbehaviors(
	ctx context.Context,
	misbehaviors []abci.Misbehavior,
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
			math.U64(misbehavior.GetHeight()),
			index,
		)
		slashingInfo[i] = t
	}
	return slashingInfo, nil
}

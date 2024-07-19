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
	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	cmtabci "github.com/cometbft/cometbft/abci/types"
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

func (c *ConsensusEngine[
	AttestationDataT,
	SlashingInfoT,
	SlotDataT,
	_,
]) convertPrepareProposalToSlotData(
	req *cmtabci.PrepareProposalRequest,
) (SlotDataT, error) {
	var t SlotDataT
	t = t.New(
		math.U64(req.Height),
		nil,
		nil,
	)
	return t, nil
}

// // attestationDataFromVotes returns a list of attestation data from the
// // comet vote info. This is used to build the attestations for the block.
// func (h *ValidatorMiddleware[
// 	AvailabilityStoreT,
// 	BeaconBlockT,
// 	BeaconBlockBodyT,
// 	BeaconStateT,
// 	BlobsSidecarsT,
// 	DepositStoreT,
// ]) attestationDataFromVotes(
// 	st BeaconStateT,
// 	root primitives.Root,
// 	votes []v1.ExtendedVoteInfo,
// 	slot uint64,
// ) ([]*types.AttestationData, error) {
// 	var err error
// 	var index math.U64
// 	attestations := make([]*types.AttestationData, len(votes))
// 	for i, vote := range votes {
// 		index, err = st.ValidatorIndexByCometBFTAddress(vote.Validator.Address)
// 		if err != nil {
// 			return nil, err
// 		}
// 		attestations[i] = &types.AttestationData{
// 			Slot:            slot,
// 			Index:           index.Unwrap(),
// 			BeaconBlockRoot: root,
// 		}
// 	}
// 	// Attestations are sorted by index.
// 	sort.Slice(attestations, func(i, j int) bool {
// 		return attestations[i].Index < attestations[j].Index
// 	})
// 	return attestations, nil
// }

// // slashingInfoFromMisbehaviors returns a list of slashing info from the
// // comet misbehaviors.
// func slashingInfoFromMisbehaviors(
// 	st BeaconStateT,
// 	misbehaviors []v1.Misbehavior,
// ) ([]*types.SlashingInfo, error) {
// 	var err error
// 	var index math.U64
// 	slashingInfo := make([]*types.SlashingInfo, len(misbehaviors))
// 	for i, misbehavior := range misbehaviors {
// 		index, err = st.ValidatorIndexByCometBFTAddress(
// 			misbehavior.Validator.Address,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		slashingInfo[i] = &types.SlashingInfo{
// 			//#nosec:G701 // safe.
// 			Slot:  uint64(misbehavior.GetHeight()),
// 			Index: index.Unwrap(),
// 		}
// 	}
// 	return slashingInfo, nil
// }

// ProcessProposalHandler

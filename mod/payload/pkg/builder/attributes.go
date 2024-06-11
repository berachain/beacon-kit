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

package builder

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// getPayloadAttributes returns the payload attributes for the given state and
// slot. The attribute is required to initiate a payload build process in the
// context of an `engine_forkchoiceUpdated` call.
func (pb *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
]) getPayloadAttribute(
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	prevHeadRoot [32]byte,
) (engineprimitives.PayloadAttributer, error) {
	var (
		prevRandao [32]byte
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		pb.logger.Error(
			"Could not get expected withdrawals to get payload attribute",
			"error",
			err,
		)
		return nil, err
	}

	epoch := pb.chainSpec.SlotToEpoch(slot)

	// Get the previous randao mix.
	prevRandao, err = st.GetRandaoMixAtIndex(
		uint64(epoch) % pb.chainSpec.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return nil, err
	}

	return engineprimitives.NewPayloadAttributes(
		pb.chainSpec.ActiveForkVersionForEpoch(epoch),
		timestamp,
		prevRandao,
		pb.cfg.SuggestedFeeRecipient,
		withdrawals,
		prevHeadRoot,
	)
}

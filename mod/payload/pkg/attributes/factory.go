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

package attributes

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Factory is a factory for creating payload attributes.
type Factory[
	BeaconStateT BeaconState[WithdrawalT], WithdrawalT any,
] struct {
	// chainSpec is the chain spec for the attributes factory.
	chainSpec common.ChainSpec
	// logger is the logger for the attributes factory.
	logger log.Logger[any]
	// suggestedFeeRecipient is the suggested fee recipient sent to
	// the execution client for the payload build.
	suggestedFeeRecipient common.ExecutionAddress
}

// NewAttributesFactory creates a new instance of AttributesFactory.
func NewAttributesFactory[
	BeaconStateT BeaconState[WithdrawalT], WithdrawalT any,
](
	chainSpec common.ChainSpec,
	logger log.Logger[any],
	suggestedFeeRecipient common.ExecutionAddress,
) *Factory[BeaconStateT, WithdrawalT] {
	return &Factory[BeaconStateT, WithdrawalT]{
		chainSpec:             chainSpec,
		logger:                logger,
		suggestedFeeRecipient: suggestedFeeRecipient,
	}
}

// CreateAttributes creates a new instance of PayloadAttributes.
func (f *Factory[BeaconStateT, WithdrawalT]) BuildPayloadAttributes(
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	prevHeadRoot [32]byte,
) (engineprimitives.PayloadAttributer, error) {
	var (
		prevRandao [32]byte
		epoch      = f.chainSpec.SlotToEpoch(slot)
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		f.logger.Error(
			"Could not get expected withdrawals to get payload attribute",
			"error",
			err,
		)
		return nil, err
	}

	// Get the previous randao mix.
	if prevRandao, err = st.GetRandaoMixAtIndex(
		uint64(epoch) % f.chainSpec.EpochsPerHistoricalVector(),
	); err != nil {
		return nil, err
	}

	return engineprimitives.NewPayloadAttributes(
		f.chainSpec.ActiveForkVersionForEpoch(epoch),
		timestamp,
		prevRandao,
		f.suggestedFeeRecipient,
		withdrawals,
		prevHeadRoot,
	)
}

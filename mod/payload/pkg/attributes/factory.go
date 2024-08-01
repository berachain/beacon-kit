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
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Factory is a factory for creating payload attributes.
type Factory[
	BeaconStateT BeaconState[WithdrawalT],
	StateProcessorT StateProcessor[BeaconStateT, WithdrawalT],
	PayloadAttributesT PayloadAttributes[PayloadAttributesT, WithdrawalT],
	WithdrawalT any,
] struct {
	// sp is the state processor for the attributes factory.
	sp StateProcessorT
	// chainSpec is the chain spec for the attributes factory.
	chainSpec common.ChainSpec
	// logger is the logger for the attributes factory.
	logger log.Logger[any]
	// suggestedFeeRecipient is the suggested fee recipient sent to
	// the execution client for the payload build.
	suggestedFeeRecipient gethprimitives.ExecutionAddress
}

// NewAttributesFactory creates a new instance of AttributesFactory.
func NewAttributesFactory[
	BeaconStateT BeaconState[WithdrawalT],
	StateProcessorT StateProcessor[BeaconStateT, WithdrawalT],
	PayloadAttributesT PayloadAttributes[PayloadAttributesT, WithdrawalT],
	WithdrawalT any,
](
	sp StateProcessorT,
	chainSpec common.ChainSpec,
	logger log.Logger[any],
	suggestedFeeRecipient gethprimitives.ExecutionAddress,
) *Factory[BeaconStateT, StateProcessorT, PayloadAttributesT, WithdrawalT] {
	return &Factory[BeaconStateT, StateProcessorT, PayloadAttributesT, WithdrawalT]{
		sp:                    sp,
		chainSpec:             chainSpec,
		logger:                logger,
		suggestedFeeRecipient: suggestedFeeRecipient,
	}
}

// CreateAttributes creates a new instance of PayloadAttributes.
func (f *Factory[
	BeaconStateT,
	StateProcessorT,
	PayloadAttributesT,
	WithdrawalT,
]) BuildPayloadAttributes(
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	prevHeadRoot [32]byte,
) (PayloadAttributesT, error) {
	var (
		prevRandao [32]byte
		attributes PayloadAttributesT
		epoch      = f.chainSpec.SlotToEpoch(slot)
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := f.sp.ExpectedWithdrawals(st)
	if err != nil {
		f.logger.Error(
			"Could not get expected withdrawals to get payload attribute",
			"error",
			err,
		)
		return attributes, err
	}

	// Get the previous randao mix.
	if prevRandao, err = st.GetRandaoMixAtIndex(
		uint64(epoch) % f.chainSpec.EpochsPerHistoricalVector(),
	); err != nil {
		return attributes, err
	}

	return attributes.New(
		f.chainSpec.ActiveForkVersionForEpoch(epoch),
		timestamp,
		prevRandao,
		f.suggestedFeeRecipient,
		withdrawals,
		prevHeadRoot,
	)
}

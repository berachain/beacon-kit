// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Factory is a factory for creating payload attributes.
type Factory struct {
	// chainSpec is the chain spec for the attributes factory.
	chainSpec ChainSpec
	// logger is the logger for the attributes factory.
	logger log.Logger
	// suggestedFeeRecipient is the suggested fee recipient sent to
	// the execution client for the payload build.
	suggestedFeeRecipient common.ExecutionAddress
}

// NewAttributesFactory creates a new instance of AttributesFactory.
func NewAttributesFactory(
	chainSpec ChainSpec,
	logger log.Logger,
	suggestedFeeRecipient common.ExecutionAddress,
) *Factory {
	return &Factory{
		chainSpec:             chainSpec,
		logger:                logger,
		suggestedFeeRecipient: suggestedFeeRecipient,
	}
}

// BuildPayloadAttributes creates a new instance of PayloadAttributes.
func (f *Factory) BuildPayloadAttributes(
	st ReadOnlyBeaconState,
	timestamp math.U64,
) (*engineprimitives.PayloadAttributes, error) {
	// Get the expected withdrawals to include in this payload.
	withdrawals, _, err := st.ExpectedWithdrawals(timestamp)
	if err != nil {
		f.logger.Error(
			"Could not get expected withdrawals to get payload attribute",
			"error",
			err,
		)
		return nil, err
	}

	// Get the previous randao mix.
	slot, err := st.GetSlot()
	if err != nil {
		f.logger.Error(
			"Could not get slot to get payload attribute",
			"error", err,
		)
		return nil, err
	}
	epoch := f.chainSpec.SlotToEpoch(slot).Unwrap()
	prevRandao, err := st.GetRandaoMixAtIndex(epoch % f.chainSpec.EpochsPerHistoricalVector())
	if err != nil {
		f.logger.Error(
			"Could not get randao mix to get payload attribute",
			"epoch", epoch,
			"error", err,
		)
		return nil, err
	}

	latestHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		f.logger.Error(
			"Could not get latest block header to get payload attribute",
			"error", err,
		)
		return nil, err
	}
	prevHeadRoot := latestHeader.HashTreeRoot()

	return engineprimitives.NewPayloadAttributes(
		f.chainSpec.ActiveForkVersionForTimestamp(timestamp),
		timestamp,
		prevRandao,
		f.suggestedFeeRecipient,
		withdrawals,
		prevHeadRoot,
	)
}

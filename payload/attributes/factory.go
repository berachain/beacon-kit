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
	"github.com/berachain/beacon-kit/chain-spec/chain"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// Factory is a factory for creating payload attributes.
type Factory struct {
	// chainSpec is the chain spec for the attributes factory.
	chainSpec chain.ChainSpec
	// logger is the logger for the attributes factory.
	logger log.Logger
	// suggestedFeeRecipient is the suggested fee recipient sent to
	// the execution client for the payload build.
	suggestedFeeRecipient common.ExecutionAddress
}

// NewAttributesFactory creates a new instance of AttributesFactory.
func NewAttributesFactory(
	chainSpec chain.ChainSpec,
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
	st *statedb.StateDB,
	slot math.Slot,
	timestamp uint64,
	prevHeadRoot [32]byte,
) (*engineprimitives.PayloadAttributes, error) {
	var (
		prevRandao [32]byte
		attributes *engineprimitives.PayloadAttributes
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
		return attributes, err
	}

	// Get the previous randao mix.
	if prevRandao, err = st.GetRandaoMixAtIndex(
		epoch.Unwrap() % f.chainSpec.EpochsPerHistoricalVector(),
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

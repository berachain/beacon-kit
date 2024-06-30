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

package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/require"
)

type Withdrawal struct{}

func TestPayloadAttributes(t *testing.T) {
	forkVersion := uint32(1)
	timestamp := uint64(123456789)
	prevRandao := common.Bytes32{1, 2, 3}
	suggestedFeeRecipient := gethprimitives.ExecutionAddress{}
	withdrawals := []Withdrawal{}
	parentBeaconBlockRoot := common.Root{}

	payloadAttributes, err := engineprimitives.NewPayloadAttributes[Withdrawal](
		forkVersion,
		timestamp,
		prevRandao,
		suggestedFeeRecipient,
		withdrawals,
		parentBeaconBlockRoot,
	)
	require.NoError(t, err)
	require.NotNil(t, payloadAttributes)

	require.False(t, payloadAttributes.IsNil())

	require.Equal(
		t,
		suggestedFeeRecipient,
		payloadAttributes.GetSuggestedFeeRecipient(),
	)
	require.Equal(t, forkVersion, payloadAttributes.Version())

	require.NoError(t, payloadAttributes.Validate())
}

func TestNewPayloadAttributes_ErrorCases(t *testing.T) {
	forkVersion := uint32(1)
	prevRandao := common.Bytes32{1, 2, 3}
	suggestedFeeRecipient := gethprimitives.ExecutionAddress{}
	withdrawals := []Withdrawal{}
	parentBeaconBlockRoot := common.Root{}

	// Test case where Timestamp is 0
	_, err := engineprimitives.NewPayloadAttributes[Withdrawal](
		forkVersion,
		0,
		prevRandao,
		suggestedFeeRecipient,
		withdrawals,
		parentBeaconBlockRoot,
	)
	require.ErrorIs(t, err, engineprimitives.ErrInvalidTimestamp)

	// Test case where PrevRandao is an empty array
	_, err = engineprimitives.NewPayloadAttributes[Withdrawal](
		forkVersion,
		123456789,
		common.Bytes32{},
		suggestedFeeRecipient,
		withdrawals,
		parentBeaconBlockRoot,
	)
	require.ErrorIs(t, err, engineprimitives.ErrEmptyPrevRandao)

	// Test case where Withdrawals is nil and version is equal to Capella
	_, err = engineprimitives.NewPayloadAttributes[Withdrawal](
		version.Capella,
		123456789,
		prevRandao,
		suggestedFeeRecipient,
		nil,
		parentBeaconBlockRoot,
	)
	require.ErrorIs(t, err, engineprimitives.ErrNilWithdrawals)
}

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

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

type testWithdrawal struct{}

type payloadAttributesInput struct {
	forkVersion           uint32
	timestamp             uint64
	prevRandao            common.Bytes32
	suggestedFeeRecipient common.ExecutionAddress
	withdrawals           []testWithdrawal
	parentBeaconBlockRoot common.Root
}

func TestPayloadAttributes(t *testing.T) {
	// default valid data
	validInput := payloadAttributesInput{
		forkVersion:           uint32(1),
		timestamp:             uint64(123456789),
		prevRandao:            common.Bytes32{1, 2, 3},
		suggestedFeeRecipient: common.ExecutionAddress{},
		withdrawals:           []testWithdrawal{},
		parentBeaconBlockRoot: common.Root{},
	}
	tests := []struct {
		name    string
		input   func() payloadAttributesInput
		wantErr error
	}{
		{
			name: "Valid payload attributes",
			input: func() payloadAttributesInput {
				return validInput
			},
			wantErr: nil,
		},
		{
			name: "Invalid timestamp",
			input: func() payloadAttributesInput {
				res := validInput
				res.timestamp = 0
				return res
			},
			wantErr: engineprimitives.ErrInvalidTimestamp,
		},
		{
			name: "Invalid PreRandao",
			input: func() payloadAttributesInput {
				res := validInput
				res.prevRandao = common.Bytes32{}
				return res
			},
			wantErr: engineprimitives.ErrEmptyPrevRandao,
		},
		{
			name: "Invalid nil withdrawals on Capella",
			input: func() payloadAttributesInput {
				res := validInput
				res.forkVersion = version.Capella
				res.withdrawals = nil
				return res
			},
			wantErr: engineprimitives.ErrNilWithdrawals,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.input()
			p := &engineprimitives.PayloadAttributes[testWithdrawal]{}
			got, err := p.New(
				in.forkVersion,
				in.timestamp,
				in.prevRandao,
				in.suggestedFeeRecipient,
				in.withdrawals,
				in.parentBeaconBlockRoot,
			)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)

				require.False(t, got.IsNil())
				require.NoError(t, got.Validate())

				require.Equal(
					t,
					in.suggestedFeeRecipient,
					got.GetSuggestedFeeRecipient(),
				)
				require.Equal(t, in.forkVersion, got.Version())
			}
		})
	}
}

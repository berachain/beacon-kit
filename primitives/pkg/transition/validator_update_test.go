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

package transition_test

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/primitives/pkg/math"
	"github.com/berachain/beacon-kit/primitives/pkg/transition"
	"github.com/stretchr/testify/require"
)

func TestValidatorUpdate_CanonicalSort(t *testing.T) {
	pubkey1 := crypto.BLSPubkey{1}
	pubkey2 := crypto.BLSPubkey{2}
	pubkey3 := crypto.BLSPubkey{3}

	type test struct {
		name  string
		input transition.ValidatorUpdates
		want  transition.ValidatorUpdates
	}

	tests := []test{
		{
			name: "RemoveDuplicates-PickLatest",
			input: transition.ValidatorUpdates{
				&transition.ValidatorUpdate{
					Pubkey:           pubkey1,
					EffectiveBalance: math.Gwei(1000),
				},
				&transition.ValidatorUpdate{
					Pubkey:           pubkey1,
					EffectiveBalance: math.Gwei(500),
				},
				&transition.ValidatorUpdate{
					Pubkey:           pubkey2,
					EffectiveBalance: math.Gwei(2000),
				},
			},
			want: transition.ValidatorUpdates{
				&transition.ValidatorUpdate{
					Pubkey:           pubkey1,
					EffectiveBalance: math.Gwei(500),
				},
				&transition.ValidatorUpdate{
					Pubkey:           pubkey2,
					EffectiveBalance: math.Gwei(2000),
				},
			},
		},
		{
			name: "SortByPubKey",
			input: transition.ValidatorUpdates{
				&transition.ValidatorUpdate{
					Pubkey:           pubkey3,
					EffectiveBalance: math.Gwei(2000),
				},
				&transition.ValidatorUpdate{
					Pubkey:           pubkey1,
					EffectiveBalance: math.Gwei(5000),
				},
				&transition.ValidatorUpdate{
					Pubkey:           pubkey2,
					EffectiveBalance: math.Gwei(1000),
				},
			},
			want: transition.ValidatorUpdates{
				&transition.ValidatorUpdate{
					Pubkey:           pubkey1,
					EffectiveBalance: math.Gwei(5000),
				},
				&transition.ValidatorUpdate{
					Pubkey:           pubkey2,
					EffectiveBalance: math.Gwei(1000),
				},
				&transition.ValidatorUpdate{
					Pubkey:           pubkey3,
					EffectiveBalance: math.Gwei(2000),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.CanonicalSort()
			require.Equal(t, tc.want, got)
		})
	}
}

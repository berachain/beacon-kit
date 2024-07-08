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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalSSZ(t *testing.T) {
	withdrawal := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(2),
		Address:   [20]byte{},
		Amount:    math.Gwei(100),
	}

	data, err := withdrawal.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	err = withdrawal.UnmarshalSSZ(data)
	require.NoError(t, err)

	size := withdrawal.SizeSSZ()
	require.Equal(t, 44, size)

	tree, errHashTree := withdrawal.HashTreeRoot()
	require.NoError(t, errHashTree)
	require.NotNil(t, tree)
}

func TestWithdrawalGetTree(t *testing.T) {
	withdrawal := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(2),
		Address:   [20]byte{},
		Amount:    math.Gwei(100),
	}

	tree, err := withdrawal.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestWithdrawalUnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *engineprimitives.Withdrawal
		wantErr bool
	}{
		{
			name: "Valid input little endian",
			input: []byte{
				1, 0, 0, 0, 0, 0, 0, 0, // Index: 1
				2, 0, 0, 0, 0, 0, 0, 0, // Validator: 2
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
				12, 13, 14, 15, 16, 17, 18, 19, 20, // Address
				100, 0, 0, 0, 0, 0, 0, 0, // Amount: 100
			},
			want: &engineprimitives.Withdrawal{
				Index:     math.U64(1),
				Validator: math.ValidatorIndex(2),
				Address: [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				Amount: math.Gwei(100),
			},
			wantErr: false,
		},
		{
			name:    "Invalid size - too short",
			input:   []byte{0, 1, 2, 3},
			want:    &engineprimitives.Withdrawal{},
			wantErr: true,
		},
		{
			name:    "Invalid size - too long",
			input:   make([]byte, 45),
			want:    &engineprimitives.Withdrawal{},
			wantErr: true,
		},
		{
			name: "Max values",
			input: []byte{
				255, 255, 255, 255, 255, 255, 255, 255, // Index: max uint64
				255, 255, 255, 255, 255, 255, 255, 255, // Validator: max uint64
				255, 255, 255, 255, 255, 255, 255, 255, 255,
				255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, // Address
				255, 255, 255, 255, 255, 255, 255, 255, // Amount: max uint64
			},
			want: &engineprimitives.Withdrawal{
				Index:     math.U64(^uint64(0)),
				Validator: math.ValidatorIndex(^uint64(0)),
				Address: [20]byte{255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
				Amount: math.Gwei(^uint64(0)),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &engineprimitives.Withdrawal{}
			err := w.UnmarshalSSZ(tt.input)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.Index, w.Index)
				require.Equal(t, tt.want.Validator, w.Validator)
				require.Equal(t, tt.want.Address, w.Address)
				require.Equal(t, tt.want.Amount, w.Amount)
			}
		})
	}
}

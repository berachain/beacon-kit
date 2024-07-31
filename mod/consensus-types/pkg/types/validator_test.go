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

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

func TestNewValidatorFromDeposit(t *testing.T) {
	tests := []struct {
		name                      string
		pubkey                    crypto.BLSPubkey
		withdrawalCredentials     types.WithdrawalCredentials
		amount                    math.Gwei
		effectiveBalanceIncrement math.Gwei
		maxEffectiveBalance       math.Gwei
		want                      *types.Validator
	}{
		{
			name:   "normal case",
			pubkey: [48]byte{0x01},
			withdrawalCredentials: types.
				NewCredentialsFromExecutionAddress(
					gethprimitives.ExecutionAddress{0x01},
				),
			amount:                    32e9,
			effectiveBalanceIncrement: 1e9,
			maxEffectiveBalance:       32e9,
			want: &types.Validator{
				Pubkey: [48]byte{0x01},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
		{
			name:   "effective balance capped at max",
			pubkey: [48]byte{0x02},
			withdrawalCredentials: types.
				NewCredentialsFromExecutionAddress(
					gethprimitives.ExecutionAddress{0x02},
				),
			amount:                    40e9,
			effectiveBalanceIncrement: 1e9,
			maxEffectiveBalance:       32e9,
			want: &types.Validator{
				Pubkey: [48]byte{0x02},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x02},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
		{
			name:   "effective balance rounded down",
			pubkey: [48]byte{0x03},
			withdrawalCredentials: types.
				NewCredentialsFromExecutionAddress(
					gethprimitives.ExecutionAddress{0x03},
				),
			amount:                    32.5e9,
			effectiveBalanceIncrement: 1e9,
			maxEffectiveBalance:       32e9,
			want: &types.Validator{
				Pubkey: [48]byte{0x03},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x03},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.NewValidatorFromDeposit(
				tt.pubkey,
				tt.withdrawalCredentials,
				tt.amount,
				tt.effectiveBalanceIncrement,
				tt.maxEffectiveBalance,
			)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidator_IsActive(t *testing.T) {
	tests := []struct {
		name      string
		epoch     math.Epoch
		validator *types.Validator
		want      bool
	}{
		{
			name:  "active",
			epoch: 10,
			validator: &types.Validator{
				ActivationEpoch: 5,
				ExitEpoch:       15,
			},
			want: true,
		},
		{
			name:  "not active, before activation",
			epoch: 4,
			validator: &types.Validator{
				ActivationEpoch: 5,
				ExitEpoch:       15,
			},
			want: false,
		},
		{
			name:  "not active, after exit",
			epoch: 16,
			validator: &types.Validator{
				ActivationEpoch: 5,
				ExitEpoch:       15,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.validator.IsActive(tt.epoch))
		})
	}
}

func TestValidator_IsEligibleForActivation(t *testing.T) {
	tests := []struct {
		name           string
		finalizedEpoch math.Epoch
		validator      *types.Validator
		want           bool
	}{
		{
			name:           "eligible",
			finalizedEpoch: 10,
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
			want: true,
		},
		{
			name:           "not eligible, activation eligibility in future",
			finalizedEpoch: 4,
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
			want: false,
		},
		{
			name:           "not eligible, already activated",
			finalizedEpoch: 10,
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				ActivationEpoch:            8,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsEligibleForActivation(tt.finalizedEpoch),
			)
		})
	}
}

func TestValidator_IsEligibleForActivationQueue(t *testing.T) {
	maxEffectiveBalance := math.Gwei(32e9)
	tests := []struct {
		name      string
		validator *types.Validator
		want      bool
	}{
		{
			name: "eligible",
			validator: &types.Validator{
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				EffectiveBalance: maxEffectiveBalance,
			},
			want: true,
		},
		{
			name: "not eligible, activation eligibility set",
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				EffectiveBalance:           maxEffectiveBalance,
			},
			want: false,
		},
		{
			name: "not eligible, effective balance too low",
			validator: &types.Validator{
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				EffectiveBalance: maxEffectiveBalance - 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsEligibleForActivationQueue(maxEffectiveBalance),
			)
		})
	}
}

func TestValidator_IsSlashable(t *testing.T) {
	tests := []struct {
		name      string
		epoch     math.Epoch
		validator *types.Validator
		want      bool
	}{
		{
			name:  "slashable",
			epoch: 10,
			validator: &types.Validator{
				Slashed:           false,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: true,
		},
		{
			name:  "not slashable, already slashed",
			epoch: 10,
			validator: &types.Validator{
				Slashed:           true,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: false,
		},
		{
			name:  "not slashable, before activation",
			epoch: 4,
			validator: &types.Validator{
				Slashed:           false,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: false,
		},
		{
			name:  "not slashable, after withdrawable",
			epoch: 16,
			validator: &types.Validator{
				Slashed:           false,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.validator.IsSlashable(tt.epoch))
		})
	}
}

func TestValidator_IsFullyWithdrawable(t *testing.T) {
	tests := []struct {
		name      string
		balance   math.Gwei
		epoch     math.Epoch
		validator *types.Validator
		want      bool
	}{
		{
			name:    "fully withdrawable",
			balance: 32e9,
			epoch:   10,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				WithdrawableEpoch: 5,
			},
			want: true,
		},
		{
			name:    "not fully withdrawable, non-eth1 credentials",
			balance: 32e9,
			epoch:   10,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					WithdrawalCredentials{0x00},
				WithdrawableEpoch: 5,
			},
			want: false,
		},
		{
			name:    "not fully withdrawable, before withdrawable epoch",
			balance: 32e9,
			epoch:   4,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				WithdrawableEpoch: 5,
			},
			want: false,
		},
		{
			name:    "not fully withdrawable, zero balance",
			balance: 0,
			epoch:   10,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				WithdrawableEpoch: 5,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsFullyWithdrawable(tt.balance, tt.epoch),
			)
		})
	}
}

func TestValidator_IsPartiallyWithdrawable(t *testing.T) {
	maxEffectiveBalance := math.Gwei(32e9)
	tests := []struct {
		name      string
		balance   math.Gwei
		validator *types.Validator
		want      bool
	}{
		{
			name:    "partially withdrawable",
			balance: 33e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				EffectiveBalance: maxEffectiveBalance,
			},
			want: true,
		},
		{
			name:    "not partially withdrawable, non-eth1 credentials",
			balance: 33e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.WithdrawalCredentials{
					0x00,
				},
				EffectiveBalance: maxEffectiveBalance,
			},
			want: false,
		},
		{
			name:    "not partially withdrawable, not at max effective balance",
			balance: 33e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				EffectiveBalance: maxEffectiveBalance - 1,
			},
			want: false,
		},
		{
			name:    "not partially withdrawable, no excess balance",
			balance: 32e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				EffectiveBalance: maxEffectiveBalance,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsPartiallyWithdrawable(
					tt.balance,
					maxEffectiveBalance,
				),
			)
		})
	}
}

func TestValidator_HasEth1WithdrawalCredentials(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
		want      bool
	}{
		{
			name: "has eth1 credentials",
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
			},
			want: true,
		},
		{
			name: "does not have eth1 credentials",
			validator: &types.Validator{
				WithdrawalCredentials: types.WithdrawalCredentials{
					0x00,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.HasEth1WithdrawalCredentials(),
			)
		})
	}
}

func TestValidator_HasMaxEffectiveBalance(t *testing.T) {
	maxEffectiveBalance := math.Gwei(32e9)
	tests := []struct {
		name      string
		validator *types.Validator
		want      bool
	}{
		{
			name: "has max effective balance",
			validator: &types.Validator{
				EffectiveBalance: maxEffectiveBalance,
			},
			want: true,
		},
		{
			name: "does not have max effective balance",
			validator: &types.Validator{
				EffectiveBalance: maxEffectiveBalance - 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.HasMaxEffectiveBalance(maxEffectiveBalance),
			)
		})
	}
}

func TestValidator_MarshalUnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name       string
		validator  *types.Validator
		invalidSSZ bool
	}{
		{
			name: "normal case",
			validator: &types.Validator{
				Pubkey: [48]byte{0x01},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
			invalidSSZ: false,
		},
		{
			name: "slashed validator",
			validator: &types.Validator{
				Pubkey: [48]byte{0x02},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x02},
					),
				EffectiveBalance:           32e9,
				Slashed:                    true,
				ActivationEligibilityEpoch: 5,
				ActivationEpoch:            6,
				ExitEpoch:                  10,
				WithdrawableEpoch:          15,
			},
			invalidSSZ: false,
		},
		{
			name: "validator with zero balance",
			validator: &types.Validator{
				Pubkey: [48]byte{0x03},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x03},
					),
				EffectiveBalance: 0,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
			invalidSSZ: false,
		},
		{
			name: "validator with non-default epochs",
			validator: &types.Validator{
				Pubkey: [48]byte{0x04},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x04},
					),
				EffectiveBalance:           16e9,
				Slashed:                    false,
				ActivationEligibilityEpoch: 10,
				ActivationEpoch:            12,
				ExitEpoch:                  20,
				WithdrawableEpoch:          25,
			},
			invalidSSZ: false,
		},
		{
			name:       "invalid SSZ size",
			validator:  &types.Validator{},
			invalidSSZ: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.invalidSSZ {
				// Create a byte slice with an invalid size (not 121)
				invalidSizeData := make([]byte, 120)
				var v types.Validator
				err := v.UnmarshalSSZ(invalidSizeData)
				require.Error(t, err, "Test case: %s", tt.name)
				require.Equal(t, io.ErrUnexpectedEOF, err,
					"Test case: %s", tt.name)
			} else {
				// Marshal the validator
				marshaled, err := tt.validator.MarshalSSZ()
				require.NoError(t, err)

				// Unmarshal into a new validator
				var unmarshaled types.Validator
				err = unmarshaled.UnmarshalSSZ(marshaled)
				require.NoError(t, err)

				// Check if the original and unmarshaled validators are equal
				require.Equal(
					t,
					tt.validator,
					&unmarshaled,
					"Test case: %s",
					tt.name,
				)

				var buf []byte
				buf, err = tt.validator.MarshalSSZTo(buf)
				require.NoError(t, err)

				// The two byte slices should be equal
				require.Equal(t, marshaled, buf)
			}
		})
	}
}

func TestValidator_HashTreeRoot(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
	}{
		{
			name: "normal case",
			validator: &types.Validator{
				Pubkey: [48]byte{0x01},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
		{
			name: "slashed validator",
			validator: &types.Validator{
				Pubkey: [48]byte{0x02},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x02},
					),
				EffectiveBalance:           32e9,
				Slashed:                    true,
				ActivationEligibilityEpoch: 5,
				ActivationEpoch:            6,
				ExitEpoch:                  10,
				WithdrawableEpoch:          15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test HashTreeRoot
			root := tt.validator.HashTreeRoot()
			require.NotEqual(t, [32]byte{}, root)

			// Test HashTreeRootWith
			hh := ssz.NewHasher()
			err := tt.validator.HashTreeRootWith(hh)
			require.NoError(t, err)

			// Test GetTree
			tree, err := tt.validator.GetTree()
			require.NoError(t, err)
			require.NotNil(t, tree)
		})
	}
}

func TestValidator_SetEffectiveBalance(t *testing.T) {
	tests := []struct {
		name      string
		balance   math.Gwei
		validator *types.Validator
		want      math.Gwei
	}{
		{
			name:    "set effective balance",
			balance: 32e9,
			validator: &types.Validator{
				EffectiveBalance: 0,
			},
			want: 32e9,
		},
		{
			name:    "update effective balance",
			balance: 16e9,
			validator: &types.Validator{
				EffectiveBalance: 32e9,
			},
			want: 16e9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validator.SetEffectiveBalance(tt.balance)
			require.Equal(t, tt.want, tt.validator.EffectiveBalance,
				"Test case: %s", tt.name)
		})
	}
}

func TestValidator_GetWithdrawableEpoch(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
		want      math.Epoch
	}{
		{
			name: "get withdrawable epoch",
			validator: &types.Validator{
				WithdrawableEpoch: 10,
			},
			want: 10,
		},
		{
			name: "get far future withdrawable epoch",
			validator: &types.Validator{
				WithdrawableEpoch: math.Epoch(constants.FarFutureEpoch),
			},
			want: math.Epoch(constants.FarFutureEpoch),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.validator.GetWithdrawableEpoch()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestValidator_GetWithdrawalCredentials(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
		want      types.WithdrawalCredentials
	}{
		{
			name: "get withdrawal credentials",
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
			},
			want: types.NewCredentialsFromExecutionAddress(
				gethprimitives.ExecutionAddress{0x01},
			),
		},
		{
			name: "get empty withdrawal credentials",
			validator: &types.Validator{
				WithdrawalCredentials: types.WithdrawalCredentials{0x00},
			},
			want: types.WithdrawalCredentials{0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.validator.GetWithdrawalCredentials()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestValidator_IsSlashed(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
		want      bool
	}{
		{
			name: "validator is slashed",
			validator: &types.Validator{
				Slashed: true,
			},
			want: true,
		},
		{
			name: "validator is not slashed",
			validator: &types.Validator{
				Slashed: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.validator.IsSlashed()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestValidator_New(t *testing.T) {
	tests := []struct {
		name                      string
		pubkey                    crypto.BLSPubkey
		withdrawalCredentials     types.WithdrawalCredentials
		amount                    math.Gwei
		effectiveBalanceIncrement math.Gwei
		maxEffectiveBalance       math.Gwei
		want                      *types.Validator
	}{
		{
			name:   "create new validator",
			pubkey: [48]byte{0x01},
			withdrawalCredentials: types.
				NewCredentialsFromExecutionAddress(
					gethprimitives.ExecutionAddress{0x01},
				),
			amount:                    32e9,
			effectiveBalanceIncrement: 1e9,
			maxEffectiveBalance:       32e9,
			want: &types.Validator{
				Pubkey: [48]byte{0x01},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						gethprimitives.ExecutionAddress{0x01},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &types.Validator{}
			got := v.New(
				tt.pubkey,
				tt.withdrawalCredentials,
				tt.amount,
				tt.effectiveBalanceIncrement,
				tt.maxEffectiveBalance,
			)
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestValidator_GetPubkey(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
		want      crypto.BLSPubkey
	}{
		{
			name: "get pubkey",
			validator: &types.Validator{
				Pubkey: [48]byte{0x01},
			},
			want: [48]byte{0x01},
		},
		{
			name: "get pubkey with multiple bytes",
			validator: &types.Validator{
				Pubkey: [48]byte{0x01, 0x02, 0x03},
			},
			want: [48]byte{0x01, 0x02, 0x03},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.validator.GetPubkey()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

func TestValidator_GetEffectiveBalance(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
		want      math.Gwei
	}{
		{
			name: "get effective balance",
			validator: &types.Validator{
				EffectiveBalance: 32e9,
			},
			want: 32e9,
		},
		{
			name: "get zero effective balance",
			validator: &types.Validator{
				EffectiveBalance: 0,
			},
			want: 0,
		},
		{
			name: "get maximum effective balance",
			validator: &types.Validator{
				EffectiveBalance: math.Gwei(1<<64 - 1),
			},
			want: math.Gwei(1<<64 - 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.validator.GetEffectiveBalance()
			require.Equal(t, tt.want, got, "Test case: %s", tt.name)
		})
	}
}

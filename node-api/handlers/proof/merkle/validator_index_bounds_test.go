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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package merkle_test

import (
	"testing"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle/mock"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// TestProveInBlock_RejectsOutOfBoundsValidatorIndex asserts a shared invariant
// that a validator index at or above the registry limit must be rejected.
func TestProveInBlock_RejectsOutOfBoundsValidatorIndex(t *testing.T) {
	t.Parallel()

	bs := mock.NewBeaconStateWith(
		4,
		ctypes.Validators{&ctypes.Validator{}},
		0,
		common.ExecutionAddress{},
		version.Electra(),
	)
	bs.Balances = []uint64{32000000000}
	bbh := ctypes.NewBeaconBlockHeader(
		4, 0, common.Root{1, 2, 3}, bs.HashTreeRoot(), common.Root{3, 2, 1},
	)

	provers := []struct {
		name string
		fn   func(math.U64) error
	}{
		{"ProveValidatorPubkeyInBlock", func(idx math.U64) error {
			_, _, err := merkle.ProveValidatorPubkeyInBlock(idx, bbh, bs)
			return err
		}},
		{"ProveWithdrawalCredentialsInBlock", func(idx math.U64) error {
			_, _, err := merkle.ProveWithdrawalCredentialsInBlock(idx, bbh, bs)
			return err
		}},
		{"ProveBalanceInBlock", func(idx math.U64) error {
			_, _, _, err := merkle.ProveBalanceInBlock(idx, bbh, bs, bs.Balances)
			return err
		}},
	}

	indices := []struct {
		name           string
		validatorIndex math.U64
	}{
		{"at registry limit", math.U64(constants.ValidatorsRegistryLimit)},
		{"above registry limit", math.U64(constants.ValidatorsRegistryLimit) + 1},
	}

	for _, p := range provers {
		t.Run(p.name, func(t *testing.T) {
			for _, idx := range indices {
				t.Run(idx.name, func(t *testing.T) {
					err := p.fn(idx.validatorIndex)
					require.ErrorContains(t, err, "exceeds registry limit")
				})
			}
		})
	}
}

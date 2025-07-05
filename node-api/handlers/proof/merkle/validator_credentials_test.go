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

package merkle_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle/mock"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// TestValidatorCredentialsProof tests the ProveWithdrawalCredentialsInBlock function
// and that the generated proof correctly verifies.
func TestValidatorCredentialsProof(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name              string
		forkVersion       common.Version
		numValidators     int
		slot              math.Slot
		valIndex          math.ValidatorIndex
		parentBlockRoot   common.Root
		bodyRoot          common.Root
		pubKey            crypto.BLSPubkey
		withdrawalAddress common.ExecutionAddress
		expectedProofFile string
	}{
		{
			name:              "1 Validator Set",
			forkVersion:       version.Electra(),
			numValidators:     1,
			slot:              4,
			valIndex:          0,
			parentBlockRoot:   common.Root{1, 2, 3},
			bodyRoot:          common.Root{3, 2, 1},
			pubKey:            [48]byte{9, 8, 7, 6, 5, 4, 3, 2, 1},
			withdrawalAddress: common.ExecutionAddress{1, 2, 3},
			expectedProofFile: "one_validator_validator_credentials_proof.json",
		},
		{
			name:              "Many Validator Set",
			forkVersion:       version.Electra(),
			numValidators:     100,
			slot:              5,
			valIndex:          95,
			parentBlockRoot:   common.Root{1, 2, 3, 4, 5, 6},
			bodyRoot:          common.Root{3, 2, 1, 9, 8, 7},
			pubKey:            [48]byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2},
			withdrawalAddress: common.ExecutionAddress{1, 2, 3},
			expectedProofFile: "many_validators_validator_credentials_proof.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vals := make(types.Validators, tc.numValidators)
			for i := range vals {
				vals[i] = &types.Validator{}
			}
			vals[tc.valIndex] = &types.Validator{
				Pubkey:                tc.pubKey,
				WithdrawalCredentials: types.NewCredentialsFromExecutionAddress(tc.withdrawalAddress),
			}

			bs := mock.NewBeaconStateWith(
				tc.slot, vals, 0, common.ExecutionAddress{}, tc.forkVersion,
			)

			stateRoot, err := bs.HashTreeRoot()
			require.NoError(t, err)
			bbh := types.NewBeaconBlockHeader(
				tc.slot,
				tc.valIndex,
				tc.parentBlockRoot,
				stateRoot,
				tc.bodyRoot,
			)

			proof, _, err := merkle.ProveWithdrawalCredentialsInBlock(
				tc.valIndex, bbh, bs,
			)
			require.NoError(t, err)
			expectedProof := ReadProofFromFile(t, tc.expectedProofFile)
			require.Equal(t, expectedProof, proof)
		})
	}
}

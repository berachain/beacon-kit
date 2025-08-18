//go:build test
// +build test

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

package genesis_test

import (
	libbytes "bytes"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/cli/utils/genesis"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

// generator for random deposits.
type TestDeposits []types.Deposit

func (TestDeposits) Generate(rand *rand.Rand, size int) reflect.Value {
	res := make(TestDeposits, size)
	for i := range size {
		var (
			pubKey crypto.BLSPubkey
			creds  types.WithdrawalCredentials
			sign   crypto.BLSSignature

			err error
		)
		_, err = rand.Read(pubKey[:])
		if err != nil {
			panic(fmt.Errorf("failed generating random pubKey: %w", err))
		}
		_, err = rand.Read(creds[:])
		if err != nil {
			panic(fmt.Errorf("failed generating random cred: %w", err))
		}
		_, err = rand.Read(sign[:])
		if err != nil {
			panic(fmt.Errorf("failed generating random sign: %w", err))
		}

		res[i] = types.Deposit{
			Pubkey:      pubKey,
			Credentials: creds,
			Amount:      math.Gwei(rand.Uint64()),
			Signature:   sign,
			Index:       0, // indexes will be set in order in the test
		}
	}
	return reflect.ValueOf(res)
}

func TestCompareGenesisCmdWithStateProcessor(t *testing.T) {
	t.Parallel()
	qc := &quick.Config{MaxCount: 1_000}
	csDev, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	csTest, err := spec.TestnetChainSpec()
	require.NoError(t, err)
	csMain, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	specs := []chain.Spec{csDev, csTest, csMain}
	for i, cs := range specs {
		t.Run(fmt.Sprintf("spec-%d", i), func(t *testing.T) {
			f := func(inputs TestDeposits) bool {
				deposits := make(types.Deposits, len(inputs))
				for i, input := range inputs {
					deposits[i] = &types.Deposit{
						Pubkey:      input.Pubkey,
						Credentials: input.Credentials,
						Amount:      input.Amount,
						Signature:   input.Signature,
						Index:       uint64(i),
					}
				}
				// genesis validators root from CLI
				cliValRoot := genesis.ComputeValidatorsRoot(deposits, cs)

				// genesis validators root from StateProcessor
				sp, st, _, _, _, _ := statetransition.SetupTestState(t, cs)
				genPayloadHeader := types.NewEmptyExecutionPayloadHeaderWithVersion(cs.GenesisForkVersion())

				_, err = sp.InitializeBeaconStateFromEth1(
					st,
					deposits,
					genPayloadHeader,
					cs.GenesisForkVersion(),
				)
				require.NoError(t, err)

				var processorRoot common.Root
				processorRoot, err = st.GetGenesisValidatorsRoot()
				require.NoError(t, err)

				// assert that they generate the same root, given the same list of deposits
				return libbytes.Equal(cliValRoot[:], processorRoot[:])
			}
			require.NoError(t, quick.Check(f, qc))
		})
	}
}

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
	"testing"

	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

var (
	pubKey1 = bytes.B48{0xff, 0xff, 0xff, 0xff}
	creds1  = types.WithdrawalCredentials{0xaa, 0xaa, 0xaa}
	amount1 = math.U64(2025)
	idx1    = int(0)

	pubKey2 = bytes.B48{0xee, 0xee, 0xee, 0xee}
	creds2  = types.WithdrawalCredentials{0xbb, 0xbb, 0xbb}
	amount2 = math.U64(5052)
	idx2    = int(1)

	emptySignature = crypto.BLSSignature{}

	genDeposits = types.Deposits{
		{
			Pubkey:      pubKey1,
			Credentials: creds1,
			Amount:      amount1,
			Signature:   emptySignature,
			Index:       uint64(idx1),
		},
		{
			Pubkey:      pubKey2,
			Credentials: creds2,
			Amount:      amount2,
			Signature:   emptySignature,
			Index:       uint64(idx2),
		},
	}

	expectedValRoot = common.Root{
		0xa3, 0xfa, 0xd, 0x97, 0x0, 0xeb, 0xdc, 0x2c,
		0x2, 0x1b, 0x51, 0xa1, 0xb, 0xcb, 0xb4, 0x80,
		0x5d, 0xe6, 0x13, 0x53, 0x9a, 0x77, 0xdc, 0x65,
		0xc7, 0x64, 0x85, 0x36, 0x6, 0xde, 0x4f, 0xa2,
	}
)

func TestOracle(t *testing.T) {
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	cliValRoot := genesis.ValidatorsRoot(genDeposits, cs)
	require.Equal(t, expectedValRoot, cliValRoot)
}

func TestStateTransitionGenesis(t *testing.T) {
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	sp, st, ds, ctx := statetransition.SetupTestState(t, cs)
	var (
		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
		genVersion       = version.FromUint32[common.Version](version.Deneb)
	)
	require.NoError(t, ds.EnqueueDeposits(ctx, genDeposits))
	_, err = sp.InitializePreminedBeaconStateFromEth1(
		st,
		genDeposits,
		genPayloadHeader,
		genVersion,
	)
	require.NoError(t, err)

	processorRoot, err := st.GetGenesisValidatorsRoot()
	require.NoError(t, err)
	require.Equal(t, expectedValRoot, processorRoot)
}

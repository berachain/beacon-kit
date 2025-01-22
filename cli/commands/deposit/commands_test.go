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

package deposit_test

import (
	"os"
	"path/filepath"
	"testing"
	"testing/quick"

	"github.com/berachain/beacon-kit/cli/commands/deposit"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtbls12381 "github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/privval"
	"github.com/stretchr/testify/require"
)

func TestCreateAndValidateCommandsDuality(t *testing.T) {
	qc := &quick.Config{MaxCount: 100}

	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	// create a tmp folder where test stores bls keys and
	// overwrite relevant files across test cases
	tmpFolder := t.TempDir()

	cometCfg := cmtcfg.DefaultConfig()
	cometCfg.RootDir = tmpFolder

	require.NoError(t, os.MkdirAll(filepath.Dir(cometCfg.PrivValidatorKeyFile()), 0o777))
	require.NoError(t, os.MkdirAll(filepath.Dir(cometCfg.PrivValidatorStateFile()), 0o777))

	f := func(
		blsKeySecret [32]byte,
		genValRoot common.Root,
		creds types.WithdrawalCredentials,
		amount math.Gwei,
	) bool {
		// generate random blsKey from the given secret
		var privKey *cmtbls12381.PrivKey
		privKey, err = cmtbls12381.GenPrivKeyFromSecret(blsKeySecret[:])
		require.NoError(t, err)

		// update relevant files and create corresponding blsSigner
		pv := privval.NewFilePV(privKey, cometCfg.PrivValidatorKeyFile(), cometCfg.PrivValidatorStateFile())
		pv.Save()
		blsSigner := signer.NewBLSSigner(cometCfg.PrivValidatorKeyFile(), cometCfg.PrivValidatorStateFile())

		// create the deposit and check that it verifies
		var (
			msg  *types.DepositMessage
			sign crypto.BLSSignature
		)
		msg, sign, err = deposit.CreateDepositMessage(cs, blsSigner, genValRoot, creds, amount)
		require.NoError(t, err)

		return deposit.ValidateDeposit(cs, msg.Pubkey, msg.Credentials, msg.Amount, genValRoot, sign) == nil
	}

	require.NoError(t, quick.Check(f, qc))
}

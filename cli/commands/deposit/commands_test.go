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
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	tmpFolder := t.TempDir()

	cometCfg := cmtcfg.DefaultConfig()
	cometCfg.RootDir = tmpFolder

	require.NoError(t, os.MkdirAll(filepath.Dir(cometCfg.PrivValidatorKeyFile()), 0o777))
	require.NoError(t, os.MkdirAll(filepath.Dir(cometCfg.PrivValidatorStateFile()), 0o777))

	for i := 1; i < 10; i++ {
		var privKey *cmtbls12381.PrivKey
		privKey, err = cmtbls12381.GenPrivKey()
		require.NoError(t, err)

		pv := privval.NewFilePV(privKey, cometCfg.PrivValidatorKeyFile(), cometCfg.PrivValidatorStateFile())
		pv.Save()

		blsSigner := signer.NewBLSSigner(cometCfg.PrivValidatorKeyFile(), cometCfg.PrivValidatorStateFile())
		genValRoot := common.Root{}
		creds := types.WithdrawalCredentials{}
		amount := math.Gwei(2025)

		var msg *types.DepositMessage
		var sign crypto.BLSSignature
		msg, sign, err = deposit.CreateDepositMessage(cs, blsSigner, genValRoot, creds, amount)
		require.NoError(t, err)

		require.NoError(t, deposit.ValidateDeposit(cs, msg.Pubkey, msg.Credentials, msg.Amount, genValRoot, sign))
	}
}

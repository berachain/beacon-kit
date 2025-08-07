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

package core

import (
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/ethereum/go-ethereum/params"
)

// TODO: confirm data
const (
	luganodesPubKey = "0xafd0ad061f698eae0d483098948e26e254f4b7089244bda897895257c668196ffd5e2ddf458fdf8bcea295b7d47a5b37"
	luganodesCreds  = "0x010000000000000000000000b0c615224a053236ac7d1c239f6c1b5fbf8f0617"
	luganodesAmount = 901_393_690_000_000 * params.Wei // 901_393.69 Gwei

	fixSmileePubKey = "0x84acfd38a13af12add8d82e1ef0842c4dfc1e4175fae5b8ab73770f9050cbf673cafdbf6d8ab679fe9ea13208f50b485"
)

//nolint:gochecknoglobals // unexported
var (
	luganodesCredentials = ctypes.WithdrawalCredentials(hex.MustToBytes(luganodesCreds))
	luganodesPubKeyBytes = crypto.BLSPubkey(hex.MustToBytes(luganodesPubKey))

	smileePubKey = crypto.BLSPubkey(hex.MustToBytes(fixSmileePubKey))
)

// processElectra1Fixes handles some fixes made necessary by accidents or wrong validator choices in mainnet
func (sp *StateProcessor) processElectra1Fixes(st *state.StateDB) error {
	if !sp.cs.IsMainnet() {
		return nil
	}

	if err := sp.processLuganodesRecovery(st); err != nil {
		return err
	}

	return sp.processSmileeFix(st)
}

func (sp *StateProcessor) processLuganodesRecovery(st *state.StateDB) error {
	sp.logger.Info("Luganodes recovery - creating deposit")

	// Make a one-time hardcoded deposit
	deposit := ctypes.Deposit{
		Pubkey:      luganodesPubKeyBytes,
		Credentials: luganodesCredentials,
		Amount:      luganodesAmount,
	}
	sp.logger.Info(
		"Luganodes deposit",
		"pubkey", deposit.GetPubkey().String(),
		"amount", deposit.GetAmount().Unwrap(),
		"credentials", deposit.GetWithdrawalCredentials().String(),
	)

	if err := sp.addValidatorToRegistry(st, &deposit); err != nil {
		return fmt.Errorf("failed to add Luganodes validator: %w", err)
	}

	sp.logger.Info("Luganodes recovery - created deposit")
	return nil
}

func (sp *StateProcessor) processSmileeFix(st *state.StateDB) error {
	sp.logger.Info("smilee fix - forcing validator exit")
	idx, err := st.ValidatorIndexByPubkey(smileePubKey)
	if err != nil {
		return fmt.Errorf("smilee fix - failed retrieving validator index: %w", err)
	}
	val, err := st.ValidatorByIndex(idx)
	if err != nil {
		return fmt.Errorf("smilee fix - failed retrieving validator: %w", err)
	}

	if err = sp.InitiateValidatorExit(st, idx); err != nil {
		return fmt.Errorf("smilee fix - failed initiating validator exit: %w", err)
	}
	sp.logger.Info(
		"smilee validator",
		"pubkey", val.Pubkey.String(),
		"amount", val.EffectiveBalance.Unwrap(),
		"credentials", val.WithdrawalCredentials.String(),
	)
	sp.logger.Info("smilee fix - successfully forced validator exit")
	return nil
}

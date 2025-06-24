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
	fixForkChainID    = uint64(80094) // avoid import cycle
	fixLuganodePubKey = "0xafd0ad061f698eae0d483098948e26e254f4b7089244bda897895257c668196ffd5e2ddf458fdf8bcea295b7d47a5b37"
	fixLuganodeAmount = 900_000 * params.GWei
)

//nolint:gochecknoglobals // unexported
var fixCredentials = ctypes.WithdrawalCredentials{0x00} // TODO: update this to a correct value

// processElectra1Fixes handles some fixes made necessary by accidents or wrong validator choices in mainnet
func (sp *StateProcessor) processElectra1Fixes(st *state.StateDB) error {
	if sp.cs.DepositEth1ChainID() == fixForkChainID {
		return nil
	}

	if err := sp.processLuganodesRecovery(st); err != nil {
		return err
	}
	return nil
}

func (sp *StateProcessor) processLuganodesRecovery(st *state.StateDB) error {
	sp.logger.Info("Luganode recovery - creating deposit")

	// Make a one-time hardcoded deposit
	deposit := ctypes.Deposit{
		Pubkey:      crypto.BLSPubkey(hex.MustToBytes(fixLuganodePubKey)),
		Credentials: fixCredentials,
		Amount:      fixLuganodeAmount,
	}
	sp.logger.Info(
		"Luganode deposit",
		"pubkey", deposit.GetPubkey().String(),
		"amount", deposit.GetAmount().Unwrap(),
		"credentials", deposit.GetWithdrawalCredentials().String(),
	)

	if err := sp.addValidatorToRegistry(st, &deposit); err != nil {
		return fmt.Errorf("failed to add Luganode validator: %w", err)
	}

	sp.logger.Info("Luganode recovery - created deposit")
	return nil
}

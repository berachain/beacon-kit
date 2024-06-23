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

package e2e_test

import (
	"context"
	"math/big"

	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types/tx"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (s *BeaconKitE2ESuite) Test4844Live() {
	// Sender account
	sender := s.TestAccounts()[1]

	ctx, cancel := context.WithTimeout(s.Ctx(), suite.DefaultE2ETestTimeout)
	defer cancel()

	chainID, err := s.JSONRPCBalancer().ChainID(ctx)
	s.Require().NoError(err)

	tip, err := s.JSONRPCBalancer().SuggestGasTipCap(ctx)
	s.Require().NoError(err)

	gasFee, err := s.JSONRPCBalancer().SuggestGasPrice(ctx)
	s.Require().NoError(err)

	tx := tx.New4844Tx(
		0, nil, 1000000,
		chainID, tip, gasFee, big.NewInt(0),
		[]byte{0x01, 0x02, 0x03, 0x04},
		big.NewInt(1), []byte{0x01, 0x02, 0x03, 0x04},
		types.AccessList{},
	)

	tx, err = sender.SignTx(chainID, tx)
	s.Require().NoError(err)

	s.Logger().Info("submitted blob transaction", "tx", tx.Hash().Hex())
	s.Require().NoError(s.JSONRPCBalancer().SendTransaction(ctx, tx))

	s.Logger().
		Info("waiting for blob transaction to be mined", "tx", tx.Hash().Hex())
	receipt, err := bind.WaitMined(ctx, s.JSONRPCBalancer(), tx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), receipt.Status)

	// Ensure Blob Tx doesn't cause liveliness issues.
	err = s.WaitForNBlockNumbers(10)
	s.Require().NoError(err)
}

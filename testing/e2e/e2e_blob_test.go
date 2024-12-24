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
	"fmt"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types/tx"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const (
	// NumBlobsLoad is the number of blob-carrying transactions to submit in
	// the Test4844Live test.
	NumBlobsLoad uint64 = 1

	// BlocksToWait4844 is the number of blocks to wait for the nodes to catch up.
	BlocksToWait4844 = 10
)

// Test4844Live tests sending a large number of blob carrying txs over the
// network.
func (s *BeaconKitE2ESuite) Test4844Live() {
	// Sender account
	sender := s.TestAccounts()[0]

	// Set the test timeout
	ctx, cancel := context.WithTimeout(s.Ctx(), suite.DefaultE2ETestTimeout)
	defer cancel()

	// Get the chain ID.
	chainID, err := s.JSONRPCBalancer().ChainID(ctx)
	s.Require().NoError(err)

	// Get the gas tip.
	tip, err := s.JSONRPCBalancer().SuggestGasTipCap(ctx)
	s.Require().NoError(err)

	// Get the gas fee
	gasFee, err := s.JSONRPCBalancer().SuggestGasPrice(ctx)
	s.Require().NoError(err)

	var blobTx *coretypes.Transaction
	for i := range NumBlobsLoad {
		// Get the block num.
		blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
		s.Require().NoError(err)

		// Get the nonce of our tx sender.
		//nolint:staticcheck // used below.
		nonce, err := s.JSONRPCBalancer().NonceAt(
			s.Ctx(), sender.Address(), new(big.Int).SetUint64(blkNum),
		)
		s.Require().NoError(err)

		// Craft the blob-carrying transaction.
		blobTx := tx.New4844Tx(
			nonce+i, nil, 1000000,
			chainID, tip, gasFee, big.NewInt(0),
			[]byte{0x01, 0x02, 0x03, 0x04},
			big.NewInt(1), []byte{0x01, 0x02, 0x03, 0x04},
			coretypes.AccessList{},
		)

		// Sign and submit the transaction.
		blobTx, err = sender.SignTx(chainID, blobTx)
		s.Require().NoError(err)
		s.Logger().Info("submitting blob transaction", "tx", blobTx.Hash().Hex())
		s.Require().NoError(s.JSONRPCBalancer().SendTransaction(ctx, blobTx))
	}

	// Wait for the last tx to be mined.
	s.Logger().
		Info("waiting for blob transaction to be mined", "tx", blobTx.Hash().Hex())
	receipt, err := bind.WaitMined(ctx, s.JSONRPCBalancer(), blobTx)
	s.Require().NoError(err)
	s.Require().Equal(coretypes.ReceiptStatusSuccessful, receipt.Status)

	// Ensure Blob Tx doesn't cause liveliness issues.
	err = s.WaitForNBlockNumbers(BlocksToWait4844)
	s.Require().NoError(err)

	client0 := s.ConsensusClients()[config.ClientValidator0]
	s.Require().NotNil(client0)

	opts := api.BlobSidecarsOpts{
		Block: receipt.BlockNumber.String(),
	}
	response, err := client0.BlobSidecars(ctx, &opts)
	s.Require().NoError(err)
	fmt.Println(response)

	// TODO: query and validate a sample (or all) of blob data from node-api
	// to ensure data availability.
	//client0.Connect()
}

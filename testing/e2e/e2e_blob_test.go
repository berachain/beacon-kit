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
	"bytes"
	"context"
	"encoding/binary"
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
	NumBlobsLoad uint64 = 10

	// BlocksToWait4844 is the number of blocks to wait for the nodes to catch up.
	BlocksToWait4844 = 1
)

// Test4844Live tests sending a large number of blob carrying txs over the
// network.
func (s *BeaconKitE2ESuite) Test4844Live() {
	// Set the test timeout
	ctx, cancel := context.WithTimeout(s.Ctx(), suite.DefaultE2ETestTimeout)
	defer cancel()

	// Connect the consensus client node-api
	client0 := s.ConsensusClients()[config.ClientValidator0]
	s.Require().NotNil(client0)
	s.Require().NoError(client0.Connect(ctx))

	// Grab values to plug into txs
	sender := s.TestAccounts()[0]
	chainID, err := s.JSONRPCBalancer().ChainID(ctx)
	s.Require().NoError(err)
	tip, err := s.JSONRPCBalancer().SuggestGasTipCap(ctx)
	s.Require().NoError(err)
	gasFee, err := s.JSONRPCBalancer().SuggestGasPrice(ctx)
	s.Require().NoError(err)
	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
	s.Require().NoError(err)
	nonce, err := s.JSONRPCBalancer().NonceAt(
		s.Ctx(), sender.Address(), new(big.Int).SetUint64(blkNum),
	)
	s.Require().NoError(err)

	// TODO: Query node-api for blobs to make sure they error

	// Craft and send each blob-carrying transaction.
	var blobTxs []*coretypes.Transaction
	for i := range NumBlobsLoad {
		blobData := make([]byte, 8)
		binary.LittleEndian.PutUint64(blobData, nonce+i)

		// Craft the blob-carrying transaction.
		blobTx := tx.New4844Tx(
			nonce+i, nil, 1000000,
			chainID, tip, gasFee, big.NewInt(0),
			[]byte{0x01, 0x02, 0x03, 0x04},
			big.NewInt(1), blobData,
			coretypes.AccessList{},
		)

		// Sign and submit the transaction.
		blobTx, err = sender.SignTx(chainID, blobTx)
		s.Require().NoError(err)
		s.Logger().Info("submitting blob transaction", "blobTx", blobTx.Hash().Hex())

		blobTxs = append(blobTxs, blobTx)

		err = s.JSONRPCBalancer().SendTransaction(ctx, blobTx)
		// TODO: Figure out why this error happens and why errors.Is(err, txpool.ErrAlreadyKnown) doesn't catch it
		if err != nil && err.Error() == "already known" {
			fmt.Println("FOUND ALREADY KNOWN")
			continue
		}
		fmt.Println(err)
		s.Require().NoError(err)
	}

	// TODO Make all node-api calls and verification asynchronous. node-api should be able to handle async calls

	// Wait for txs to be mined and group them by the block they get mined in.
	sidecarsByBlockNumber := make(map[string][]*coretypes.BlobTxSidecar)
	for _, blobTx := range blobTxs {
		s.Logger().
			Info("waiting for blob transaction to be mined", "blobTx", blobTx.Hash().Hex())
		receipt, err := bind.WaitMined(ctx, s.JSONRPCBalancer(), blobTx)
		s.Require().NoError(err)
		s.Require().Equal(coretypes.ReceiptStatusSuccessful, receipt.Status)

		blockNum := receipt.BlockNumber.String()

		if _, exists := sidecarsByBlockNumber[blockNum]; !exists {
			sidecarsByBlockNumber[blockNum] = []*coretypes.BlobTxSidecar{blobTx.BlobTxSidecar()}
		} else {
			sidecarsByBlockNumber[blockNum] = append(
				sidecarsByBlockNumber[blockNum], blobTx.BlobTxSidecar(),
			)
		}
	}

	// Validate each blob via the node-api.
	for blockNum := range sidecarsByBlockNumber {
		sidecarsInBlock := sidecarsByBlockNumber[blockNum]

		// Fetch blobs from node-api
		response, err := client0.BlobSidecars(ctx, &api.BlobSidecarsOpts{Block: blockNum})
		s.Require().NoError(err, "unable to fetch blob sidecars from node-api")

		// Verify blob data from each transaction is published by the node-api.
		for _, sidecar := range sidecarsInBlock {
			for i := range sidecar.Commitments {
				verified := false
				for _, blob := range response.Data {
					if bytes.Equal(blob.KZGCommitment[:], sidecar.Commitments[i][:]) {
						s.Require().Equal(sidecar.Blobs[i][:], blob.Blob[:], "blob data not equal")
						verified = true
					}
				}
				s.Require().True(verified, "unable to find blob commitment")
			}
		}
	}

	// Ensure Blob Tx doesn't cause liveliness issues.
	err = s.WaitForNBlockNumbers(BlocksToWait4844)
	s.Require().NoError(err)
}

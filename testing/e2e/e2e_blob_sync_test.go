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

package e2e_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types/tx"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	// NumBlocksWithBlobs is the number of blocks with blob transactions to create before restarting the syncing node.
	NumBlocksWithBlobs = 10
)

// TestBlobSync validates that a node can sync from behind and fetch blobs from other peers via the blob reactor.
// This test does the following steps:
// 1. Stop a full node to simulate it being offline
// 2. Produce several blocks with blob transactions while the node is down
// 3. Restart the full node so it needs to catch up
// 4. Verify that the node successfully fetches blobs from peers via P2P
func (s *BeaconKitE2ESuite) TestBlobSync() {
	ctx, cancel := context.WithTimeout(s.Ctx(), suite.DefaultE2ETestTimeout)
	defer cancel()

	// We use full node 0 for this test
	fullNodeELService := "el-full-reth-0"
	fullNodeCLService := config.ClientFullNode0

	// 1. Stop a full node to simulate it being offline
	//
	s.Logger().Info("Stopping full node to simulate being offline", "service", fullNodeCLService)
	err := s.StopService(ctx, fullNodeCLService)
	s.Require().NoError(err, "failed to stop full node consensus client")
	err = s.StopService(ctx, fullNodeELService)
	s.Require().NoError(err, "failed to stop full node execution client")
	s.Logger().Info("Full node stopped, now producing blocks with blobs while it's offline")

	// Set up connection to a validator's consensus client to produce blocks with blobs
	// while the full node is offline. Lets use validator 0 for this.
	//
	client0 := s.ConsensusClients()[config.ClientValidator0]
	s.Require().NotNil(client0)
	s.Require().NoError(client0.Connect(ctx))

	// Get initial block number before submitting blob transactions
	initialBlockNum, err := s.JSONRPCBalancer().BlockNumber(ctx)
	s.Require().NoError(err)
	s.Logger().Info("Initial block number", "block", initialBlockNum)

	// Prepare transaction parameters
	sender := s.TestAccounts()[0]
	chainID, err := s.JSONRPCBalancer().ChainID(ctx)
	s.Require().NoError(err)
	tip, err := s.JSONRPCBalancer().SuggestGasTipCap(ctx)
	s.Require().NoError(err)
	gasFee, err := s.JSONRPCBalancer().SuggestGasPrice(ctx)
	s.Require().NoError(err)
	nonce, err := s.JSONRPCBalancer().NonceAt(ctx, sender.Address(), new(big.Int).SetUint64(initialBlockNum))
	s.Require().NoError(err)

	// 2. Produce several blocks with blob transactions while the node is down
	//
	var (
		blobTxs          = make([]*coretypes.Transaction, 0)
		receipts         = make([]*coretypes.Receipt, 0)
		currentNonce     = nonce
		lastBlobBlockNum uint64
	)
	for blockIdx := range NumBlocksWithBlobs {
		// Each block can have 1-6 blob sidecars
		numBlobsInBlock := uint64((blockIdx % 6) + 1)

		s.Logger().Info("Creating block with blobs", "block_index", blockIdx, "num_blobs", numBlobsInBlock)
		for blobIdx := range numBlobsInBlock {
			// Create unique blob data for each transaction
			blobData := make([]byte, 8)
			binary.LittleEndian.PutUint64(blobData, currentNonce)

			// Craft blob-carrying transaction
			blobTx := tx.New4844Tx(
				currentNonce, nil, 1000000,
				chainID, tip, gasFee, big.NewInt(0),
				[]byte{0x01, 0x02, 0x03, 0x04},
				big.NewInt(1), blobData,
				coretypes.AccessList{},
			)

			// Sign and submit the transaction
			blobTx, err = sender.SignTx(chainID, blobTx)
			s.Require().NoError(err)
			s.Logger().Info("Submitting blob transaction",
				"tx_hash", blobTx.Hash().Hex(),
				"nonce", currentNonce,
				"block_index", blockIdx,
				"blob_index", blobIdx)

			err = s.JSONRPCBalancer().SendTransaction(ctx, blobTx)
			s.Require().NoError(err)
			blobTxs = append(blobTxs, blobTx)

			// Wait for this transaction to be mined
			receipt, errWait := bind.WaitMined(ctx, s.JSONRPCBalancer(), blobTx)
			s.Require().NoError(errWait)
			s.Require().Equal(coretypes.ReceiptStatusSuccessful, receipt.Status)
			receipts = append(receipts, receipt)

			// Track the highest block number
			if receipt.BlockNumber.Uint64() > lastBlobBlockNum {
				lastBlobBlockNum = receipt.BlockNumber.Uint64()
			}

			s.Logger().Info("Blob transaction mined",
				"tx_hash", blobTx.Hash().Hex(),
				"block", receipt.BlockNumber.Uint64(),
				"block_index", blockIdx,
				"blob_index", blobIdx)

			currentNonce++
		}
	}

	// 3. Restart the full node so it needs to catch up
	//
	err = s.StartService(ctx, fullNodeCLService)
	s.Require().NoError(err, "failed to start full node consensus client")
	err = s.StartService(ctx, fullNodeELService)
	s.Require().NoError(err, "failed to start full node execution client")
	s.Logger().Info("Full node restarted, waiting for it to sync to last blob block...", "last_blob_block", lastBlobBlockNum)
	s.Require().NoError(s.WaitForFinalizedBlockNumber(lastBlobBlockNum))

	// After catching up, the full node may need to wait a bit more for the blob fetcher to detect missing blobs
	s.Logger().Info("Waiting for blob fetcher to process queued blob requests...")
	time.Sleep(20 * time.Second)

	// 4. Verify that the node successfully fetches blobs from peers via P2P
	//
	s.Logger().Info("Setting up full node consensus client to verify blob sync...")
	err = s.SetupFullNodeConsensusClients()
	s.Require().NoError(err, "failed to setup full node consensus clients")
	fullNodeClient := s.FullNodeClients()[fullNodeCLService]
	s.Require().NotNil(fullNodeClient, "full node consensus client is nil")
	s.Logger().Info("Connecting to full node's consensus client...")
	err = fullNodeClient.Connect(ctx)
	s.Require().NoError(err, "failed to connect to full node consensus client")

	// Verify all blobs are accessible from the full node's node-api
	for i, receipt := range receipts {
		s.Logger().Info("Verifying blob availability on full node",
			"block", receipt.BlockNumber.Uint64(),
			"tx_index", i,
			"full_node", fullNodeCLService)

		// Fetch blob sidecars from the full node's node-api
		response, errAPI := fullNodeClient.BlobSidecars(ctx, &api.BlobSidecarsOpts{Block: receipt.BlockNumber.String()})
		s.Require().NoError(errAPI, "failed to fetch blob sidecars from full node for block %s", receipt.BlockNumber.String())
		s.Require().NotNil(response)
		s.Require().NotEmpty(response.Data, "no blob sidecars found on full node for block %s", receipt.BlockNumber.String())

		// Verify each blob commitment matches what was originally submitted
		sidecar := blobTxs[i].BlobTxSidecar()
		s.Require().NotNil(sidecar, "blob transaction %d missing sidecar", i)
		for j, commitment := range sidecar.Commitments {
			found := false
			for _, blob := range response.Data {
				if bytes.Equal(blob.KZGCommitment[:], commitment[:]) {
					s.Require().Equal(sidecar.Blobs[j][:], blob.Blob[:], "blob data mismatch on full node for tx %d blob %d", i, j)
					found = true
					break
				}
			}
			s.Require().True(found, "blob commitment not found on full node for tx %d blob %d", i, j)
		}

		s.Logger().Info("Blob verified successfully on full node",
			"block", receipt.BlockNumber.Uint64(),
			"tx_hash", blobTxs[i].Hash().Hex(),
			"num_blobs", len(sidecar.Commitments))
	}

	s.Logger().Info("Blob sync test completed successfully - full node fetched all blobs via P2P",
		"blocks_with_blobs", NumBlocksWithBlobs,
		"total_transactions", len(receipts),
		"full_node_synced", fullNodeCLService)
}

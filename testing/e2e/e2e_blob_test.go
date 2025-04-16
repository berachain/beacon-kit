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

package e2e_test

// import (
// 	"bytes"
// 	"context"
// 	"encoding/binary"
// 	"math/big"
// 	"sync"

// 	"github.com/attestantio/go-eth2-client/api"
// 	"github.com/berachain/beacon-kit/primitives/encoding/hex"
// 	"github.com/berachain/beacon-kit/testing/e2e/config"
// 	"github.com/berachain/beacon-kit/testing/e2e/suite"
// 	"github.com/berachain/beacon-kit/testing/e2e/suite/types/tx"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/core/txpool"
// 	coretypes "github.com/ethereum/go-ethereum/core/types"
// )

// const (
// 	// NumBlobsLoad is the number of blob-carrying transactions to submit in
// 	// the Test4844Live test. Cannot pool more than 16 txs currently.
// 	NumBlobsLoad uint64 = 16

// 	// BlocksToWait4844 is the number of blocks to wait for the nodes to catch up.
// 	BlocksToWait4844 = 5
// )

// // Test4844Live tests sending a large number of blob carrying txs over the
// // network.
// func (s *BeaconKitE2ESuite) Test4844Live() {
// 	ctx, cancel := context.WithTimeout(s.Ctx(), suite.DefaultE2ETestTimeout)
// 	defer cancel()

// 	// Connect the consensus client node-api
// 	client0 := s.ConsensusClients()[config.ClientValidator0]
// 	s.Require().NotNil(client0)
// 	s.Require().NoError(client0.Connect(ctx))

// 	// Grab values to plug into txs
// 	sender := s.TestAccounts()[0]
// 	chainID, err := s.JSONRPCBalancer().ChainID(ctx)
// 	s.Require().NoError(err)
// 	tip, err := s.JSONRPCBalancer().SuggestGasTipCap(ctx)
// 	s.Require().NoError(err)
// 	gasFee, err := s.JSONRPCBalancer().SuggestGasPrice(ctx)
// 	s.Require().NoError(err)
// 	blkNum, err := s.JSONRPCBalancer().BlockNumber(s.Ctx())
// 	s.Require().NoError(err)
// 	nonce, err := s.JSONRPCBalancer().NonceAt(
// 		s.Ctx(), sender.Address(), new(big.Int).SetUint64(blkNum),
// 	)
// 	s.Require().NoError(err)

// 	// Craft and send each blob-carrying transaction.
// 	blobTxs := make([]*coretypes.Transaction, 0, NumBlobsLoad)
// 	for i := range NumBlobsLoad {
// 		blobData := make([]byte, 8)

// 		// For the first 5 transactions, submit duplicate blobs in separate
// 		// transactions by leaving blobData empty. This is allowed by the
// 		// protocol.
// 		if i >= 5 {
// 			binary.LittleEndian.PutUint64(blobData, nonce+i)
// 		}

// 		// Craft the blob-carrying transaction.
// 		blobTx := tx.New4844Tx(
// 			nonce+i, nil, 1000000,
// 			chainID, tip, gasFee, big.NewInt(0),
// 			[]byte{0x01, 0x02, 0x03, 0x04},
// 			big.NewInt(1), blobData,
// 			coretypes.AccessList{},
// 		)

// 		// Sign and submit the transaction.
// 		blobTx, err = sender.SignTx(chainID, blobTx)
// 		s.Require().NoError(err)
// 		s.Logger().Info("submitting blob transaction", "blobTx", blobTx.Hash().Hex())
// 		blobTxs = append(blobTxs, blobTx)

// 		err = s.JSONRPCBalancer().SendTransaction(ctx, blobTx)
// 		// TODO: Figure out what is causing this to happen.
// 		// Also, `errors.Is(err, txpool.ErrAlreadyKnown)` doesn't catch it.
// 		if err != nil && err.Error() == txpool.ErrAlreadyKnown.Error() {
// 			continue
// 		}
// 		s.Require().NoError(err)
// 	}

// 	// All node-api calls and verification are asynchronous. node-api should
// 	// be able to handle async calls.
// 	var wg sync.WaitGroup
// 	for _, blobTx := range blobTxs {
// 		wg.Add(1)
// 		go func(blobTx *coretypes.Transaction) {
// 			defer wg.Done()

// 			// Wait for the blob transaction to be mined before making request.
// 			s.Logger().
// 				Info("waiting for blob transaction to be mined", "blobTx", blobTx.Hash().Hex())
// 			receipt, errWait := bind.WaitMined(ctx, s.JSONRPCBalancer(), blobTx)
// 			s.Require().NoError(errWait)
// 			s.Require().Equal(coretypes.ReceiptStatusSuccessful, receipt.Status)

// 			// WaitMined only waits until the tx is included in a block. This
// 			// gets triggered whenever beacon-kit sends a FCU including the
// 			// block. In the optimistic builder, this happens before the block
// 			// is finalized. Meaning the data has not yet been stored. Let's
// 			// just wait 1 block.
// 			//
// 			//nolint:contextcheck // uses the service context.
// 			s.Require().NoError(s.WaitForNBlockNumbers(1))

// 			// Fetch blobs from node-api.
// 			response, errAPI := client0.BlobSidecars(ctx, &api.BlobSidecarsOpts{Block: receipt.BlockNumber.String()})
// 			s.Require().NoError(errAPI, "unable to fetch blob sidecars from node-api")

// 			// Verify blob data from each transaction is published by the node-api.
// 			sidecar := blobTx.BlobTxSidecar()
// 			for i, commitment := range sidecar.Commitments {
// 				verified := false
// 				for _, blob := range response.Data {
// 					if bytes.Equal(blob.KZGCommitment[:], commitment[:]) {
// 						s.Require().Equal(sidecar.Blobs[i][:], blob.Blob[:], "blob data not equal")
// 						verified = true
// 					}
// 				}
// 				s.Require().True(verified, "blob data was not made available by node-api")
// 				s.Logger().Info("verified blob data availability", "KzgCommitment", hex.EncodeBytes(commitment[:]))
// 			}
// 		}(blobTx)
// 	}

// 	// Wait for DA validation to finish
// 	wg.Wait()

// 	// Ensure Blob Tx doesn't cause liveliness issues.
// 	err = s.WaitForNBlockNumbers(BlocksToWait4844)
// 	s.Require().NoError(err)
// }

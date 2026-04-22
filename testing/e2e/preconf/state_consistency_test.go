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

//go:build e2e

package preconf_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	stateConsistencyTxsPerWorker = 20
	stateConsistencyPollInterval = 25 * time.Millisecond
	stateConsistencyPerTxTimeout = 30 * time.Second
)

type workerResult struct {
	pendingBalMismatches     int
	pendingNonceMismatches   int
	inclusionBalMismatches   int
	inclusionNonceMismatches int
	err                      error
}

// TestFlashblockStateConsistency verifies that the preconf RPC exposes *exactly correct* recipient balance
// and sender nonce at both the "pending" tag and at the inclusion block after each individual transaction.
func (s *PreconfE2ESuite) TestFlashblockStateConsistency() {
	// Run test with single sender, fresh recipient, N sequential TXs.
	s.Run("Serial", s.runStateConsistencySerial)
	// Run test with one worker per test account, each with its own recipient.
	s.Run("Parallel", s.runStateConsistencyParallel)
}

func (s *PreconfE2ESuite) runStateConsistencySerial() {
	ctx := s.Ctx()
	preconf := s.PreconfRPCClients(0)
	gasTipCap, gasFeeCap := s.suggestGasCaps(preconf.Client)

	r := s.runStateConsistencyWorker(
		ctx, preconf, s.TestAccounts()[0], newRandomAddress(),
		stateConsistencyTxsPerWorker, gasTipCap, gasFeeCap,
	)
	s.Require().NoError(r.err)
	s.Require().Zero(r.pendingBalMismatches)
	s.Require().Zero(r.pendingNonceMismatches)
	s.Require().Zero(r.inclusionBalMismatches)
	s.Require().Zero(r.inclusionNonceMismatches)
}

func (s *PreconfE2ESuite) runStateConsistencyParallel() {
	ctx := s.Ctx()
	preconf := s.PreconfRPCClients(0)
	gasTipCap, gasFeeCap := s.suggestGasCaps(preconf.Client)
	senders := s.TestAccounts()

	results := make([]workerResult, len(senders))
	var wg sync.WaitGroup
	for i := range senders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = s.runStateConsistencyWorker(
				ctx, preconf, senders[i], newRandomAddress(),
				stateConsistencyTxsPerWorker, gasTipCap, gasFeeCap,
			)
		}()
	}
	wg.Wait()

	for i, r := range results {
		s.Require().NoError(r.err, "worker %d fatal error", i)
		s.Require().Zero(r.pendingBalMismatches)
		s.Require().Zero(r.pendingNonceMismatches)
		s.Require().Zero(r.inclusionBalMismatches)
		s.Require().Zero(r.inclusionNonceMismatches)
	}
}

// runStateConsistencyWorker sends numTxs ETH transfers and, for each TX, asserts exact recipient
// balance and sender nonce on the preconf RPC at both the "pending" tag and the inclusion block.
func (s *PreconfE2ESuite) runStateConsistencyWorker(
	ctx context.Context,
	preconf *types.ExecutionClient,
	sender *types.EthAccount,
	recipient common.Address,
	numTxs int,
	gasTipCap, gasFeeCap *big.Int,
) workerResult {
	var r workerResult
	txValue := new(big.Int).SetUint64(suite.Ether / 1000) //nolint:mnd // 0.001 ETH
	senderAddr := sender.Address()

	startNonce, err := preconf.PendingNonceAt(ctx, senderAddr)
	if err != nil {
		r.err = fmt.Errorf("get initial pending nonce for %s: %w", senderAddr, err)
		return r
	}
	lastScanned, err := preconf.BlockNumber(ctx)
	if err != nil {
		r.err = fmt.Errorf("get initial block number: %w", err)
		return r
	}

	var firstBlock, lastBlock uint64
	for i := range numTxs {
		expectedBal := new(big.Int).Mul(txValue, big.NewInt(int64(i+1)))
		expectedNonce := startNonce + uint64(i+1)

		tx, sErr := s.sendETHTransfer(transferParams{
			client:    preconf.Client,
			sender:    sender,
			to:        recipient,
			nonce:     startNonce + uint64(i),
			amount:    txValue,
			gasTipCap: gasTipCap,
			gasFeeCap: gasFeeCap,
		})
		if sErr != nil {
			r.err = fmt.Errorf("send tx %d from %s: %w", i+1, senderAddr, sErr)
			return r
		}

		if !waitForTxInPending(ctx, preconf, tx.Hash()) {
			r.err = fmt.Errorf("tx %s from %s timed out waiting for pending inclusion", tx.Hash().Hex(), senderAddr)
			return r
		}
		pendingBal, bErr := preconf.PendingBalanceAt(ctx, recipient)
		if bErr != nil {
			r.err = fmt.Errorf("get pending balance for %s: %w", recipient, bErr)
			return r
		}
		if pendingBal.Cmp(expectedBal) != 0 {
			r.pendingBalMismatches++
			s.T().Errorf("pending balance mismatch: tx=%d recipient=%s got=%s want=%s", i+1, recipient, pendingBal, expectedBal)
		}
		pendingNonce, nErr := preconf.PendingNonceAt(ctx, senderAddr)
		if nErr != nil {
			r.err = fmt.Errorf("get pending nonce for %s: %w", senderAddr, nErr)
			return r
		}
		if pendingNonce != expectedNonce {
			r.pendingNonceMismatches++
			s.T().Errorf("pending nonce mismatch: tx=%d sender=%s got=%d want=%d", i+1, senderAddr, pendingNonce, expectedNonce)
		}

		inclBlock := waitForTxInclusion(ctx, preconf, tx.Hash(), lastScanned)
		if inclBlock == 0 {
			r.err = fmt.Errorf("tx %s from %s timed out waiting for inclusion", tx.Hash().Hex(), senderAddr)
			return r
		}
		if firstBlock == 0 {
			firstBlock = inclBlock
		}
		lastBlock = inclBlock
		lastScanned = inclBlock

		inclBlockBig := new(big.Int).SetUint64(inclBlock)
		inclBal, bErr := preconf.BalanceAt(ctx, recipient, inclBlockBig)
		if bErr != nil {
			r.err = fmt.Errorf("get balance for %s at block %d: %w", recipient, inclBlock, bErr)
			return r
		}
		if inclBal.Cmp(expectedBal) != 0 {
			r.inclusionBalMismatches++
			s.T().Errorf("inclusion balance mismatch: tx=%d recipient=%s block=%d got=%s want=%s",
				i+1, recipient, inclBlock, inclBal, expectedBal)
		}
		inclNonce, nErr := preconf.NonceAt(ctx, senderAddr, inclBlockBig)
		if nErr != nil {
			r.err = fmt.Errorf("get nonce for %s at block %d: %w", senderAddr, inclBlock, nErr)
			return r
		}
		if inclNonce != expectedNonce {
			r.inclusionNonceMismatches++
			s.T().Errorf("inclusion nonce mismatch: tx=%d sender=%s block=%d got=%d want=%d",
				i+1, senderAddr, inclBlock, inclNonce, expectedNonce)
		}
	}

	s.T().Logf("verified %d txs from %s to %s in blocks %d-%d",
		numTxs, senderAddr.Hex()[:10], recipient.Hex()[:10], firstBlock, lastBlock)
	return r
}

// waitForTxInPending polls the preconf pending block until hash appears (or times out)
func waitForTxInPending(ctx context.Context, preconf *types.ExecutionClient, hash common.Hash) bool {
	ctx, cancel := context.WithTimeout(ctx, stateConsistencyPerTxTimeout)
	defer cancel()
	tick := time.NewTicker(stateConsistencyPollInterval)
	defer tick.Stop()
	for {
		pb, err := preconf.BlockByNumber(ctx, big.NewInt(-1))
		if err == nil && pb != nil {
			for _, tx := range pb.Transactions() {
				if tx.Hash() == hash {
					return true
				}
			}
		}
		select {
		case <-ctx.Done():
			return false
		case <-tick.C:
		}
	}
}

// waitForTxInclusion scans newly-committed blocks (from lastScanned+1) on the given client until the tx is found.
func waitForTxInclusion(ctx context.Context, client *types.ExecutionClient, hash common.Hash, lastScanned uint64) uint64 {
	ctx, cancel := context.WithTimeout(ctx, stateConsistencyPerTxTimeout)
	defer cancel()
	tick := time.NewTicker(stateConsistencyPollInterval)
	defer tick.Stop()
	scanFrom := lastScanned + 1
	for {
		block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(scanFrom))
		if err != nil || block == nil {
			// Block not yet committed (or transient RPC error); wait then retry.
			select {
			case <-ctx.Done():
				return 0
			case <-tick.C:
			}
			continue
		}
		for _, tx := range block.Transactions() {
			if tx.Hash() == hash {
				return scanFrom
			}
		}
		scanFrom++
	}
}

// newRandomAddress returns a random 20-byte address.
func newRandomAddress() common.Address {
	var b [common.AddressLength]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("crypto/rand.Read: %w", err))
	}
	return common.Address(b)
}

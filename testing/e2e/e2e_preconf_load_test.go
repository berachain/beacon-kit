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

//go:build e2e_preconf

package e2e_test

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	// Kurtosis service names.
	loadTestPreconfRPCEL   = "el-preconf-rpc-reth-0"
	loadTestPreconfRPCPort = "eth-json-rpc"

	// Load parameters.
	loadTestDuration           = 60 * time.Second       // how long the load test phase runs
	loadTestTxPerFlashblock    = 40                     // txs per flashblock (0 to disable)
	loadTestSpamTxPerSec       = 100                    // background tx to spam on fullnode rpc (0 to disable)
	loadTestFlashblockInterval = 200 * time.Millisecond // the flashblock cadence

	// Performance thresholds.
	loadTestP50LatencyThreshold = 250 * time.Millisecond // 250ms to give us some room on slow CI
	loadTestP99LatencyThreshold = 750 * time.Millisecond
	loadTestMaxInclusionDelay   = 2
)

// transferParams groups the parameters for an ETH transfer.
type transferParams struct {
	client    *ethclient.Client
	sender    *types.EthAccount
	to        common.Address
	nonce     uint64
	amount    *big.Int
	gasTipCap *big.Int
	gasFeeCap *big.Int
}

// txResult captures the preconf latency and block inclusion info for a tx.
type txResult struct {
	latency        time.Duration
	sendBlock      uint64
	inclusionBlock uint64
	err            error
}

// pendingTx tracks a sent transaction awaiting block inclusion.
type pendingTx struct {
	hash      common.Hash
	sendTime  time.Time
	sendBlock uint64
}

// PreconfLoadE2ESuite tests the preconf system by sending real ETH
// transactions through the preconf RPC node, measuring flashblock
// latency, and verifying state consistency with the standard RPC.
type PreconfLoadE2ESuite struct {
	suite.KurtosisE2ESuite
	preconfClient *ethclient.Client
	chainID       *big.Int
}

// TestPreconfLoadE2ESuite runs the preconf load test suite.
func TestPreconfLoadE2ESuite(t *testing.T) {
	suite.Run(t, new(PreconfLoadE2ESuite))
}

// SetupSuite initializes the network with a dedicated sequencer and
// preconf RPC node, then discovers the preconf RPC endpoint.
func (s *PreconfLoadE2ESuite) SetupSuite() {
	s.SetupSuiteWithOptions(suite.WithPreconfLoadConfig())

	// Discover preconf RPC EL node via Kurtosis port mapping.
	sCtx, err := s.Enclave().GetServiceContext(loadTestPreconfRPCEL)
	s.Require().NoError(err, "Should get preconf RPC EL service context")

	port, ok := sCtx.GetPublicPorts()[loadTestPreconfRPCPort]
	s.Require().True(ok, "Preconf RPC EL should expose eth-json-rpc port")

	preconfURL := fmt.Sprintf("http://0.0.0.0:%d", port.GetNumber())
	s.T().Logf("Preconf RPC EL URL: %s", preconfURL)

	s.preconfClient, err = types.DialWithPooling(preconfURL)
	s.Require().NoError(err, "Should connect to preconf RPC EL")
	s.T().Cleanup(func() { s.preconfClient.Close() })

	s.chainID, err = s.RPCClient().ChainID(s.Ctx())
	s.Require().NoError(err, "Should get chain ID")

	// Brief warmup: confirm network is producing blocks after funding.
	err = s.WaitForNBlockNumbers(1)
	s.Require().NoError(err, "Network should produce warmup blocks")
}

// TestPreconfTransactions fires bursts of parallel ETH transfers through
// the preconf RPC node every flashblock interval and waits for receipts.
// This measures per-transaction flashblock inclusion latency under
// realistic load while a background spammer creates mempool pressure
// on the standard full node RPC.
//
//nolint:funlen // load test with multiple phases
func (s *PreconfLoadE2ESuite) TestPreconfTransactions() {
	ctx := s.Ctx()
	sender := s.TestAccounts()[0]
	receiver := s.TestAccounts()[1]
	senderAddr := sender.Address()
	receiverAddr := receiver.Address()

	spammer := s.TestAccounts()[2]
	spammerAddr := spammer.Address()

	initialNonce, err := s.preconfClient.NonceAt(ctx, senderAddr, nil)
	s.Require().NoError(err, "Should get initial sender nonce")

	initialBalance, err := s.preconfClient.BalanceAt(ctx, receiverAddr, nil)
	s.Require().NoError(err, "Should get initial receiver balance")

	spammerInitialNonce, err := s.RPCClient().NonceAt(ctx, spammerAddr, nil)
	s.Require().NoError(err, "Should get initial spammer nonce")

	// Pre-compute gas caps once to avoid per-tx RPC calls.
	gasTipCap, gasFeeCap := s.suggestGasCaps(s.preconfClient)

	// Start background load on the standard RPC to create mempool pressure.
	stopSpammer := s.startBackgroundLoad(ctx)
	defer func() { stopSpammer() }()

	transferAmt := new(big.Int).SetUint64(suite.Ether / 100)
	nonce := initialNonce

	maxExpectedTxs := int(loadTestDuration/loadTestFlashblockInterval) * loadTestTxPerFlashblock

	// Channel for sender goroutines to report successfully sent txs.
	pendingTxCh := make(chan pendingTx, maxExpectedTxs)

	// Start block-scanning collector with a generous timeout.
	collectorCtx, collectorCancel := context.WithTimeout(ctx, loadTestDuration+30*time.Second)
	defer collectorCancel()
	collectorDone := make(chan []txResult, 1)
	go func() {
		collectorDone <- s.collectResults(collectorCtx, pendingTxCh)
	}()

	ticker := time.NewTicker(loadTestFlashblockInterval)
	defer ticker.Stop()
	deadline := time.After(loadTestDuration)

loop:
	for {
		select {
		case <-deadline:
			break loop
		case <-ticker.C:
			sendBlock, bErr := s.preconfClient.BlockNumber(ctx)
			s.Require().NoError(bErr, "Should get block number at send time")

			// Send txs sequentially to avoid nonce gaps. The collector
			// handles receipt detection concurrently via block scanning.
			for range loadTestTxPerFlashblock {
				sendTime := time.Now()
				tx, sErr := s.sendETHTransfer(transferParams{
					client:    s.preconfClient,
					sender:    sender,
					to:        receiverAddr,
					nonce:     nonce,
					amount:    transferAmt,
					gasTipCap: gasTipCap,
					gasFeeCap: gasFeeCap,
				})
				if sErr != nil {
					break // stop burst, retry same nonce on next tick
				}
				nonce++
				pendingTxCh <- pendingTx{
					hash:      tx.Hash(),
					sendTime:  sendTime,
					sendBlock: sendBlock,
				}
			}
		}
	}

	// Signal collector that no more txs are coming.
	close(pendingTxCh)

	// Wait for collector to discover all sent txs in blocks.
	results := <-collectorDone

	// Check for collection errors.
	for _, r := range results {
		s.Require().NoError(r.err, "All txs should be included")
	}

	// Stop the background spammer and collect sent tx hashes.
	spamHashes := stopSpammer()
	s.T().Logf("Background spammer sent %d txs total", len(spamHashes))

	// Determine the max inclusion block from preconf results OR wait for
	// the current head when no preconf txs were sent (spam-only mode).
	var maxInclusionBlock uint64

	if loadTestTxPerFlashblock > 0 {
		// Compute and assert latency stats.
		latencies := make([]time.Duration, len(results))
		var totalLatency time.Duration
		for i, r := range results {
			latencies[i] = r.latency
			totalLatency += r.latency
		}
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

		avg := totalLatency / time.Duration(len(latencies))
		p50 := durationPercentile(latencies, 0.50)
		p99 := durationPercentile(latencies, 0.99)

		s.T().Logf("Preconf latency: avg=%v p50=%v p99=%v min=%v max=%v (n=%d)",
			avg, p50, p99, latencies[0], latencies[len(latencies)-1], len(latencies))

		s.Require().LessOrEqual(p50, loadTestP50LatencyThreshold,
			"p50 preconf latency %v exceeds %v threshold", p50, loadTestP50LatencyThreshold)
		s.Require().LessOrEqual(p99, loadTestP99LatencyThreshold,
			"p99 preconf latency %v exceeds %v threshold", p99, loadTestP99LatencyThreshold)

		// Assert inclusion delay and find max inclusion block.
		for i, r := range results {
			delay := r.inclusionBlock - r.sendBlock
			s.Require().LessOrEqual(delay, uint64(loadTestMaxInclusionDelay),
				"tx %d included %d blocks late (sendBlock=%d inclusionBlock=%d)",
				i, delay, r.sendBlock, r.inclusionBlock)
			if r.inclusionBlock > maxInclusionBlock {
				maxInclusionBlock = r.inclusionBlock
			}
		}
	}

	// Wait for the full node to reach the max inclusion block (or current
	// head when no preconf txs) so that state is settled before verification.
	if maxInclusionBlock == 0 {
		maxInclusionBlock, err = s.RPCClient().BlockNumber(ctx)
		s.Require().NoError(err)
		maxInclusionBlock += 3 //nolint:mnd // extra blocks for pending spam txs to settle
	}
	s.T().Logf("Waiting for full node to reach block %d", maxInclusionBlock)
	err = s.WaitForFinalizedBlockNumber(maxInclusionBlock)
	s.Require().NoError(err, "full node should reach block %d", maxInclusionBlock)

	// Verify receipts for all background spam transactions.
	// Done after finalization wait so that the last pending spam txs
	// have time to be included.
	spamVerifyStart := time.Now()
	for i, h := range spamHashes {
		receipt, rErr := s.RPCClient().TransactionReceipt(ctx, h)
		s.Require().NoError(rErr, "spam tx %d (%s) should have a receipt", i, h.Hex())
		s.Require().Equal(ethtypes.ReceiptStatusSuccessful, receipt.Status,
			"spam tx %d (%s) should succeed", i, h.Hex())
	}
	s.T().Logf("All %d spam tx receipts verified in %v", len(spamHashes), time.Since(spamVerifyStart))

	// Verify final nonce and balance at the max inclusion block.
	verifyBlock := new(big.Int).SetUint64(maxInclusionBlock)

	finalNonce, err := s.RPCClient().NonceAt(ctx, senderAddr, verifyBlock)
	s.Require().NoError(err)
	s.Require().Equal(initialNonce+uint64(len(results)), finalNonce,
		"Sender nonce should increment by %d", len(results))

	expectedIncrease := new(big.Int).Mul(transferAmt, big.NewInt(int64(len(results))))
	finalBalance, err := s.RPCClient().BalanceAt(ctx, receiverAddr, verifyBlock)
	s.Require().NoError(err)
	actualIncrease := new(big.Int).Sub(finalBalance, initialBalance)
	s.Require().Equal(expectedIncrease.String(), actualIncrease.String(),
		"Receiver balance should increase by %s", expectedIncrease)

	// Verify spammer nonce is at least the number of tracked sends.
	// Under high load, SendTransaction may return an error for a tx that
	// was actually accepted by the node, so the on-chain nonce can be
	// slightly higher than len(spamHashes).
	spammerFinalNonce, err := s.RPCClient().NonceAt(ctx, spammerAddr, nil)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(spammerFinalNonce, spammerInitialNonce+uint64(len(spamHashes)),
		"Spammer nonce should be at least %d", len(spamHashes))

	// Verify state consistency between both RPCs at the same block.
	for _, addr := range []common.Address{senderAddr, receiverAddr, spammerAddr} {
		stdBal, bErr := s.RPCClient().BalanceAt(ctx, addr, verifyBlock)
		s.Require().NoError(bErr)
		preconfBal, bErr := s.preconfClient.BalanceAt(ctx, addr, verifyBlock)
		s.Require().NoError(bErr)
		s.Require().Equal(stdBal.String(), preconfBal.String(),
			"Balance mismatch for %s at block %d", addr, maxInclusionBlock)

		stdNonce, nErr := s.RPCClient().NonceAt(ctx, addr, verifyBlock)
		s.Require().NoError(nErr)
		preconfNonce, nErr := s.preconfClient.NonceAt(ctx, addr, verifyBlock)
		s.Require().NoError(nErr)
		s.Require().Equal(stdNonce, preconfNonce,
			"Nonce mismatch for %s at block %d", addr, maxInclusionBlock)
	}

	s.T().Logf("State consistency verified at block %d", maxInclusionBlock)

	// Log eth transaction count per block from genesis to end of test.
	s.T().Log("Eth transactions count in block:")
	for bn := uint64(0); bn <= maxInclusionBlock; bn++ {
		block, bErr := s.RPCClient().BlockByNumber(ctx, new(big.Int).SetUint64(bn))
		if bErr != nil {
			s.T().Logf("  block %d: error fetching: %v", bn, bErr)
			continue
		}
		s.T().Logf("  block %d: %d txs", bn, len(block.Transactions()))
	}
}

// collectResults polls finalized and pending blocks on the preconf RPC
// to discover when sent transactions are included. This replaces per-tx
// receipt polling with a single scanning loop that makes far fewer RPC
// calls, avoiding 429 rate-limiting under heavy load.
func (s *PreconfLoadE2ESuite) collectResults(
	ctx context.Context,
	pendingTxCh <-chan pendingTx,
) []txResult {
	unseen := make(map[common.Hash]pendingTx)
	var results []txResult
	var lastScannedBlock uint64
	channelOpen := true
	lastLogTime := time.Now()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for channelOpen || len(unseen) > 0 {
		select {
		case <-ctx.Done():
			s.T().Logf("Collector timeout: %d unseen, %d collected, lastScanned=%d",
				len(unseen), len(results), lastScannedBlock)
			for hash, ptx := range unseen {
				s.T().Logf("  unseen tx=%s sendBlock=%d age=%v",
					hash.Hex(), ptx.sendBlock, time.Since(ptx.sendTime))
				results = append(results, txResult{
					err: fmt.Errorf("timeout: tx %s not included", hash.Hex()),
				})
			}
			return results
		case <-ticker.C:
		}

		// Drain newly sent txs (non-blocking).
	drain:
		for {
			select {
			case ptx, ok := <-pendingTxCh:
				if !ok {
					channelOpen = false
					s.T().Logf("Collector: send phase done, %d unseen, %d collected, lastScanned=%d",
						len(unseen), len(results), lastScannedBlock)
					break drain
				}
				unseen[ptx.hash] = ptx
			default:
				break drain
			}
		}

		if len(unseen) == 0 {
			continue
		}

		// Periodic progress log (every 10s).
		if time.Since(lastLogTime) > 10*time.Second {
			s.T().Logf("Collector: %d unseen, %d collected, lastScanned=%d",
				len(unseen), len(results), lastScannedBlock)
			lastLogTime = time.Now()
		}

		// Scan newly finalized blocks.
		currentBlock, err := s.preconfClient.BlockNumber(ctx)
		if err != nil {
			s.T().Logf("Collector: BlockNumber failed: %v", err)
			continue
		}
		if lastScannedBlock == 0 && currentBlock > 0 {
			lastScannedBlock = currentBlock - 1
		}
		for bn := lastScannedBlock + 1; bn <= currentBlock; bn++ {
			block, bErr := s.preconfClient.BlockByNumber(ctx, new(big.Int).SetUint64(bn))
			if bErr != nil {
				s.T().Logf("Collector: BlockByNumber(%d) failed: %v (lastScanned=%d current=%d unseen=%d)",
					bn, bErr, lastScannedBlock, currentBlock, len(unseen))
				break
			}
			now := time.Now()
			for _, tx := range block.Transactions() {
				if ptx, found := unseen[tx.Hash()]; found {
					results = append(results, txResult{
						latency:        now.Sub(ptx.sendTime),
						sendBlock:      ptx.sendBlock,
						inclusionBlock: bn,
					})
					delete(unseen, tx.Hash())
				}
			}
			lastScannedBlock = bn
		}

		// Check pending block for txs preconfirmed via flashblocks
		// but not yet finalized.
		if len(unseen) > 0 {
			pb, pErr := s.preconfClient.BlockByNumber(ctx, big.NewInt(-1))
			if pErr == nil && pb != nil {
				now := time.Now()
				for _, tx := range pb.Transactions() {
					if ptx, found := unseen[tx.Hash()]; found {
						results = append(results, txResult{
							latency:        now.Sub(ptx.sendTime),
							sendBlock:      ptx.sendBlock,
							inclusionBlock: pb.NumberU64(),
						})
						delete(unseen, tx.Hash())
					}
				}
			}
		}
	}

	return results
}

// startBackgroundLoad spams the standard RPC with self-transfers
// from testAccounts[2] to create realistic mempool pressure. Returns a
// stop func that cancels the spammer, waits for it to exit, and returns
// the hashes of all successfully sent transactions.
func (s *PreconfLoadE2ESuite) startBackgroundLoad(
	ctx context.Context,
) (stop func() []common.Hash) {
	if loadTestSpamTxPerSec == 0 {
		return func() []common.Hash { return nil }
	}

	ctx, cancel := context.WithCancel(ctx)

	spammer := s.TestAccounts()[2]
	spammerAddr := spammer.Address()
	client := s.RPCClient()

	nonce, err := client.NonceAt(ctx, spammerAddr, nil)
	s.Require().NoError(err, "Should get spammer nonce")

	gasTipCap, gasFeeCap := s.suggestGasCaps(client.Client)
	amount := new(big.Int).SetUint64(1) // 1 wei self-transfers

	var (
		wg     sync.WaitGroup
		hashes []common.Hash
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		rate := time.Duration(loadTestSpamTxPerSec) // variable to avoid compile-time division by zero
		ticker := time.NewTicker(time.Second / rate)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			tx, sErr := spammer.SignTx(
				s.chainID, ethtypes.NewTx(&ethtypes.DynamicFeeTx{
					ChainID:   s.chainID,
					Nonce:     nonce,
					GasTipCap: gasTipCap,
					GasFeeCap: gasFeeCap,
					Gas:       suite.EtherTransferGasLimit,
					To:        &spammerAddr,
					Value:     amount,
				}),
			)
			if sErr != nil {
				continue
			}

			if sErr = client.SendTransaction(ctx, tx); sErr != nil {
				// Nonce may be stale, refresh it.
				if refreshed, nErr := client.NonceAt(ctx, spammerAddr, nil); nErr == nil {
					nonce = refreshed
				}
				continue
			}

			hashes = append(hashes, tx.Hash())
			nonce++
		}
	}()

	return func() []common.Hash {
		cancel()
		wg.Wait()
		return hashes
	}
}

// sendETHTransfer creates, signs, and sends an ETH transfer via the given
// client. It is goroutine-safe: errors are returned, not asserted.
func (s *PreconfLoadE2ESuite) sendETHTransfer(
	p transferParams,
) (*ethtypes.Transaction, error) {
	signedTx, err := p.sender.SignTx(
		s.chainID, ethtypes.NewTx(&ethtypes.DynamicFeeTx{
			ChainID:   s.chainID,
			Nonce:     p.nonce,
			GasTipCap: p.gasTipCap,
			GasFeeCap: p.gasFeeCap,
			Gas:       suite.EtherTransferGasLimit,
			To:        &p.to,
			Value:     p.amount,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("sign tx nonce=%d: %w", p.nonce, err)
	}

	if err = p.client.SendTransaction(s.Ctx(), signedTx); err != nil {
		return nil, fmt.Errorf("send tx nonce=%d: %w", p.nonce, err)
	}

	return signedTx, nil
}

// suggestGasCaps queries the client for gas tip and fee caps, falling back
// to defaults if the RPC method is unsupported.
func (s *PreconfLoadE2ESuite) suggestGasCaps(
	client *ethclient.Client,
) (gasTipCap, gasFeeCap *big.Int) {
	gasTipCap, err := client.SuggestGasTipCap(s.Ctx())
	if err != nil {
		var rpcErr rpc.Error
		if errors.As(err, &rpcErr) && rpcErr.ErrorCode() == -32601 {
			gasTipCap = new(big.Int).SetUint64(suite.TenGwei)
		} else {
			s.Require().NoError(err, "Should get gas tip cap")
		}
	}
	gasFeeCap = new(big.Int).Add(gasTipCap, new(big.Int).SetUint64(suite.TenGwei))
	return gasTipCap, gasFeeCap
}

// durationPercentile returns the p-th percentile from a sorted slice of durations.
func durationPercentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}

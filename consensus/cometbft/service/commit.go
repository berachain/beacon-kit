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
//

package cometbft

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"cosmossdk.io/store/rootmulti"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
)

func (s *Service) commit(
	*cmtabci.CommitRequest,
) (*cmtabci.CommitResponse, error) {
	if _, _, err := s.cachedStates.GetFinal(); err != nil {
		// This is unexpected since CometBFT should call Commit only
		// after FinalizeBlock has been called. Panic appeases nilaway.
		panic(fmt.Errorf("commit: %w", err))
	}

	// The cached state context carries an empty block header, so use the height and time captured from FinalizeBlock instead.
	retainHeight := s.GetBlockRetentionHeight(s.finalizedHeight)

	rms, ok := s.sm.GetCommitMultiStore().(*rootmulti.Store)
	if ok {
		rms.SetCommitHeader(cmtproto.Header{ChainID: s.chainID, Height: s.finalizedHeight, Time: s.finalizedTime})
	}
	s.sm.GetCommitMultiStore().Commit()

	s.cachedStates.Reset()

	if s.blockDelay != nil {
		if err := s.sm.SaveBlockDelay(s.blockDelay.ToBytes()); err != nil {
			panic(fmt.Errorf("failed to save block delay: %w", err))
		}
	}

	s.haltIfReached()

	return &cmtabci.CommitResponse{
		RetainHeight: retainHeight,
	}, nil
}

// haltPointReached reports whether a block at the given height and time has reached the configured halt-height
// or halt-time. It is the single halt predicate, applied to the last finalized block by ensureNotHalted and haltIfReached.
func haltPointReached(haltHeight, haltTime uint64, height int64, blockTime time.Time) bool {
	unixTime := blockTime.Unix()
	switch {
	case haltHeight > 0 && height >= 0 && uint64(height) >= haltHeight:
		return true
	case haltTime > 0 && unixTime >= 0 && uint64(unixTime) >= haltTime:
		return true
	default:
		return false
	}
}

// ensureNotHalted returns an error once the last finalized block has reached the halt point. Service start
// refuses to run with it and the ABCI handlers gate on it (see abci.go), so a node with the halt flags still
// set neither advances state past the halt block nor creeps one block per restart.
func (s *Service) ensureNotHalted() error {
	if !haltPointReached(s.haltHeight, s.haltTime, s.finalizedHeight, s.finalizedTime) {
		return nil
	}
	return fmt.Errorf(
		"chain reached the configured halt point (halt-height %d, halt-time %d) at committed height %d, unset the halt flags to resume",
		s.haltHeight, s.haltTime, s.finalizedHeight,
	)
}

// haltShutdownSlack bounds how long a parked ABCI call waits after the app context is cancelled. The halt
// shutdown normally exits the process within it, so the caller's error return (and the CometBFT panic it
// causes) only happens on a wedged shutdown.
//
//nolint:gochecknoglobals // var instead of const so tests can shorten it
var haltShutdownSlack = 30 * time.Second

// waitForHaltShutdown parks an ABCI call that must not proceed past the halt point while the halt shutdown
// brings the process down.
func (s *Service) waitForHaltShutdown() {
	<-s.ctx.Done()
	time.Sleep(haltShutdownSlack)
}

// haltGracePeriod keeps Commit blocked after the halt block so vote gossip can deliver the halt-block
// precommits to validators still one short. Exiting immediately can wedge those peers at the previous height,
// and a wedged set larger than 1/3 cannot recover after the restart (individual precommit signatures do not
// survive commit aggregation).
const haltGracePeriod = 5 * time.Second

// haltIfReached gracefully shuts down the node once the committed block reaches the configured halt-height or
// halt-time. It runs after the block has been fully committed, so a restarted node resumes consensus at the
// next height with no replay needed.
func (s *Service) haltIfReached() {
	if !haltPointReached(s.haltHeight, s.haltTime, s.finalizedHeight, s.finalizedTime) {
		return
	}

	s.logger.Info("halting node per configuration",
		"halt_height", s.haltHeight, "halt_time", s.haltTime, "committed_height", s.finalizedHeight, "grace_period", haltGracePeriod)

	// Sleeping here blocks the consensus state machine inside Commit, so no halting node can advance to the
	// next height while its peer gossip routines keep serving the halt-block precommits from the live vote set.
	time.Sleep(haltGracePeriod)

	// Signal our own process so the node's regular shutdown path runs, the same mechanism a cosmos-sdk baseapp
	// uses for halt-height.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		os.Exit(0)
	}
	if err = p.Signal(syscall.SIGINT); err != nil {
		if err = p.Signal(syscall.SIGTERM); err != nil {
			os.Exit(0)
		}
	}
}

// GetBlockRetentionHeight returns the height for which all blocks below this
// height
// are pruned from CometBFT. Given a commitment height and a non-zero local
// minRetainBlocks configuration, the retentionHeight is the smallest height
// that
// satisfies:
//
// - Unbonding (safety threshold) time: The block interval in which validators
// can be economically punished for misbehavior. Blocks in this interval must be
// auditable e.g. by the light client.
//
// - Logical store snapshot interval: The block interval at which the underlying
// logical store database is persisted to disk, e.g. every 10000 heights. Blocks
// since the last IAVL snapshot must be available for replay on application
// restart.
//
// - State sync snapshots: Blocks since the oldest available snapshot must be
// available for state sync nodes to catch up (oldest because a node may be
// restoring an old snapshot while a new snapshot was taken).
//
// - Local (minRetainBlocks) config: Archive nodes may want to retain more or
// all blocks, e.g. via a local config option min-retain-blocks. There may also
// be a need to vary retention for other nodes, e.g. sentry nodes which do not
// need historical blocks.
func (s *Service) GetBlockRetentionHeight(commitHeight int64) int64 {
	// pruning is disabled if minRetainBlocks is zero
	if s.minRetainBlocks == 0 {
		return 0
	}

	minNonZero := func(x, y int64) int64 {
		switch {
		case x == 0:
			return y

		case y == 0:
			return x

		case x < y:
			return x

		default:
			return y
		}
	}

	// Define retentionHeight as the minimum value that satisfies all non-zero
	// constraints. All blocks below (commitHeight-retentionHeight) are pruned
	// from CometBFT.
	var retentionHeight int64

	// Define the number of blocks needed to protect against misbehaving
	// validators
	// which allows light clients to operate safely. Note, we piggy back of the
	// evidence parameters instead of computing an estimated number of blocks
	// based
	// on the unbonding period and block commitment time as the two should be
	// equivalent.
	if _, _, err := s.cachedStates.GetFinal(); err != nil {
		return 0
	}
	cp := s.cmtConsensusParams.ToProto()
	if cp.Evidence != nil && cp.Evidence.MaxAgeNumBlocks > 0 {
		retentionHeight = commitHeight - cp.Evidence.MaxAgeNumBlocks
	}

	v := commitHeight - int64(s.minRetainBlocks) // #nosec G115
	retentionHeight = minNonZero(retentionHeight, v)

	if retentionHeight <= 0 {
		// prune nothing in the case of a non-positive height
		return 0
	}

	return retentionHeight
}

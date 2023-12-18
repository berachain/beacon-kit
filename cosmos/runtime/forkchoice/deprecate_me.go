// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package forkchoice

import (
	"context"
	"time"

	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

func (m *Service) SyncEl(ctx context.Context) error {
	// Trigger the execution client to begin building the block, and update
	// the proposers forkchoice state accordingly.

	// TODO `block` needs to come from the latest blocked stored IAVL tree
	// on the consensus client.
	// block, _ := m.EngineCaller.EarliestBlock(ctx)

	m.getForkchoiceFromExecutionClient(ctx)
	// genesisHash := m.ek.RetrieveGenesis(ctx)
	m.logger.Info("waiting for execution client to finish sync")
	for {
		var err error
		fc := m.curForkchoiceState
		payloadID, lastValidHash, err := m.EngineCaller.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
		if err == nil {
			break
		}

		m.curForkchoiceState = fc
		m.logger.Info("waiting for execution client to sync", "error", err)
		m.logger.Info("waiting for execution client to sync", "payloadID", payloadID, "lastValidHash", lastValidHash)
		time.Sleep(1 * time.Second)
	}

	return nil
}

// TODO DEPRECATE ME
func (m *Service) getForkchoiceFromExecutionClient(ctx context.Context) error {
	var latestBlock *enginev1.ExecutionBlock
	var err error
	latestBlock, err = m.EngineCaller.LatestExecutionBlock(ctx)
	if err != nil {
		m.logger.Error("failed to get block number", "err", err)
		return err
	}

	m.curForkchoiceState.HeadBlockHash = latestBlock.Hash.Bytes()
	m.logger.Info("forkchoice state", "head", latestBlock.Header.Hash())

	safe, err := m.EngineCaller.LatestSafeBlock(ctx)
	if err != nil {
		m.logger.Error("failed to get safe block", "err", err)
		safe = latestBlock
	}

	m.curForkchoiceState.SafeBlockHash = safe.Hash.Bytes()

	final, err := m.EngineCaller.LatestFinalizedBlock(ctx)
	m.logger.Info("forkchoice state", "finalized", safe.Hash)
	if err != nil {
		m.logger.Error("failed to get final block", "err", err)
		final = latestBlock
	}

	m.curForkchoiceState.FinalizedBlockHash = final.Hash.Bytes()
	m.logger.Info("forkchoice state", "finalized", final.Hash)

	return nil
}

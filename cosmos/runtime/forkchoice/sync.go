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
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

const (
	syncCheckDelay = 3
)

// WaitforExecutionClientSync waits for the execution client to sync up with the beacon chain.
func (m *Service) WaitforExecutionClientSync(ctx context.Context) error {
	var (
		err error
	)

	// Wait for the sync to complete before starting the consensus client.
	sp := &ethereum.SyncProgress{}
	for sp.CurrentBlock < sp.HighestBlock {
		sp, err = m.EngineCaller.SyncProgress(ctx)
		if err != nil {
			m.logger.Error("failed to get sync progress", "err", err)
			return err
		}
		m.logger.Info("waiting for sync to complete", "current_block", sp.CurrentBlock,
			"highest_block", sp.HighestBlock)
		time.Sleep(syncCheckDelay * time.Second)
	}

	m.logger.Info("execution client is synced", "current_block", sp.CurrentBlock)

	blk, err := m.EngineCaller.BlockByNumber(ctx, big.NewInt(int64(sp.CurrentBlock)))
	if err != nil {
		m.logger.Error("failed to get block by number", "err", err)
		return err
	}

	fc := &enginev1.ForkchoiceState{
		HeadBlockHash:      blk.Hash().Bytes(),
		SafeBlockHash:      common.Hash(m.ek.ForkChoiceStore(ctx).GetSafeBlockHash()).Bytes(),
		FinalizedBlockHash: common.Hash(m.ek.ForkChoiceStore(ctx).GetFinalizedBlockHash()).Bytes(),
	}

	id, x, err := m.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
	if err != nil {
		return err
	}

	m.logger.Info("forkchoice updated", "id", id, "x", x)

	return nil
}

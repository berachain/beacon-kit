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

package block_sync

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"time"

	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	"github.com/itsdevbear/bolaris/types/config"
	v1 "github.com/itsdevbear/bolaris/types/v1"
)

type HeadSubscriber interface {
	engine.Caller
	// SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereub.Subscription, error)
}

type ForkChoiceStoreProvider interface {
	ForkChoiceStore(ctx context.Context) v1.ForkChoiceStore
}

// BlockSync is responsible for managing the synchornization of the execution client.
type BlockSync struct {
	logger         log.Logger
	beaconCfg      *config.Beacon
	fcsp           ForkChoiceStoreProvider
	headSubscriber HeadSubscriber
}

func New(opts ...Option) *BlockSync {
	b := &BlockSync{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *BlockSync) Start(ctx context.Context) {
	b.logger.Info("starting block sync service for execution client...")
	go b.loop(ctx)
}

func (b *BlockSync) loop(ctx context.Context) {
	// ch := make(chan *types.Header)
	// b.headSubscriber.SubscribeNewHead(ctx, ch)
	// for {
	// 	select {
	// 	case header := <-ch:
	// 		b.logger.Info("received new header"n, "header", header)
	// 	case <-ctx.Done():
	// 		return
	// 	}
	// }
}

// WaitforExecutionClientSync waits for the execution client to sync up with the beacon chain.
func (b *BlockSync) WaitforExecutionClientSync(ctx context.Context) error {
	var err error
	blk := b.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	if bytes.Equal(blk[:], (common.Hash{}).Bytes()) {
		blk = b.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
		if bytes.Equal(blk[:], (common.Hash{}).Bytes()) {
			blk_, err := b.headSubscriber.BlockByNumber(ctx, new(big.Int))
			if err != nil {
				return err
			}
			blk = [32]byte(blk_.Hash())
		}
	}

	fc := &enginev1.ForkchoiceState{
		HeadBlockHash:      blk[:],
		SafeBlockHash:      common.Hash(b.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()).Bytes(),
		FinalizedBlockHash: common.Hash(b.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()).Bytes(),
	}

	b.logger.Info("setting execution client forkchoice startup",
		"head", common.Bytes2Hex(fc.HeadBlockHash),
		"safe", common.Bytes2Hex(fc.SafeBlockHash),
		"finalized", common.Bytes2Hex(fc.FinalizedBlockHash))
retry:
	_, latestValidHash, err := b.headSubscriber.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
	if err != nil {
		if errors.Is(err, engine.ErrSyncingPayloadStatus) {
			b.logger.Info("payload on sync is acepted or syncing, retrying....",
				"latestValidHash", latestValidHash)
			time.Sleep(1 * time.Second)
			goto retry
		}
		b.logger.Error("invalid forkchoice at startup",
			"head", common.Bytes2Hex(fc.HeadBlockHash),
			"safe", common.Bytes2Hex(fc.SafeBlockHash),
			"finalized", common.Bytes2Hex(fc.FinalizedBlockHash))
		return err
	}

	b.logger.Info("execution client forkchoice startup complete", "latestValidHash", latestValidHash)
	return nil
}

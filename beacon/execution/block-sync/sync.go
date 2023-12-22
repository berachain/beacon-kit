package block_sync

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"time"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	"github.com/itsdevbear/bolaris/types/config"
	v1 "github.com/itsdevbear/bolaris/types/v1"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

type HeadSubscriber interface {
	engine.Caller
	// SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
}

type ForkChoiceStoreProvider interface {
	ForkChoiceStore(ctx context.Context) v1.ForkChoiceStore
}

// BlockSync is responsible for managing the synchornization of the execution client
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
	// 		b.logger.Info("received new header", "header", header)
	// 	case <-ctx.Done():
	// 		return
	// 	}
	// }
}

// WaitforExecutionClientSync waits for the execution client to sync up with the beacon chain.
func (m *BlockSync) WaitforExecutionClientSync(ctx context.Context) error {
	var err error

	// // Wait for the sync to complete before starting the consensus client.
	// sp := &ethereum.SyncProgress{}
	// for sp.CurrentBlock < sp.HighestBlock {
	// 	sp, err = m.Caller.SyncProgress(ctx)
	// 	if err != nil {
	// 		m.logger.Error("failed to get sync progress", "err", err)
	// 		return err
	// 	}
	// 	m.logger.Info("waiting for sync to complete", "current_block", sp.CurrentBlock,
	// 		"highest_block", sp.HighestBlock)
	// 	time.Sleep(syncCheckDelay * time.Second)
	// }

	// m.logger.Info("execution client is reporting sync'd", "current_block", sp.CurrentBlock)

	blk := m.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	if bytes.Equal(blk[:], (common.Hash{}).Bytes()) {
		blk = m.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
		if bytes.Equal(blk[:], (common.Hash{}).Bytes()) {
			blk_, err := m.headSubscriber.BlockByNumber(ctx, new(big.Int))
			if err != nil {
				return err
			}
			blk = [32]byte(common.Hash(blk_.Hash()).Bytes())
		}
	}

	fc := &enginev1.ForkchoiceState{
		HeadBlockHash:      blk[:],
		SafeBlockHash:      common.Hash(m.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()).Bytes(),
		FinalizedBlockHash: common.Hash(m.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()).Bytes(),
	}

	m.logger.Info("setting execution client forkchoice startup",
		"head", common.Bytes2Hex(fc.HeadBlockHash),
		"safe", common.Bytes2Hex(fc.SafeBlockHash),
		"finalized", common.Bytes2Hex(fc.FinalizedBlockHash))
retry:
	_, latestValidHash, err := m.headSubscriber.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
	if err != nil {
		if errors.Is(err, execution.ErrAcceptedSyncingPayloadStatus) {
			m.logger.Info("payload on sync is acepted or syncing, retrying....",
				"latestValidHash", latestValidHash)
			time.Sleep(1 * time.Second)
			goto retry
		}
		m.logger.Error("invalid forkchoice at startup",
			"head", common.Bytes2Hex(fc.HeadBlockHash),
			"safe", common.Bytes2Hex(fc.SafeBlockHash),
			"finalized", common.Bytes2Hex(fc.FinalizedBlockHash))
		return err
	}

	m.logger.Info("execution client forkchoice startup complete", "latestValidHash", latestValidHash)
	return nil
}

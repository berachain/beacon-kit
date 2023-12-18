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

// Package forker implements the Ethereum forker.
package forkchoice

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	prysmexecution "github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	pb "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/execution"
	BeaconKeeper "github.com/itsdevbear/bolaris/cosmos/x/beacon/keeper"
)

// Service implements the baseapp.TxSelector interface.
type Service struct {
	execution.EngineCaller
	logger    log.Logger
	ek        *BeaconKeeper.Keeper
	etherbase common.Address
	bk        *blockchain.Service

	// TODO DEPRECATE ME
	curForkchoiceState *pb.ForkchoiceState

	// TODO: handle this better with a proper LRU cache.
	cachedPayload *pb.PayloadIDBytes
}

// New produces a cosmos forker from a geth forker.
func New(gm execution.EngineCaller, bk *blockchain.Service, ek *BeaconKeeper.Keeper, logger log.Logger) *Service {
	return &Service{
		EngineCaller:       gm,
		curForkchoiceState: &pb.ForkchoiceState{},
		logger:             logger,
		ek:                 ek,
		bk:                 bk,
	}
}

// TODO DEPRECATE ME
func (m *Service) getForkchoiceFromExecutionClient(ctx context.Context) error {
	var latestBlock *pb.ExecutionBlock
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

		fmt.Println("PAYLOAD ID")
		m.curForkchoiceState = fc
		m.logger.Info("waiting for execution client to sync", "error", err)
		m.logger.Info("waiting for execution client to sync", "payloadID", payloadID, "lastValidHash", lastValidHash)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (m *Service) BuildBlockV2(ctx sdk.Context) (interfaces.ExecutionData, error) {
	// Reference: https://hackmd.io/@danielrachi/engine_api#Block-Building
	// builder := (&execution.Builder{EngineCaller: m.EngineCaller.(*execution.Service)})
	m.logger.Info("entering build-block-v2")
	defer m.logger.Info("exiting build-block-v2")

	attrs, err := m.getPayloadAttributes(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: SHOULD THIS BE LATEST OR FINALIZED????
	// IN THEORY LATEST MEANS THAT THE PROPOSER CAN APPEND ARBITRARY BLOCKS
	// AND SET THE CANONICAL CHAIN TO WHATEVER IT WANTS (i.e) could
	// apply a deep deep reorg?
	// On the flip side, if we force LatestFinalizedBlock(), we can
	// at the Consensus Layer ensure the reorg will always be 1 deep.
	b, err := m.EngineCaller.LatestFinalizedBlock(ctx)
	if err != nil {
		m.logger.Error("failed to get block number", "err", err)
	}

	// We start by setting the head of our execution client to the
	// latest block that we have seen.
	zeroHash := [32]byte{}
	fc := &pb.ForkchoiceState{
		HeadBlockHash:      b.Hash.Bytes(),
		SafeBlockHash:      zeroHash[:],
		FinalizedBlockHash: zeroHash[:],
	}

	var payloadID *pb.PayloadIDBytes
	if m.cachedPayload == nil {
		m.logger.Info("No cached payload found, building new payload")
		// Trigger the execution client to begin building the block, and update
		// the proposers forkchoice state accordingly. By setting the HeadBlockHash
		// to the last finalized block that we have seen, we are telling the execution client
		// to begin building a block on top of that block.
		payloadID, _, err = m.EngineCaller.ForkchoiceUpdated(ctx, fc, attrs)
		if err != nil {
			m.logger.Error("failed to get forkchoice updated", "err", err)
			return nil, err
		}

		// TODO: this should be something that is 80% of proposal timeout, or so?
		// TODO: maybe this should be some sort of event that we wait for?
		// But the TLDR is that we need to wait for the execution client to
		// build the payload before we can include it in the beacon block.
		time.Sleep(8000 * time.Millisecond) //nolint:gomnd // temp.
	} else {
		payloadID = m.cachedPayload
	}
	m.logger.Info("calling getPayload", "id", payloadID)

	// Get the Payload From the Execution Client
	builtPayload, _, _, err := m.EngineCaller.GetPayload(
		ctx, *payloadID, primitives.Slot(ctx.BlockHeight()))
	if err != nil {
		m.logger.Error("failed to get previosuly queued payload", "err", err, "payloadID", payloadID)
		return nil, err
	}

	// This CL node's execution client head is now at the head of this
	// payload, so we can update our forkchoice state accordingly.
	var latestValidHash []byte
	_, latestValidHash, err = m.EngineCaller.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
	if errors.Is(err, prysmexecution.ErrInvalidPayloadStatus) {
		// If we
		m.logger.Error("invalid payload, proposing last valid", "err", err)
		fc := &pb.ForkchoiceState{
			HeadBlockHash:      latestValidHash,
			SafeBlockHash:      zeroHash[:],
			FinalizedBlockHash: zeroHash[:],
		}
		_, _, err2 := m.EngineCaller.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
		if err2 != nil {
			return nil, err
		}
	} else if err != nil {
		m.logger.Error("failed to get forkchoice updated", "err", err)
		return nil, err
	}

	// Go and include the builtPayload on the BeaconBlock
	return builtPayload, nil
}

func (m *Service) ValidateBlock(ctx sdk.Context, builtPayload interfaces.ExecutionData) error {
	// // NewPayload calls "InsertBlockWithoutSetHead"
	// lastValidHash, _ := m.NewPayload(ctx, builtPayload, []common.Hash{}, &common.Hash{} /*empty version hashes and root before Deneb*/)
	// fmt.Println("LAST VALID HASH FOUND ON ETH ONE", common.Bytes2Hex(lastValidHash))
	payload, err := m.bk.ProcessBlock(ctx, ctx.HeaderInfo(), builtPayload)
	if err != nil {
		return err
	}
	m.cachedPayload = payload
	// // TODO FIX, rn we are just blindly finalizing whatever the proposer has sent us.
	// // The blind finalization is "sorta safe" cause we will get an STATUS_INVALID From the forkchoice update
	// // if it is deemed ot break the rules of the execution layer.
	// // still needs to be addressed of course.
	// fc := &pb.ForkchoiceState{
	// 	HeadBlockHash:      builtPayload.BlockHash(),
	// 	SafeBlockHash:      builtPayload.BlockHash(),
	// 	FinalizedBlockHash: builtPayload.BlockHash(),
	// }

	// attrs, err := m.getPayloadAttributes(ctx)
	// if err != nil {
	// 	return err
	// }

	// // Update forkchoice and start building the next block.
	// m.cachedPayload, _, err = m.EngineCaller.ForkchoiceUpdated(ctx, fc, attrs)
	// if errors.Is(err, prysmexecution.ErrAcceptedSyncingPayloadStatus) {
	// 	// retry as we are waiting for the import?
	// }
	// if err != nil {
	// 	m.logger.Error("failed to get forkchoice updated", "err", err)
	// 	return err
	// }

	// m.logger.Info("successfully validated execution layer block", "hash", common.Bytes2Hex(builtPayload.BlockHash()))
	return nil
}

func (m *Service) getPayloadAttributes(ctx sdk.Context) (payloadattribute.Attributer, error) {
	// TODO: modularize andn make better.
	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		return nil, err
	}

	return payloadattribute.New(&pb.PayloadAttributesV2{
		Timestamp:             uint64(time.Now().Unix()),
		SuggestedFeeRecipient: m.etherbase.Bytes(),
		Withdrawals:           nil,
		PrevRandao:            append([]byte{}, random[:]...),
	})
}

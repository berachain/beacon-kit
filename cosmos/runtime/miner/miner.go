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

// Package miner implements the Ethereum miner.
package miner

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	pb "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/execution"
	BeaconKeeper "github.com/itsdevbear/bolaris/cosmos/x/beacon/keeper"
)

// emptyHash is a common.Hash initialized to all zeros.
// var emptyHash = common.Hash{}

// EnvelopeSerializer is used to convert an envelope into a byte slice that represents
// a cosmos sdk.Tx.
type EnvelopeSerializer interface {
	ToSdkTxBytes(interfaces.ExecutionData, uint64) ([]byte, error)
}

// Miner implements the baseapp.TxSelector interface.
type Miner struct {
	execution.EngineCaller
	ek                 *BeaconKeeper.Keeper
	serializer         EnvelopeSerializer
	etherbase          common.Address
	curForkchoiceState *pb.ForkchoiceState
	// lastBlockTime      uint64
	logger log.Logger
}

// New produces a cosmos miner from a geth miner.
func New(gm execution.EngineCaller, ek *BeaconKeeper.Keeper, logger log.Logger) *Miner {
	return &Miner{
		EngineCaller:       gm,
		curForkchoiceState: &pb.ForkchoiceState{},
		logger:             logger,
		ek:                 ek,
	}
}

// Init sets the transaction serializer.
func (m *Miner) Init(serializer EnvelopeSerializer) {
	m.serializer = serializer
}

func (m *Miner) getForkchoiceFromExecutionClient(ctx context.Context) error {
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

func (m *Miner) SyncEl(ctx context.Context) error {
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
		fmt.Println(common.Bytes2Hex(m.curForkchoiceState.HeadBlockHash))
		fmt.Println(common.Bytes2Hex(m.curForkchoiceState.SafeBlockHash))
		fmt.Println(common.Bytes2Hex(m.curForkchoiceState.FinalizedBlockHash))
		// fmt.Println(common.Bytes2Hex(genesisHash.Bytes()))
		// fc := &pb.ForkchoiceState{
		// 	HeadBlockHash:      genesisHash.Bytes(),
		// 	SafeBlockHash:      genesisHash.Bytes(),
		// 	FinalizedBlockHash: genesisHash.Bytes(),
		// }
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

func (m *Miner) BuildBlockV2(ctx sdk.Context) (interfaces.ExecutionData, error) {
	// Reference: https://hackmd.io/@danielrachi/engine_api#Block-Building
	// builder := (&execution.Builder{EngineCaller: m.EngineCaller.(*execution.Service)})
	m.logger.Info("entering build-block-v2")
	defer m.logger.Info("exiting build-block-v2")

	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		return nil, err
	}

	attrs, err := payloadattribute.New(&pb.PayloadAttributesV2{
		Timestamp:             uint64(time.Now().Unix()),
		SuggestedFeeRecipient: m.etherbase.Bytes(),
		Withdrawals:           nil,
		PrevRandao:            append([]byte{}, random[:]...),
	})
	if err != nil {
		fmt.Println("attribute erorr")
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

	// Trigger the execution client to begin building the block, and update
	// the proposers forkchoice state accordingly. By setting the HeadBlockHash
	// to the last finalized block that we have seen, we are telling the execution client
	// to begin building a block on top of that block.
	payloadID, _, err := m.EngineCaller.ForkchoiceUpdated(ctx, fc, attrs)
	if err != nil {
		m.logger.Error("failed to get forkchoice updated", "err", err)
		return nil, err
	}

	// TODO: this should be something that is 80% of proposal timeout, or so?
	// TODO: maybe this should be some sort of event that we wait for?
	// But the TLDR is that we need to wait for the execution client to
	// build the payload before we can include it in the beacon block.
	time.Sleep(3000 * time.Millisecond) //nolint:gomnd // temp.

	// Get the Payload From the Execution Client
	builtPayload, _, _, err := m.EngineCaller.GetPayload(
		ctx, *payloadID, primitives.Slot(ctx.BlockHeight()))
	if err != nil {
		m.logger.Error("failed to get payload", "err", err)
		return nil, err
	}

	// This CL node's execution client head is now at the head of this
	// payload, so we can update our forkchoice state accordingly.
	m.curForkchoiceState.HeadBlockHash = builtPayload.BlockHash()

	// Go and include the builtPayload on the BeaconBlock
	return builtPayload, nil
}

func (m *Miner) ValidateBlock(ctx sdk.Context, builtPayload interfaces.ExecutionData) error {
	lastValidHash, _ := m.NewPayload(ctx, builtPayload, nil, nil) // last param here is nil pre-Deneb. must be specified post.
	fmt.Println("LAST VALID HASH FOUND ON ETH ONE", common.Bytes2Hex(lastValidHash))

	// TODO FIX, rn we are just blindly finalizing whatever the proposer has sent us.
	m.curForkchoiceState.HeadBlockHash = builtPayload.BlockHash()
	m.curForkchoiceState.FinalizedBlockHash = builtPayload.BlockHash()
	m.curForkchoiceState.SafeBlockHash = builtPayload.BlockHash()

	// The blind finalization is "sorta safe" cause we will get an STATUS_INVALID From the forkchoice update
	// if it is deemed ot break the rules of the execution layer.
	// still needs to be addressed of course.
	_, _, err := m.EngineCaller.ForkchoiceUpdated(ctx, m.curForkchoiceState, payloadattribute.EmptyWithVersion(3))
	if err != nil {
		m.logger.Error("failed to get forkchoice updated", "err", err)
		return err
	}

	m.logger.Info("successfully validated execution layer block", "hash", common.Bytes2Hex(builtPayload.BlockHash()))
	return nil
}

// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

// Package miner implements the Ethereum miner.
package miner

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	pb "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/eth"
	"github.com/itsdevbear/bolaris/beacon/prysm"
	evmkeeper "github.com/itsdevbear/bolaris/cosmos/x/evm/keeper"
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
	eth.EngineAPI
	ek                 *evmkeeper.Keeper
	serializer         EnvelopeSerializer
	etherbase          common.Address
	curForkchoiceState *pb.ForkchoiceState
	// lastBlockTime      uint64
	logger log.Logger
}

// New produces a cosmos miner from a geth miner.
func New(gm eth.EngineAPI, ek *evmkeeper.Keeper, logger log.Logger) *Miner {
	return &Miner{
		EngineAPI:          gm,
		curForkchoiceState: &pb.ForkchoiceState{},
		logger:             logger,
		ek:                 ek,
	}
}

// Init sets the transaction serializer.
func (m *Miner) Init(serializer EnvelopeSerializer) {
	m.serializer = serializer
}

func (m *Miner) BuildVoteExtension(ctx sdk.Context, _ int64) ([]byte, error) {
	var data interfaces.ExecutionData
	var err error
	if data, err = m.BuildBlock(ctx); err != nil {
		return nil, err
	}

	return data.MarshalSSZ()
}

// finalizedBlockHash returns the block hash of the finalized block corresponding to the given
// number or nil if doesn't exist in the chain.
func (m *Miner) finalizedBlockHash(number uint64) *common.Hash {
	var finalizedNumber = number
	// The code below is basically faking only updating the finalized block once per epoch.
	// if number%devEpochLength == 0 {
	// } else {
	// 	finalizedNumber = (number - 1) / devEpochLength * devEpochLength
	// }

	if finalizedBlock, err := m.EngineAPI.HeaderByNumber(context.Background(),
		big.NewInt(int64(finalizedNumber))); finalizedBlock != nil && err == nil {
		fh := finalizedBlock.Hash
		return &fh
	}
	return nil
}

func (m *Miner) getForkchoiceFromExecutionClient(ctx context.Context) error {
	var latestBlock *pb.ExecutionBlock
	var err error
	latestBlock, err = m.EngineAPI.LatestExecutionBlock(ctx)
	if err != nil {
		m.logger.Error("failed to get block number", "err", err)
	}

	m.curForkchoiceState.HeadBlockHash = latestBlock.Hash.Bytes()
	m.logger.Info("forkchoice state", "head", latestBlock.Header.Hash())

	safe, err := m.EngineAPI.LatestSafeBlock(ctx)
	if err != nil {
		m.logger.Error("failed to get safe block", "err", err)
		safe = latestBlock
	}

	m.curForkchoiceState.SafeBlockHash = safe.Hash.Bytes()

	final, err := m.EngineAPI.LatestFinalizedBlock(ctx)
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
	// block, _ := m.EngineAPI.EarliestBlock(ctx)

	m.getForkchoiceFromExecutionClient(ctx)
	genesisHash := m.ek.RetrieveGenesis(ctx)
	m.logger.Info("waiting for execution layer to sync")
	for {
		var err error
		fmt.Println(common.Bytes2Hex(m.curForkchoiceState.HeadBlockHash))
		fmt.Println(common.Bytes2Hex(m.curForkchoiceState.SafeBlockHash))
		fmt.Println(common.Bytes2Hex(m.curForkchoiceState.FinalizedBlockHash))
		fmt.Println(common.Bytes2Hex(genesisHash.Bytes()))
		// fc := &pb.ForkchoiceState{
		// 	HeadBlockHash:      genesisHash.Bytes(),
		// 	SafeBlockHash:      genesisHash.Bytes(),
		// 	FinalizedBlockHash: genesisHash.Bytes(),
		// }
		fc := m.curForkchoiceState
		payloadID, lastValidHash, err := m.EngineAPI.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
		if err == nil {
			break
		}
		m.curForkchoiceState = fc
		m.logger.Info("waiting for execution layer to sync", "error", err)
		m.logger.Info("waiting for execution layer to sync", "payloadID", payloadID, "lastValidHash", lastValidHash)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (m *Miner) BuildBlockV2(ctx sdk.Context) (interfaces.ExecutionData, error) {
	// Reference: https://hackmd.io/@danielrachi/engine_api#Block-Building
	builder := (&prysm.Builder{EngineCaller: m.EngineAPI.(*prysm.Service)})
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

	// Trigger the execution client to begin building the block, and update
	// the proposers forkchoice state accordingly.

	payloadID, _, err := m.EngineAPI.ForkchoiceUpdated(ctx, m.curForkchoiceState, attrs)
	if err != nil {
		m.logger.Error("failed to get forkchoice updated", "err", err)
		return nil, err
	}

	// TODO: this should be something that is 80% of proposal timeout, or so?
	// TODO: maybe this should be some sort of event that we wait for?
	// But the TLDR is that we need to wait for the execution client to
	// build the payload before we can include it in the beacon block.
	time.Sleep(2000 * time.Millisecond) //nolint:gomnd // temp.

	// Get the Payload From the Execution Client
	builtPayload, _, _, err := builder.GetPayload(
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
	builder := (&prysm.Builder{EngineCaller: m.EngineAPI.(*prysm.Service)})
	// todo parentBlockRoot should maybe be beacOn?
	// _, err := builder.NewPayload(ctx, builtPayload, nil, (*common.Hash)(builtPayload.ParentHash()))
	_, err := builder.NewPayload(ctx, builtPayload, nil, (*common.Hash)(m.curForkchoiceState.FinalizedBlockHash))
	m.curForkchoiceState.HeadBlockHash = builtPayload.BlockHash()
	if err != nil {
		// We need to resync.
		fmt.Println("ERROR", err, "RESYNCING")
		_, _, err := m.EngineAPI.ForkchoiceUpdated(ctx, m.curForkchoiceState, payloadattribute.EmptyWithVersion(3))
		if err != nil {
			m.logger.Error("failed to get forkchoice updated", "err", err)
			return err
		}

		return err
	}
	m.curForkchoiceState.FinalizedBlockHash = builtPayload.BlockHash()
	m.curForkchoiceState.SafeBlockHash = builtPayload.BlockHash()
	m.logger.Info("successfully validated block", "block", builtPayload.BlockHash())
	// m.getForkchoiceFromExecutionClient(ctx)
	return nil
}

// buildBlock builds and submits a payload, it also waits for the txs
// to resolve from the underying worker.
func (m *Miner) BuildBlock(ctx sdk.Context) (interfaces.ExecutionData, error) {
	builder := (&prysm.Builder{EngineCaller: m.EngineAPI.(*prysm.Service)})
	var (
		err error
		// envelope *engine.ExecutionPayloadEnvelope
		// sCtx = sdk.UnwrapSDKContext(ctx)
	)

	var payloadID *pb.PayloadIDBytes
	// Reset to CurrentBlock in case of the chain was rewound
	// ALL THIS CODE DOES IS FORCES RESETTING TO THE LATEST EXECUTION BLOCK
	// CALLS JSON RPC with "latest" block param
	{
		if err := m.getForkchoiceFromExecutionClient(ctx); err != nil {
			return nil, err
		}

		// tstamp := sCtx.BlockTime()
		var random [32]byte
		if _, err = rand.Read(random[:]); err != nil {
			return nil, err
		}
		var attrs payloadattribute.Attributer
		attrs, err = payloadattribute.New(&pb.PayloadAttributesV2{
			Timestamp:             uint64(time.Now().Unix()),
			SuggestedFeeRecipient: m.etherbase.Bytes(),
			Withdrawals:           nil,
			PrevRandao:            append([]byte{}, random[:]...),
		})
		if err != nil {
			fmt.Println("attribute erorr")
			return nil, err
		}
		// the proposer updates it's forkchoice and then
		// I think this needs to actually happen in FinalizeBlock and everyone should
		// update forkchoice together?
		payloadID, _, err = m.EngineAPI.ForkchoiceUpdated(ctx,
			m.curForkchoiceState, attrs)
		if err != nil {
			m.logger.Error("failed to get forkchoice updated", "err", err)
		}
		// TODO: this should be something that is 80% of proposal timeout, or so?
		time.Sleep(2500 * time.Millisecond) //nolint:gomnd // temp.
	}

	// builds a payload
	builtPayload, _, _, err := builder.GetPayload(
		ctx, *payloadID, primitives.Slot(ctx.BlockHeight()),
	)
	if err != nil {
		fmt.Println("PAYLOAD ERROR", err)
		return nil, err
	}

	hash, err := builder.NewPayload(ctx, builtPayload, nil, (*common.Hash)(m.curForkchoiceState.FinalizedBlockHash))
	if err != nil {
		fmt.Println("NEW PAYLOAD ERROR", err)
	}
	fmt.Println("PAYLOAD HASH", hash, "cosmos block height", ctx.BlockHeight(), "PayloadID", payloadID,
		"payload block number", builtPayload.BlockNumber())
	// finalizedHash := m.finalizedBlockHash(builtPayload.BlockNumber())
	// if finalizedHash == nil {
	// 	fmt.Println("EMPTY FINALIZED HASH")
	// 	finalizedHash = new(common.Hash)
	// 	x := common.BytesToHash(builtPayload.BlockHash())
	// 	finalizedHash = &x
	// 	m.curForkchoiceState.HeadBlockHash = finalizedHash.Bytes()
	// }
	m.curForkchoiceState.HeadBlockHash = builtPayload.BlockHash()
	// finalizedHash, when there is epochs, could be in the past. But since
	// we are finalizing every block, the builtPayload and the finalized Hash are the same.

	// // finalizedHash := builtPayload.BlockHash()
	// fmt.Println("finalizedHash", common.Bytes2Hex(finalizedHash))
	// fmt.Println("payload()", common.Bytes2Hex(builtPayload.BlockHash()))
	// m.setCurrentState(builtPayload.BlockHash(), finalizedHash)

	if err != nil {
		fmt.Println("ERROR", err)
		return nil, err
	}

	fmt.Println(
		"finalized", common.Bytes2Hex(m.curForkchoiceState.FinalizedBlockHash),
		"head", common.Bytes2Hex(m.curForkchoiceState.HeadBlockHash),
		"safe", common.Bytes2Hex(m.curForkchoiceState.SafeBlockHash),
	)
	// time.Sleep(*time.Second)
	// _, _, err = b.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3)) //nolint:gomnd // okay for now.
	_, _, err = builder.EngineCaller.ForkchoiceUpdated(ctx, m.curForkchoiceState, payloadattribute.EmptyWithVersion(3))
	fmt.Println("UPDATE ERROR", err)
	// _, err = builder.BlockValidation(ctx, builtPayload)
	return builtPayload, err
}

// // setCurrentState sets the current forkchoice state.
// func (m *Miner) setCurrentState(headHash, finalizedHash []byte) {
// 	m.curForkchoiceState = &pb.ForkchoiceState{
// 		HeadBlockHash:      headHash,
// 		SafeBlockHash:      finalizedHash,
// 		FinalizedBlockHash: finalizedHash,
// 	}
// }

// // constructPayloadArgs builds a payload to submit to the miner.
// func (m *Miner) constructPayloadArgs(
// 	ctx sdk.Context, parent *types.Block) *miner.BuildPayloadArgs {
// 	// etherbase, err := m.Etherbase(ctx)
// 	// if err != nil {
// 	// 	ctx.Logger().Error("failed to get etherbase", "err", err)
// 	// 	return nil
// 	// }

// 	return &miner.BuildPayloadArgs{
// 		Timestamp: parent.Header().Time + 2, //nolint:gomnd // todo fix this arbitrary number.
// 		// FeeRecipient: etherbase,
// 		Random:      common.Hash{}, /* todo: generated random */
// 		Withdrawals: make(types.Withdrawals, 0),
// 		BeaconRoot:  &emptyHash,
// 		Parent:      parent.Hash(),
// 	}
// }

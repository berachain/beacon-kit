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
	"bytes"
	"context"
	"crypto/rand"
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
	serializer         EnvelopeSerializer
	etherbase          common.Address
	curForkchoiceState *pb.ForkchoiceState
	// lastBlockTime      uint64
	logger log.Logger
}

// New produces a cosmos miner from a geth miner.
func New(gm eth.EngineAPI, logger log.Logger) *Miner {
	return &Miner{
		EngineAPI:          gm,
		curForkchoiceState: &pb.ForkchoiceState{},
		logger:             logger,
	}
}

// Init sets the transaction serializer.
func (m *Miner) Init(serializer EnvelopeSerializer) {
	m.serializer = serializer
}

func (m *Miner) BuildVoteExtension(ctx sdk.Context, _ int64) ([]byte, error) {
	var data interfaces.ExecutionData
	var err error
	if data, err = m.buildBlock(ctx); err != nil {
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

// buildBlock builds and submits a payload, it also waits for the txs
// to resolve from the underying worker.
func (m *Miner) buildBlock(ctx sdk.Context) (interfaces.ExecutionData, error) {
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
		var latestBlock *pb.ExecutionBlock
		latestBlock, err = m.EngineAPI.LatestExecutionBlock(ctx)
		if err != nil {
			m.logger.Error("failed to get block number", "err", err)
		}

		if !bytes.Equal(m.curForkchoiceState.HeadBlockHash, latestBlock.Hash.Bytes()) {
			finalizedHash := m.finalizedBlockHash(latestBlock.Number.Uint64())
			m.setCurrentState(latestBlock.Hash.Bytes(), finalizedHash.Bytes())
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
			return nil, err
		}

		payloadID, _, err = m.EngineAPI.ForkchoiceUpdated(ctx,
			m.curForkchoiceState, attrs)
		if err != nil {
			m.logger.Error("failed to get forkchoice updated", "err", err)
		}
		// TODO: this should be something that is 80% of proposal timeout, or so?
		time.Sleep(2500 * time.Millisecond) //nolint:gomnd // temp.
	}

	builtPayload, _, _, err := builder.GetPayload(
		ctx, *payloadID, primitives.Slot(ctx.BlockHeight()),
	)
	if err != nil {
		return nil, err
	}

	finalizedHash := builtPayload.BlockHash()
	// finalizedHash, when there is epochs, could be in the past. But since
	// we are finalizing every block, the builtPayload and the finalized Hash are the same.
	m.setCurrentState(builtPayload.BlockHash(), finalizedHash)

	// _, _, err = builder.BlockValidation(ctx, builtPayload)
	return builtPayload, err
}

// setCurrentState sets the current forkchoice state.
func (m *Miner) setCurrentState(headHash, finalizedHash []byte) {
	m.curForkchoiceState = &pb.ForkchoiceState{
		HeadBlockHash:      headHash,
		SafeBlockHash:      headHash,
		FinalizedBlockHash: finalizedHash,
	}
}

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

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

package prysm

import (
	"context"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

type EVMKeeper interface {
	LatestForkChoice(ctx context.Context) (*enginev1.ForkchoiceState, error)
	SetLatestForkChoice(ctx context.Context, fc *enginev1.ForkchoiceState) error
}

type Builder struct {
	Keeper EVMKeeper
	EngineCaller
}

func (b *Builder) BlockProposal(
	_ context.Context, _ interfaces.ExecutionData, _ payloadattribute.Attributer,
) (interfaces.ExecutionData, error) {
	// payloadID, latestValidHash, err := b.BlockValidation(ctx, payload)
	// if err != nil e
	// 	return nil, err
	// }
	// if latestValidHash == nil {
	// 	return nil, err
	// }

	// builtPayload, _, _, err := b.GetPayload(ctx, *payloadID, primitives.Slot(0))

	return nil, nil
}

// BlockValidation builds a payload from the provided execution data, and submits it to
// the execution client. It then submits a forkchoice update with the updated
// valid hash returned by the execution client.1
// This should be called by a node when it receives Execution Data as part of
// beacon chain consensus.
// receives payload -> get latestValidHash from our execution client -> forkchoice locally.
func (b *Builder) BlockValidation(
	ctx context.Context, payload interfaces.ExecutionData,
) ([]byte, error) {
	// pull the lastest valid hash from the execution node.
	// this payload should have been p2p'd to this node, callinpro NewPayload will execute it
	// adding it to the local canonical chain. the Forkchoice update can then check to see if
	// everything is gucci gang.
	// This should probably be ran in prepare proposal in order to reject the proposal if
	// // the node receives a payload that produces a bad block and/or errors.
	// latestValidHash, err := b.NewPayload(ctx, payload, nil, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println("BLOCK VALIDATION")
	// fmt.Println("latestValidHash: ", common.BytesToHash(latestValidHash).Hex())
	// // fc, err := b.Keeper.LatestForkChoice(ctx)
	// if fc == nil || err != nil {
	// 	return nil, err
	// }

	fc := &enginev1.ForkchoiceState{
		HeadBlockHash:      payload.BlockHash(),
		SafeBlockHash:      payload.BlockHash(),
		FinalizedBlockHash: payload.BlockHash(),
	}

	// if fc == nil {
	// 	fc = &enginev1.ForkchoiceState{
	// 		HeadBlockHash:      latestValidHash,
	// 		SafeBlockHash:      latestValidHash,
	// 		FinalizedBlockHash: latestValidHash,
	// 	}
	// } else {
	// 	fc = &enginev1.ForkchoiceState{
	// 		HeadBlockHash: latestValidHash,
	// 		// The two below are technically incorrect? These should be set later imo.
	// 		// Should we set head block hash in prepare proposal?
	// 		// Then update safe and finalized in end block / beginning of the following block?
	// 		SafeBlockHash:      fc.SafeBlockHash,
	// 		FinalizedBlockHash: fc.FinalizedBlockHash,
	// 	}
	// }

	// b.Keeper.SetLatestForkChoice(ctx, fc)

	_, _, err := b.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3)) //nolint:gomnd // okay for now.

	return payload.BlockHash(), err
}

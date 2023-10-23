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

package ve

import (
	"context"
	"fmt"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/log"

	cometabci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/prysm"
	"github.com/itsdevbear/bolaris/cosmos/runtime/miner"
)

type Processor struct {
	logger log.Logger

	miner *miner.Miner
}

func NewProcessor(miner *miner.Miner, logger log.Logger) *Processor {
	return &Processor{
		logger: logger,
		miner:  miner,
	}
}

func (h *Processor) ProcessCommitInfo(ctx context.Context, height int64,
	extendedCommitInfo cometabci.ExtendedCommitInfo) error {
	for _, vote := range extendedCommitInfo.Votes {
		address := sdk.ConsAddress{}
		if err := address.Unmarshal(vote.Validator.Address); err != nil {
			h.logger.Error(
				"failed to unmarshal validator address",
				"height", height,
			)

			return err
		}

		if err := h.ValidateOracleVoteExtension(ctx, vote.VoteExtension,
			uint64(height)); err != nil {
			// h.logger.Error(
			// 	"failed to validate oracle vote extension",
			// 	"height", height,
			// 	"validator", address.String(),
			// 	"err", err,
			// )

			return err
		}
	}
	return nil
}

func (h *Processor) ValidateOracleVoteExtension(
	ctx context.Context, voteExtension []byte, height uint64) error {
	builder := (&prysm.Builder{EngineCaller: h.miner})

	payload := new(enginev1.ExecutionPayloadCapellaWithValue)
	payload.Payload = new(enginev1.ExecutionPayloadCapella)
	if err := payload.Payload.UnmarshalSSZ(voteExtension); err != nil {
		h.logger.Error(
			"failed to unmarshal vote extension",
			"height", height,
			"err", err,
		)

		return err
	}
	// todo handle hardforks without needing codechange.
	data, err := blocks.WrappedExecutionPayloadCapella(
		payload.Payload, blocks.PayloadValueToGwei(payload.Value),
	)
	if err != nil {
		h.logger.Error(
			"failed to wrap vote extension",
			"height", height,
			"err", err,
		)
		return err
	}

	// TODO: switch evmtypes.WrappedPayloadEnvelope to just use the proto type
	// from prysm.
	_, err = builder.BlockValidation(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to validate payload: %w", err)
	}

	return nil
}

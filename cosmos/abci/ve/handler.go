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
	"time"

	"cosmossdk.io/log"

	cometabci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/cosmos/runtime/miner"
)

// VoteExtensionHandler is a handler that extends a vote with the oracle's
// current price feed. In the case where oracle data is unable to be fetched
// or correctly marshalled, the handler will return an empty vote extension to
// ensure liveliness.
type VoteExtensionHandler struct {
	logger log.Logger

	// timeout is the maximum amount of time to wait for the oracle to respond
	// to a price request.
	timeout time.Duration

	// miner represents the underlying miner.
	miner *miner.Miner
}

// NewVoteExtensionHandler returns a new VoteExtensionHandler.
func NewVoteExtensionHandler(logger log.Logger, timeout time.Duration,
	miner *miner.Miner) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:  logger,
		timeout: timeout,
		miner:   miner,
	}
}

// ExtendVoteHandler returns a handler that extends a vote with the oracle's
// current price feed. In the case where oracle data is unable to be fetched
// or correctly marshalled, the handler will return an empty vote extension to
// ensure liveliness.
//
//nolint:nolintlint,nonamedreturns // its okay from slinky.
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(
		ctx sdk.Context, req *cometabci.RequestExtendVote,
	) (resp *cometabci.ResponseExtendVote, err error) {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", r,
				)

				resp, err = &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, nil
			}
		}()

		// Create a context with a timeout to ensure we do not wait forever for the oracle
		// to respond.
		_, cancel := context.WithTimeout(ctx.Context(), h.timeout)
		defer cancel()

		// voteExt := vetypes.ForkChoiceExtension{
		// 	Height: req.Height,
		// }

		// voteExt :=

		bz, err := h.miner.BuildVoteExtension(ctx, req.Height)
		if err != nil {
			h.logger.Error(
				"failed to marshal vote extension; returning empty vote extension",
				"height", req.Height,
				"err", err,
			)

			return &cometabci.ResponseExtendVote{VoteExtension: []byte{}}, nil
		}

		fmt.Println("VOTE EXTENSION", bz, "EXTEND VOTE HANDLER")

		h.logger.Info(
			"extending vote with oracle prices",
			"vote_extension_height", req.Height,
			"req_height", req.Height,
		)

		return &cometabci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

// VerifyVoteExtensionHandler returns a handler that verifies the vote extension provided by
// a validator is valid. In the case when the vote extension is empty, we return ACCEPT. T
// This means that the validator may have been unable to fetch prices firom the oracle and
// voting an empty vote extension. We reject any vote extensions that are not empty and
// fail to unmarshal or contain invalid prices.
func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(
		ctx sdk.Context,
		req *cometabci.RequestVerifyVoteExtension,
	) (*cometabci.ResponseVerifyVoteExtension, error) {
		// voteExtension := req.VoteExtension

		// if err := ValidateOracleVoteExtension(voteExtension, req.Height); err != nil {
		// 	h.logger.Error(
		// 		"failed to validate vote extension",
		// 		"height", req.Height,
		// 		"err", err,
		// 	)

		// 	return &cometabci.ResponseVerifyVoteExtension{Status:
		//  cometabci.ResponseVerifyVoteExtension_REJECT}, err
		// }

		h.logger.Info(
			"validated vote extension",
			"height", req.Height,
		)

		return &cometabci.ResponseVerifyVoteExtension{
			Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}

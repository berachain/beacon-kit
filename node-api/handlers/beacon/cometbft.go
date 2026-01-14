// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package beacon

import (
	"fmt"

	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
)

// GetCometBFTBlock returns the CometBFT block at the specified height.
// GET /eth/v1/beacon/cometbft/block/:height
func (h *Handler) GetCometBFTBlock(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.CometBFTHeightRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	height := req.Height
	block := h.backend.GetCometBFTBlock(height)
	if block == nil {
		return nil, fmt.Errorf("block not found at height %d", height)
	}

	// Return CometBFT block directly - it has native JSON serialization
	return beacontypes.NewResponse(block), nil
}

// GetCometBFTCommit returns the CometBFT commit at the specified height.
// GET /eth/v1/beacon/cometbft/commit/:height
func (h *Handler) GetCometBFTCommit(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.CometBFTHeightRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	height := req.Height
	commit := h.backend.GetCometBFTCommit(height)
	if commit == nil {
		return nil, fmt.Errorf("commit not found at height %d", height)
	}

	// Return CometBFT commit directly - it has native JSON serialization
	return beacontypes.NewResponse(commit), nil
}

// GetCometBFTValidators returns the CometBFT validator set at the specified height.
// GET /eth/v1/beacon/cometbft/validators/:height
func (h *Handler) GetCometBFTValidators(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.CometBFTHeightRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	height := req.Height
	validators, err := h.backend.GetCometBFTValidators(height)
	if err != nil {
		return nil, fmt.Errorf("failed to get validators at height %d: %w", height, err)
	}
	if validators == nil {
		return nil, fmt.Errorf("validators not found at height %d", height)
	}

	// Return CometBFT ValidatorSet directly - it has native JSON serialization
	return beacontypes.NewResponse(validators), nil
}
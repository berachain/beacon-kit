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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package cometbft

import (
	errorsmod "cosmossdk.io/errors"
	servercmtlog "github.com/berachain/beacon-kit/consensus/cometbft/service/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CreateQueryContext creates a new sdk.Context for a query, taking as args
// the block height and whether the query needs a proof or not.
func (s *Service) CreateQueryContext(
	height int64,
	prove bool,
) (sdk.Context, error) {
	// use custom query multi-store if provided
	lastBlockHeight := s.sm.GetCommitMultiStore().LatestVersion()
	if lastBlockHeight == 0 {
		return sdk.Context{}, errorsmod.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"%s is not ready; please wait for first block",
			AppName,
		)
	}

	if height > lastBlockHeight {
		return sdk.Context{},
			errorsmod.Wrap(
				sdkerrors.ErrInvalidHeight,
				"cannot query with height in the future; please provide a valid height",
			)
	}

	// when a client did not provide a query height, manually inject the latest
	if height == 0 {
		height = lastBlockHeight
	}

	if height <= 1 && prove {
		return sdk.Context{},
			errorsmod.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			)
	}

	cacheMS, err := s.sm.GetCommitMultiStore().CacheMultiStoreWithVersion(height)
	if err != nil {
		return sdk.Context{},
			errorsmod.Wrapf(
				sdkerrors.ErrNotFound,
				"failed to load state at height %d; %s (latest height: %d)",
				height,
				err,
				lastBlockHeight,
			)
	}

	return sdk.NewContext(
		cacheMS,
		true,
		servercmtlog.WrapSDKLogger(s.logger),
	), nil
}

// LastBlockHeight returns the last committed block height.
func (s *Service) LastBlockHeight() int64 {
	return s.sm.GetCommitMultiStore().LastCommitID().Version
}

func (s *Service) GetSyncData() (int64, int64) {
	latestHeight := s.sm.GetCommitMultiStore().LastCommitID().Version
	return latestHeight, s.syncToHeight
}

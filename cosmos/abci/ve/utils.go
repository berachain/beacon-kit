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
	cometabci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// // ValidateOracleVoteExtension validates the vote extension provided by a validator.
// func ValidateOracleVoteExtension(voteExtension []byte, _ int64) error {
// 	if len(voteExtension) == 0 {
// 		return nil
// 	}

// 	return nil
// }

// VoteExtensionsEnabled determines if vote extensions are enabled for the current block.
func VoteExtensionsEnabled(ctx sdk.Context) bool {
	cp := ctx.ConsensusParams()
	if cp.Abci == nil || cp.Abci.VoteExtensionsEnableHeight == 0 {
		return false
	}

	// Per the cosmos sdk, the first block should not utilize the latest finalize block state.
	// This means vote extensions should NOT be making state changes.
	if ctx.BlockHeight() <= 1 {
		return false
	}

	// We do a +1 here because the vote extensions are enabled at height h
	// but a proposer will only receive vote extensions in height h+1.
	return cp.Abci.VoteExtensionsEnableHeight+1 < ctx.BlockHeight()
}

type (
	// ValidateVoteExtensionsFn defines the function for validating vote extensions. This
	// function is not explicitly used to validate the oracle data but rather that
	// the signed vote extensions included in the proposal are valid and provide
	// a supermajority of vote extensions for the current block. This method is
	// expected to be used in PrepareProposal and ProcessProposal.
	ValidateVoteExtensionsFn func(
		ctx sdk.Context,
		height int64,
		extInfo cometabci.ExtendedCommitInfo,
	) error
)

// NewDefaultValidateVoteExtensionsFn returns a new DefaultValidateVoteExtensionsFn.
func NewDefaultValidateVoteExtensionsFn(
	chainID string, validatorStore baseapp.ValidatorStore,
) ValidateVoteExtensionsFn {
	return func(ctx sdk.Context, height int64, info cometabci.ExtendedCommitInfo) error {
		return baseapp.ValidateVoteExtensions(ctx, validatorStore, height, chainID, info)
	}
}

// NoOpValidateVoteExtensions is a no-op validation method (purely used for testing).
func NoOpValidateVoteExtensions(
	_ sdk.Context,
	_ int64,
	_ cometabci.ExtendedCommitInfo,
) error {
	return nil
}

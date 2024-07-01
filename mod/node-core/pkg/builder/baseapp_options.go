// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package builder

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This file contains Options that extend our default baseapp options to be
// called by cosmos when building the app.
// TODO: refactor into consensus_options for serverv2 migration.

// WithCometParamStore sets the param store to the comet consensus engine.
func WithCometParamStore(
	chainSpec common.ChainSpec,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetParamStore(comet.NewConsensusParamsStore(chainSpec))
	}
}

// WithPrepareProposal sets the prepare proposal handler to the baseapp.
func WithPrepareProposal(
	handler sdk.PrepareProposalHandler,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetPrepareProposal(handler)
	}
}

// WithProcessProposal sets the process proposal handler to the baseapp.
func WithProcessProposal(
	handler sdk.ProcessProposalHandler,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetProcessProposal(handler)
	}
}

// WithPreBlocker sets the pre-blocker to the baseapp.
func WithPreBlocker(
	preBlocker sdk.PreBlocker,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetPreBlocker(preBlocker)
	}
}

func WithStreamingManager(
	streamingManager storetypes.StreamingManager,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetStreamingManager(streamingManager)
	}
}

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

package runtime

import (
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "runtime"

// AppI implements the common methods for a Cosmos SDK-based application
// specific blockchain.
type AppI interface {
	// Name the assigned name of the app.
	Name() string

	// BeginBlocker updates every begin block.
	BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error)

	// EndBlocker updates every end block.
	EndBlocker(ctx sdk.Context) (sdk.EndBlock, error)

	// InitChainer update at chain (i.e app) initialization.
	InitChainer(
		ctx sdk.Context,
		req *abci.InitChainRequest,
	) (*abci.InitChainResponse, error)

	// LoadHeight load the app at a given height.
	LoadHeight(height int64) error

	// ExportAppStateAndValidators exports the state of the application for a
	// genesis file.
	ExportAppStateAndValidators(
		forZeroHeight bool,
		jailAllowedAddrs, modulesToExport []string,
	) (types.ExportedApp, error)
}

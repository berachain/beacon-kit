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

package components

import (
	"path/filepath"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/store/v2"
	"cosmossdk.io/store/v2/commitment/iavl"
	"cosmossdk.io/store/v2/db"
	"cosmossdk.io/store/v2/root"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

type StoreConfigInput struct {
	depinject.In

	Logger  log.Logger
	AppOpts servertypes.AppOptions
}

func ProvideStoreOptions(in StoreConfigInput) *root.FactoryOptions {
	homeDir := cast.ToString(in.AppOpts.Get(serverv2.FlagHome))
	scRawDb, err := db.NewDB(
		db.DBTypePebbleDB,
		"application",
		filepath.Join(homeDir, "data"),
		nil,
	)
	if err != nil {
		panic(err)
	}

	return &root.FactoryOptions{
		Logger:          in.Logger,
		RootDir:         homeDir,
		SSType:          root.SSTypeSQLite,
		SCType:          root.SCTypeIavl,
		SCPruningOption: store.NewPruningOption(store.PruningNothing),
		IavlConfig:      iavl.DefaultConfig(),
		SCRawDB:         scRawDb,
	}
}

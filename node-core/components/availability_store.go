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
	"os"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/config"
	dastore "github.com/berachain/beacon-kit/da/store"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/storage/filedb"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
)

// AvailabilityStoreInput is the input for the ProviderAvailabilityStore
// function for the depinject framework.
type AvailabilityStoreInput[LoggerT any] struct {
	depinject.In
	AppOpts   config.AppOptions
	ChainSpec chain.ChainSpec
	Logger    LoggerT
}

// ProvideAvailibilityStore provides the availability store.
func ProvideAvailibilityStore[
	LoggerT log.AdvancedLogger[LoggerT],
](
	in AvailabilityStoreInput[LoggerT],
) (*dastore.Store, error) {
	return dastore.New(
		filedb.NewRangeDB(
			filedb.NewDB(
				filedb.WithRootDirectory(
					cast.ToString(
						in.AppOpts.Get(flags.FlagHome),
					)+"/data/blobs",
				),
				filedb.WithFileExtension("ssz"),
				filedb.WithDirectoryPermissions(os.ModePerm),
				filedb.WithLogger(in.Logger),
			),
		),
		in.Logger.With("service", "da-store"),
		in.ChainSpec,
	), nil
}

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

package deposit

import (
	"fmt"

	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	clicontext "github.com/berachain/beacon-kit/cli/context"
	servercmtlog "github.com/berachain/beacon-kit/consensus/cometbft/service/log"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/storage/db"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetDBCheckCmd returns a command for checking that the deposit store
// is in sync with the beacon state.
func GetDBCheckCmd(appCreator servertypes.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-check",
		Short: `Checks if the deposit store is in sync with the beacon state. Fails if either of the beacon or deposit DBs are not available.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Create the application from home directory configs and data.
			v := clicontext.GetViperFromCmd(cmd)
			logger := clicontext.GetLoggerFromCmd(cmd)
			cfg := clicontext.GetConfigFromCmd(cmd)
			db, err := db.OpenDB(cfg.RootDir, dbm.PebbleDBBackend)
			if err != nil {
				return err
			}
			app := appCreator(logger, db, nil, cfg, v)

			// Setup the state to check.
			ctx := sdk.NewContext(
				app.CommitMultiStore().CacheMultiStore(), false, servercmtlog.WrapSDKLogger(logger),
			).WithContext(cmd.Context())
			beaconState := app.StorageBackend().StateFromContext(ctx)
			depositStore := app.StorageBackend().DepositStore()

			// First, check the deposit store contains the correct amount of deposits.
			eth1DepositIndex, err := beaconState.GetEth1DepositIndex()
			if err != nil {
				return err
			}
			deposits, err := depositStore.GetDepositsByIndex(
				ctx, constants.FirstDepositIndex, eth1DepositIndex,
			)
			if err != nil {
				return err
			}
			if uint64(len(deposits)) != eth1DepositIndex {
				return fmt.Errorf(
					"deposits in Deposit store (length: %d) are not in sync with the Beacon state (Eth1DepositIndex: %d)",
					len(deposits),
					eth1DepositIndex,
				)
			}

			// Second, verify that the deposit store contains the correct deposits.
			eth1Data, err := beaconState.GetEth1Data()
			if err != nil {
				return err
			}
			if deposits.HashTreeRoot() != eth1Data.DepositRoot {
				return fmt.Errorf(
					"deposits from Deposit store (HTR: %s) are not consistent with the Beacon state (HTR: %s)",
					deposits.HashTreeRoot(),
					eth1Data.DepositRoot,
				)
			}

			logger.Info("✅ Deposit store is in sync with the Beacon state!")
			return nil
		},
	}

	return cmd
}

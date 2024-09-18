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

package server

import (
	"context"
	"fmt"

	"cosmossdk.io/store"
	types "github.com/berachain/beacon-kit/mod/cli/pkg/commands/server/types"
	clicontext "github.com/berachain/beacon-kit/mod/cli/pkg/context"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/storage/pkg/db"
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/cobra"
)

// NewRollbackCmd creates a command to rollback CometBFT and multistore state by
// one height.
func NewRollbackCmd[
	T interface {
		Start(context.Context) error
		CommitMultiStore() store.CommitMultiStore
	},
	LoggerT log.AdvancedLogger[LoggerT],
](
	appCreator types.AppCreator[T, LoggerT],
) *cobra.Command {
	var removeBlock bool

	//nolint:lll // its okay.
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "rollback Cosmos SDK and CometBFT state by one height",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when CometBFT has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also rolls back to height n - 1. No blocks are removed, so upon
restarting CometBFT the transactions in block n will be re-executed against the
application.
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			v := clicontext.GetViperFromCmd(cmd)
			logger := clicontext.GetLoggerFromCmd[LoggerT](cmd)
			cfg := clicontext.GetConfigFromCmd(cmd)

			db, err := db.OpenDB(cfg.RootDir, dbm.PebbleDBBackend)
			if err != nil {
				return err
			}
			app := appCreator(logger, db, nil, cfg, v)
			// rollback CometBFT state
			height, hash, err := cmtcmd.RollbackState(cfg, removeBlock)
			if err != nil {
				return fmt.Errorf("failed to rollback CometBFT state: %w", err)
			}
			// rollback the multistore

			if err = app.CommitMultiStore().RollbackToVersion(height); err != nil {
				return fmt.Errorf("failed to rollback to version: %w", err)
			}

			logger.Info(
				"Rolled back state to height %d and hash %X\n",
				height,
				hash,
			)
			return nil
		},
	}

	cmd.Flags().
		BoolVar(&removeBlock, "hard", false, "remove last block as well as state")
	return cmd
}

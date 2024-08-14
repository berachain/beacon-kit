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

package genesis

import (
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// Commands builds the genesis-related command. Users may
// provide application specific commands as a parameter.
func Commands[
	LegacyKeyT ~[constants.BLSSecretKeyLength]byte,
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	cs common.ChainSpec,
	cmds ...*cobra.Command,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "genesis",
		Short:                      "Application's genesis-related subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2, //nolint:mnd // from sdk.
		RunE:                       client.ValidateCmd,
	}

	// Adding subcommands for genesis-related operations.
	cmd.AddCommand(
		AddGenesisDepositCmd[LegacyKeyT, LoggerT](cs),
		CollectGenesisDepositsCmd[LoggerT](),
		AddExecutionPayloadCmd[LoggerT](cs),
		GetGenesisValidatorRootCmd(cs),
	)

	// Add additional commands
	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}

	return cmd
}

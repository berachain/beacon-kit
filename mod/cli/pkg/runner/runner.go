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

package runner

import (
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// Runner is a type that runs the root command.
type Runner struct {
	// nodeHome is the home directory of the node.
	nodeHome string
	// rootCmd is the root command.
	rootCmd *commands.Root
}

// New returns a new Runner with the given root command.
func New[NodeT servertypes.Application](
	rootCmd *commands.Root,
	appCreator servertypes.AppCreator[NodeT],
) *Runner {
	commands.SetupRootCmdWithNode[NodeT](
		rootCmd,
		appCreator,
	)
	return &Runner{rootCmd: rootCmd}
}

// Run runs the root command.
func (runner *Runner) Run() error {
	return runner.rootCmd.Run(runner.nodeHome)
}

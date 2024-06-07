// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package main

import (
	"log/slog"
	"os"

	nodebuilder "github.com/berachain/beacon-kit/mod/node-core/pkg/builder"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config/spec"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"go.uber.org/automaxprocs/maxprocs"
)

// run runs the beacon node.
func run() error {
	// Set the uber max procs
	if _, err := maxprocs.Set(); err != nil {
		return err
	}

	// Build the node using the node-core.
	nb := nodebuilder.New[types.NodeI](
		nodebuilder.WithName[types.NodeI]("beacond"),
		nodebuilder.WithDescription[types.NodeI](
			"beacond is a beacon node for any beacon-kit chain",
		),
		nodebuilder.WithDepInjectConfig[types.NodeI](Config()),
		// TODO: Don't hardcode the default chain spec.
		nodebuilder.WithChainSpec[types.NodeI](spec.TestnetChainSpec()),
	)

	node, err := nb.Build()
	if err != nil {
		return err
	}
	return node.Run()
}

// main is the entry point.
func main() {
	if err := run(); err != nil {
		//nolint:sloglint // todo fix.
		slog.Error("startup failure", "error", err)
		os.Exit(1)
	}
}

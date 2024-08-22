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

package node

import (
	"context"
	"os"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server"
	service "github.com/berachain/beacon-kit/mod/node-core/pkg/services/registry"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

// Compile-time assertion that node implements the NodeI interface.
var _ types.Node = (*node)(nil)

// node is the hard-type representation of the beacon-kit node.
type node struct {
	types.Application
	// registry is the node's service registry.
	registry *service.Registry
}

// New returns a new node.
func New[NodeT types.Node](registry *service.Registry) NodeT {
	return types.Node(&node{registry: registry}).(NodeT)
}

// Start starts the node.
func (n *node) Start(
	ctx context.Context,
) error {
	g, ctx := server.GetCtx(ctx, log.NewLogger(os.Stdout), true)
	if err := n.registry.StartAll(ctx); err != nil {
		return err
	}
	return g.Wait()
}

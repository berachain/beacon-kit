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

package node

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/berachain/beacon-kit/log"
	service "github.com/berachain/beacon-kit/node-core/services/registry"
	"github.com/berachain/beacon-kit/node-core/types"
)

// Compile-time assertion that node implements the NodeI interface.
var _ types.Node = (*node)(nil)

// node is the hard-type representation of the beacon-kit node.
type node struct {
	// logger is the node's logger.
	logger log.Logger
	// registry is the node's service registry.
	registry *service.Registry

	// TODO: FIX, HACK TO MAKE CLI HAPPY FOR NOW.
	// THIS SHOULD BE REMOVED EVENTUALLY.
	types.Node
}

// New returns a new node.
func New[NodeT types.Node](
	registry *service.Registry, logger log.Logger) NodeT {
	//nolint:errcheck // should be safe
	return types.Node(&node{registry: registry, logger: logger}).(NodeT)
}

// Start starts the node.
func (n *node) Start(
	ctx context.Context,
) error {
	// Make the context cancellable.
	cctx, cancelFn := context.WithCancel(ctx)

	stop := make(chan struct{})

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		sig := <-sigc
		n.logger.Info("caught exit signal", "signal", sig.String())
		cancelFn()
		go func() {
			n.registry.StopAll()
			close(stop)
		}()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				n.logger.Info("Already shutting down, interrupt more to panic")
			}
		}
		panic("Panic closing the beacon node")
	}()

	n.registry.StartAll(cctx)

	// Wait for stop channel to be closed.
	<-stop

	return nil
}

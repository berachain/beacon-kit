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
	"os/signal"
	"syscall"

	"github.com/berachain/beacon-kit/log"
	service "github.com/berachain/beacon-kit/node-core/pkg/services/registry"
	"github.com/berachain/beacon-kit/node-core/pkg/types"
	"golang.org/x/sync/errgroup"
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

	// Create an errgroup to manage the lifecycle of all the services.
	g, gctx := errgroup.WithContext(cctx)

	// listen for quit signals so the calling parent process can gracefully exit
	n.listenForQuitSignals(g, true, cancelFn)

	// Start all the registered services.
	if err := n.registry.StartAll(gctx); err != nil {
		return err
	}

	// Wait for those aforementioned exit signals.
	return g.Wait()
}

// listenForQuitSignals listens for SIGINT and SIGTERM. When a signal is
// received,
// the cleanup function is called, indicating the caller can gracefully exit or
// return.
//
// Note, the blocking behavior of this depends on the block argument.
// The caller must ensure the corresponding context derived from the cancelFn is
// used correctly.
func (n *node) listenForQuitSignals(
	g *errgroup.Group,
	block bool,
	cancelFn context.CancelFunc,
) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	f := func() {
		sig := <-sigCh
		cancelFn()

		n.logger.Info("caught exit signal", "signal", sig.String())
	}

	if block {
		g.Go(func() error {
			f()
			return nil
		})
	} else {
		go f()
	}
}

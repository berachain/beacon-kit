// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package comet

import (
	"errors"
	"net/http"

	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/light/mod/provider/comet/types"
	cometOs "github.com/cometbft/cometbft/libs/os"
	lproxy "github.com/cometbft/cometbft/light/proxy"
	lrpc "github.com/cometbft/cometbft/light/rpc"
)

// StartProxy starts both the light client and the corresponding proxy
// with the given configuration.
func StartProxy(cfg *Config) error {
	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	serverCfg := initServerConfig(cfg.MaxOpenConnections)

	// Set options for the light rpc proxy
	opts := []lrpc.Option{
		lrpc.KeyPathFn(lrpc.DefaultMerkleKeyPathFn()),
		func(c *lrpc.Client) {
			c.RegisterOpDecoder(
				storetypes.ProofOpIAVLCommitment,
				storetypes.CommitmentOpDecoder,
			)
			c.RegisterOpDecoder(
				storetypes.ProofOpSimpleMerkleCommitment,
				storetypes.CommitmentOpDecoder,
			)
		},
	}
	proxy, err := lproxy.NewProxy(
		client,
		cfg.ListeningAddr,
		cfg.PrimaryAddr,
		serverCfg,
		cfg.Logger,
		opts...,
	)
	if err != nil {
		return err
	}

	// Stop upon receiving SIGTERM or CTRL-C.
	cometOs.TrapSignal(cfg.Logger, func() {
		proxy.Listener.Close()
	})

	cfg.Logger.Info("Starting proxy...", "listening at", cfg.ListeningAddr)
	go func() {
		if err = proxy.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			cfg.Logger.Error(types.ListenAndServeError, err)
		}
	}()

	return nil
}

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

package provider

import (
	"context"
	"time"

	lightprovider "github.com/cometbft/cometbft/light/provider"
	httpprovider "github.com/cometbft/cometbft/light/provider/http"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cometbft/cometbft/types"
)

// Provider is a struct which provides all full node data to the light client.
type Provider struct {
	config        Config
	client        rpcclient.Client
	cometProvider lightprovider.Provider
}

// New returns a new provider with the given configuration.
func New(config Config) *Provider {
	return &Provider{
		config: config,
	}
}

// Start starts the provider.
func (p *Provider) Start() error {
	httpClient, err := httpclient.New(
		p.config.HTTPEndpoint,
		p.config.WSEndpoint,
	)
	if err != nil {
		return err
	}
	p.client = httpClient
	p.cometProvider = httpprovider.NewWithClient(p.config.ChainID, httpClient)
	return nil
}

// Subscribe subscribes to the provider.
// TODO: Remove this subscription; Once we have a lightchain this will no longer
// be useful.
func (p *Provider) SubscribeToLightBlock(
	ctx context.Context,
) (chan *types.LightBlock, error) {
	ch := make(
		chan *types.LightBlock,
		256, //nolint:gomnd // will be removed later
	)
	ticker := time.NewTicker(
		5 * time.Second, //nolint:gomnd // will be removed later
	)

	latestBlockHeight := int64(-1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				res, err := p.client.ABCIInfo(ctx)
				if err != nil {
					panic(err)
				}

				var lb *types.LightBlock
				if res.Response.GetLastBlockHeight() > latestBlockHeight {
					lb, err = p.cometProvider.LightBlock(
						ctx,
						res.Response.GetLastBlockHeight(),
					)
					if err != nil {
						panic(err)
					}
					ch <- lb
					latestBlockHeight = lb.Height - 1
				}
			}
		}
	}()

	return ch, nil
}

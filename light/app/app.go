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

package app

import (
	"context"
	"log"

	"github.com/berachain/beacon-kit/light/mod/provider"
)

// RunLightNode starts the light node with the given configuration.
// In the future, this will build the light node runtime and start it.
func RunLightNode(ctx context.Context, config *Config) {
	// builds the runtime and starts the node
	client := provider.New(*config.Provider)
	if err := client.Start(); err != nil {
		panic(err)
	}

	// subscribe to light block
	ch, err := client.SubscribeToLightBlock(ctx)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case block := <-ch:
			// process the light block
			log.Println(block)
			log.Println("trustedEth1Hash", client.GetEth1BlockHash(ctx))
			log.Println("latestBlockHeader", client.GetLatestBlockHeader(ctx))

			// if block.Height > 4 {
			// 	fmt.Println("blockRootAtIndex", client.GetBlockRootAtIndex(ctx,
			// 2))
			// }
		case <-ctx.Done():
			// context cancelled, exit the loop
			return
		}
	}
}

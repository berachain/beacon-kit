package app

import (
	"context"
	"fmt"

	"github.com/berachain/beacon-kit/light/provider"
)

func RunLightNode(ctx context.Context, config *Config) {
	// builds the runtime and starts the node
	client := provider.New(*config.Provider)
	client.Start()

	// subscribe to light block
	ch, err := client.SubscribeToLightBlock(ctx)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case block := <-ch:
			// process the light block
			fmt.Println(block)
			fmt.Println("trustedEth1Hash", client.GetTrustedEth1Hash())
			fmt.Println("latestBlockHeader", client.GetLatestBlockHeader())

			if block.Height > 4 {
				fmt.Println("blockRootAtIndex", client.GetBlockRootAtIndex(2))
			}
		case <-ctx.Done():
			// context cancelled, exit the loop
			return
		}
	}
}

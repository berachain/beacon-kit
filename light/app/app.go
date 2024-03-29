package app

import (
	"context"
	"fmt"

	"github.com/berachain/beacon-kit/light/client"
)

func RunLightNode(ctx context.Context, config *Config) {
	// builds the runtime and starts the node
	lightClient := client.New(*config.Provider)
	lightClient.Start(ctx)

	// subscribe to light block
	ch, err := lightClient.SubscribeToLightBlock(ctx)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case block := <-ch:
			// process the light block
			fmt.Println(block)
		case <-ctx.Done():
			// context cancelled, exit the loop
			return
		}
	}
}

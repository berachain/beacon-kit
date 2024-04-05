package app

import (
	"context"

	nodebuilder "github.com/berachain/beacon-kit/mod/node-builder"
	beaconkitruntime "github.com/berachain/beacon-kit/mod/runtime"
	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/berachain/beacon-kit/light/provider"
)

var _ nodebuilder.BeaconApp = (*LightApp)(nil)

type LightApp struct {
	runtime *beaconkitruntime.BeaconKitRuntime
}

func NewLightApp(ctx context.Context, config *Config) *LightApp {
	return &LightApp{
		// runtime: runtime,
	}
}

func (app *LightApp) PostStartup(ctx context.Context, _ client.Context) error {
	// Initial check for execution client sync.
	app.runtime.StartServices(
		ctx,
	)
	return nil
}

// func RunLightNode(ctx context.Context, config *Config) {
// 	// builds the runtime and starts the node
// 	client := provider.New(*config.Provider)
// 	client.Start()

// 	// subscribe to light block
// 	ch, err := client.SubscribeToLightBlock(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for {
// 		select {
// 		case block := <-ch:
// 			// process the light block
// 			fmt.Println(block)
// 			fmt.Println("trustedEth1Hash", client.GetTrustedEth1Hash())
// 			fmt.Println("latestBlockHeader", client.GetLatestBlockHeader())

// 			if block.Height > 4 {
// 				fmt.Println("blockRootAtIndex", client.GetBlockRootAtIndex(2))
// 			}
// 		case <-ctx.Done():
// 			// context cancelled, exit the loop
// 			return
// 		}
// 	}
// }

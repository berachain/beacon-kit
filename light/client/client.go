package client

import (
	"context"

	"github.com/cometbft/cometbft/types"

	"github.com/berachain/beacon-kit/light/client/provider"
)

// Client is a client for the Light API.
type Client struct {
	provider *provider.Provider

	trustedHeight int64
}

// New creates a new light client.
func New(config provider.Config) *Client {
	return &Client{
		provider: provider.New(config),
	}
}

func (c *Client) GetTrustedEth1Hash() []byte {
	hash, err := c.provider.QueryWithProof(context.Background(), finalized_key, c.trustedHeight)
	if err != nil {
		panic(err)
	}
	return hash.Bytes()
}

// TODO: get trusted block execution payload

// func (c *Client) GetTrustedEthRpcURL() string {
// 	return c.provider.GetTrustedEthRpcURL
// }

// func (c *Client) GetClient() rpcclient.Client {
// 	return p.client
// }

// func (c *Client) LatestBlockHeight() int64 {
// 	return p.latestBlockHeight
// }

// Start starts the provider.
func (c *Client) Start(_ context.Context) error {
	return c.provider.Start()
}

// Subscribe subscribes to the provider.
func (c *Client) SubscribeToLightBlock(ctx context.Context) (chan *types.LightBlock, error) {
	return c.provider.SubscribeToLightBlock(ctx)
}

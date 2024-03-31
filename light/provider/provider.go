package provider

import (
	"context"
	"time"

	lightprovider "github.com/cometbft/cometbft/light/provider"
	httpprovider "github.com/cometbft/cometbft/light/provider/http"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cometbft/cometbft/types"

	// "github.com/berachain/cometbft/crypto/merkle"
	"github.com/cometbft/cometbft/crypto/merkle"

	storetypes "cosmossdk.io/store/types"
	// storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

// Provider is a struct which provides all full node data to the light client
type Provider struct {
	config        Config
	client        rpcclient.Client
	cometProvider lightprovider.Provider
	prt           *merkle.ProofRuntime

	latestBlockHeight int64
}

func New(config Config) *Provider {
	prt := merkle.DefaultProofRuntime()
	prt.RegisterOpDecoder(storetypes.ProofOpIAVLCommitment, storetypes.CommitmentOpDecoder)
	prt.RegisterOpDecoder(storetypes.ProofOpSimpleMerkleCommitment, storetypes.CommitmentOpDecoder)

	return &Provider{
		prt:               prt,
		config:            config,
		latestBlockHeight: -1,
	}
}

func (p *Provider) Start() error {
	httpClient, err := httpclient.New(p.config.HttpEndpoint, p.config.WSEndpoint)
	if err != nil {
		return err
	}
	p.client = httpClient
	p.cometProvider = httpprovider.NewWithClient(p.config.ChainID, httpClient)
	return nil
}

// Subscribe subscribes to the provider.
func (p *Provider) SubscribeToLightBlock(ctx context.Context) (chan *types.LightBlock, error) {
	ch := make(chan *types.LightBlock, 256)
	ticker := time.NewTicker(5 * time.Second)
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

				if res.Response.GetLastBlockHeight() > p.latestBlockHeight {
					lb, err := p.cometProvider.LightBlock(ctx, res.Response.GetLastBlockHeight())
					if err != nil {
						panic(err)
					}
					ch <- lb
					p.latestBlockHeight = lb.Height - 1
				}
			}
		}
	}()
	return ch, nil
}

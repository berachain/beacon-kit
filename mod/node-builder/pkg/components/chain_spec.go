package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/spec"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/chain"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

const ChainSpecFileName = "chainspec.toml"

type ChainSpecInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

func ProvideChainSpec(in ChainSpecInput) (primitives.ChainSpec, error) {
	data, err := spec.ReadFromAppOpts(in.AppOpts)
	if err != nil {
		return nil, err
	}

	return chain.NewChainSpec(*data), nil
}

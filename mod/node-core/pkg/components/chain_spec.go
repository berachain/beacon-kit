package components

import (
	"fmt"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config/spec"
	"github.com/berachain/beacon-kit/mod/primitives"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

type ChainSpecIn struct {
	depinject.In

	AppOpts servertypes.AppOptions
}

// ProvideChainSpec provides the chain spec.
func ProvideChainSpec(in ChainSpecIn) (primitives.ChainSpec, error) {
	cs, err := spec.ReadFromAppOpts(in.AppOpts)
	if err != nil {
		panic("REEE")
	}

	fmt.Println("CHAIN SPEC FROM PROVIDE", cs)
	return cs, nil
}

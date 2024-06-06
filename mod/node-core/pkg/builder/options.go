package builder

import (
	"cosmossdk.io/depinject"
	node "github.com/berachain/beacon-kit/mod/node-core/pkg"
	"github.com/berachain/beacon-kit/mod/primitives"
)

type Opt[NodeT node.NodeI] func(*Builder[NodeT]) error

func WithName[NodeT node.NodeI](name string) Opt[NodeT] {
	return func(b *Builder[NodeT]) error {
		b.node.SetAppName(name)
		b.name = name
		return nil
	}
}

func WithDescription[NodeT node.NodeI](description string) Opt[NodeT] {
	return func(b *Builder[NodeT]) error {
		b.node.SetAppDescription(description)
		b.description = description
		return nil
	}
}

func WithDepInjectConfig[NodeT node.NodeI](config depinject.Config) Opt[NodeT] {
	return func(b *Builder[NodeT]) error {
		b.depInjectCfg = config
		return nil
	}
}

func WithChainSpec[NodeT node.NodeI](chainSpec primitives.ChainSpec) Opt[NodeT] {
	return func(b *Builder[NodeT]) error {
		b.chainSpec = chainSpec
		return nil
	}
}

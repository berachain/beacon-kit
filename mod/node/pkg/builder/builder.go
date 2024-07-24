package builder

import "github.com/berachain/beacon-kit/mod/node/pkg/types"

type NodeBuilder[NodeT types.Node] struct {
	node NodeT
}

package components

import (
	"github.com/berachain/beacon-kit/mod/node-core/pkg/node"
	service "github.com/berachain/beacon-kit/mod/node-core/pkg/services/registry"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

// ProvideNode is a function that provides the module to the
func ProvideNode(
	registry *service.Registry,
) types.Node {
	return node.New[types.Node](registry)
}

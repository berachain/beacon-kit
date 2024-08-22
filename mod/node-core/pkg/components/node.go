package components

import (
	"github.com/berachain/beacon-kit/mod/node-core/pkg/node"
	service "github.com/berachain/beacon-kit/mod/node-core/pkg/services/registry"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

// ProvideNode is a function that provides the module to the
func ProvideNode(registry *service.Registry) types.Node {
	n := node.New[types.Node]()
	n.SetServiceRegistry(registry)
	return n
}

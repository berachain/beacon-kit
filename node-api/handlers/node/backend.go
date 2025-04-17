package node

import "github.com/cometbft/cometbft/node"

// Backend is the interface for backend of the node API.
type Backend interface {
	GetNode() *node.Node
}

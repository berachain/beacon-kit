package ethclient2

import (
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient2/rpc"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// EthRPC - Ethereum rpc client
type EthRPC[ExecutionPayloadT interface {
	constraints.JSONMarshallable
	Empty(uint32) ExecutionPayloadT
}] struct {
	*rpc.Client
}

// New create new rpc client with given url
func New[
	ExecutionPayloadT interface {
		constraints.JSONMarshallable
		Empty(uint32) ExecutionPayloadT
	},
](url string, options ...func(rpc *rpc.Client)) *EthRPC[ExecutionPayloadT] {
	rpc := &EthRPC[ExecutionPayloadT]{
		Client: rpc.NewClient(url, options...),
	}

	return rpc
}

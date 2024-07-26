package ethclient

import (
	"context"
	"math/big"

	hexutil "github.com/ethereum/go-ethereum/common/hexutil"
)

// ChainID retrieves the current chain ID for transaction replay protection.
func (s *Eth1Client[ExecutionPayloadT]) ChainID(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := s.Client.CallContext(ctx, &result, "eth_chainId")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), nil
}

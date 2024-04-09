package cosmos

import (
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c ChainQuerier) callABCIQuery(ctx sdk.Context, path string) (*cometabci.ResponseQuery, error) {
	req := cometabci.RequestQuery{Path: path}
	responseQuery, err := c.ABCI.Query(ctx, &req)
	if err != nil {
		return nil, err
	}

	return responseQuery, nil
}

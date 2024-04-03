package provider

import (
	"context"
	"strconv"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/metadata"
)

func (p *Provider) QueryWithProof(ctx context.Context, key string, height int64) (common.Hash, error) {
	resp, _, err := p.RunGRPCQuery(
		context.Background(),
		beaconStoreKey,
		[]byte(key),
		height,
	)

	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(resp.Value), nil
}

func (p *Provider) RunGRPCQuery(ctx context.Context, method string, reqBz []byte, height int64) (abci.ResponseQuery, metadata.MD, error) {
	// // parse height header
	// if heights := md.Get(grpctypes.GRPCBlockHeightHeader); len(heights) > 0 {
	// 	height, err := strconv.ParseInt(heights[0], 10, 64)
	// 	if err != nil {
	// 		return abci.ResponseQuery{}, nil, err
	// 	}
	// 	if height < 0 {
	// 		return abci.ResponseQuery{}, nil, sdkerrors.Wrapf(
	// 			legacyerrors.ErrInvalidRequest,
	// 			"client.Context.Invoke: height (%d) from %q must be >= 0", height, grpctypes.GRPCBlockHeightHeader)
	// 	}

	// }

	// height, err := GetHeightFromMetadata(md)
	// if err != nil {
	// 	return abci.ResponseQuery{}, nil, err
	// }

	// prove, err := GetProveFromMetadata(md)
	// if err != nil {
	// 	return abci.ResponseQuery{}, nil, err
	// }

	abciReq := abci.RequestQuery{
		Path:   method,
		Data:   reqBz,
		Height: height,
	}

	abciRes, err := p.QueryABCI(ctx, abciReq)
	if err != nil {
		return abci.ResponseQuery{}, nil, err
	}

	// Create header metadata. For now the headers contain:
	// - block height
	// We then parse all the call options, if the call option is a
	// HeaderCallOption, then we manually set the value of that header to the
	// metadata.
	md := metadata.Pairs(grpctypes.GRPCBlockHeightHeader, strconv.FormatInt(abciRes.Height, 10))

	return abciRes, md, nil
}

// QueryABCI performs an ABCI query and returns the appropriate response and error sdk error code.
func (p *Provider) QueryABCI(ctx context.Context, req abci.RequestQuery) (abci.ResponseQuery, error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: req.Height,
	}

	// Note: ABCIQueryWithOptions verifies proofs by default. Thus we do not have to
	// check the proof validity ourselves.
	result, err := p.client.ABCIQueryWithOptions(ctx, req.Path, req.Data, opts)
	if err != nil {
		return abci.ResponseQuery{}, err
	} else if !result.Response.IsOK() {
		return abci.ResponseQuery{}, sdkErrorToGRPCError(result.Response)
	}

	return result.Response, nil
}

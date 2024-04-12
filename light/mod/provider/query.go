// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package provider

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
)

// Query performs a query to the store at path at a given height and
// returns the result.
func (p *Provider) Query(
	ctx context.Context,
	path string,
	key []byte,
	height int64,
) ([]byte, error) {
	// If height is not provided, query the latest height.
	if height == 0 {
		res, err := p.client.ABCIInfo(ctx)
		if err != nil {
			return nil, err
		}
		height = res.Response.LastBlockHeight - 1
	}

	resp, err := p.QueryABCI(
		ctx,
		abci.RequestQuery{
			Path:   path,
			Data:   key,
			Height: height,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// QueryABCI performs an ABCI query and returns the appropriate response and
// error sdk error code.
func (p *Provider) QueryABCI(ctx context.Context, req abci.RequestQuery) (
	abci.ResponseQuery, error,
) {
	opts := rpcclient.ABCIQueryOptions{
		Height: req.Height,
	}

	// Note: ABCIQueryWithOptions verifies proofs by default. Thus we do not
	// have to
	// check the proof validity ourselves.
	result, err := p.client.ABCIQueryWithOptions(ctx, req.Path, req.Data, opts)
	if err != nil {
		return abci.ResponseQuery{}, err
	} else if !result.Response.IsOK() {
		return abci.ResponseQuery{}, abciErrorToSdkError(result.Response)
	}

	return result.Response, nil
}

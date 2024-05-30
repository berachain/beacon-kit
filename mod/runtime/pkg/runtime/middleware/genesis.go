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

package middleware

import (
	"context"
	"encoding/json"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/sourcegraph/conc/iter"
)

// InitGenesis is called by the base app to initialize the beacon state from
// the eth1 genesis data.
func (h *FinalizeBlockMiddleware[
	BeaconBlockT, BeaconStateT, BlobsSidecarsT,
]) InitGenesis(
	ctx context.Context,
	bz []byte,
) ([]appmodulev2.ValidatorUpdate, error) {
	data := new(
		genesis.Genesis[*types.Deposit, *types.ExecutionPayloadHeaderDeneb],
	)
	if err := json.Unmarshal(bz, data); err != nil {
		return nil, err
	}

	updates, err := h.chainService.ProcessGenesisData(
		ctx,
		data,
	)
	if err != nil {
		return nil, err
	}

	// Convert updates into the Cosmos SDK format.
	return iter.MapErr(updates, convertValidatorUpdate)
}

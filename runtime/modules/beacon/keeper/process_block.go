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

package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/itsdevbear/bolaris/config/version"
	"github.com/itsdevbear/bolaris/runtime/modules/beacon/types"
	"github.com/itsdevbear/bolaris/types/consensus"
)

func (k *Keeper) ProcessBeaconBlock(
	ctx context.Context,
	beaconBlockBz *types.BeaconKitBlockContainer,
) (*types.ProcessBeaconBlockResponse, error) {
	cctx := sdk.UnwrapSDKContext(ctx)
	if cctx.BlockHeight()%2 == 0 || cctx.BlockHeight()%4 == 0 ||
		cctx.BlockHeight()%5 == 0 ||
		cctx.BlockHeight()%6 == 0 {
		k.chainService.Logger().Error("MISSED SLOT", "slot", cctx.BlockHeight())
		return nil, errors.New("$$$$ MISSED SLOT $$$$$$$")
	}

	// Create a new beacon block from the given SSZ bytes.
	blk, err := consensus.BeaconKitBlockFromSSZ(
		beaconBlockBz.SszBlockBytes,
		version.Deneb,
	)
	if err != nil {
		return nil, err
	}

	// Process the finalization of the beacon block.
	if err = k.chainService.FinalizeBeaconBlock(
		ctx, blk,
	); err != nil {
		return nil, err
	}

	return &types.ProcessBeaconBlockResponse{}, nil
}

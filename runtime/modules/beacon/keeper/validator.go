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

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/runtime/modules/beacon/types"
	"github.com/ethereum/go-ethereum/common"
)

// CreateValidator creates a new validator in the beacon state.
func (k *Keeper) CreateValidator(
	ctx context.Context,
	msg *types.MsgCreateValidatorX,
) (*types.MsgCreateValidatorResponse, error) {
	address := common.HexToAddress(msg.Credentials)
	val := &beacontypes.Validator{
		Pubkey: [48]byte(msg.Pubkey),
		Credentials: beacontypes.
			NewCredentialsFromExecutionAddress(address),
		EffectiveBalance: 1, // todo: should be zero at creation.
		Slashed:          false,
	}
	if err := k.beaconStore.AddValidator(ctx, val); err != nil {
		return nil, err
	}

	return &types.MsgCreateValidatorResponse{}, nil
}

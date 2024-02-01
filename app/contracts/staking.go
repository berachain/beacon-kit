// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package contracts

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// TODO: figure out if we can skip this, I feel as if we can?
//
// The is_valid_merkle_branch() checks to ensure that it's not possible to fake a deposit in
// the process_deposit() function. The eth1data.deposit_root from the deposit contract has
// been agreed by the beacon chain and includes all pending deposits visible to the beacon
// chain. The deposit itself contains a Merkle proof that it is included in that root. The
// state.eth1_deposit_index counter ensures that deposits are processed in order. In short,
// the proposer provides leaf and branch, but neither index nor root.
type StakingCallbacks struct {
	stakingtypes.MsgServer
}

func (s *StakingCallbacks) ABIEvents() map[string]abi.Event {
	x, _ := StakingMetaData.GetAbi()
	return x.Events
}

func (s *StakingCallbacks) Delegate(
	ctx context.Context, operatorAddress string, amount *big.Int,
) error {
	sdk.UnwrapSDKContext(ctx).Logger().Info("delegating from execution layer",
		"operatorAddress", operatorAddress, "amt", amount)

	return nil
}

func (s *StakingCallbacks) Undelegate(
	ctx context.Context, operatorAddress string, amount *big.Int,
) error {
	sdk.UnwrapSDKContext(ctx).Logger().Info("undelegating from execution layer",
		"operatorAddress", operatorAddress, "amt", amount)
	return nil
}

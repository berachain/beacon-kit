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

package logs

import (
	"context"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/itsdevbear/bolaris/beacon/execution/logs/callback"
	"github.com/itsdevbear/bolaris/contracts/abi"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus"
	"github.com/itsdevbear/bolaris/types/engine"
)

var _ callback.ContractHandler = &Handler{}

// Handler is a struct that implements the callback Handler interface.
type Handler struct {
	service.BaseService

	// sks is the staking service.
	sks StakingService
}

// ABIEvents returns the events defined in the staking contract ABI.
func (s *Handler) ABIEvents() map[string]ethabi.Event {
	stakingAbi, err := abi.StakingMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	return stakingAbi.Events
}

// Deposit is a callback function that is called
// when a Deposit event is emitted from the staking contract.
func (s *Handler) Deposit(
	ctx context.Context,
	validatorPubkey []byte,
	withdrawalCredentials []byte,
	amount uint64,
) error {
	deposit := consensus.NewDeposit(
		validatorPubkey,
		amount,
		withdrawalCredentials,
	)
	return s.sks.AcceptDepositIntoQueue(ctx, deposit)
}

// Withdrawal is a callback function that is called
// when a Withdrawal event is emitted from the staking contract.
func (s *Handler) Withdrawal(
	ctx context.Context,
	validatorPubkey []byte,
	_ []byte,
	amount uint64,
) error {
	return s.sks.ProcessWithdrawal(
		ctx,
		engine.NewWithdrawal(validatorPubkey, amount),
	)
}

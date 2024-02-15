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

package staking

import (
	"context"
	"math/big"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/itsdevbear/bolaris/beacon/execution/logs/callback"
	stakingabi "github.com/itsdevbear/bolaris/beacon/execution/staking/abi"
	"github.com/itsdevbear/bolaris/runtime/modules/beacon/keeper/store"
	"github.com/itsdevbear/bolaris/runtime/service"
)

var _ callback.Handler = &Handler{}

// Handler is a struct that implements the callback Handler interface.
type Handler struct {
	service.BaseService
	logger log.Logger
}

// ABIEvents returns the events defined in the staking contract ABI.
func (s *Handler) ABIEvents() map[string]abi.Event {
	x, err := stakingabi.StakingMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	return x.Events
}

func NewHandler(
	base service.BaseService,
	opts ...Option,
) (*Handler, error) {
	s := &Handler{
		BaseService: base,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Delegate is a callback function that is called
// when a Delegate event is emitted from the staking contract.
func (s *Handler) Delegate(
	ctx context.Context,
	validatorPubkey []byte,
	withdrawalCredentials []byte,
	amount *big.Int,
	nonce *big.Int,
) error {
	beaconState := s.BeaconState(ctx)
	expectedNonce := beaconState.GetStakingNonce()
	logNonce := nonce.Uint64()
	// We may receive the same deposit log twice from the execution layer, just ignore it.
	if logNonce < expectedNonce {
		return nil
	}
	// The deposit log does not come in order.
	if logNonce != expectedNonce {
		return ErrInvalidNonce
	}
	deposit := store.NewDeposit(validatorPubkey, amount, withdrawalCredentials)
	err := beaconState.AddDeposit(deposit)
	if err != nil {
		return err
	}
	beaconState.SetStakingNonce(expectedNonce + 1)
	s.Logger().Info("delegating from execution layer",
		"validatorPubkey", validatorPubkey, "amount", amount, "nonce", nonce)
	return nil
}

func (s *Handler) Undelegate(
	ctx context.Context,
	validatorPubkey []byte,
	amount *big.Int,
	nonce *big.Int,
) error {
	s.Logger().Info("undelegating from execution layer",
		"validatorPubkey", validatorPubkey, "amount", amount, "nonce", nonce)
	return nil
}

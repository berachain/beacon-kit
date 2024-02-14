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

package store

import (
	"cosmossdk.io/collections"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

type Deposit struct {
	*consensusv1.Deposit
}

// GetAmount returns the amount of the deposit.
func (d *Deposit) GetAmount() uint64 {
	return d.GetData().GetAmount()
}

// GetPubkey returns the public key of the validator in the deposit.
func (d *Deposit) GetPubkey() []byte {
	return d.GetData().GetPubkey()
}

// AddDeposit adds a deposit to the staking module.
func (s *BeaconStore) AddDeposit(deposit *Deposit) error {
	s.deposits = append(s.deposits, deposit)
	return nil
}

// NextDeposit returns the next deposit in the queue.
func (s *BeaconStore) NextDeposit() (*Deposit, error) {
	if len(s.deposits) == 0 {
		return nil, collections.ErrNotFound
	}
	deposit := s.deposits[0]
	s.deposits = s.deposits[1:]
	return deposit, nil
}

// ProcessDeposit processes a deposit with the staking keeper.
func (s *BeaconStore) ProcessDeposit(deposit *Deposit) error {
	_, err := s.stakingKeeper.Delegate(s.sdkCtx, deposit)
	return err
}

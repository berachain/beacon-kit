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

package store

import (
	"cosmossdk.io/collections/codec"
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

type DepositValue struct{}

var _ codec.ValueCodec[*Deposit] = DepositValue{}

func (DepositValue) Encode(value *Deposit) ([]byte, error) {
	return value.MarshalSSZ()
}

func (DepositValue) Decode(b []byte) (*Deposit, error) {
	value := &Deposit{&consensusv1.Deposit{}}
	if err := value.UnmarshalSSZ(b); err != nil {
		return nil, err
	}
	return value, nil
}

func (DepositValue) EncodeJSON(_ *Deposit) ([]byte, error) {
	panic("not implemented")
}

func (DepositValue) DecodeJSON(_ []byte) (*Deposit, error) {
	panic("not implemented")
}

func (DepositValue) Stringify(value *Deposit) string {
	return value.String()
}

func (d DepositValue) ValueType() string {
	return "Deposit"
}

// AddDeposit adds a deposit to the staking module.
func (s *BeaconStore) AddDeposit(deposit *Deposit) error {
	return s.deposits.Push(s.sdkCtx, deposit)
}

// NextDeposit returns the next deposit in the queue.
func (s *BeaconStore) NextDeposit() (*Deposit, error) {
	return s.deposits.Pop(s.sdkCtx)
}

// ProcessDeposit processes a deposit with the staking keeper.
func (s *BeaconStore) ProcessDeposit(deposit *Deposit) error {
	_, err := s.stakingKeeper.Delegate(s.sdkCtx, deposit)
	return err
}

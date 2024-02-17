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
	"github.com/itsdevbear/bolaris/lib/encoding"
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

// NewDeposit creates a new deposit.
func NewDeposit(pubkey []byte, amount uint64, withdrawalCredentials []byte) *Deposit {
	depositData := &consensusv1.Deposit_Data{
		Pubkey:                pubkey,
		Amount:                amount,
		WithdrawalCredentials: withdrawalCredentials,
	}
	return &Deposit{&consensusv1.Deposit{Data: depositData}}
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

// CacheDeposit caches a deposit.
func (s *BeaconStore) CacheDeposit(deposit *Deposit) error {
	s.depositCache = append(s.depositCache, deposit)
	return nil
}

// CommitDeposits commits the cached deposits to the queue.
func (s *BeaconStore) CommitDeposits() error {
	err := s.deposits.PushMulti(s.sdkCtx, s.depositCache)
	if err != nil {
		return err
	}
	s.depositCache = nil
	return nil
}

// PersistDeposits pops the next deposits, up to n,
// from the queue and delegate them with staking keeper.
func (s *BeaconStore) PersistDeposits(n uint64) ([]*Deposit, error) {
	var err error
	depositsToProcess, err := s.deposits.PopMulti(s.sdkCtx, n)
	if err != nil {
		return nil, err
	}
	for _, deposit := range depositsToProcess {
		// TODO: If an error occurs in the middle of processing deposits,
		// should we continue to process the remaining deposits?
		if err = s.processDeposit(deposit); err != nil {
			return nil, err
		}
	}
	return depositsToProcess, nil
}

// processDeposit processes a deposit with the staking keeper.
func (s *BeaconStore) processDeposit(deposit *Deposit) error {
	_, err := s.stakingKeeper.Delegate(s.sdkCtx, deposit)
	return err
}

// SetStakingNonce sets the staking nonce.
func (s *BeaconStore) SetStakingNonce(nonce uint64) {
	bz := encoding.EncodeUint64(nonce)
	s.Set([]byte(stakingNonceKey), bz)
}

// GetStakingNonce returns the staking nonce.
func (s *BeaconStore) GetStakingNonce() uint64 {
	bz := s.Get([]byte(stakingNonceKey))
	if bz == nil {
		return 0
	}
	return encoding.DecodeUint64(bz)
}

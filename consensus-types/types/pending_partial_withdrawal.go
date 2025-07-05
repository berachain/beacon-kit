// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//go:generate sszgen -path pending_partial_withdrawal.go -objs PendingPartialWithdrawal -output pending_partial_withdrawal_sszgen.go -include ../../primitives/common,../../primitives/math

package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
)


var (
	_ constraints.SSZMarshallableRootable = (*PendingPartialWithdrawal)(nil)
	_ constraints.SSZMarshallable = (*PendingPartialWithdrawals)(nil)
)

// PendingPartialWithdrawal reflects the following spec:
//
//	class PendingPartialWithdrawal(Container):
//	    validator_index: ValidatorIndex
//	    amount: Gwei
//	    withdrawable_epoch: Epoch
type PendingPartialWithdrawal struct {
	ValidatorIndex    math.ValidatorIndex
	Amount            math.Gwei
	WithdrawableEpoch math.Epoch
}

// ValidateAfterDecodingSSZ validates the PendingPartialWithdrawal object
// after decoding from SSZ. Customize further validation as needed.
func (p *PendingPartialWithdrawal) ValidateAfterDecodingSSZ() error {
	return nil
}

// PendingPartialWithdrawals is a SSZ list of PendingPartialWithdrawal containers.
type PendingPartialWithdrawals []*PendingPartialWithdrawal

// NewEmptyPendingPartialWithdrawals returns a new empty PendingPartialWithdrawals list.
func NewEmptyPendingPartialWithdrawals() *PendingPartialWithdrawals {
	return &PendingPartialWithdrawals{}
}

// SizeSSZ returns the size of the PendingPartialWithdrawals list.
func (p *PendingPartialWithdrawals) SizeSSZ() int {
	return 4 + len(*p)*24 // 24 bytes per PendingPartialWithdrawal
}

// MarshalSSZ returns the SSZ encoding of the PendingPartialWithdrawals list.
func (p *PendingPartialWithdrawals) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, p.SizeSSZ())
	return p.MarshalSSZTo(buf)
}

// MarshalSSZTo marshals the PendingPartialWithdrawals list into a pre-allocated byte slice.
func (p *PendingPartialWithdrawals) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Write offset
	offset := 4
	dst = fastssz.MarshalUint32(dst, uint32(offset))

	// Write elements
	for _, elem := range *p {
		var err error
		dst, err = elem.MarshalSSZTo(dst)
		if err != nil {
			return nil, err
		}
	}

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the PendingPartialWithdrawals list.
func (p *PendingPartialWithdrawals) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 4 {
		return fastssz.ErrSize
	}

	// Read offset
	offset := fastssz.UnmarshallUint32(buf[0:4])
	if offset != 4 {
		return fastssz.ErrInvalidVariableOffset
	}

	// Read elements
	if (len(buf)-4)%24 != 0 {
		return fastssz.ErrSize
	}

	numItems := (len(buf) - 4) / 24
	if numItems > constants.PendingPartialWithdrawalsLimit {
		return errors.New("pending partial withdrawals too large")
	}

	*p = make([]*PendingPartialWithdrawal, numItems)
	for i := 0; i < numItems; i++ {
		(*p)[i] = &PendingPartialWithdrawal{}
		start := 4 + i*24
		end := start + 24
		if err := (*p)[i].UnmarshalSSZ(buf[start:end]); err != nil {
			return err
		}
	}

	return nil
}

// HashTreeRoot computes the hash tree root for PendingPartialWithdrawals.
func (p *PendingPartialWithdrawals) HashTreeRoot() ([32]byte, error) {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	if err := p.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()

}

// HashTreeRootWith ssz hashes the PendingPartialWithdrawals object with a hasher.
func (p *PendingPartialWithdrawals) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(*p))
	if num > constants.PendingPartialWithdrawalsLimit {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range *p {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(indx, num, constants.PendingPartialWithdrawalsLimit)
	return nil
}

// GetTree ssz hashes the PendingPartialWithdrawals object.
func (p *PendingPartialWithdrawals) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(p)
}

// ValidateAfterDecodingSSZ validates the PendingPartialWithdrawals list after decoding from SSZ.
func (p *PendingPartialWithdrawals) ValidateAfterDecodingSSZ() error {
	if p == nil {
		return errors.New("nil PendingPartialWithdrawals")
	}
	if len(*p) > constants.PendingPartialWithdrawalsLimit {
		return errors.New("pending partial withdrawals too large")
	}
	return nil
}

// PendingBalanceToWithdraw implements get_pending_balance_to_withdraw from the ETH2.0 spec.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#new-get_pending_balance_to_withdraw
func (p PendingPartialWithdrawals) PendingBalanceToWithdraw(validatorIndex math.ValidatorIndex) math.Gwei {
	var total math.Gwei
	for _, withdrawal := range p {
		if withdrawal.ValidatorIndex == validatorIndex {
			total += withdrawal.Amount
		}
	}
	return total
}

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

// TODO: PendingPartialWithdrawal ready for full fastssz migration (already has all needed methods)
// go:generate sszgen -path . -objs PendingPartialWithdrawal -output pending_partial_withdrawal_sszgen.go -include ../../primitives/common,../../primitives/math

package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
)

// sszPendingPartialWithdrawalSize defines the total SSZ serialized size for
// PendingPartialWithdrawal. The fields are assumed to be encoded as follows:
// - ValidatorIndex: 8 bytes (uint64)
// - Amount:         8 bytes (math.Gwei)
// - WithdrawableEpoch: 8 bytes (uint64)
// Total = 8 + 8 + 8 = 24 bytes.
const sszPendingPartialWithdrawalSize = 24

// TODO: Re-enable interface assertions once constraints are updated
// var (
// 	_ constraints.SSZMarshallable = (*PendingPartialWithdrawal)(nil)
// 	_ constraints.SSZMarshallable = (*PendingPartialWithdrawals)(nil)
// )

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

/* -------------------------------------------------------------------------- */
/*                      PendingPartialWithdrawal SSZ                          */
/* -------------------------------------------------------------------------- */

// ValidateAfterDecodingSSZ validates the PendingPartialWithdrawal object
// after decoding from SSZ. Customize further validation as needed.
func (p *PendingPartialWithdrawal) ValidateAfterDecodingSSZ() error {
	return nil
}


// SizeSSZ returns the fixed size of the SSZ serialization for PendingPartialWithdrawal.
func (p *PendingPartialWithdrawal) SizeSSZ() int {
	return sszPendingPartialWithdrawalSize
}

// MarshalSSZ returns the SSZ encoding of the PendingPartialWithdrawal.
func (p *PendingPartialWithdrawal) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, sszPendingPartialWithdrawalSize)
	return p.MarshalSSZTo(buf)
}

// HashTreeRoot computes and returns the hash tree root for the PendingPartialWithdrawal.
func (p *PendingPartialWithdrawal) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	p.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

// HashTreeRootWith SSZ hashes the Deposit object with a hasher. Needed for BeaconState SSZ.
func (p *PendingPartialWithdrawal) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'ValidatorIndex'
	hh.PutUint64(uint64(p.ValidatorIndex))

	// Field (1) 'Amount'
	hh.PutUint64(uint64(p.Amount))

	// Field (2) 'WithdrawableEpoch'
	hh.PutUint64(uint64(p.WithdrawableEpoch))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the PendingPartialWithdrawal object.
func (p *PendingPartialWithdrawal) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(p)
}

// MarshalSSZTo marshals the PendingPartialWithdrawal object into a pre-allocated byte slice.
func (p *PendingPartialWithdrawal) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Field (0) 'ValidatorIndex'
	dst = fastssz.MarshalUint64(dst, uint64(p.ValidatorIndex))

	// Field (1) 'Amount'
	dst = fastssz.MarshalUint64(dst, uint64(p.Amount))

	// Field (2) 'WithdrawableEpoch'
	dst = fastssz.MarshalUint64(dst, uint64(p.WithdrawableEpoch))

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the PendingPartialWithdrawal object.
func (p *PendingPartialWithdrawal) UnmarshalSSZ(buf []byte) error {
	if len(buf) != sszPendingPartialWithdrawalSize {
		return fastssz.ErrSize
	}

	// Field (0) 'ValidatorIndex'
	p.ValidatorIndex = math.ValidatorIndex(fastssz.UnmarshallUint64(buf[0:8]))

	// Field (1) 'Amount'
	p.Amount = math.Gwei(fastssz.UnmarshallUint64(buf[8:16]))

	// Field (2) 'WithdrawableEpoch'
	p.WithdrawableEpoch = math.Epoch(fastssz.UnmarshallUint64(buf[16:24]))

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
	return 4 + len(*p)*sszPendingPartialWithdrawalSize
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
	if (len(buf)-4)%sszPendingPartialWithdrawalSize != 0 {
		return fastssz.ErrSize
	}

	numItems := (len(buf) - 4) / sszPendingPartialWithdrawalSize
	if numItems > constants.PendingPartialWithdrawalsLimit {
		return errors.New("pending partial withdrawals too large")
	}

	*p = make([]*PendingPartialWithdrawal, numItems)
	for i := 0; i < numItems; i++ {
		(*p)[i] = &PendingPartialWithdrawal{}
		start := 4 + i*sszPendingPartialWithdrawalSize
		end := start + sszPendingPartialWithdrawalSize
		if err := (*p)[i].UnmarshalSSZ(buf[start:end]); err != nil {
			return err
		}
	}

	return nil
}

// HashTreeRoot computes the hash tree root for PendingPartialWithdrawals.
func (p *PendingPartialWithdrawals) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	p.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
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

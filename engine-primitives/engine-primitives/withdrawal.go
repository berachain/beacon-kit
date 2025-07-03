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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package engineprimitives

import (
	"io"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	fastrlp "github.com/umbracle/fastrlp"
)

// withdrawalSize is the size of the Withdrawal in bytes.
const withdrawalSize = 44

// TODO: Re-enable interface assertion once constraints are updated
// var (
// 	_ constraints.SSZMarshallableRootable = (*Withdrawal)(nil)
// )

// Withdrawal represents a validator withdrawal from the consensus layer.
type Withdrawal struct {
	// Index is the unique identifier for the withdrawal.
	Index math.U64 `json:"index"`
	// Validator is the index of the validator initiating the withdrawal.
	Validator math.ValidatorIndex `json:"validatorIndex"`
	// Address is the execution address where the withdrawal will be sent.
	// It has a fixed size of 20 bytes.
	Address common.ExecutionAddress `json:"address"`
	// Amount is the amount of Gwei to be withdrawn.
	Amount math.Gwei `json:"amount"`
}

func NewWithdrawal(
	index math.U64,
	validator math.ValidatorIndex,
	address common.ExecutionAddress,
	amount math.Gwei,
) *Withdrawal {
	return &Withdrawal{
		Index:     index,
		Validator: validator,
		Address:   address,
		Amount:    amount,
	}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the Withdrawal in bytes when SSZ encoded.
func (*Withdrawal) SizeSSZ() int {
	return withdrawalSize
}


// HashTreeRoot.
func (w *Withdrawal) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	w.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

// MarshalSSZ marshals the Withdrawal object to SSZ format.
func (w *Withdrawal) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, withdrawalSize)
	return w.MarshalSSZTo(buf)
}

func (*Withdrawal) ValidateAfterDecodingSSZ() error { return nil }

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the Withdrawal object to a target array.
func (w *Withdrawal) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Field (0) 'Index'
	dst = fastssz.MarshalUint64(dst, uint64(w.Index))

	// Field (1) 'Validator'
	dst = fastssz.MarshalUint64(dst, uint64(w.Validator))

	// Field (2) 'Address'
	dst = append(dst, w.Address[:]...)

	// Field (3) 'Amount'
	dst = fastssz.MarshalUint64(dst, uint64(w.Amount))

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the Withdrawal object.
func (w *Withdrawal) UnmarshalSSZ(buf []byte) error {
	if len(buf) != withdrawalSize {
		return fastssz.ErrSize
	}

	// Field (0) 'Index'
	w.Index = math.U64(fastssz.UnmarshallUint64(buf[0:8]))

	// Field (1) 'Validator'
	w.Validator = math.ValidatorIndex(fastssz.UnmarshallUint64(buf[8:16]))

	// Field (2) 'Address'
	copy(w.Address[:], buf[16:36])

	// Field (3) 'Amount'
	w.Amount = math.Gwei(fastssz.UnmarshallUint64(buf[36:44]))

	return nil
}

// HashTreeRootWith ssz hashes the Withdrawal object with a hasher.
func (w *Withdrawal) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Index'
	hh.PutUint64(uint64(w.Index))

	// Field (1) 'Validator'
	hh.PutUint64(uint64(w.Validator))

	// Field (2) 'Address'
	hh.PutBytes(w.Address[:])

	// Field (3) 'Amount'
	hh.PutUint64(uint64(w.Amount))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the Withdrawal object.
func (w *Withdrawal) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(w)
}

/* -------------------------------------------------------------------------- */
/*                                     RLP                                    */
/* -------------------------------------------------------------------------- */

// SetIndex sets the unique identifier for the withdrawal.
func (w Withdrawal) EncodeRLP(_w io.Writer) error {
	a := fastrlp.DefaultArenaPool.Get()
	defer fastrlp.DefaultArenaPool.Put(a)
	v := a.NewArray()
	v.Set(a.NewUint(uint64(w.Index)))
	v.Set(a.NewUint(uint64(w.Validator)))
	v.Set(a.NewCopyBytes(w.Address[:]))
	v.Set(a.NewUint(uint64(w.Amount)))
	_, err := _w.Write(v.MarshalTo(nil))
	return err
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// Equals returns true if the Withdrawal is equal to the other.
func (w *Withdrawal) Equals(rhs *Withdrawal) bool {
	switch {
	case w == nil && rhs == nil:
		return true
	case w != nil && rhs != nil:
		return w.Index == rhs.Index &&
			w.Validator == rhs.Validator &&
			w.Address == rhs.Address &&
			w.Amount == rhs.Amount
	default:
		return false
	}
}

// GetIndex returns the unique identifier for the withdrawal.
func (w *Withdrawal) GetIndex() math.U64 {
	return w.Index
}

// GetValidatorIndex returns the index of the validator initiating the
// withdrawal.
func (w *Withdrawal) GetValidatorIndex() math.ValidatorIndex {
	return w.Validator
}

// GetAddress returns the execution address where the withdrawal will be sent.
func (w *Withdrawal) GetAddress() common.ExecutionAddress {
	return w.Address
}

// GetAmount returns the amount of Gwei to be withdrawn.
func (w *Withdrawal) GetAmount() math.Gwei {
	return w.Amount
}

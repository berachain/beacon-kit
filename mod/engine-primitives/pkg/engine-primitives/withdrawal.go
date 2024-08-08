// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
	fastrlp "github.com/umbracle/fastrlp"
)

// WithdrawalSize is the size of the Withdrawal in bytes.
const WithdrawalSize = 44

var (
	_ ssz.StaticObject                    = (*Withdrawal)(nil)
	_ constraints.SSZMarshallableRootable = (*Withdrawal)(nil)
)

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

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

func (w *Withdrawal) New(
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
func (*Withdrawal) SizeSSZ() uint32 {
	return WithdrawalSize
}

// MarshalSSZ marshals the Withdrawal into SSZ format.
func (w *Withdrawal) DefineSSZ(c *ssz.Codec) {
	ssz.DefineUint64(c, &w.Index)        // Field  (0) -     Index -  8 bytes
	ssz.DefineUint64(c, &w.Validator)    // Field  (1) - Validator -  8 bytes
	ssz.DefineStaticBytes(c, &w.Address) // Field  (2) -   Address - 20 bytes
	ssz.DefineUint64(c, &w.Amount)       // Field  (3) -    Amount -  8 bytes
}

// HashTreeRoot.
func (w *Withdrawal) HashTreeRoot() common.Root {
	return ssz.HashSequential(w)
}

// MarshalSSZ marshals the Withdrawal object to SSZ format.
func (w *Withdrawal) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, w.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, w)
}

// UnmarshalSSZ unmarshals the SSZ encoded data to a Withdrawal object.
func (w *Withdrawal) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, w)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the Withdrawal object to a target array.
func (w *Withdrawal) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := w.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
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
func (w *Withdrawal) Equals(other *Withdrawal) bool {
	return w.Index == other.Index &&
		w.Validator == other.Validator &&
		w.Address == other.Address &&
		w.Amount == other.Amount
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

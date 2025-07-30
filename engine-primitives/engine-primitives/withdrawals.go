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
	"bytes"

	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	fastssz "github.com/ferranbt/fastssz"
)

var (
	_ constraints.SSZRootable = (*Withdrawals)(nil)
)

// Withdrawals represents a list of withdrawals.
type Withdrawals []*Withdrawal

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Withdrawals.
func (w Withdrawals) SizeSSZ() int {
	return 4 + len(w)*44 // offset + each withdrawal is 44 bytes
}

// HashTreeRoot returns the hash tree root of the Withdrawals.
func (w Withdrawals) HashTreeRoot() ([32]byte, error) {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	if err := w.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the Withdrawals object with a hasher.
func (w Withdrawals) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(w))
	if num > constants.MaxWithdrawalsPerPayload {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range w {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(indx, num, constants.MaxWithdrawalsPerPayload)
	return nil
}

// GetTree ssz hashes the Withdrawals object.
func (w Withdrawals) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(w)
}

/* -------------------------------------------------------------------------- */
/*                                     RLP                                    */
/* -------------------------------------------------------------------------- */

// Len returns the length of s.
func (w Withdrawals) Len() int { return len(w) }

// EncodeIndex encodes the i'th withdrawal to w. Note that this does not check
// for errors because we assume that *Withdrawal will only ever contain valid
// withdrawals that were either
// constructed by decoding or via public API in this package.
func (w Withdrawals) EncodeIndex(i int, _w *bytes.Buffer) {
	// #nosec:G703 // its okay.
	_ = w[i].EncodeRLP(_w)
}

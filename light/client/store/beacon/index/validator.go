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

package index

import (
	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
)

// Collection prefixes.
const (
	validatorByIndexPrefix                 = "val_idx_to_pk"
	validatorPubkeyToIndexPrefix           = "val_pk_to_idx"
	validatorConsAddrToIndexPrefix         = "val_cons_addr_to_idx"
	validatorEffectiveBalanceToIndexPrefix = "val_eff_bal_to_idx"
)

// ValidatorsIndex is a struct that holds a unique index for validators based
// on their public key.
type ValidatorsIndex struct {
	// Pubkey is a unique index mapping a validator's public key to their
	// numeric ID and vice versa.
	Pubkey *indexes.Unique[[]byte, uint64, *beacontypes.Validator]
	// EffectiveBalance is a multi-index mapping a validator's effective balance
	// to their numeric ID.
	EffectiveBalance *indexes.Multi[uint64, uint64, *beacontypes.Validator]
}

// IndexesList returns a list of all indexes associated with the
// validatorsIndex.
func (a ValidatorsIndex) IndexesList() []sdkcollections.Index[
	uint64, *beacontypes.Validator,
] {
	return []sdkcollections.Index[uint64, *beacontypes.Validator]{
		a.Pubkey,
		a.EffectiveBalance,
	}
}

// NewValidatorsIndex creates a new validatorsIndex with a unique index for
// validator public keys.
func NewValidatorsIndex(sb *sdkcollections.SchemaBuilder) ValidatorsIndex {
	return ValidatorsIndex{
		Pubkey: indexes.NewUnique(
			sb,
			sdkcollections.NewPrefix(validatorPubkeyToIndexPrefix),
			validatorPubkeyToIndexPrefix,
			sdkcollections.BytesKey,
			sdkcollections.Uint64Key,

			func(_ uint64, validator *beacontypes.Validator) ([]byte, error) {
				return validator.Pubkey[:], nil
			},
		),
		EffectiveBalance: indexes.NewMulti(
			sb,
			sdkcollections.NewPrefix(validatorEffectiveBalanceToIndexPrefix),
			validatorEffectiveBalanceToIndexPrefix,
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Key,
			func(_ uint64, validator *beacontypes.Validator) (uint64, error) {
				return uint64(validator.EffectiveBalance), nil
			},
		),
	}
}

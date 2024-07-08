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

package index

import (
	sdkcollections "cosmossdk.io/collections"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections/indexes"
	cmtcrypto "github.com/cometbft/cometbft/crypto"
)

// Collection prefixes.
const (
	validatorByIndexPrefix                 = "val_idx_to_pk"
	validatorPubkeyToIndexPrefix           = "val_pk_to_idx"
	validatorConsAddrToIndexPrefix         = "val_cons_addr_to_idx"
	validatorEffectiveBalanceToIndexPrefix = "val_eff_bal_to_idx"
	moduleName                             = "beacon"
)

// Validator is an interface that combines the ssz.Marshaler and
// ssz.Unmarshaler interfaces.
type Validator interface {
	// MarshalSSZTo marshals the object into the provided byte slice and returns
	// it along with any error.
	MarshalSSZTo([]byte) ([]byte, error)
	// MarshalSSZ marshals the object into a new byte slice and returns it along
	// with any error.
	MarshalSSZ() ([]byte, error)
	// UnmarshalSSZ unmarshals the object from the provided byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalSSZ([]byte) error
	// SizeSSZ returns the size in bytes that the object would take when
	// marshaled.
	SizeSSZ() int
	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetEffectiveBalance returns the effective balance of the validator.
	GetEffectiveBalance() math.Gwei
}

// ValidatorsIndex is a struct that holds a unique index for validators based
// on their public key.
type ValidatorsIndex[ValidatorT Validator] struct {
	// Pubkey is a unique index mapping a validator's public key to their
	// numeric ID and vice versa.
	Pubkey *indexes.Unique[[]byte, uint64, ValidatorT]
	// EffectiveBalance is a multi-index mapping a validator's effective balance
	// to their numeric ID.
	EffectiveBalance *indexes.Multi[uint64, uint64, ValidatorT]
	// CometBFTAddress is a unique index mapping a validator's Comet BFT address
	// to their numeric ID.
	CometBFTAddress *indexes.Unique[[]byte, uint64, ValidatorT]
}

// IndexesList returns a list of all indexes associated with the
// validatorsIndex.
func (a ValidatorsIndex[ValidatorT]) IndexesList() []collections.Index[
	uint64, ValidatorT,
] {
	return []collections.Index[uint64, ValidatorT]{
		a.Pubkey,
		a.EffectiveBalance,
		a.CometBFTAddress,
	}
}

// NewValidatorsIndex creates a new validatorsIndex with a unique index for
// validator public keys.
func NewValidatorsIndex[ValidatorT Validator](
	sa collections.StoreAccessor,
) ValidatorsIndex[ValidatorT] {
	// lol temp
	storeKey := []byte(moduleName)
	return ValidatorsIndex[ValidatorT]{
		Pubkey: indexes.NewUnique(
			storeKey,
			sdkcollections.NewPrefix(validatorPubkeyToIndexPrefix),
			sdkcollections.BytesKey,
			sdkcollections.Uint64Key,
			sa,
			func(_ uint64, validator ValidatorT) ([]byte, error) {
				pk := validator.GetPubkey()
				return pk[:], nil
			},
		),
		EffectiveBalance: indexes.NewMulti(
			storeKey,
			sdkcollections.NewPrefix(validatorEffectiveBalanceToIndexPrefix),
			sdkcollections.Uint64Key,
			sdkcollections.Uint64Key,
			sa,
			func(_ uint64, validator ValidatorT) (uint64, error) {
				return uint64(validator.GetEffectiveBalance()), nil
			},
		),
		CometBFTAddress: indexes.NewUnique(
			storeKey,
			sdkcollections.NewPrefix(validatorConsAddrToIndexPrefix),
			sdkcollections.BytesKey,
			sdkcollections.Uint64Key,
			sa,
			func(_ uint64, validator ValidatorT) ([]byte, error) {
				pk := validator.GetPubkey()
				return cmtcrypto.AddressHash(pk[:]).Bytes(), nil
			},
		),
	}
}

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

//nolint:dupl // separate file for ease of future implementation
package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	ssz "github.com/ferranbt/fastssz"
)

// Compile-time assertions to ensure AttesterSlashing implements necessary interfaces.
var (
	_ constraints.SSZMarshallableRootable = (*AttesterSlashing)(nil)
	_ common.UnusedEnforcer               = (*AttesterSlashings)(nil)
)

type (
	AttesterSlashing  = common.UnusedType
	AttesterSlashings []*AttesterSlashing
)

// SizeSSZ returns the SSZ encoded size in bytes for the AttesterSlashings.
func (ass AttesterSlashings) SizeSSZ() int {
	return 4 + len(ass)*16 // offset + each attester slashing size (UnusedType is 16 bytes)
}

// HashTreeRoot returns the hash tree root of the AttesterSlashings.
func (ass AttesterSlashings) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)
	if err := ass.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()

}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the AttesterSlashings object with a hasher.
func (ass AttesterSlashings) HashTreeRootWith(hh ssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(ass))
	if num > constants.MaxAttesterSlashings {
		return ssz.ErrIncorrectListSize
	}
	for _, elem := range ass {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(indx, num, constants.MaxAttesterSlashings)
	return nil
}

// GetTree ssz hashes the AttesterSlashings object.
func (ass AttesterSlashings) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(ass)
}

// EnforceUnused return true if the length of the AttestererSlashings is 0.
// As long as this type remains unimplemented and unvalidated by consensus,
// we must enforce that it contains no data.
func (ass AttesterSlashings) EnforceUnused() error {
	if len(ass) != 0 {
		return errors.New("AttesterSlashings must be unused")
	}
	return nil
}

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
	fastssz "github.com/ferranbt/fastssz"
)

// Compile-time assertions to ensure ProposerSlashing implements necessary interfaces.
var (
	_ constraints.SSZMarshallableRootable = (*ProposerSlashing)(nil)
	_ common.UnusedEnforcer               = (*ProposerSlashings)(nil)
)

type (
	ProposerSlashing  = common.UnusedType
	ProposerSlashings []*ProposerSlashing
)

// SizeSSZ returns the SSZ encoded size in bytes for the ProposerSlashings.
func (ps ProposerSlashings) SizeSSZ() int {
	return 4 + len(ps)*16 // offset + each proposer slashing size (UnusedType is 16 bytes)
}


// HashTreeRoot returns the hash tree root of the ProposerSlashings.
func (ps ProposerSlashings) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	ps.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the ProposerSlashings object with a hasher.
func (ps ProposerSlashings) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(ps))
	if num > constants.MaxProposerSlashings {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range ps {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(indx, num, constants.MaxProposerSlashings)
	return nil
}

// GetTree ssz hashes the ProposerSlashings object.
func (ps ProposerSlashings) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(ps)
}

// EnforceUnused return true if the length of the ProposerSlashings is 0.
// As long as this type remains unimplemented and unvalidated by consensus,
// we must enforce that it contains no data.
func (ps ProposerSlashings) EnforceUnused() error {
	if len(ps) != 0 {
		return errors.New("ProposerSlashings must be unused")
	}
	return nil
}

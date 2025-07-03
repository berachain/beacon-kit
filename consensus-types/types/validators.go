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

//nolint:dupl // False positive detected.
package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	fastssz "github.com/ferranbt/fastssz"
)

// Validators is a type alias for a SSZ list of Validator containers.
type Validators []*Validator

// SizeSSZ returns the SSZ encoded size in bytes for the Validators.
func (vs Validators) SizeSSZ() int {
	return 4 + len(vs)*121 // offset + each validator is 121 bytes
}


// HashTreeRoot returns the SSZ hash tree root for the Validators object.
func (vs Validators) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	vs.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the Validators object with a hasher.
func (vs Validators) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(vs))
	if num > constants.ValidatorsRegistryLimit {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range vs {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(indx, num, constants.ValidatorsRegistryLimit)
	return nil
}

// GetTree ssz hashes the Validators object.
func (vs Validators) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(vs)
}

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
	"github.com/karalabe/ssz"
)

// Compile-time assertions to ensure BlsToExecutionChange implements necessary interfaces.
var (
	_ ssz.StaticObject                    = (*BlsToExecutionChange)(nil)
	_ constraints.SSZMarshallableRootable = (*BlsToExecutionChange)(nil)
	_ common.UnusedEnforcer               = (*BlsToExecutionChanges)(nil)
)

type (
	BlsToExecutionChange  = common.UnusedType
	BlsToExecutionChanges []*BlsToExecutionChange
)

// SizeSSZ returns the SSZ encoded size in bytes for the BlsToExecutionChanges.
func (bs BlsToExecutionChanges) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, bs)
}

// DefineSSZ defines the SSZ encoding for the BlsToExecutionChanges object.
func (bs BlsToExecutionChanges) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*BlsToExecutionChange)(&bs), constants.MaxBlsToExecutionChanges)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*BlsToExecutionChange)(&bs), constants.MaxBlsToExecutionChanges)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*BlsToExecutionChange)(&bs), constants.MaxBlsToExecutionChanges)
	})
}

// HashTreeRoot returns the hash tree root of the BlsToExecutionChanges.
func (bs BlsToExecutionChanges) HashTreeRoot() common.Root {
	return ssz.HashSequential(bs)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the BlsToExecutionChanges object with a hasher.
func (bs BlsToExecutionChanges) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()
	num := uint64(len(bs))
	if num > constants.MaxBlsToExecutionChanges {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range bs {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(indx, num, constants.MaxBlsToExecutionChanges)
	return nil
}

// GetTree ssz hashes the BlsToExecutionChanges object.
func (bs BlsToExecutionChanges) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(bs)
}

// EnforceUnused return true if the length of the BlsToExecutionChanges is 0.
// As long as this type remains unimplemented and unvalidated by consensus,
// we must enforce that it contains no data.
func (bs BlsToExecutionChanges) EnforceUnused() error {
	if len(bs) != 0 {
		return errors.New("BlsToExecutionChanges must be unused")
	}
	return nil
}

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

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// SSZListBasic is a list of basic types.
type SSZListBasic[T Basic[T]] []T

/* -------------------------------------------------------------------------- */
/*                                    Size                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the list in bytes.
func (l SSZListBasic[T]) SizeSSZ() int {
	// The same for SSZListBasic as for SSZVectorBasic.
	return SSZVectorBasic[T](l).SizeSSZ()
}

/* -------------------------------------------------------------------------- */
/*                                HashTreeRoot                                */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith returns the Merkle root of the SSZListBasic
// with a given merkleizer.
func (l SSZListBasic[T]) HashTreeRootWith(
	merkleizer Merkleizer[common.ChainSpec, [32]byte, T],
) ([32]byte, error) {
	return merkleizer.MerkleizeListBasic(l)
}

// HashTreeRoot returns the Merkle root of the SSZListBasic.
func (l SSZListBasic[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(ssz.NewMerkleizer[
		common.ChainSpec, [32]byte, T,
	]())
}

/* -------------------------------------------------------------------------- */
/*                                 Marshalling                                */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo marshals the SSZListBasic into SSZ format.
func (l SSZListBasic[T]) MarshalSSZTo(out []byte) ([]byte, error) {
	// The same for SSZListBasic as for SSZVectorBasic.
	return SSZVectorBasic[T](l).MarshalSSZTo(out)
}

// MarshalSSZ marshals the SSZListBasic into SSZ format.
func (l SSZListBasic[T]) MarshalSSZ() ([]byte, error) {
	// The same for SSZListBasic as for SSZVectorBasic.
	return SSZVectorBasic[T](l).MarshalSSZ()
}

// NewFromSSZ creates a new SSZListBasic from SSZ format.
func (SSZListBasic[T]) NewFromSSZ(buf []byte) (SSZListBasic[T], error) {
	// The same for SSZListBasic as for SSZVectorBasic
	var (
		t   SSZVectorBasic[T]
		err error
	)
	t, err = t.NewFromSSZ(buf)
	return SSZListBasic[T](t), err
}

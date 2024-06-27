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
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// SSZVectorBasic is a vector of basic types.
type SSZVectorBasic[T Basic[T]] []T

// VectorBasicFromElements creates a new SSZListComposite from elements.
// TODO: Deprecate once off of FastSSZTypes
func VectorBasicFromElements[T Basic[T]](elements ...T) SSZVectorBasic[T] {
	return elements
}

/* -------------------------------------------------------------------------- */
/*                                    Size                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the list in bytes.
func (l SSZVectorBasic[T]) SizeSSZ() int {
	var t T
	return t.SizeSSZ() * len(l)
}

/* -------------------------------------------------------------------------- */
/*                                HashTreeRoot                                */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith returns the Merkle root of the SSZVectorBasic
// with a given merkleizer.
func (l SSZVectorBasic[T]) HashTreeRootWith(
	merkleizer Merkleizer[common.ChainSpec, [32]byte, T],
) ([32]byte, error) {
	return merkleizer.MerkleizeVecBasic(l)
}

// HashTreeRoot returns the Merkle root of the SSZVectorBasic.
func (l SSZVectorBasic[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(ssz.NewMerkleizer[
		common.ChainSpec, [32]byte, T,
	]())
}

/* -------------------------------------------------------------------------- */
/*                                 Marshalling                                */
/* -------------------------------------------------------------------------- */

// MarshalSSZToBytes marshals the SSZVectorBasic into SSZ format.
func (l SSZVectorBasic[T]) MarshalSSZTo(out []byte) ([]byte, error) {
	// From the Spec:
	// fixed_parts = [
	// 		serialize(element)
	// 			if not is_variable_size(element)
	//			else None for element in value,
	// 		]
	// VectorBasic has all fixed types, so we simply
	// serialize each element and pack them together.
	for _, v := range l {
		bytes, err := v.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		out = append(out, bytes...)
	}
	return out, nil
}

// MarshalSSZ marshals the SSZVectorBasic into SSZ format.
func (l SSZVectorBasic[T]) MarshalSSZ() ([]byte, error) {
	return l.MarshalSSZTo(make([]byte, 0, l.SizeSSZ()))
}

// NewFromSSZ creates a new SSZVectorBasic from SSZ format.
func (SSZVectorBasic[T]) NewFromSSZ(buf []byte) (SSZVectorBasic[T], error) {
	var (
		err error
		t   T
	)
	elementSize := t.SizeSSZ()
	if len(buf)%elementSize != 0 {
		return nil, fmt.Errorf(
			"invalid buffer length %d for element size %d",
			len(buf),
			elementSize,
		)
	}

	result := make(SSZVectorBasic[T], 0, len(buf)/elementSize)
	for i := 0; i < len(buf); i += elementSize {
		if t, err = t.NewFromSSZ(buf[i : i+elementSize]); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

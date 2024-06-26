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

// SizeSSZ returns the size of the list in bytes.
func (l SSZVectorBasic[T]) SizeSSZ() int {
	var t T
	return t.SizeSSZ() * len(l)
}

// HashTreeRootWith returns the Merkle root of the SSZVectorBasic
// with a given merkleizer.
func (l SSZVectorBasic[T]) HashTreeRootWith(
	merkleizer interface {
		MerkleizeByteSlice([]byte) ([32]byte, error)
	},
) ([32]byte, error) {
	packedBytes, err := l.MarshalSSZ()
	if err != nil {
		return [32]byte{}, err
	}
	return merkleizer.MerkleizeByteSlice(packedBytes)
}

// HashTreeRoot returns the Merkle root of the SSZVectorBasic.
func (l SSZVectorBasic[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(ssz.NewMerkleizer[
		common.ChainSpec, [32]byte, common.Root,
	]())
}

// MarshalSSZ marshals the SSZVectorBasic into SSZ format.
func (l SSZVectorBasic[T]) MarshalSSZ() ([]byte, error) {
	packedBytes := make([]byte, 0, l.SizeSSZ())
	for _, v := range l {
		bytes, err := v.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		packedBytes = append(packedBytes, bytes...)
	}
	return packedBytes, nil
}

// UnmarshalSSZ unmarshals the SSZVectorBasic from SSZ format.
func (l *SSZVectorBasic[T]) UnmarshalSSZ(buf []byte) error {
	var (
		err error
		t   T
	)
	elementSize := t.SizeSSZ()
	if len(buf)%elementSize != 0 {
		return fmt.Errorf(
			"invalid buffer length %d for element size %d",
			len(buf),
			elementSize,
		)
	}

	if l == nil {
		l = new(SSZVectorBasic[T])
		*l = make([]T, 0, len(buf)/elementSize)
	}

	for i := 0; i < len(buf); i += elementSize {
		if t, err = t.NewFromSSZ(buf[i : i+elementSize]); err != nil {
			return err
		}
		*l = append(*l, t)
	}

	return nil
}

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
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

type SSZMarshallable interface {
	SizeSSZ() int
}

type Basic interface {
	~bool | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	MarshalSSZ() ([]byte, error)
}

type SSZVectorBasic[T Basic] []T

// SizeSSZ returns the size of the list in bytes.
func (l SSZVectorBasic[T]) SizeSSZ() int {
	elementSize := reflect.TypeOf((*T)(nil)).Elem().Size()
	fmt.Println(elementSize)
	return int(elementSize) * len(l)
}

// HashTreeRoot returns the Merkle root of the SSZVectorBasic.
func (l SSZVectorBasic[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	m := ssz.NewMerkleizer[
		common.ChainSpec, [32]byte, common.Root,
	]()
	packedBytes := make([]byte, l.SizeSSZ())
	for _, v := range l {
		v, err := v.MarshalSSZ()
		if err != nil {
			return [32]byte{}, err
		}

		packedBytes = append(packedBytes, v...)
	}
	return m.MerkleizeByteSlice(packedBytes)
}

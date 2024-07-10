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

package compare_test

import (
	"testing"
	"testing/quick"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	ztyp "github.com/protolambda/ztyp/tree"
)

var hFn = ztyp.GetHashFn()

func TestListHashTreeRootZtyp(t *testing.T) {
	f := func(s []byte, limit uint64) bool {
		a := ssz.ByteListFromBytes(s, limit)

		root1, err := a.HashTreeRoot()
		root2 := hFn.ByteListHTR(s, limit)
		if err != nil {
			return false
		}
		return root1 == [32]byte(root2)
	}
	c := quick.Config{MaxCount: 1000000}
	if err := quick.Check(f, &c); err != nil {
		t.Error(err)
	}
}

func TestVectorHashTreeRootZTyp(t *testing.T) {
	f := func(s []byte) bool {
		a := ssz.ByteVectorFromBytes(s)

		root1, err := a.HashTreeRoot()
		root2 := hFn.ByteVectorHTR(s)
		if err != nil {
			return false
		}
		return root1 == [32]byte(root2)
	}
	c := quick.Config{MaxCount: 1000000}
	if err := quick.Check(f, &c); err != nil {
		t.Error(err)
	}
}

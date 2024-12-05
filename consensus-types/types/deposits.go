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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/karalabe/ssz"
)

// Deposits is a typealias for a list of Deposits.
type Deposits []*Deposit

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (ds Deposits) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ([]*Deposit)(ds))
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
func (ds Deposits) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			c, (*[]*Deposit)(&ds), constants.MaxDepositsPerBlock)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			c, (*[]*Deposit)(&ds), constants.MaxDepositsPerBlock)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(
			c, (*[]*Deposit)(&ds), constants.MaxDepositsPerBlock)
	})
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (ds Deposits) HashTreeRoot() common.Root {
	return ssz.HashSequential(ds)
}

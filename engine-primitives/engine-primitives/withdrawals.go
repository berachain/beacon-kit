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

package engineprimitives

import (
	"bytes"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
)

var (
	_ ssz.StaticObject        = (*Withdrawals)(nil)
	_ constraints.SSZRootable = (*Withdrawals)(nil)
)

// Withdrawals represents a list of withdrawals.
type Withdrawals []*Withdrawal

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Withdrawals.
func (w Withdrawals) SizeSSZ(siz *ssz.Sizer) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, w)
}

// DefineSSZ defines the SSZ encoding for the Withdrawals object.
func (w Withdrawals) DefineSSZ(codec *ssz.Codec) {
	codec.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			codec, (*[]*Withdrawal)(&w), constants.MaxWithdrawalsPerPayload)
	})
	codec.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			codec, (*[]*Withdrawal)(&w), constants.MaxWithdrawalsPerPayload)
	})
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(
			codec, (*[]*Withdrawal)(&w), constants.MaxWithdrawalsPerPayload)
	})
}

// HashTreeRoot returns the hash tree root of the Withdrawals.
func (w Withdrawals) HashTreeRoot() common.Root {
	return ssz.HashSequential(w)
}

/* -------------------------------------------------------------------------- */
/*                                     RLP                                    */
/* -------------------------------------------------------------------------- */

// Len returns the length of s.
func (w Withdrawals) Len() int { return len(w) }

// EncodeIndex encodes the i'th withdrawal to w. Note that this does not check
// for errors because we assume that *Withdrawal will only ever contain valid
// withdrawals that were either
// constructed by decoding or via public API in this package.
func (w Withdrawals) EncodeIndex(i int, _w *bytes.Buffer) {
	// #nosec:G703 // its okay.
	_ = w[i].EncodeRLP(_w)
}

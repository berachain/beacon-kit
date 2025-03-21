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
	"github.com/karalabe/ssz"
)

// Compile-time assertions to ensure VoluntaryExit implements necessary interfaces.
var (
	_ ssz.StaticObject                    = (*VoluntaryExit)(nil)
	_ constraints.SSZMarshallableRootable = (*VoluntaryExit)(nil)
	_ common.UnusedEnforcer               = (*VoluntaryExits)(nil)
)

type (
	VoluntaryExit  = common.UnusedType
	VoluntaryExits []*VoluntaryExit
)

// SizeSSZ returns the SSZ encoded size in bytes for the VoluntaryExits.
func (vs VoluntaryExits) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, vs)
}

// DefineSSZ defines the SSZ encoding for the VoluntaryExits object.
func (vs VoluntaryExits) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*VoluntaryExit)(&vs), constants.MaxVoluntaryExits)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*VoluntaryExit)(&vs), constants.MaxVoluntaryExits)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*VoluntaryExit)(&vs), constants.MaxVoluntaryExits)
	})
}

// HashTreeRoot returns the hash tree root of the VoluntaryExits.
func (vs VoluntaryExits) HashTreeRoot() common.Root {
	return ssz.HashSequential(vs)
}

// EnforceUnused return true if the length of the VoluntaryExits is 0.
// As long as this type remains unimplemented and unvalidated by consensus,
// we must enforce that it contains no data.
func (vs VoluntaryExits) EnforceUnused() error {
	if len(vs) != 0 {
		return errors.New("VoluntaryExits must be unused")
	}
	return nil
}

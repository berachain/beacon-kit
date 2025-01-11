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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
)

// Compile-time assertions to ensure ProposerSlashing implements necessary interfaces.
var (
	_ ssz.StaticObject                    = (*ProposerSlashing)(nil)
	_ constraints.SSZMarshallableRootable = (*ProposerSlashing)(nil)
	_ UnusedEnforcer                      = (*ProposerSlashings)(nil)
)

type (
	ProposerSlashing  = UnusedType
	ProposerSlashings []*ProposerSlashing
)

// SizeSSZ returns the SSZ encoded size in bytes for the ProposerSlashings.
func (ps ProposerSlashings) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ps)
}

// DefineSSZ defines the SSZ encoding for the ProposerSlashings object.
func (ps ProposerSlashings) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*ProposerSlashing)(&ps), constants.MaxProposerSlashings)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*ProposerSlashing)(&ps), constants.MaxProposerSlashings)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*ProposerSlashing)(&ps), constants.MaxProposerSlashings)
	})
}

// HashTreeRoot returns the hash tree root of the ProposerSlashings.
func (ps ProposerSlashings) HashTreeRoot() common.Root {
	return ssz.HashSequential(ps)
}

// EnforceUnused return true if the length of the ProposerSlashings is 0.
// As long as this type remains unimplemented and unvalidated by consensus,
// we must enforce that it contains no data.
func (ps ProposerSlashings) EnforceUnused() error {
	if len(ps) != 0 {
		return errors.New("ProposerSlashings must be unused")
	}
	return nil
}

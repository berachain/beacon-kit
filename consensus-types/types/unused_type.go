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

package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
)

// UnusedEnforcer is an interface that asserts that a type is unused.
type UnusedEnforcer interface {
	EnforceUnused() error
}

// Compile-time assertions to ensure UnusedType implements necessary interfaces.
var (
	_ ssz.StaticObject                                 = (*UnusedType)(nil)
	_ constraints.SSZMarshallableRootable[*UnusedType] = (*UnusedType)(nil)
	_ UnusedEnforcer                                   = (*UnusedType)(nil)
)

type UnusedType uint8

// SizeSSZ returns the SSZ encoded size in bytes for the UnusedType.
func (ut *UnusedType) SizeSSZ(_ *ssz.Sizer) uint32 {
	return 1
}

// DefineSSZ defines the SSZ encoding for the UnusedType object.
func (ut *UnusedType) DefineSSZ(c *ssz.Codec) {
	ssz.DefineUint8(c, ut)
}

// MarshalSSZ marshals the UnusedType object to SSZ format.
func (ut *UnusedType) MarshalSSZ() ([]byte, error) {
	return []byte{uint8(*ut)}, nil
}

func (*UnusedType) IsUnusedFromSZZ() bool      { return true }
func (*UnusedType) EnsureSyntaxFromSSZ() error { return nil }

// HashTreeRoot returns the hash tree root of the Deposits.
func (ut *UnusedType) HashTreeRoot() common.Root {
	return ssz.HashSequential(ut)
}

// EnforceUnused return true if the UnusedType contains all zero values.
// As long as this type remains unused and unvalidated by consensus,
// we must enforce that it contains no data.
func (ut *UnusedType) EnforceUnused() error {
	if *ut != UnusedType(0) {
		return errors.New("UnusedType must be unused")
	}
	return nil
}

func EnforceAllUnused(enforcers ...UnusedEnforcer) error {
	for _, enforcer := range enforcers {
		if err := enforcer.EnforceUnused(); err != nil {
			return err
		}
	}
	return nil
}

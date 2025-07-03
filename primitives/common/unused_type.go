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

package common

import (
	"github.com/berachain/beacon-kit/errors"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// UnusedEnforcer is an interface that asserts that a type is unused.
type UnusedEnforcer interface {
	EnforceUnused() error
}

// Compile-time assertions to ensure UnusedType implements necessary interfaces.
var (
	_ ssz.StaticObject        = (*UnusedType)(nil)
	_ SSZMarshallableRootable = (*UnusedType)(nil)
	_ UnusedEnforcer          = (*UnusedType)(nil)
	// Note: fastssz interfaces are not enforced due to method signature conflicts during migration
)

type UnusedType uint8

// SizeSSZInternal returns the SSZ encoded size in bytes for the UnusedType.
// This method is kept for compatibility with karalabe/ssz.
func (ut *UnusedType) SizeSSZ(_ *ssz.Sizer) uint32 {
	return 1
}

// DefineSSZ defines the SSZ encoding for the UnusedType object.
func (ut *UnusedType) DefineSSZ(c *ssz.Codec) {
	ssz.DefineUint8(c, ut)
}

// MarshalSSZ marshals the UnusedType object to SSZ format.
func (ut *UnusedType) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(ut))
	return buf, ssz.EncodeToBytes(buf, ut)
}

func (ut *UnusedType) ValidateAfterDecodingSSZ() error {
	return ut.EnforceUnused()
}

// HashTreeRoot returns the hash tree root of the UnusedType.
func (ut *UnusedType) HashTreeRoot() Root {
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

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the UnusedType object to a target array.
func (ut *UnusedType) MarshalSSZTo(dst []byte) ([]byte, error) {
	dst = append(dst, byte(*ut))
	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the UnusedType object.
func (ut *UnusedType) UnmarshalSSZ(buf []byte) error {
	if len(buf) != 1 {
		return errors.Wrapf(fastssz.ErrSize, "expected buffer of length 1, received %d", len(buf))
	}
	*ut = UnusedType(buf[0])
	// Validate after unmarshaling
	return ut.EnforceUnused()
}

// SizeSSZ returns the ssz encoded size in bytes for the UnusedType (fastssz).
func (ut *UnusedType) SizeSSZFastSSZ() (size int) {
	return 1
}

// HashTreeRoot calculates the SSZ hash tree root of the UnusedType.
// This method satisfies the fastssz.HashRoot interface.
func (ut *UnusedType) HashTreeRootFastSSZ() ([32]byte, error) {
	// UnusedType is a uint8, so its hash tree root is the value padded to 32 bytes
	var root [32]byte
	root[0] = byte(*ut)
	return root, nil
}

// HashTreeRootWith ssz hashes the UnusedType object with a hasher.
func (ut *UnusedType) HashTreeRootWith(hh fastssz.HashWalker) error {
	hh.PutUint8(uint8(*ut))
	return nil
}

// GetTree ssz hashes the UnusedType object.
func (ut *UnusedType) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(ut)
}

func EnforceAllUnused(enforcers ...UnusedEnforcer) error {
	for _, enforcer := range enforcers {
		if err := enforcer.EnforceUnused(); err != nil {
			return err
		}
	}
	return nil
}

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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// Compile-time assertions to ensure SignedBeaconBlockHeader implements necessary interfaces.
var (
	_ ssz.StaticObject                    = (*SignedBeaconBlockHeader)(nil)
	_ constraints.SSZMarshallableRootable = (*SignedBeaconBlockHeader)(nil)
)

// SignedBeaconBlockHeader is a struct that contains a BeaconBlockHeader and a BLSSignature.
//
// NOTE: This struct is only ever (un)marshalled with SSZ and NOT with JSON.
type SignedBeaconBlockHeader struct {
	Header    *BeaconBlockHeader
	Signature crypto.BLSSignature
}

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

// NewSignedBeaconBlockHeader creates a new BeaconBlockHeader.
func NewSignedBeaconBlockHeader(
	header *BeaconBlockHeader,
	signature crypto.BLSSignature,
) *SignedBeaconBlockHeader {
	return &SignedBeaconBlockHeader{
		header, signature,
	}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the SignedBeaconBlockHeader object
// in SSZ encoding. Total size: Header (112) + Signature (96).
func (b *SignedBeaconBlockHeader) SizeSSZ(sizer *ssz.Sizer) uint32 {
	//nolint:mnd // no magic
	size := (*BeaconBlockHeader)(nil).SizeSSZ(sizer) + 96
	return size
}

// DefineSSZ defines the SSZ encoding for the SignedBeaconBlockHeader object.
func (b *SignedBeaconBlockHeader) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticObject(codec, &b.Header)
	ssz.DefineStaticBytes(codec, &b.Signature)
}

// MarshalSSZ marshals the SignedBeaconBlockHeader object to SSZ format.
func (b *SignedBeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

func (*SignedBeaconBlockHeader) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot computes the SSZ hash tree root of the
// SignedBeaconBlockHeader object.
func (b *SignedBeaconBlockHeader) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}

/* -------------------------------------------------------------------------- */
/*                            Getters and Setters                             */
/* -------------------------------------------------------------------------- */

// Getheader retrieves the header of the SignedBeaconBlockHeader.
func (b *SignedBeaconBlockHeader) GetHeader() *BeaconBlockHeader {
	return b.Header
}

// GetSignature retrieves the Signature of the SignedBeaconBlockHeader.
func (b *SignedBeaconBlockHeader) GetSignature() crypto.BLSSignature {
	return b.Signature
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the SignedBeaconBlockHeader object to a target array.
func (b *SignedBeaconBlockHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := b.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the SignedBeaconBlockHeader object.
func (b *SignedBeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	// For now, delegate to karalabe/ssz for unmarshaling
	return ssz.DecodeFromBytes(buf, b)
}

// SizeSSZFastSSZ returns the ssz encoded size in bytes for the SignedBeaconBlockHeader (fastssz).
// TODO: Rename to SizeSSZ() once karalabe/ssz is fully removed.
func (b *SignedBeaconBlockHeader) SizeSSZFastSSZ() (size int) {
	// Use the existing karalabe/ssz Size function to get the size
	size = int(ssz.Size(b))
	return
}

// HashTreeRootWith ssz hashes the SignedBeaconBlockHeader object with a hasher.
func (b *SignedBeaconBlockHeader) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Header'
	if err := b.Header.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (1) 'Signature'
	hh.PutBytes(b.Signature[:])

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the SignedBeaconBlockHeader object.
func (b *SignedBeaconBlockHeader) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}

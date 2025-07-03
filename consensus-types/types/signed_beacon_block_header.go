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
	"github.com/berachain/beacon-kit/primitives/crypto"
	fastssz "github.com/ferranbt/fastssz"
)

// TODO: Re-enable interface assertion once constraints are updated
// var (
// 	_ constraints.SSZMarshallableRootable = (*SignedBeaconBlockHeader)(nil)
// )

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
func (b *SignedBeaconBlockHeader) SizeSSZ() int {
	return 112 + 96 // BeaconBlockHeaderSize + SignatureSize
}


// MarshalSSZ marshals the SignedBeaconBlockHeader object to SSZ format.
func (b *SignedBeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, b.SizeSSZ())
	return b.MarshalSSZTo(buf)
}

func (*SignedBeaconBlockHeader) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot computes the SSZ hash tree root of the
// SignedBeaconBlockHeader object.
func (b *SignedBeaconBlockHeader) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	b.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
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
	// Header
	dst, err := b.Header.MarshalSSZTo(dst)
	if err != nil {
		return nil, err
	}

	// Signature
	dst = append(dst, b.Signature[:]...)

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the SignedBeaconBlockHeader object.
func (b *SignedBeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	if len(buf) != 208 { // 112 + 96
		return fastssz.ErrSize
	}

	// Header
	if b.Header == nil {
		b.Header = &BeaconBlockHeader{}
	}
	if err := b.Header.UnmarshalSSZ(buf[0:112]); err != nil {
		return err
	}

	// Signature
	copy(b.Signature[:], buf[112:208])

	return nil
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

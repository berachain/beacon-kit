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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/karalabe/ssz"
)

// SignedBeaconBlockHeaderSize is the size of the SignedBeaconBlockHeader object in bytes.
//
// Total size: header (112) + Signature (96)
const SignedBeaconBlockHeaderSize = 208

var (
	_ ssz.StaticObject                    = (*SignedBeaconBlockHeader)(nil)
	_ constraints.SSZMarshallableRootable = (*SignedBeaconBlockHeader)(nil)
)

type SignedBeaconBlockHeader struct {
	Header    *BeaconBlockHeader  `json:"header"`
	Signature crypto.BLSSignature `json:"signature"`
}

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

// NewBeaconBlockHeader creates a new BeaconBlockHeader.
func NewSignedBeaconBlockHeader(
	header *BeaconBlockHeader,
	signature crypto.BLSSignature,
) *SignedBeaconBlockHeader {
	return &SignedBeaconBlockHeader{
		header, signature,
	}
}

// Empty creates an empty BeaconBlockHeader instance.
func (*SignedBeaconBlockHeader) Empty() *SignedBeaconBlockHeader {
	return &SignedBeaconBlockHeader{}
}

// New creates a new SignedBeaconBlockHeader.
func (b *SignedBeaconBlockHeader) New(
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

// SizeSSZ returns the size of the SignedBeaconBlockHeader object in SSZ encoding.
func (b *SignedBeaconBlockHeader) SizeSSZ(*ssz.Sizer) uint32 {
	return SignedBeaconBlockHeaderSize
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

// UnmarshalSSZ unmarshals the SignedBeaconBlockHeader object from SSZ format.
func (b *SignedBeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the SignedBeaconBlockHeader object.
func (b *SignedBeaconBlockHeader) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// TODO

/* -------------------------------------------------------------------------- */
/*                            Getters and Setters                             */
/* -------------------------------------------------------------------------- */

// Getheader retrieves the header of the SignedBeaconBlockHeader.
func (b *SignedBeaconBlockHeader) GetHeader() *BeaconBlockHeader {
	return b.Header
}

// Setheader sets the header of the BeaconBlockHeader.
func (b *SignedBeaconBlockHeader) SetHeader(header *BeaconBlockHeader) {
	b.Header = header
}

// GetSignature retrieves the Signature of the SignedBeaconBlockHeader.
func (b *SignedBeaconBlockHeader) GetSignature() crypto.BLSSignature {
	return b.Signature
}

// SetSignature sets the Signature of the BeaconBlockHeader.
func (b *SignedBeaconBlockHeader) SetSignature(signature crypto.BLSSignature) {
	b.Signature = signature
}

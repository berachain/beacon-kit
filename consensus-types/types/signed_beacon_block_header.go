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
	"github.com/karalabe/ssz"
)

// Compile-time assertions to ensure SignedBeaconBlockHeader implements necessary interfaces.
var (
	_ ssz.StaticObject                                              = (*SignedBeaconBlockHeader)(nil)
	_ constraints.SSZMarshallableRootable[*SignedBeaconBlockHeader] = (*SignedBeaconBlockHeader)(nil)
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

// NewFromSSZ creates a new SignedBeaconBlockHeader from SSZ format.
func (*SignedBeaconBlockHeader) NewFromSSZ(buf []byte) (*SignedBeaconBlockHeader, error) {
	b := &SignedBeaconBlockHeader{}
	return b, ssz.DecodeFromBytes(buf, b)
}

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

func (b *SignedBeaconBlockHeader) IsNil() bool {
	return b == nil
}

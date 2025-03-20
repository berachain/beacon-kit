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
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/karalabe/ssz"
)

// Compile-time assertions to ensure SignedBeaconBlock implements necessary interfaces.
var (
	_ ssz.DynamicObject                            = (*SignedBeaconBlock)(nil)
	_ constraints.SSZVersionedMarshallableRootable = (*SignedBeaconBlock)(nil)
)

// SignedBeaconBlock is a struct that contains a BeaconBlock and a BLSSignature.
//
// NOTE: This struct is only ever (un)marshalled with SSZ and NOT with JSON.
type SignedBeaconBlock struct {
	*BeaconBlock
	Signature crypto.BLSSignature
}

/* -------------------------------------------------------------------------- */
/*                                 Constructors                               */
/* -------------------------------------------------------------------------- */

// NewSignedBeaconBlock signs the provided BeaconBlock and populates the receiver.
//
// NOTE: will panic if any provided argument is nil. Only errors if signing fails.
func NewSignedBeaconBlock(
	blk *BeaconBlock, forkData *ForkData, cs ProposerDomain, signer crypto.BLSSigner,
) (*SignedBeaconBlock, error) {
	domain := forkData.ComputeDomain(cs.DomainTypeProposer())
	signingRoot := ComputeSigningRoot(blk, domain)
	signature, err := signer.Sign(signingRoot[:])
	if err != nil {
		return nil, err
	}

	return &SignedBeaconBlock{
		BeaconBlock: blk,
		Signature:   signature,
	}, nil
}

func NewEmptySignedBeaconBlockWithVersion(forkVersion common.Version) (*SignedBeaconBlock, error) {
	switch forkVersion {
	case version.Deneb(), version.Deneb1():
		return &SignedBeaconBlock{
			BeaconBlock: NewEmptyBeaconBlockWithVersion(forkVersion),
		}, nil
	default:
		// We return a non-nil block here to appease nilaway.
		return nil, errors.Wrapf(ErrForkVersionNotSupported, "fork %d", forkVersion)
	}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the SignedBeaconBlockHeader object
// in SSZ encoding.
// Total size: MessageOffset (4) + Signature (96) + MessageContentDynamic.
func (b *SignedBeaconBlock) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	var size = uint32(constants.SSZOffsetSize + bytes.B96Size)
	if fixed {
		return size
	}
	size += ssz.SizeDynamicObject(siz, b.BeaconBlock)
	return size
}

// DefineSSZ defines the SSZ encoding for the SignedBeaconBlockHeader object.
func (b *SignedBeaconBlock) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineDynamicObjectOffset(codec, &b.BeaconBlock)
	ssz.DefineStaticBytes(codec, &b.Signature)

	// Define the dynamic data (fields)
	ssz.DefineDynamicObjectContent(codec, &b.BeaconBlock)
}

// MarshalSSZ marshals the SignedBeaconBlockHeader object to SSZ format.
func (b *SignedBeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

func (b *SignedBeaconBlock) EnsureSyntaxFromSSZ() error {
	return b.BeaconBlock.EnsureSyntaxFromSSZ()
}

// HashTreeRoot computes the SSZ hash tree root of the
// SignedBeaconBlockHeader object.
func (b *SignedBeaconBlock) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}

/* -------------------------------------------------------------------------- */
/*                                 Getters                                    */
/* -------------------------------------------------------------------------- */

func (b *SignedBeaconBlock) GetBeaconBlock() *BeaconBlock {
	return b.BeaconBlock
}

func (b *SignedBeaconBlock) GetSignature() crypto.BLSSignature {
	return b.Signature
}

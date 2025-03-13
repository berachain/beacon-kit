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
	"fmt"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/karalabe/ssz"
)

var (
	_ ssz.DynamicObject                   = (*SignedBeaconBlock)(nil)
	_ constraints.SSZMarshallableRootable = (*SignedBeaconBlock)(nil)
)

type SignedBeaconBlock struct {
	Message   *BeaconBlock        `json:"message"`
	Signature crypto.BLSSignature `json:"signature"`
}

/* -------------------------------------------------------------------------- */
/*                                 Constructors                               */
/* -------------------------------------------------------------------------- */

// NewSignedBeaconBlockFromSSZ creates a new beacon block from the given SSZ bytes.
func NewSignedBeaconBlockFromSSZ(
	bz []byte,
	forkVersion common.Version,
) (*SignedBeaconBlock, error) {
	block := &SignedBeaconBlock{}
	switch forkVersion {
	case version.Deneb(), version.Deneb1():
		if err := block.UnmarshalSSZ(bz); err != nil {
			return block, err
		}

		// make sure Withdrawals in execution payload are not nil
		EnsureNotNilWithdrawals(block.Message.Body.ExecutionPayload)

		// duly setup fork version in every relevant block member
		block.Message.forkVersion = forkVersion
		block.Message.Body.Versionable = block.Message
		block.Message.Body.ExecutionPayload.Versionable = block.Message
		return block, nil
	default:
		// we return block here to appease nilaway
		return block, errors.Wrap(
			ErrForkVersionNotSupported,
			fmt.Sprintf("fork %d", forkVersion),
		)
	}
}

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
		Message:   blk,
		Signature: signature,
	}, nil
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
	size += ssz.SizeDynamicObject(siz, b.Message)
	return size
}

// DefineSSZ defines the SSZ encoding for the SignedBeaconBlockHeader object.
func (b *SignedBeaconBlock) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineDynamicObjectOffset(codec, &b.Message)
	ssz.DefineStaticBytes(codec, &b.Signature)

	// Define the dynamic data (fields)
	ssz.DefineDynamicObjectContent(codec, &b.Message)
}

// MarshalSSZ marshals the SignedBeaconBlockHeader object to SSZ format.
func (b *SignedBeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ unmarshals the SignedBeaconBlockHeader object from SSZ format.
func (b *SignedBeaconBlock) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the
// SignedBeaconBlockHeader object.
func (b *SignedBeaconBlock) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}

/* -------------------------------------------------------------------------- */
/*                            Getters and Setters                             */
/* -------------------------------------------------------------------------- */

func (b *SignedBeaconBlock) GetMessage() *BeaconBlock {
	return b.Message
}

func (b *SignedBeaconBlock) SetHeader(message *BeaconBlock) {
	b.Message = message
}

func (b *SignedBeaconBlock) GetSignature() crypto.BLSSignature {
	return b.Signature
}

func (b *SignedBeaconBlock) SetSignature(signature crypto.BLSSignature) {
	b.Signature = signature
}

func (b *SignedBeaconBlock) IsNil() bool {
	return b == nil
}

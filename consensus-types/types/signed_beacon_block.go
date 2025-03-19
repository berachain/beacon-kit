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

var (
	_ ssz.DynamicObject                                                = (*SignedBeaconBlock)(nil)
	_ constraints.SSZVersionedMarshallableRootable[*SignedBeaconBlock] = (*SignedBeaconBlock)(nil)
)

type SignedBeaconBlock struct {
	*BeaconBlock `json:"message"`
	Signature    crypto.BLSSignature `json:"signature"`
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
		var err error
		block, err = block.NewFromSSZ(bz, forkVersion)
		if err != nil {
			return nil, err
		}

		// Make sure Withdrawals in execution payload are not nil.
		block.Body.ExecutionPayload.EnsureNotNilWithdrawals()

		return block, nil
	case version.Electra():
		var err error
		block, err = block.NewFromSSZ(bz, forkVersion)
		if err != nil {
			return nil, err
		}
		// TODO(REZ): Come back here and add decoding
		blockBody := block.GetBody()
		// Make sure Withdrawals in execution payload are not nil.
		blockBody.GetExecutionPayload().EnsureNotNilWithdrawals()
		requests, err := blockBody.GetExecutionRequests()
		if err != nil {
			return nil, err
		}
		if requests == nil {
			return nil, errors.New("execution requests was nil")
		}
		return block, nil
	default:
		// We return block here to appease nilaway.
		return block, errors.Wrapf(ErrForkVersionNotSupported, "fork %d", forkVersion)
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
		BeaconBlock: blk,
		Signature:   signature,
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

// empty creates a new SignedBeaconBlock with empty values.
func (*SignedBeaconBlock) empty(version common.Version) *SignedBeaconBlock {
	return &SignedBeaconBlock{
		BeaconBlock: (&BeaconBlock{}).empty(version),
	}
}

// NewFromSSZ creates a new SignedBeaconBlock from SSZ format.
func (*SignedBeaconBlock) NewFromSSZ(
	buf []byte, version common.Version,
) (*SignedBeaconBlock, error) {
	b := (&SignedBeaconBlock{}).empty(version)
	return b, ssz.DecodeFromBytes(buf, b)
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

func (b *SignedBeaconBlock) IsNil() bool {
	return b == nil
}

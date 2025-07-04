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
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/version"
	fastssz "github.com/ferranbt/fastssz"
)

// TODO: Re-enable interface assertion once constraints are updated
// var (
// 	_ constraints.SSZVersionedMarshallableRootable = (*SignedBeaconBlock)(nil)
// )

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
	case version.Deneb(), version.Deneb1(), version.Electra(), version.Electra1():
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

// SizeSSZ returns the size of the SignedBeaconBlock object in SSZ encoding.
// Total size: MessageOffset (4) + Signature (96) + MessageContentDynamic.
func (b *SignedBeaconBlock) SizeSSZ() int {
	return 4 + 96 + b.BeaconBlock.SizeSSZ()
}


// MarshalSSZ marshals the SignedBeaconBlock object to SSZ format.
func (b *SignedBeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, b.SizeSSZ())
	return b.MarshalSSZTo(buf)
}

func (b *SignedBeaconBlock) ValidateAfterDecodingSSZ() error {
	return b.BeaconBlock.ValidateAfterDecodingSSZ()
}

// HashTreeRoot computes the SSZ hash tree root of the
// SignedBeaconBlock object.
func (b *SignedBeaconBlock) HashTreeRoot() ([32]byte, error) {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	if err := b.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()
	
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

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the SignedBeaconBlock object to a target array.
func (b *SignedBeaconBlock) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Offset for BeaconBlock
	offset := 100 // 4 + 96
	dst = fastssz.MarshalUint32(dst, uint32(offset))

	// Signature
	dst = append(dst, b.Signature[:]...)

	// BeaconBlock
	dst, err := b.BeaconBlock.MarshalSSZTo(dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the SignedBeaconBlock object.
func (b *SignedBeaconBlock) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 100 {
		return fastssz.ErrSize
	}

	// Read offset
	offset := fastssz.UnmarshallUint32(buf[0:4])
	if offset != 100 {
		return fastssz.ErrInvalidVariableOffset
	}

	// Signature
	copy(b.Signature[:], buf[4:100])

	// BeaconBlock
	if b.BeaconBlock == nil {
		b.BeaconBlock = &BeaconBlock{}
	}
	if err := b.BeaconBlock.UnmarshalSSZ(buf[100:]); err != nil {
		return err
	}

	return b.ValidateAfterDecodingSSZ()
}


// HashTreeRootWith ssz hashes the SignedBeaconBlock object with a hasher.
func (b *SignedBeaconBlock) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Message' (BeaconBlock)
	if err := b.BeaconBlock.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (1) 'Signature'
	hh.PutBytes(b.Signature[:])

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the SignedBeaconBlock object.
func (b *SignedBeaconBlock) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}

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
	"encoding/binary"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/karalabe/ssz"
)

// SigningData as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#signingdata
//
//nolint:lll // link.
type SigningData struct {
	// ObjectRoot is the hash tree root of the object being signed.
	ObjectRoot common.Root
	// Domain is the domain the object is being signed in.
	Domain common.Domain
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the SigningData object in SSZ encoding.
func (*SigningData) SizeSSZ() uint32 {
	//nolint:mnd // 32*2 = 64.
	return 64
}

// DefineSSZ defines the SSZ encoding for the SigningData object.
func (s *SigningData) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &s.ObjectRoot)
	ssz.DefineStaticBytes(codec, &s.Domain)
}

// HashTreeRoot computes the SSZ hash tree root of the SigningData object.
func (s *SigningData) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(s), nil
}

// MarshalSSZTo marshals the SigningData object to SSZ format into the provided
// buffer.
func (s *SigningData) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, s)
}

// MarshalSSZ marshals the SigningData object to SSZ format.
func (s *SigningData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, s.SizeSSZ())
	return s.MarshalSSZTo(buf)
}

// UnmarshalSSZ unmarshals the SigningData object from SSZ format.
func (s *SigningData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, s)
}

// ComputeSigningRoot as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_signing_root
//
//nolint:lll // link.
func ComputeSigningRoot(
	sszObject interface{ HashTreeRoot() ([32]byte, error) },
	domain common.Domain,
) (common.Root, error) {
	objectRoot, err := sszObject.HashTreeRoot()
	if err != nil {
		return common.Root{}, err
	}
	return (&SigningData{
		ObjectRoot: objectRoot,
		Domain:     domain,
	}).HashTreeRoot()
}

// ComputeSigningRootUInt64 computes the signing root of a uint64 value.
func ComputeSigningRootUInt64(
	value uint64,
	domain common.Domain,

) (common.Root, error) {
	bz := make([]byte, constants.RootLength)
	binary.LittleEndian.PutUint64(bz, value)
	return (&SigningData{
		ObjectRoot: common.Root(bz),
		Domain:     domain,
	}).HashTreeRoot()
}

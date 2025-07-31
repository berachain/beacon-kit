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
	"encoding/binary"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
)

//go:generate sszgen --path . --include ../../primitives/common,../../primitives/bytes --objs SigningData --output signing_data_sszgen.go

// SigningData as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#signingdata
type SigningData struct {
	// ObjectRoot is the hash tree root of the object being signed.
	ObjectRoot common.Root `ssz-size:"32"`
	// Domain is the domain the object is being signed in.
	Domain common.Domain `ssz-size:"32"`
}

func (*SigningData) ValidateAfterDecodingSSZ() error { return nil }

// ComputeSigningRoot as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_signing_root
func ComputeSigningRoot(
	sszObject interface{ HashTreeRoot() ([32]byte, error) },
	domain common.Domain,
) common.Root {
	htr, err := sszObject.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	root, err := (&SigningData{
		ObjectRoot: common.Root(htr),
		Domain:     domain,
	}).HashTreeRoot()
	if err != nil {
		panic(err)
	}
	return common.Root(root)
}

// ComputeSigningRootFromHTR computes the signing root from a pre-computed hash tree root.
func ComputeSigningRootFromHTR(
	htr [32]byte,
	domain common.Domain,
) common.Root {
	root, err := (&SigningData{
		ObjectRoot: common.Root(htr),
		Domain:     domain,
	}).HashTreeRoot()
	if err != nil {
		panic(err)
	}
	return common.Root(root)
}

// ComputeSigningRootUInt64 computes the signing root of a uint64 value.
func ComputeSigningRootUInt64(
	value uint64,
	domain common.Domain,

) common.Root {
	bz := make([]byte, constants.RootLength)
	binary.LittleEndian.PutUint64(bz, value)
	root, err := (&SigningData{
		ObjectRoot: common.Root(bz),
		Domain:     domain,
	}).HashTreeRoot()
	if err != nil {
		panic(err)
	}
	return common.Root(root)
}

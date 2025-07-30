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
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
)

// NOTE: sszgen generation requires special handling to avoid generating duplicate BeaconBlockHeader methods.
// The generated file needs manual cleanup to remove BeaconBlockHeader methods.
// go:generate sszgen -path . -objs SignedBeaconBlockHeader -output signed_beacon_block_header_sszgen.go -include ../../primitives/common,../../primitives/crypto,../../primitives/math,../../primitives/bytes -exclude-objs BeaconBlockBody,BeaconState,ExecutionPayload,ExecutionPayloadHeader,BeaconBlock,SignedBeaconBlock

var (
	_ constraints.SSZMarshallableRootable = (*SignedBeaconBlockHeader)(nil)
)

// SignedBeaconBlockHeader is a struct that contains a BeaconBlockHeader and a BLSSignature.
//
// NOTE: This struct is only ever (un)marshalled with SSZ and NOT with JSON.
type SignedBeaconBlockHeader struct {
	Header    *BeaconBlockHeader  `ssz-size:"112"`
	Signature crypto.BLSSignature `ssz-size:"96"`
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

func (*SignedBeaconBlockHeader) ValidateAfterDecodingSSZ() error { return nil }

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

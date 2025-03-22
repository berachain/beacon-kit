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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/eip7685"
	"github.com/karalabe/ssz"
)

// DepositRequest is introduced in EIP6110 which is currently not processed.
type DepositRequest = Deposit

// Compile-time checks for the DepositRequests list.
var (
	_ ssz.StaticObject            = (*DepositRequests)(nil)
	_ constraints.SSZMarshallable = (*DepositRequests)(nil)
)

// DepositRequests is used for SSZ unmarshalling a list of DepositRequests.
type DepositRequests []*DepositRequest

// DefineSSZ defines the SSZ encoding for the DepositRequests object.
func (dr DepositRequests) DefineSSZ(codec *ssz.Codec) {
	codec.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			codec, (*[]*DepositRequest)(&dr), constants.MaxDepositRequestsPerPayload,
		)
	})
	codec.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			codec, (*[]*DepositRequest)(&dr), constants.MaxDepositRequestsPerPayload,
		)
	})
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(
			codec, (*[]*DepositRequest)(&dr), constants.MaxDepositRequestsPerPayload,
		)
	})
}

// SizeSSZ returns the SSZ encoded size of the DepositRequests object.
func (dr DepositRequests) SizeSSZ(siz *ssz.Sizer) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ([]*DepositRequest)(dr))
}

// MarshalSSZ marshals the DepositRequests list with SSZ format.
func (dr DepositRequests) MarshalSSZ() ([]byte, error) {
	// return eip7685.MarshalItems(dr)
	buf := make([]byte, ssz.Size(dr))
	return buf, ssz.EncodeToBytes(buf, dr)
}

// ValidateAfterDecodingSSZ validates the DepositRequests object after decoding.
func (dr DepositRequests) ValidateAfterDecodingSSZ() error { return nil }

// DecodeDepositRequests decodes SSZ data by decoding each request individually.
func DecodeDepositRequests(data []byte) (*DepositRequests, error) {
	maxSize := constants.MaxDepositRequestsPerPayload * DepositSize
	if len(data) > int(maxSize) {
		return nil, fmt.Errorf(
			"invalid deposit requests SSZ size, requests should not be more than the max per "+
				"payload, got %d max %d", len(data), maxSize,
		)
	}
	depositSize := int(ssz.Size(&Deposit{}))
	if len(data) < depositSize {
		return nil, fmt.Errorf("invalid deposit requests SSZ size, got %d expected at least %d", len(data), depositSize)
	}
	// Use the generic unmarshalItems helper.
	items, err := eip7685.UnmarshalItems(data, depositSize, func() *DepositRequest { return new(DepositRequest) })
	if err != nil {
		return nil, err
	}
	deposits := DepositRequests(items)
	return &deposits, nil
}

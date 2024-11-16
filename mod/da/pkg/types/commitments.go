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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"github.com/karalabe/ssz"
)

// COMMITMENT_SIZE is the length of a KZG commitment in bytes.
const COMMITMENT_SIZE = 48

// BLOBS_PER_SLOT is the number of blobs that can be included in a slot.
const BLOBS_PER_SLOT = 6

// SlotCommitments represents a list of blob commitments for a slot.
// Used to store the blob commitments for a slot in the DB, because
// we need the commitments to retrieve the blobs from the DB.
type SlotCommitments struct {
	Commitments [][]byte
}

// SizeSSZ returns the size of the SSZ encoding
func (sc *SlotCommitments) SizeSSZ(sizer *ssz.Sizer, fixed bool) uint32 {
	size := uint32(4) // offset
	if fixed {
		return size
	}

	size += ssz.SizeSliceOfDynamicBytes(sizer, sc.Commitments)
	return size
}

// DefineSSZ defines the SSZ encoding for SlotCommitments
func (sc *SlotCommitments) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfDynamicBytesOffset(codec, &sc.Commitments, BLOBS_PER_SLOT, COMMITMENT_SIZE)
	ssz.DefineSliceOfDynamicBytesContent(codec, &sc.Commitments, BLOBS_PER_SLOT, COMMITMENT_SIZE)
}

// MarshalSSZ marshals SlotCommitments into SSZ format
func (sc *SlotCommitments) MarshalSSZ() ([]byte, error) {
	sizer := &ssz.Sizer{}
	size := sc.SizeSSZ(sizer, false)
	buf := make([]byte, size)

	return buf, ssz.EncodeToBytes(buf, sc)
}

// UnmarshalSSZ unmarshals SlotCommitments from SSZ format
func (sc *SlotCommitments) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, sc)
}

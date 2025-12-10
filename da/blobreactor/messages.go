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

package blobreactor

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cosmos/gogoproto/proto"
	karalabessz "github.com/karalabe/ssz"
)

// BlobMessage wraps our messages for CometBFT
// This implements proto.Message interface for CometBFT compatibility
type BlobMessage struct {
	Data []byte `json:"data"`
}

// Ensure BlobMessage implements proto.Message
var _ proto.Message = (*BlobMessage)(nil)

// Reset implements proto.Message
func (m *BlobMessage) Reset() {
	if m != nil {
		m.Data = nil
	}
}

// String implements proto.Message
func (m *BlobMessage) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("BlobMessage{Data: %d bytes}", len(m.Data))
}

// ProtoMessage implements proto.Message
func (m *BlobMessage) ProtoMessage() {}

// Marshal implements encoding for CometBFT
func (m *BlobMessage) Marshal() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return m.Data, nil
}

// Unmarshal implements decoding for CometBFT
func (m *BlobMessage) Unmarshal(data []byte) error {
	if m == nil {
		return errors.New("cannot unmarshal into nil BlobMessage")
	}
	m.Data = make([]byte, len(data))
	copy(m.Data, data)
	return nil
}

// Size returns the size of the message
func (m *BlobMessage) Size() int {
	if m == nil {
		return 0
	}
	return len(m.Data)
}

// NewBlobMessage creates a new blob message from JSON data
func NewBlobMessage(data []byte) *BlobMessage {
	return &BlobMessage{Data: data}
}

// ============================================================================
// SSZ Message Types for BlobReactor Protocol
// ============================================================================

// Ensure our types implement the necessary interfaces
var (
	_ karalabessz.StaticObject            = (*BlobRequest)(nil)
	_ constraints.SSZMarshallableRootable = (*BlobRequest)(nil)
	_ constraints.SSZMarshallable         = (*BlobResponse)(nil)
)

// MessageType identifies the type of message being sent
type MessageType uint8

const (
	MessageTypeRequest MessageType = iota
	MessageTypeResponse
)

// BlobRequest requests all blobs for a specific slot
type BlobRequest struct {
	Slot      math.Slot
	RequestID uint64 // Unique ID for request/response matching
}

// DefineSSZ defines the SSZ encoding for BlobRequest
func (r *BlobRequest) DefineSSZ(c *karalabessz.Codec) {
	karalabessz.DefineUint64(c, &r.Slot)
	karalabessz.DefineUint64(c, &r.RequestID)
}

// SizeSSZ returns the size of BlobRequest in SSZ encoding
//
//nolint:mnd // ok for now
func (*BlobRequest) SizeSSZ(*karalabessz.Sizer) uint32 {
	return 16 // uint64 slot + uint64 requestID
}

// MarshalSSZ marshals the BlobRequest to SSZ format
func (r *BlobRequest) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabessz.Size(r))
	return buf, karalabessz.EncodeToBytes(buf, r)
}

// HashTreeRoot computes the SSZ hash tree root
func (r *BlobRequest) HashTreeRoot() common.Root {
	return karalabessz.HashSequential(r)
}

// UnmarshalSSZ unmarshals BlobRequest from SSZ format
//
//nolint:mnd // ok for now
func (r *BlobRequest) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 16 {
		return fmt.Errorf("insufficient data for BlobRequest: need 16 bytes, got %d", len(buf))
	}
	r.Slot = math.Slot(binary.LittleEndian.Uint64(buf[0:8]))
	r.RequestID = binary.LittleEndian.Uint64(buf[8:16])
	return nil
}

// ValidateAfterDecodingSSZ validates the BlobRequest after SSZ decoding
func (*BlobRequest) ValidateAfterDecodingSSZ() error {
	return nil
}

const BlobResponseStaticSize uint32 = 28

// BlobResponse contains all blobs for the requested slot
type BlobResponse struct {
	Slot        math.Slot
	RequestID   uint64 // Echo back the request ID for matching
	HeadSlot    math.Slot
	SidecarData []byte // Raw SSZ-encoded BlobSidecars (avoiding double marshal/unmarshal)
}

// DefineSSZ defines the SSZ encoding for BlobResponse using karalabe/ssz codec
func (r *BlobResponse) DefineSSZ(c *karalabessz.Codec) {
	// Define fixed-size fields
	karalabessz.DefineUint64(c, &r.Slot)
	karalabessz.DefineUint64(c, &r.RequestID)
	karalabessz.DefineUint64(c, &r.HeadSlot)

	// Define dynamic field - offset first, then content
	karalabessz.DefineDynamicBytesOffset(c, &r.SidecarData, defaultRecvMessageCapacity)
	karalabessz.DefineDynamicBytesContent(c, &r.SidecarData, defaultRecvMessageCapacity)
}

// SizeSSZ returns the SSZ encoded size in bytes
func (r *BlobResponse) SizeSSZ(_ *karalabessz.Sizer, fixed bool) uint32 {
	var size = BlobResponseStaticSize
	if fixed {
		return size
	}

	// Dynamic part: actual sidecar data length
	size += uint32(len(r.SidecarData)) // #nosec G115 // length validated in ValidateAfterDecodingSSZ
	return size
}

// MarshalSSZ marshals the BlobResponse to SSZ format using the codec
func (r *BlobResponse) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabessz.Size(r))
	return buf, karalabessz.EncodeToBytes(buf, r)
}

// UnmarshalSSZ unmarshals BlobResponse from SSZ format using the codec
func (r *BlobResponse) UnmarshalSSZ(buf []byte) error {
	return karalabessz.DecodeFromBytes(buf, r)
}

// ValidateAfterDecodingSSZ validates the BlobResponse after SSZ decoding
func (r *BlobResponse) ValidateAfterDecodingSSZ() error {
	// Validate sidecar data size
	if len(r.SidecarData) > defaultRecvMessageCapacity {
		return fmt.Errorf("sidecar data too large: %d bytes (max %d)", len(r.SidecarData), defaultRecvMessageCapacity)
	}
	return nil
}

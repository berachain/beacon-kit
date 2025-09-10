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

// BlobResponse contains all blobs for the requested slot
type BlobResponse struct {
	Slot        math.Slot
	RequestID   uint64    // Echo back the request ID for matching
	Error       string    // Error message if blob fetch failed
	SidecarData []byte    // Raw SSZ-encoded BlobSidecars (avoiding double marshal/unmarshal)
	HeadSlot    math.Slot // Sender's current head slot (for status updates)
}

// DefineSSZ defines the SSZ encoding for BlobResponse (needed for constraints.SSZMarshallable)
func (r *BlobResponse) DefineSSZ(*karalabessz.Codec) {
	// For dynamic objects, we need custom encoding
	// This is handled in MarshalSSZ/UnmarshalSSZ
}

// MarshalSSZ marshals the BlobResponse to SSZ format
// We manually encode: slot (8) + request_id (8) + head_slot (8) + error_len (2) + error_str + offset (4) + sidecar data
//
//nolint:mnd // ok for now
func (r *BlobResponse) MarshalSSZ() ([]byte, error) {
	// Limit error message to 256 bytes to prevent buffer issues
	errorBytes := []byte(r.Error)
	if len(errorBytes) > 256 {
		errorBytes = errorBytes[:256]
	}

	// Total size: slot (8) + request_id (8) + head_slot (8) + error_len (2) + error_str + offset (4) + sidecar data
	totalSize := 8 + 8 + 8 + 2 + len(errorBytes) + 4 + len(r.SidecarData)
	buf := make([]byte, totalSize)

	// Write slot (8 bytes)
	binary.LittleEndian.PutUint64(buf[0:8], uint64(r.Slot))

	// Write request ID (8 bytes)
	binary.LittleEndian.PutUint64(buf[8:16], r.RequestID)

	// Write head slot (8 bytes)
	binary.LittleEndian.PutUint64(buf[16:24], uint64(r.HeadSlot))

	// Write error length (2 bytes)
	binary.LittleEndian.PutUint16(buf[24:26], uint16(len(errorBytes))) // #nosec G115

	// Write error string
	pos := 26
	if len(errorBytes) > 0 {
		copy(buf[pos:pos+len(errorBytes)], errorBytes)
		pos += len(errorBytes)
	}

	// Write offset to sidecars (4 bytes) - points after the fixed part
	// Current position + 4 bytes for offset itself
	//nolint:mnd // ok for now
	offset := uint32(pos + 4) // #nosec G115
	binary.LittleEndian.PutUint32(buf[pos:pos+4], offset)
	pos += 4

	// Write sidecar data (already SSZ-encoded BlobSidecars)
	if len(r.SidecarData) > 0 {
		copy(buf[pos:], r.SidecarData)
	}

	return buf, nil
}

// UnmarshalSSZ unmarshals BlobResponse from SSZ format
//
//nolint:mnd // ok for now
func (r *BlobResponse) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 26 {
		return fmt.Errorf("insufficient data for BlobResponse: need at least 26 bytes, got %d", len(buf))
	}

	// Read slot (8 bytes)
	r.Slot = math.Slot(binary.LittleEndian.Uint64(buf[0:8]))

	// Read request ID (8 bytes)
	r.RequestID = binary.LittleEndian.Uint64(buf[8:16])

	// Read head slot (8 bytes)
	r.HeadSlot = math.Slot(binary.LittleEndian.Uint64(buf[16:24]))

	// Read error length (2 bytes)
	errorLen := binary.LittleEndian.Uint16(buf[24:26])

	pos := 26

	// Read error string if present
	if errorLen > 0 {
		if len(buf) < pos+int(errorLen)+4 {
			return fmt.Errorf("insufficient data for error string: need %d bytes, got %d", pos+int(errorLen)+4, len(buf))
		}
		r.Error = string(buf[pos : pos+int(errorLen)])
		pos += int(errorLen)
	} else {
		r.Error = ""
	}

	// Read offset (4 bytes)
	if len(buf) < pos+4 {
		return fmt.Errorf("insufficient data for offset: need %d bytes, got %d", pos+4, len(buf))
	}
	offset := binary.LittleEndian.Uint32(buf[pos : pos+4])

	// Read sidecar data if present (already SSZ-encoded BlobSidecars)
	// #nosec G115
	if uint32(len(buf)) > offset {
		r.SidecarData = make([]byte, len(buf)-int(offset))
		copy(r.SidecarData, buf[offset:])
	}

	return nil
}

// ValidateAfterDecodingSSZ validates the BlobResponse after SSZ decoding
func (r *BlobResponse) ValidateAfterDecodingSSZ() error {
	// No validation needed for raw byte arrays
	return nil
}

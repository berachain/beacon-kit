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

package blobreactor

import (
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cosmos/gogoproto/proto"
	karalabessz "github.com/karalabe/ssz"
)

// BlobMessage is the raw envelope payload exchanged on the blob channel. The first byte of Data is the MessageType, the rest is the SSZ
// encoding of the corresponding message struct. It implements proto.Message so it can be used as a CometBFT channel message type without
// a dedicated protobuf schema.
type BlobMessage struct {
	Data []byte `json:"data"`
}

var _ proto.Message = (*BlobMessage)(nil)

func (m *BlobMessage) Reset() {
	if m != nil {
		m.Data = nil
	}
}

func (m *BlobMessage) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("BlobMessage{Data: %d bytes}", len(m.Data))
}

func (m *BlobMessage) ProtoMessage() {}

func (m *BlobMessage) Marshal() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return m.Data, nil
}

func (m *BlobMessage) Unmarshal(data []byte) error {
	if m == nil {
		return errors.New("cannot unmarshal into nil BlobMessage")
	}
	m.Data = make([]byte, len(data))
	copy(m.Data, data)
	return nil
}

func (m *BlobMessage) Size() int {
	if m == nil {
		return 0
	}
	return len(m.Data)
}

// newBlobMessage builds the envelope payload for a typed message.
func newBlobMessage(msgType MessageType, msgBz []byte) *BlobMessage {
	data := make([]byte, 0, 1+len(msgBz))
	data = append(data, byte(msgType))
	data = append(data, msgBz...)
	return &BlobMessage{Data: data}
}

// MessageType discriminates the messages multiplexed on the blob channel. Every message travels as a
// BlobMessage envelope whose first byte is the MessageType and whose remainder is the SSZ encoding of the
// corresponding struct. Pushes and haves are one-way gossip. By-root and by-range requests are answered with a
// SidecarsResponse correlated by a crypto-random RequestID, and the requester only accepts the response from
// the exact peer it asked. A response carries one chunk per complete slot, a slot the responder cannot fully
// serve is simply absent, and the requester treats absence as a miss rather than a success.
type MessageType uint8

const (
	// MessageTypePush carries unsolicited sidecars for a block at the tip.
	MessageTypePush MessageType = iota
	// MessageTypeHave announces that the sender holds all sidecars for a root.
	MessageTypeHave
	// MessageTypeByRootRequest asks for the sidecars of a single block.
	MessageTypeByRootRequest
	// MessageTypeByRangeRequest asks for the sidecars of a slot range.
	MessageTypeByRangeRequest
	// MessageTypeResponse answers a by-root or by-range request.
	MessageTypeResponse
)

func (t MessageType) String() string {
	switch t {
	case MessageTypePush:
		return "push"
	case MessageTypeHave:
		return "have"
	case MessageTypeByRootRequest:
		return "by_root_request"
	case MessageTypeByRangeRequest:
		return "by_range_request"
	case MessageTypeResponse:
		return "response"
	default:
		return fmt.Sprintf("unknown(%d)", uint8(t))
	}
}

const (
	// maxSidecarsChunkBytes bounds one slot's SSZ-encoded BlobSidecars. 6 sidecars x ~132 KiB each stays well below 1 MiB.
	//
	// TODO: a full block's sidecars (~768 KiB) currently travel as a single message per block (one SidecarsPush,
	// or one chunk per slot in a by-range response). A message this large monopolizes the blob channel while it
	// is in flight, causing head-of-line blocking against other traffic on that connection. Evaluate splitting
	// into smaller units (e.g. one message per sidecar, ~128 KiB) or Merkle-verified parts, as PR #2938 flagged.
	maxSidecarsChunkBytes = 1 << 20
	// maxChunksPerResponse bounds the number of per-slot chunks in a response.
	maxChunksPerResponse = 64
	// MaxRequestedSlots bounds the Count of a by-range request.
	MaxRequestedSlots = 64
)

var (
	_ karalabessz.StaticObject  = (*SidecarsByRootRequest)(nil)
	_ karalabessz.StaticObject  = (*SidecarsByRangeRequest)(nil)
	_ karalabessz.StaticObject  = (*SidecarsHave)(nil)
	_ karalabessz.DynamicObject = (*SidecarsPush)(nil)
	_ karalabessz.DynamicObject = (*SidecarsResponse)(nil)
)

// SidecarsPush is an unsolicited delivery of one block's sidecars.
type SidecarsPush struct {
	// BlockRoot is the hash tree root of the beacon block the sidecars belong to.
	BlockRoot common.Root
	// SidecarData is the SSZ encoding of the block's BlobSidecars.
	SidecarData []byte
}

func (p *SidecarsPush) SizeSSZ(_ *karalabessz.Sizer, fixed bool) uint32 {
	size := uint32(32 + 4) //nolint:mnd // root + offset
	if !fixed {
		size += uint32(len(p.SidecarData)) // #nosec G115 -- bounded by codec max
	}
	return size
}

func (p *SidecarsPush) DefineSSZ(c *karalabessz.Codec) {
	karalabessz.DefineStaticBytes(c, &p.BlockRoot)
	karalabessz.DefineDynamicBytesOffset(c, &p.SidecarData, maxSidecarsChunkBytes)
	karalabessz.DefineDynamicBytesContent(c, &p.SidecarData, maxSidecarsChunkBytes)
}

func (p *SidecarsPush) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabessz.Size(p))
	return buf, karalabessz.EncodeToBytes(buf, p)
}

func (p *SidecarsPush) UnmarshalSSZ(buf []byte) error {
	if err := karalabessz.DecodeFromBytes(buf, p); err != nil {
		return err
	}
	return p.ValidateAfterDecodingSSZ()
}

func (p *SidecarsPush) ValidateAfterDecodingSSZ() error {
	if len(p.SidecarData) == 0 {
		return errors.New("push without sidecar data")
	}
	return nil
}

// SidecarsHave announces that the sender holds the complete sidecar set for a block, so peers suppress
// redundant full-payload pushes toward it and by-root fetchers know where to ask.
type SidecarsHave struct {
	BlockRoot common.Root
}

func (h *SidecarsHave) SizeSSZ(*karalabessz.Sizer) uint32 { return 32 } //nolint:mnd // root

func (h *SidecarsHave) DefineSSZ(c *karalabessz.Codec) {
	karalabessz.DefineStaticBytes(c, &h.BlockRoot)
}

func (h *SidecarsHave) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabessz.Size(h))
	return buf, karalabessz.EncodeToBytes(buf, h)
}

func (h *SidecarsHave) UnmarshalSSZ(buf []byte) error {
	return karalabessz.DecodeFromBytes(buf, h)
}

// SidecarsByRootRequest asks a peer for all sidecars of one block. The slot is a lookup hint (the availability store is slot-indexed);
// the responder and the requester both check the root against the returned sidecar headers.
type SidecarsByRootRequest struct {
	RequestID uint64
	Slot      math.Slot
	BlockRoot common.Root
}

func (r *SidecarsByRootRequest) SizeSSZ(*karalabessz.Sizer) uint32 { return 48 } //nolint:mnd // id + slot + root

func (r *SidecarsByRootRequest) DefineSSZ(c *karalabessz.Codec) {
	karalabessz.DefineUint64(c, &r.RequestID)
	karalabessz.DefineUint64(c, &r.Slot)
	karalabessz.DefineStaticBytes(c, &r.BlockRoot)
}

func (r *SidecarsByRootRequest) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabessz.Size(r))
	return buf, karalabessz.EncodeToBytes(buf, r)
}

func (r *SidecarsByRootRequest) UnmarshalSSZ(buf []byte) error {
	if err := karalabessz.DecodeFromBytes(buf, r); err != nil {
		return err
	}
	return r.ValidateAfterDecodingSSZ()
}

func (r *SidecarsByRootRequest) ValidateAfterDecodingSSZ() error { return nil }

// SidecarsByRangeRequest asks a peer for the sidecars of slots [StartSlot, StartSlot+Count).
type SidecarsByRangeRequest struct {
	RequestID uint64
	StartSlot math.Slot
	Count     uint64
}

func (r *SidecarsByRangeRequest) SizeSSZ(*karalabessz.Sizer) uint32 { return 24 } //nolint:mnd // three uint64s

func (r *SidecarsByRangeRequest) DefineSSZ(c *karalabessz.Codec) {
	karalabessz.DefineUint64(c, &r.RequestID)
	karalabessz.DefineUint64(c, &r.StartSlot)
	karalabessz.DefineUint64(c, &r.Count)
}

func (r *SidecarsByRangeRequest) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabessz.Size(r))
	return buf, karalabessz.EncodeToBytes(buf, r)
}

func (r *SidecarsByRangeRequest) UnmarshalSSZ(buf []byte) error {
	if err := karalabessz.DecodeFromBytes(buf, r); err != nil {
		return err
	}
	return r.ValidateAfterDecodingSSZ()
}

func (r *SidecarsByRangeRequest) ValidateAfterDecodingSSZ() error {
	if r.Count == 0 || r.Count > MaxRequestedSlots {
		return fmt.Errorf("invalid by-range count %d (max %d)", r.Count, MaxRequestedSlots)
	}
	return nil
}

// SidecarsResponse answers a by-root or by-range request. Each chunk is the SSZ encoding of one block's complete BlobSidecars; a by-root
// response has at most one chunk. Slots the responder has no data for are simply absent.
type SidecarsResponse struct {
	RequestID uint64
	// SidecarChunks holds one complete slot's SSZ-encoded BlobSidecars each.
	SidecarChunks [][]byte
}

func (r *SidecarsResponse) SizeSSZ(sizer *karalabessz.Sizer, fixed bool) uint32 {
	const sszOffsetSize = 4
	size := uint32(8 + sszOffsetSize) //nolint:mnd // uint64 + offset
	if !fixed {
		for _, chunk := range r.SidecarChunks {
			size += sszOffsetSize + uint32(len(chunk)) // #nosec G115 -- chunk length bounded by the codec max
		}
	}
	_ = sizer
	return size
}

func (r *SidecarsResponse) DefineSSZ(c *karalabessz.Codec) {
	karalabessz.DefineUint64(c, &r.RequestID)
	karalabessz.DefineSliceOfDynamicBytesOffset(c, &r.SidecarChunks, maxChunksPerResponse, maxSidecarsChunkBytes)
	karalabessz.DefineSliceOfDynamicBytesContent(c, &r.SidecarChunks, maxChunksPerResponse, maxSidecarsChunkBytes)
}

func (r *SidecarsResponse) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabessz.Size(r))
	return buf, karalabessz.EncodeToBytes(buf, r)
}

func (r *SidecarsResponse) UnmarshalSSZ(buf []byte) error {
	if err := karalabessz.DecodeFromBytes(buf, r); err != nil {
		return err
	}
	return r.ValidateAfterDecodingSSZ()
}

func (r *SidecarsResponse) ValidateAfterDecodingSSZ() error {
	if len(r.SidecarChunks) > maxChunksPerResponse {
		return fmt.Errorf("too many sidecar chunks: %d (max %d)", len(r.SidecarChunks), maxChunksPerResponse)
	}
	return nil
}

// SPDX-License-Identifier: MIT
//
// SPDX-License-Identifier: MIT
//
// Copyright (c) 2025 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WdeHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/karalabe/ssz"
)

// ConsolidationRequest is introduced in Pectra but not used by us.
// We keep it so we can maintain parity tests with other SSZ implementations.
type ConsolidationRequest struct {
	SourceAddress common.ExecutionAddress
	SourcePubKey  crypto.BLSPubkey
	TargetPubKey  crypto.BLSPubkey
}

// ConsolidationRequests is used for SSZ unmarshalling a list of ConsolidationRequest
type ConsolidationRequests []*ConsolidationRequest

/* -------------------------------------------------------------------------- */
/*                       Consolidation Requests SSZ                           */
/* -------------------------------------------------------------------------- */

func (c *ConsolidationRequest) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &c.SourceAddress)
	ssz.DefineStaticBytes(codec, &c.SourcePubKey)
	ssz.DefineStaticBytes(codec, &c.TargetPubKey)
}

func (c *ConsolidationRequest) SizeSSZ(_ *ssz.Sizer) uint32 {
	return sszConsolidationRequestSize
}

func (c *ConsolidationRequest) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(c))
	return buf, ssz.EncodeToBytes(buf, c)
}

func (c *ConsolidationRequest) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, c)
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (c *ConsolidationRequest) HashTreeRoot() common.Root {
	return ssz.HashSequential(c)
}

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (cr ConsolidationRequests) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ([]*ConsolidationRequest)(cr))
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
// TODO: get from accessible chainspec field params.
func (cr ConsolidationRequests) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*ConsolidationRequest)(&cr), maxConsolidationRequestsPerPayload)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*ConsolidationRequest)(&cr), maxConsolidationRequestsPerPayload)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*ConsolidationRequest)(&cr), maxConsolidationRequestsPerPayload)
	})
}

// MarshalSSZ marshals the BlobSidecars object to SSZ format.
func (cr *ConsolidationRequests) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(cr))
	return buf, ssz.EncodeToBytes(buf, cr)
}

func (cr *ConsolidationRequests) NewFromSSZ(data []byte) (*ConsolidationRequests, error) {
	if cr == nil {
		cr = &ConsolidationRequests{}
	}
	return cr, ssz.DecodeFromBytes(data, cr)
}

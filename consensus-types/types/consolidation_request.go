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
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip7685"
	"github.com/karalabe/ssz"
)

const maxConsolidationRequestsPerPayload = 2
const sszConsolidationRequestSize = 116

// ConsolidationRequest is introduced in Pectra but not used by us.
// We keep it so we can maintain parity tests with other SSZ implementations.
type ConsolidationRequest struct {
	SourceAddress common.ExecutionAddress
	SourcePubKey  crypto.BLSPubkey
	TargetPubKey  crypto.BLSPubkey
}

func (c *ConsolidationRequest) EnsureSyntaxFromSSZ() error {
	return nil
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

// HashTreeRoot returns the hash tree root of the Deposits.
func (c *ConsolidationRequest) HashTreeRoot() common.Root {
	return ssz.HashSequential(c)
}

// MarshalSSZ marshals the ConsolidationRequests object to SSZ format by encoding each consolidation request individually.
func (cr *ConsolidationRequests) MarshalSSZ() ([]byte, error) {
	return eip7685.MarshalItems[*ConsolidationRequest](*cr)
}

// DecodeConsolidationRequests decodes SSZ data by decoding each request individually.
func DecodeConsolidationRequests(data []byte) (*ConsolidationRequests, error) {
	maxSize := maxConsolidationRequestsPerPayload * sszConsolidationRequestSize
	if len(data) > maxSize {
		return nil, fmt.Errorf(
			"invalid consolidation requests SSZ size, requests should not be more than the "+
				"max per payload, got %d max %d", len(data), maxSize,
		)
	}
	requestSize := int(ssz.Size(&ConsolidationRequest{}))
	if len(data) < requestSize {
		return nil, fmt.Errorf(
			"invalid consolidation requests SSZ size, got %d expected at least %d", len(data), requestSize,
		)
	}
	if len(data)%requestSize != 0 {
		return nil, fmt.Errorf(
			"invalid data length: %d is not a multiple of consolidation request size %d", len(data), requestSize,
		)
	}
	items, err := eip7685.UnmarshalItems[*ConsolidationRequest](
		data,
		requestSize,
		func() *ConsolidationRequest { return new(ConsolidationRequest) },
	)
	if err != nil {
		return nil, err
	}
	consolidationRequests := ConsolidationRequests(items)
	return &consolidationRequests, nil
}

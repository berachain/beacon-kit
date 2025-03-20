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

import "github.com/karalabe/ssz"

const maxDepositRequestsPerPayload = 8192

// DepositRequest is introduced in EIP6110 which is currently not processed.
type DepositRequest = Deposit

// DepositRequests is used for SSZ unmarshalling a list of DepositRequest
type DepositRequests []*DepositRequest

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (dr DepositRequests) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ([]*DepositRequest)(dr))
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
func (dr DepositRequests) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*DepositRequest)(&dr), maxDepositRequestsPerPayload)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*DepositRequest)(&dr), maxDepositRequestsPerPayload)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*DepositRequest)(&dr), maxDepositRequestsPerPayload)
	})
}

// MarshalSSZ marshals the BlobSidecars object to SSZ format.
func (dr *DepositRequests) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(dr))
	return buf, ssz.EncodeToBytes(buf, dr)
}

func (dr *DepositRequests) NewFromSSZ(data []byte) (*DepositRequests, error) {
	if dr == nil {
		dr = &DepositRequests{}
	}
	return dr, ssz.DecodeFromBytes(data, dr)
}

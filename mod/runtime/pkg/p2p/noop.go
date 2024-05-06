// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package p2p

import (
	"context"

	ssz "github.com/ferranbt/fastssz"
)

// NoopGossipHandler is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopGossipHandler[DataT interface {
	ssz.Marshaler
	ssz.Unmarshaler
}, BytesT ~[]byte] struct{}

// Publish creates a new NoopGossipHandler.
func (n NoopGossipHandler[DataT, BytesT]) Publish(
	_ context.Context,
	data DataT,
) (BytesT, error) {
	return data.MarshalSSZ()
}

// Request simply returns the reference it receives.
func (n NoopGossipHandler[DataT, BytesT]) Request(
	_ context.Context,
	ref BytesT,
	out DataT,
) error {
	return out.UnmarshalSSZ(ref)
}

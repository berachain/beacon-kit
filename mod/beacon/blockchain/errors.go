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

package blockchain

import "github.com/cockroachdb/errors"

var (
	// ErrInvalidPayload indicates that the payload of a beacon block is
	// invalid.
	ErrInvalidPayload = errors.New("invalid payload")
	// ErrNoPayloadInBeaconBlock indicates that a beacon block was expected to
	// have a payload, but none was found.
	ErrNoPayloadInBeaconBlock = errors.New("no payload in beacon block")
	// ErrNilBlkBody is an error for when the block body is nil.
	ErrNilBlkBody = errors.New("nil block body")
	// ErrNilBlockHeader is returned when a block header from a block is nil.
	ErrNilBlockHeader = errors.New("nil block header")
	// ErrNilBlk is an error for when the beacon block is nil.
	ErrNilBlk = errors.New("nil beacon block")
)

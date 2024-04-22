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

package builder

import "github.com/cockroachdb/errors"

var (
	// ErrNilPayloadOnValidResponse is returned when a nil payload ID is
	// received on a VALID engine response.
	ErrNilPayloadOnValidResponse = errors.New(
		"received nil payload ID on VALID engine response",
	)

	// ErrNilPayloadID is returned when a nil payload ID is received.
	ErrNilPayloadID = errors.New("received nil payload ID")

	// ErrPayloadIDNotFound is returned when a payload ID is not found in the
	// cache.
	ErrPayloadIDNotFound = errors.New("unable to find payload ID in cache")

	// ErrCachedPayloadNotFoundOnExecutionClient is returned when a cached
	// payloadID is not found on the execution client.
	ErrCachedPayloadNotFoundOnExecutionClient = errors.New(
		"cached payload ID could not be resolved on execution client",
	)

	// ErrLocalBuildingDisabled is returned when local building is disabled.
	ErrLocalBuildingDisabled = errors.New("local building is disabled")

	// ErrNilPayloadEnvelope is returned when a nil payload envelope is
	// received.
	ErrNilPayloadEnvelope = errors.New("received nil payload envelope")
)

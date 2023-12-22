// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package consensus_types

import (
	"errors"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/math"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// BytesToExecutionData converts a byte array to an ExecutionData object based on the requested version.
func BytesToExecutionData(bz []byte, value math.Gwei,
	requestedVersion int) (interfaces.ExecutionData, error) {
	if bz == nil {
		return nil, errors.New("nil execution data")
	}

	switch requestedVersion {
	case version.Bellatrix:
		return wrapExecutionPayload(bz)
	case version.Capella:
		return wrapExecutionPayloadCapella(bz, value)
	case version.Deneb:
		return wrapExecutionPayloadDeneb(bz, value)
	default:
		return nil, errors.New("version not supported")
	}
}

// wrapExecutionPayload wraps a byte array into an ExecutionPayload object.
func wrapExecutionPayload(bz []byte) (interfaces.ExecutionData, error) {
	payload := &enginev1.ExecutionPayload{}
	if err := payload.UnmarshalSSZ(bz); err != nil {
		return nil, err
	}
	return blocks.WrappedExecutionPayload(payload)
}

// wrapExecutionPayloadCapella wraps a byte array and value into an ExecutionPayloadCapella object.
func wrapExecutionPayloadCapella(bz []byte, value math.Gwei) (interfaces.ExecutionData, error) {
	payload := &enginev1.ExecutionPayloadCapella{}
	if err := payload.UnmarshalSSZ(bz); err != nil {
		return nil, err
	}
	return blocks.WrappedExecutionPayloadCapella(payload, value)
}

// wrapExecutionPayloadDeneb wraps a byte array and value into an ExecutionPayloadDeneb object.
func wrapExecutionPayloadDeneb(bz []byte, value math.Gwei) (interfaces.ExecutionData, error) {
	payload := &enginev1.ExecutionPayloadDeneb{}
	if err := payload.UnmarshalSSZ(bz); err != nil {
		return nil, err
	}
	return blocks.WrappedExecutionPayloadDeneb(payload, value)
}

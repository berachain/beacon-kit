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

package engine

import (
	"github.com/itsdevbear/bolaris/math"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	capella "github.com/itsdevbear/bolaris/types/engine/v1/capella"
	deneb "github.com/itsdevbear/bolaris/types/engine/v1/deneb"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// WrappedExecutionPayloadDeneb is a constructor which wraps a protobuf execution payload
// into an interface.
func WrappedExecutionPayloadDeneb(
	p *enginev1.ExecutionPayloadDeneb, wei math.Wei,
) (ExecutionPayload, error) {
	return deneb.NewWrappedExecutionPayloadDeneb(
		p, wei,
	), nil
}

// WrappedExecutionPayloadCapella is a constructor which wraps a protobuf execution payload
// into an interface.
func WrappedExecutionPayloadCapella(
	p *enginev1.ExecutionPayloadCapella, wei math.Wei,
) (ExecutionPayload, error) {
	return capella.NewWrappedExecutionPayloadCapella(
		p, wei,
	), nil
}

// EmptyExecutionPayloadWithVersion returns an empty execution payload for the given version.
func EmptyPayloadAttributesWithVersion(v int) PayloadAttributer {
	switch v {
	case version.Deneb:
		return &deneb.WrappedPayloadAttributesV3{}
	case version.Capella:
		return &capella.WrappedPayloadAttributesV2{}
	default:
		return nil
	}
}

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
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
)

// We re-export the interfaces from the engine package to aid in importing them
// in other packages.
type (
	// ExecutionPayloadBody is the interface for the execution data of a block.
	// It contains all the fields that are part of both an execution payload header
	// and a full execution payload.
	ExecutionPayloadBody = interfaces.ExecutionPayloadBody

	// ExecutionPayload is the interface for the execution data of a block.
	ExecutionPayload = interfaces.ExecutionPayload

	// ExecutionPayloadHeader is the interface representing an execution payload header.
	ExecutionPayloadHeader = interfaces.ExecutionPayloadHeader

	// PayloadAttributer is the interface for the payload attributes of a block.
	PayloadAttributer = interfaces.PayloadAttributer
)

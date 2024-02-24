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

package ethclient

import (
	"context"

	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
)

// GethRPCClient is an interface for the Geth RPC client.
type GethRPCClient interface {
	CallContext(
		ctx context.Context,
		result interface{},
		method string,
		args ...interface{},
	) error
}

// ForkchoiceUpdatedResponse is the response kind received by the
// engine_forkchoiceUpdatedV1 endpoint.
type ForkchoiceUpdatedResponse struct {
	Status          *enginev1.PayloadStatus  `json:"payloadStatus"`
	PayloadID       *enginev1.PayloadIDBytes `json:"payloadId"`
	ValidationError string                   `json:"validationError"`
}

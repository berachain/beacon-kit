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

package enginev1

import enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

const (
	// General
	PayloadStatus_UNKNOWN            = enginev1.PayloadStatus_UNKNOWN
	PayloadStatus_VALID              = enginev1.PayloadStatus_VALID
	PayloadStatus_INVALID            = enginev1.PayloadStatus_INVALID
	PayloadStatus_SYNCING            = enginev1.PayloadStatus_SYNCING
	PayloadStatus_ACCEPTED           = enginev1.PayloadStatus_ACCEPTED
	PayloadStatus_INVALID_BLOCK_HASH = enginev1.PayloadStatus_INVALID_BLOCK_HASH
)

type (
	// General
	PayloadStatus    = enginev1.PayloadStatus
	ForkchoiceState  = enginev1.ForkchoiceState
	PayloadIDBytes   = enginev1.PayloadIDBytes
	ExecutionBlock   = enginev1.ExecutionBlock
	ExecutionPayload = enginev1.ExecutionPayload
	Withdrawal       = enginev1.Withdrawal

	// Deneb
	PayloadAttributesV3                          = enginev1.PayloadAttributesV3
	BlobsBundle                                  = enginev1.BlobsBundle
	ExecutionPayloadDeneb                        = enginev1.ExecutionPayloadDeneb
	ExecutionPayloadDenebWithValueAndBlobsBundle = enginev1.ExecutionPayloadDenebWithValueAndBlobsBundle

	// Capella
	PayloadAttributesV2              = enginev1.PayloadAttributesV2
	ExecutionPayloadCapella          = enginev1.ExecutionPayloadCapella
	ExecutionPayloadCapellaWithValue = enginev1.ExecutionPayloadCapellaWithValue
)

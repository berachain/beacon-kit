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

package enginev1

import enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

type (
	// ExecutionPayloadCapella is an alias for the Prysm ExecutionPayloadCapella.
	ExecutionPayloadCapella = enginev1.ExecutionPayloadCapella
	// ExecutionPayloadCapellaWithValue is an alias for the Prysm ExecutionPayloadCapellaWithValue.
	ExecutionPayloadCapellaWithValue = enginev1.ExecutionPayloadCapellaWithValue
	// ExecutionPayloadHeader is an alias for the Prysm ExecutionPayloadHeader.
	ExecutionPayloadHeaderCapella = enginev1.ExecutionPayloadHeaderCapella
	// ExecutionPayloadDeneb is an alias for the Prysm ExecutionPayloadDeneb.
	ExecutionPayloadDeneb = enginev1.ExecutionPayloadDeneb
	// ExecutionPayloadDenebWithValueAndBlobsBundle is an alias for the Prysm
	// ExecutionPayloadDenebWithValueAndBlobsBundle.
	//nolint:lll // alias.
	ExecutionPayloadDenebWithValueAndBlobsBundle = enginev1.ExecutionPayloadDenebWithValueAndBlobsBundle
	// ExecutionPayloadHeaderDeneb is an alias for the Prysm ExecutionPayloadHeaderDeneb.
	ExecutionPayloadHeaderDeneb = enginev1.ExecutionPayloadHeaderDeneb
	// BlobsBundle is an alias for the Prysm BlobsBundle.
	BlobsBundle = enginev1.BlobsBundle
	// Withdrawal is an alias for the Prysm Withdrawal.
	Withdrawal = enginev1.Withdrawal
	// ForkchoiceState is an alias for the Prysm ForkchoiceState.
	ForkchoiceState = enginev1.ForkchoiceState
	// PayloadAttributesContainer is an alias for the Prysm PayloadAttributesV2.
	PayloadAttributesV2 = enginev1.PayloadAttributesV2
	// PayloadAttributesContainer is an alias for the Prysm PayloadAttributesV3.
	PayloadAttributesV3 = enginev1.PayloadAttributesV3
	// PayloadAttributesContainer is an alias for the Prysm PayloadAttributesContainer.
	PayloadIDBytes = enginev1.PayloadIDBytes
	// PayloadStatus is an alias for the Prysm PayloadStatus.
	PayloadStatus = enginev1.PayloadStatus
	// ExecutionBlock is an alias for the Prysm ExecutionBlock.
	ExecutionBlock = enginev1.ExecutionBlock
)

const (
	//nolint:stylecheck,revive // alias.
	PayloadStatus_INVALID_BLOCK_HASH = enginev1.PayloadStatus_INVALID_BLOCK_HASH
	//nolint:stylecheck,revive // alias.
	PayloadStatus_ACCEPTED = enginev1.PayloadStatus_ACCEPTED
	//nolint:stylecheck,revive // alias.
	PayloadStatus_SYNCING = enginev1.PayloadStatus_SYNCING
	//nolint:stylecheck,revive // alias.
	PayloadStatus_INVALID = enginev1.PayloadStatus_INVALID
	//nolint:stylecheck,revive // alias.
	PayloadStatus_VALID = enginev1.PayloadStatus_VALID
	//nolint:stylecheck,revive // alias.
	PayloadStatus_UNKNOWN = enginev1.PayloadStatus_UNKNOWN
)

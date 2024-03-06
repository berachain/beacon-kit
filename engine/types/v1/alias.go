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

import enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"

type (
	// ExecutionPayloadDeneb alias for Prysm's version.
	ExecutionPayloadDeneb = enginev1.ExecutionPayloadDeneb
	// ExecutionPayloadDenebWithValueAndBlobsBundle alias for Prysm's version.
	//nolint:lll // alias.
	ExecutionPayloadDenebWithValueAndBlobsBundle = enginev1.ExecutionPayloadDenebWithValueAndBlobsBundle
	// ExecutionPayloadHeaderDeneb alias for Prysm's version.
	ExecutionPayloadHeaderDeneb = enginev1.ExecutionPayloadHeaderDeneb
	// BlobsBundle alias for Prysm's version.
	BlobsBundle = enginev1.BlobsBundle
	// Withdrawal alias for Prysm's version.
	Withdrawal = enginev1.Withdrawal
	// PayloadStatus alias for Prysm's version.
	PayloadStatus = enginev1.PayloadStatus
	// ExecutionBlock alias for Prysm's version.
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

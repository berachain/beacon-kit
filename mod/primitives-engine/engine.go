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

//nolint:gochecknoglobals // alias.
package engineprimitives

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/ethereum/go-ethereum/beacon/engine"
)

// There are some types we can borrow from geth.
type (
	ClientVersionV1 = engine.ClientVersionV1
	ExecutableData  = engine.ExecutableData
)

var (
	// ExecutableDataToBlock constructs a block from executable data.
	ExecutableDataToBlock = engine.ExecutableDataToBlock
)

type PayloadStatusStr = string

var (
	// PayloadStatusValid is the status of a valid payload.
	PayloadStatusValid PayloadStatusStr = "VALID"
	// PayloadStatusInvalid is the status of an invalid payload.
	PayloadStatusInvalid PayloadStatusStr = "INVALID"
	// PayloadStatusSyncing is the status returned when the EL is syncing.
	PayloadStatusSyncing PayloadStatusStr = "SYNCING"
	// PayloadStatusAccepted is the status returned when the EL has accepted the
	// payload.
	PayloadStatusAccepted PayloadStatusStr = "ACCEPTED"
)

// ForkchoiceResponseV1 as per the EngineAPI Specification:
// https://github.com/ethereum/execution-apis/blob/main/src/engine/paris.md#response-2
//
//nolint:lll // link.
type ForkchoiceResponseV1 struct {
	// PayloadStatus is the payload status.
	PayloadStatus PayloadStatusV1 `json:"payloadStatus"`

	// PayloadID isthe identifier of the payload build process, it
	// can also be `nil`.
	PayloadID *PayloadID `json:"payloadId"`
}

// ForkchoicStateV1 as per the EngineAPI Specification:
// https://github.com/ethereum/execution-apis/blob/main/src/engine/paris.md#forkchoicestatev1
//
//nolint:lll // link.
type ForkchoiceStateV1 struct {
	// HeadBlockHash is the desired block hash of the head of the canonical
	// chain.
	HeadBlockHash common.ExecutionHash `json:"headBlockHash"`

	// SafeBlockHash is  the "safe" block hash of the canonical chain under
	// certain
	// synchrony and honesty assumptions. This value MUST be either equal to
	// or an ancestor of `HeadBlockHash`.
	SafeBlockHash common.ExecutionHash `json:"safeBlockHash"`

	// FinalizedBlockHash is the desired block hash of the most recent finalized
	// block
	FinalizedBlockHash common.ExecutionHash `json:"finalizedBlockHash"`
}

// PayloadStatusV1 represents the status of a payload as per the EngineAPI
// Specification. For more details, see:
// https://github.com/ethereum/execution-apis/blob/main/src/engine/paris.md#payloadstatusv1
//
//nolint:lll // link.
type PayloadStatusV1 struct {
	// Status string of the payload.
	Status string `json:"status"`

	// LatestValidHash is the hash of the most recent valid block
	// in the branch defined by payload and its ancestors
	LatestValidHash *common.ExecutionHash `json:"latestValidHash"`

	// ValidationError is a message providing additional details on
	// the validation error if the payload is classified as
	// INVALID or INVALID_BLOCK_HASH
	ValidationError *string `json:"validationError"`
}

// PayloadID is an identifier for the payload build process.
type PayloadID = bytes.B8

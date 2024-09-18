// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//nolint:gochecknoglobals // alias.
package engineprimitives

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// ClientVersionV1 contains information which identifies a client
// implementation.
type ClientVersionV1 struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func (v *ClientVersionV1) String() string {
	return fmt.Sprintf("%s-%s-%s-%s", v.Code, v.Name, v.Version, v.Commit)
}

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

// ForkchoiceStateV1 as per the EngineAPI Specification:
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

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

package enginetypes

import (
	"fmt"
	"math/big"

	"github.com/berachain/beacon-kit/math"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

type ExecutionPayloadEnvelope interface {
	GetExecutionPayload() ExecutionPayload
	GetValue() math.Wei
	GetBlobsBundle() *BlobsBundleV1
	ShouldOverrideBuilder() bool
}

//go:generate go run github.com/fjl/gencodec -type ExecutionPayloadEnvelopeDeneb -field-override executionPayloadEnvelopeMarshaling -out payload_env.json.go
//nolint:lll
type ExecutionPayloadEnvelopeDeneb struct {
	ExecutionPayload *ExecutableData `json:"executionPayload"      gencodec:"required"`
	BlockValue       *big.Int        `json:"blockValue"            gencodec:"required"`
	BlobsBundle      *BlobsBundleV1  `json:"blobsBundle"`
	Override         bool            `json:"shouldOverrideBuilder"`
}

func (e *ExecutionPayloadEnvelopeDeneb) GetExecutionPayload() ExecutionPayload {
	return e.ExecutionPayload
}

func (e *ExecutionPayloadEnvelopeDeneb) GetValue() math.Wei {
	val, ok := uint256.FromBig(e.BlockValue)
	if !ok {
		return math.Wei{}
	}
	return math.Wei{Int: val}
}

func (e *ExecutionPayloadEnvelopeDeneb) GetBlobsBundle() *BlobsBundleV1 {
	return e.BlobsBundle
}

func (e *ExecutionPayloadEnvelopeDeneb) ShouldOverrideBuilder() bool {
	return e.Override
}

func (e *ExecutionPayloadEnvelopeDeneb) String() string {
	return fmt.Sprintf(`
ExecutionPayloadEnvelopeDeneb{
	ExecutionPayload: %s,
	BlockValue: %s,
	BlobsBundle: %s,
	Override: %v,
}`, e.ExecutionPayload.String(),
		e.BlockValue.String(),
		e.GetBlobsBundle().Blobs,
		e.Override,
	)
}

// JSON type overrides for ExecutionPayloadEnvelope.
type executionPayloadEnvelopeMarshaling struct {
	BlockValue *hexutil.Big
}

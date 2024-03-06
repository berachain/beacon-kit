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

package typesv1

import (
	"errors"

	"github.com/itsdevbear/bolaris/config/version"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
)

// IsNil checks if the BeaconBlock is nil or not.
func (b *BeaconBlockDeneb) IsNil() bool {
	return b == nil
}

// Version returns the version of the block.
func (b *BeaconBlockDeneb) Version() int {
	return version.Deneb
}

// AttachExecution attaches the given execution data to the block.
func (b *BeaconBlockDeneb) AttachExecution(
	executionData enginetypes.ExecutionPayload,
) error {
	var ok bool
	if executionData == nil {
		b.PayloadValue = make([]byte, 32) //nolint:gomnd // fix later.
		return nil
	}
	b.Body.ExecutionPayload, ok = executionData.
		ToProto().(*enginev1.ExecutionPayloadDeneb)
	if !ok {
		return errors.New(
			"failed to convert execution data to capella payload")
	}

	// TODO: this needs to be done better, really hood rn.
	payloadValueBz := make([]byte, 32) //nolint:gomnd // 32 bytes for uint256.
	copy(payloadValueBz, executionData.GetValue().Bytes())
	b.PayloadValue = payloadValueBz
	return nil
}

// Execution returns the execution data of the block.
func (b *BeaconBlockDeneb) ExecutionPayload() (
	enginetypes.ExecutionPayload, error,
) {
	return &enginev1.ExecutionPayloadEnvelope{
		Payload: &enginev1.ExecutionPayloadEnvelope_Deneb{
			Deneb: b.GetBody().GetExecutionPayload(),
		},
		PayloadValue:          b.GetPayloadValue(),
		BlobsBundle:           &enginev1.BlobsBundle{},
		ShouldOverrideBuilder: false,
	}, nil
}

// GetSlot returns the slot of the block.
func (b *BeaconBlockDeneb) GetBlobKzgCommitments() [][]byte {
	if b.GetBody().GetExecutionPayload() == nil {
		return nil
	}
	return b.GetBody().GetBlobKzgCommitments()
}

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

package consensusv1

import (
	"errors"

	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	github_com_itsdevbear_bolaris_types_consensus_primitives "github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
	v1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	// TODO: fix jank sszgen import naming requirement
)

// BeaconKitBlock implements the BeaconKitBlock interface.
var _ interfaces.BeaconKitBlock = (*BeaconKitBlockCapella)(nil)

// BeaconKitBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconKitBlock(
	slot github_com_itsdevbear_bolaris_types_consensus_primitives.Slot,
	executionData engine.ExecutionPayload,
) (interfaces.BeaconKitBlock, error) {
	block := &BeaconKitBlockCapella{
		Slot: slot,
		Body: &BeaconKitBlockBodyCapella{
			RandaoReveal: make([]byte, 96), //nolint:gomnd // 96 bytes for RandaoReveal.
			Graffiti:     make([]byte, 32), //nolint:gomnd // 32 bytes for Graffiti.
		},
	}
	if executionData != nil {
		if err := block.AttachExecution(executionData); err != nil {
			return nil, err
		}
	}
	return block, nil
}

// IsNil checks if the BeaconKitBlock is nil or not.
func (b *BeaconKitBlockCapella) IsNil() bool {
	return b == nil
}

// Version returns the version of the block.
func (b *BeaconKitBlockCapella) Version() int {
	return version.Capella
}

// AttachExecution attaches the given execution data to the block.
func (b *BeaconKitBlockCapella) AttachExecution(
	executionData engine.ExecutionPayload,
) error {
	var ok bool
	b.Body.ExecutionPayload, ok = executionData.ToProto().(*v1.ExecutionPayloadCapella)
	// b.Body.ExecutionPayload, err = executionData.PbCapella()
	if !ok {
		return errors.New("failed to convert execution data to capella payload")
	}

	// TODO: this needs to be done better, really hood rn.
	payloadValueBz := make([]byte, 32) //nolint:gomnd // 32 bytes for uint256.
	copy(payloadValueBz, executionData.GetValue().Bytes())
	b.PayloadValue = payloadValueBz
	return nil
}

// Execution returns the execution data of the block.
func (b *BeaconKitBlockCapella) ExecutionPayload() (engine.ExecutionPayload, error) {
	return &enginev1.ExecutionPayloadContainer{
		Payload: &enginev1.ExecutionPayloadContainer_Capella{
			Capella: b.GetBody().GetExecutionPayload(),
		},
		PayloadValue: b.GetPayloadValue(),
	}, nil
}

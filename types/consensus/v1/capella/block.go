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

package capella

import (
	"math/big"

	"github.com/itsdevbear/bolaris/beacon/state"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	github_com_prysmaticlabs_prysm_v4_consensus_types_primitives "github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// BeaconKitBlock implements the BeaconKitBlock interface.
var _ interfaces.BeaconKitBlock = (*BeaconKitBlockCapella)(nil)

// BeaconKitBlockFromState assembles a new beacon block
// from the given state and execution data.
func BeaconKitBlockFromState(
	beaconState state.ReadOnlyBeaconState,
	executionData interfaces.ExecutionData,
) (interfaces.BeaconKitBlock, error) {
	return NewBeaconKitBlock(
		beaconState.Slot(),
		executionData,
	)
}

// BeaconKitBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconKitBlock(
	slot github_com_prysmaticlabs_prysm_v4_consensus_types_primitives.Slot,
	executionData interfaces.ExecutionData,
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
	executionData interfaces.ExecutionData,
) error {
	value, err := executionData.ValueInWei()
	if err != nil {
		return err
	}

	b.Body.ExecutionPayload, err = executionData.PbCapella()
	if err != nil {
		return err
	}

	// TODO: this needs to be done better, really hood rn.
	payloadValueBz := make([]byte, 32)     //nolint:gomnd // 32 bytes for uint256.
	copy(payloadValueBz, (*value).Bytes()) //nolint:gocritic // we need to copy the bytes.
	b.PayloadValue = payloadValueBz
	return err
}

// Execution returns the execution data of the block.
func (b *BeaconKitBlockCapella) Execution() (interfaces.ExecutionData, error) {
	return blocks.WrappedExecutionPayloadCapella(b.GetBody().GetExecutionPayload(),
		new(big.Int).SetBytes(b.GetPayloadValue()))
}

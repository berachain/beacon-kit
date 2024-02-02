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

package v1

import (
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/itsdevbear/bolaris/types/state"
)

// BaseBeaconKitBlock implements the BeaconKitBlock interface.
var _ interfaces.BeaconKitBlock = (*BaseBeaconKitBlock)(nil)

// BaseBeaconKitBlockFromState assembles a new beacon block
// from the given state and execution data.
func BaseBeaconKitBlockFromState(
	beaconState state.ReadOnlyBeaconState,
	executionData interfaces.ExecutionData,
) (interfaces.BeaconKitBlock, error) {
	return NewBaseBeaconKitBlock(
		beaconState.Slot(),
		beaconState.Time(),
		executionData,
		beaconState.Version(),
	)
}

// BaseBeaconKitBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBaseBeaconKitBlock(
	slot Slot,
	time uint64,
	executionData interfaces.ExecutionData,
	version int,
) (interfaces.BeaconKitBlock, error) {
	block := &BaseBeaconKitBlock{
		Slot:    slot,
		Time:    time,
		Version: uint64(version),
	}
	if executionData != nil {
		if err := block.AttachPayload(executionData); err != nil {
			return nil, err
		}
	}
	return block, nil
}

// NewEmptyBeaconKitBlockFromState assembles a new beacon block
// with no execution data from the given state.
func NewEmptyBeaconKitBlockFromState(
	beaconState state.BeaconState,
) (interfaces.BeaconKitBlock, error) {
	return NewEmptyBeaconKitBlock(
		beaconState.Slot(),
		beaconState.Time(),
		beaconState.Version(),
	)
}

// NewEmptyBeaconKitBlock assembles a new beacon block
// with no execution data.
func NewEmptyBeaconKitBlock(
	slot Slot,
	time uint64,
	version int,
) (interfaces.BeaconKitBlock, error) {
	return NewBaseBeaconKitBlock(slot, time, nil, version)
}

// ReadOnlyBeaconKitBlockFromABCIRequest assembles a
// new read-only beacon block by extracting a marshalled
// block out of an ABCI request.
func ReadOnlyBeaconKitBlockFromABCIRequest(
	req interfaces.ABCIRequest,
	bzIndex uint,
) (interfaces.ReadOnlyBeaconKitBlock, error) {
	// Extract the marshalled payload from the proposal
	txs := req.GetTxs()
	lenTxs := len(txs)
	if lenTxs == 0 {
		return nil, ErrNoBeaconBlockInProposal
	}
	if bzIndex >= uint(len(txs)) {
		return nil, ErrBzIndexOutOfBounds
	}
	block := BaseBeaconKitBlock{}
	if err := block.Unmarshal(txs[bzIndex]); err != nil {
		return nil, err
	}
	return &block, nil
}

// AttachPayload attaches the given execution data to the block.
func (b *BaseBeaconKitBlock) AttachPayload(
	executionData interfaces.ExecutionData,
) error {
	execData, err := executionData.MarshalSSZ()
	if err != nil {
		return err
	}

	value, err := executionData.ValueInGwei()
	if err != nil {
		return err
	}

	b.ExecData = execData
	b.Value = Gwei(value)
	return nil
}

// IsNil checks if the BaseBeaconKitBlock is nil or not.
func (b *BaseBeaconKitBlock) IsNil() bool {
	return b == nil
}

// SetExecutionData sets the execution data of the block.
func (b *BaseBeaconKitBlock) SetExecutionData(executionData interfaces.ExecutionData) error {
	var err error
	b.ExecData, err = executionData.MarshalSSZ()
	return err
}

// ExecutionData returns the execution data of the block.
func (b *BaseBeaconKitBlock) ExecutionData() interfaces.ExecutionData {
	// Safe to ignore the error since we successfully marshalled the data before.
	data, _ := BytesToExecutionData(b.ExecData, b.Value, int(b.Version))
	return data
}

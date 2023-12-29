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
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/math"
)

// BaseBeaconKitBlock implements the BeaconKitBlock interface.
var _ interfaces.BeaconKitBlock = (*BaseBeaconKitBlock)(nil)

// NewBaseBeaconKitBlock creates a new beacon block.
func NewBaseBeaconKitBlock(
	slot primitives.Slot,
	time uint64,
	executionData interfaces.ExecutionData,
	version int,
) (interfaces.BeaconKitBlock, error) {
	execData, err := executionData.MarshalSSZ()
	if err != nil {
		return nil, err
	}

	value, err := executionData.ValueInGwei()
	if err != nil {
		return nil, err
	}

	return &BaseBeaconKitBlock{
		Slot:     slot,
		Time:     time,
		ExecData: execData,
		Value:    math.Gwei(value),
		Version:  uint64(version),
	}, nil
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

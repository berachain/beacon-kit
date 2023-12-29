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
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/math"
)

// BeaconKitBlock is an interface for beacon blocks.
type BeaconKitBlock interface {
	// SetExecutionData sets the execution data of the block.
	SetExecutionData(executionData interfaces.ExecutionData) error
	// ExecutionData returns the execution data of the block.
	ExecutionData() (interfaces.ExecutionData, error)
}

// Ensure BaseBeaconKitBlock implements BeaconKitBlock interface.
var _ BeaconKitBlock = (*BaseBeaconKitBlock)(nil)

func NewBaseBeaconKitBlock(
	executionData interfaces.ExecutionData,
	value math.Gwei,
	version int,
) (*BaseBeaconKitBlock, error) {
	execData, err := executionData.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return &BaseBeaconKitBlock{
		ExecData: execData,
		Value:    uint64(value),
		Version:  uint64(version),
	}, nil
}

// SetExecutionData sets the execution data of the block.
func (b *BaseBeaconKitBlock) SetExecutionData(executionData interfaces.ExecutionData) error {
	var err error
	b.ExecData, err = executionData.MarshalSSZ()
	return err
}

// ExecutionData returns the execution data of the block.
func (b *BaseBeaconKitBlock) ExecutionData() (interfaces.ExecutionData, error) {
	return BytesToExecutionData(b.ExecData, math.Gwei(b.Value), int(b.Version))
}

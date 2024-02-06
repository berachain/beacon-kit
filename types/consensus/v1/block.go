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
	"errors"
	"math/big"

	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/itsdevbear/bolaris/types/state"
)

// BeaconKitBlock implements the BeaconKitBlock interface.
var _ interfaces.BeaconKitBlock = (*BeaconKitBlock)(nil)

// BeaconKitBlockFromState assembles a new beacon block
// from the given state and execution data.
func BeaconKitBlockFromState(
	beaconState state.ReadOnlyBeaconState,
	executionData interfaces.ExecutionData,
) (interfaces.BeaconKitBlock, error) {
	return NewBeaconKitBlock(
		beaconState.Slot(),
		executionData,
		beaconState.Version(),
	)
}

// BeaconKitBlock assembles a new beacon block from
// the given slot, time, execution data, and version.
func NewBeaconKitBlock(
	slot Slot,
	executionData interfaces.ExecutionData,
	version int,
) (interfaces.BeaconKitBlock, error) {
	block := &BeaconKitBlock{
		Slot: slot,
		Body: &BeaconKitBlock_BlockBodyGeneric{
			BlockBodyGeneric: &BeaconBlockBody{
				//#nosec:G701 // won't overflow, version is never negative.
				Version: int64(version),
			},
		},
	}
	if executionData != nil {
		if err := block.AttachExecution(executionData); err != nil {
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
		beaconState.Version(),
	)
}

// NewEmptyBeaconKitBlock assembles a new beacon block
// with no execution data.
func NewEmptyBeaconKitBlock(
	slot Slot,
	version int,
) (interfaces.BeaconKitBlock, error) {
	return NewBeaconKitBlock(slot, nil, version)
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
	block := BeaconKitBlock{}
	if err := block.Unmarshal(txs[bzIndex]); err != nil {
		return nil, err
	}
	return &block, nil
}

// IsNil checks if the BeaconKitBlock is nil or not.
func (b *BeaconKitBlock) IsNil() bool {
	return b == nil
}

// AttachExecution attaches the given execution data to the block.
func (b *BeaconKitBlock) AttachExecution(
	executionData interfaces.ExecutionData,
) error {
	execData, err := executionData.MarshalSSZ()
	if err != nil {
		return err
	}

	value, err := executionData.ValueInWei()
	if err != nil {
		return err
	}

	b.Body.(*BeaconKitBlock_BlockBodyGeneric).BlockBodyGeneric.ExecutionPayload = execData
	b.PayloadValue = (*value).String() //nolint:gocritic // suggestion doesn't compile.
	return nil
}

// Execution returns the execution data of the block.
func (b *BeaconKitBlock) Execution() (interfaces.ExecutionData, error) {
	// Safe to ignore the error since we successfully marshalled the data before.
	value, ok := big.NewInt(0).SetString(b.PayloadValue, 10) //nolint:gomnd // base 10.
	if !ok {
		return nil, errors.New("failed to convert payload value to big.Int")
	}
	return BytesToExecutionData(
		b.GetBlockBodyGeneric().ExecutionPayload,
		Wei(value),
		int(b.GetBlockBodyGeneric().Version))
}

func (b *BeaconKitBlock) Version() int {
	return int(b.GetBlockBodyGeneric().Version)
}

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

package types

import (
	enginetypes "github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/berachain/beacon-kit/mod/trie"
	merkleize "github.com/berachain/beacon-kit/mod/trie/merkleize"
	"github.com/cockroachdb/errors"
)

// BeaconBlockBodyDeneb represents the body of a beacon block in the Deneb
// chain.
type BeaconBlockBodyDeneb struct {
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal primitives.BLSSignature `ssz-size:"96"`

	// Graffiti is for a fun message or meme.
	Graffiti [32]byte `ssz-size:"32"`

	// Deposits is the list of deposits included in the body.
	Deposits []*Deposit `ssz-max:"16"`

	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *enginetypes.ExecutableDataDeneb

	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments kzg.Commitments `ssz-size:"?,48" ssz-max:"16"`
}

// IsNil checks if the BeaconBlockBodyDeneb is nil.
func (b *BeaconBlockBodyDeneb) IsNil() bool {
	return b == nil
}

// GetBlobKzgCommitments returns the BlobKzgCommitments of the Body.
func (b *BeaconBlockBodyDeneb) GetBlobKzgCommitments() kzg.Commitments {
	return b.BlobKzgCommitments
}

// GetGraffiti returns the Graffiti of the Body.
func (b *BeaconBlockBodyDeneb) GetGraffiti() [32]byte {
	return b.Graffiti
}

// GetRandaoReveal returns the RandaoReveal of the Body.
func (b *BeaconBlockBodyDeneb) GetRandaoReveal() primitives.BLSSignature {
	return b.RandaoReveal
}

// GetExecutionPayload returns the ExecutionPayload of the Body.
//
//nolint:lll
func (b *BeaconBlockBodyDeneb) GetExecutionPayload() enginetypes.ExecutionPayload {
	return b.ExecutionPayload
}

// GetDeposits returns the Deposits of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) GetDeposits() []*Deposit {
	return b.Deposits
}

// SetDeposits sets the Deposits of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetDeposits(deposits []*Deposit) {
	b.Deposits = deposits
}

// SetExecutionData sets the ExecutionData of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetExecutionData(
	executionData enginetypes.ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.(*enginetypes.ExecutableDataDeneb)
	if !ok {
		return errors.New("invalid execution data type")
	}
	return nil
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetBlobKzgCommitments(
	commitments kzg.Commitments,
) {
	b.BlobKzgCommitments = commitments
}

func GetTopLevelRoots(b BeaconBlockBody) ([][]byte, error) {
	layer := make([][]byte, BodyLength)
	for i := range layer {
		layer[i] = make([]byte, constants.RootLength)
	}

	randao := b.GetRandaoReveal()
	root, err := merkleize.ByteSliceSSZ(randao[:])
	if err != nil {
		return nil, err
	}
	copy(layer[0], root[:])

	// graffiti
	root = b.GetGraffiti()
	copy(layer[1], root[:])

	// Deposits
	dep := b.GetDeposits()
	//nolint:gomnd // TODO: Config
	maxDepositsPerBlock := uint64(16)
	root, err = merkleize.ListSSZ(dep, maxDepositsPerBlock)
	if err != nil {
		return nil, err
	}
	copy(layer[2], root[:])

	// Execution Payload
	rt, err := b.GetExecutionPayload().HashTreeRoot()
	if err != nil {
		return nil, err
	}
	copy(layer[3], rt[:])

	// KZG commitments is not needed
	return layer, nil
}

func GetBlobKzgCommitmentsRoot(commitments [][48]byte) ([32]byte, error) {
	commitmentsLeaves := LeavesFromCommitments(commitments)
	commitmentsSparse, err := trie.NewFromItems(
		commitmentsLeaves,
		LogMaxBlobCommitments,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return commitmentsSparse.HashTreeRoot()
}

func (b *BeaconBlockBodyDeneb) AttachExecution(
	executionData enginetypes.ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.(*enginetypes.ExecutableDataDeneb)
	if !ok {
		return errors.New("invalid execution data type")
	}
	return nil
}

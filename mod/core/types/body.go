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
	stdmath "math"
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/ssz"
	"github.com/cockroachdb/errors"
)

//nolint:gochecknoglobals // I'd prefer globals over magic numbers.
var (
	// BodyLengthDeneb is the number of fields in the BeaconBlockBodyDeneb
	// struct.
	//#nosec:G701 // realistically won't exceed 255 fields.
	BodyLengthDeneb = uint8(reflect.TypeOf(BeaconBlockBodyDeneb{}).NumField())

	// LogBodyLengthDeneb is the Log_2 of BodyLength (6).
	//#nosec:G701 // realistically won't exceed 255 fields.
	LogBodyLengthDeneb = uint8(
		stdmath.Ceil(stdmath.Log2(float64(BodyLengthDeneb))),
	)

	// KZGPosition is the position of BlobKzgCommitments in the block body.
	KZGPositionDeneb = uint64(BodyLengthDeneb - 1)
)

// BeaconBlockBodyDeneb represents the body of a beacon block in the Deneb
// chain.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path body.go -objs BeaconBlockBodyDeneb -include ../../primitives,../../primitives/math,../../primitives/kzg,../../primitives-engine,../../primitives,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output body.ssz.go
type BeaconBlockBodyDeneb struct {
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal primitives.BLSSignature `ssz-size:"96"`

	// Graffiti is for a fun message or meme.
	Graffiti [32]byte `ssz-size:"32"`

	// Deposits is the list of deposits included in the body.
	Deposits []*primitives.Deposit `ssz-max:"16"`

	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *engineprimitives.ExecutableDataDeneb

	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []kzg.Commitment `ssz-size:"?,48" ssz-max:"16"`
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
func (b *BeaconBlockBodyDeneb) GetGraffiti() primitives.Bytes32 {
	return b.Graffiti
}

// GetRandaoReveal returns the RandaoReveal of the Body.
func (b *BeaconBlockBodyDeneb) GetRandaoReveal() primitives.BLSSignature {
	return b.RandaoReveal
}

// GetExecutionPayload returns the ExecutionPayload of the Body.
//
//nolint:lll
func (b *BeaconBlockBodyDeneb) GetExecutionPayload() engineprimitives.ExecutionPayload {
	return b.ExecutionPayload
}

// GetDeposits returns the Deposits of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) GetDeposits() []*primitives.Deposit {
	return b.Deposits
}

// SetDeposits sets the Deposits of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetDeposits(
	deposits []*primitives.Deposit,
) {
	b.Deposits = deposits
}

// SetExecutionData sets the ExecutionData of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetExecutionData(
	executionData engineprimitives.ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.(*engineprimitives.ExecutableDataDeneb)
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

// GetTopLevelRoots returns the top-level roots of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) GetTopLevelRoots() ([][32]byte, error) {
	layer := make([][32]byte, BodyLengthDeneb)
	var err error
	randao := b.GetRandaoReveal()
	layer[0], err = ssz.MerkleizeByteSlice[math.U64, [32]byte](randao[:])
	if err != nil {
		return nil, err
	}

	// graffiti
	layer[1] = b.GetGraffiti()

	//nolint:mnd // TODO: Config
	maxDepositsPerBlock := uint64(16)
	// root, err = dep.HashTreeRoot()
	layer[2], err = ssz.MerkleizeListComposite[any, math.U64](
		b.GetDeposits(),
		maxDepositsPerBlock,
	)
	if err != nil {
		return nil, err
	}

	// Execution Payload
	layer[3], err = b.GetExecutionPayload().HashTreeRoot()
	if err != nil {
		return nil, err
	}

	// KZG commitments is not needed
	return layer, nil
}

func (b *BeaconBlockBodyDeneb) AttachExecution(
	executionData engineprimitives.ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.(*engineprimitives.ExecutableDataDeneb)
	if !ok {
		return errors.New("invalid execution data type")
	}
	return nil
}

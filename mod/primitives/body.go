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

package primitives

import (
	"errors"
	stdmath "math"
	"reflect"

	// engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine".
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/ssz"
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
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./body.go -objs BeaconBlockBodyDeneb -include ./primitives.go,./payload.go,./kzg.go,./bytes.go,./eth1data.go,./math,./execution.go,./deposit.go,./withdrawal_credentials.go,./withdrawal.go,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output body.ssz.go
type BeaconBlockBodyDeneb struct {
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal BLSSignature `ssz-size:"96"`

	// Eth1Data is the data from the Eth1 chain.
	Eth1Data *Eth1Data

	// Graffiti is for a fun message or meme.
	Graffiti [32]byte `ssz-size:"32"`

	// Deposits is the list of deposits included in the body.
	Deposits []*Deposit `ssz-max:"16"`

	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutableDataDeneb

	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []Commitment `ssz-size:"?,48" ssz-max:"16"`
}

// IsNil checks if the BeaconBlockBodyDeneb is nil.
func (b *BeaconBlockBodyDeneb) IsNil() bool {
	return b == nil
}

// GetBlobKzgCommitments returns the BlobKzgCommitments of the Body.
func (b *BeaconBlockBodyDeneb) GetBlobKzgCommitments() Commitments {
	return b.BlobKzgCommitments
}

// GetGraffiti returns the Graffiti of the Body.
func (b *BeaconBlockBodyDeneb) GetGraffiti() Bytes32 {
	return b.Graffiti
}

// GetRandaoReveal returns the RandaoReveal of the Body.
func (b *BeaconBlockBodyDeneb) GetRandaoReveal() BLSSignature {
	return b.RandaoReveal
}

// GetEth1Data returns the Eth1Data of the Body.
func (b *BeaconBlockBodyDeneb) GetEth1Data() *Eth1Data {
	return b.Eth1Data
}

// GetExecutionPayload returns the ExecutionPayload of the Body.
//

func (b *BeaconBlockBodyDeneb) GetExecutionPayload() ExecutionPayload {
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
	executionData ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.(*ExecutableDataDeneb)
	if !ok {
		return errors.New("invalid execution data type")
	}
	return nil
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetBlobKzgCommitments(
	commitments Commitments,
) {
	b.BlobKzgCommitments = commitments
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetEth1Data(eth1Data *Eth1Data) {
	b.Eth1Data = eth1Data
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

	layer[1], err = b.Eth1Data.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	// graffiti
	layer[2] = b.GetGraffiti()

	layer[3], err = Deposits(b.GetDeposits()).HashTreeRoot()
	if err != nil {
		return nil, err
	}

	// Execution Payload
	layer[4], err = b.GetExecutionPayload().HashTreeRoot()
	if err != nil {
		return nil, err
	}

	// KZG commitments is not needed
	return layer, nil
}

func (b *BeaconBlockBodyDeneb) AttachExecution(
	executionData ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.(*ExecutableDataDeneb)
	if !ok {
		return errors.New("invalid execution data type")
	}
	return nil
}

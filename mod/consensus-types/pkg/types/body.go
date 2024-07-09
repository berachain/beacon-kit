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
// WdeHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

const (
	// BodyLengthDeneb is the number of fields in the BeaconBlockBodyDeneb
	// struct.
	BodyLengthDeneb uint64 = 6

	// KZGPositionDeneb is the position of BlobKzgCommitments in the block body.
	KZGPositionDeneb = BodyLengthDeneb - 1

	// KZGMerkleIndexDeneb is the merkle index of BlobKzgCommitments' root
	// in the merkle tree built from the block body.
	KZGMerkleIndexDeneb = 26

	// LogsBloomSize is the size of LogsBloom in bytes.
	LogsBloomSize = 256

	// ExtraDataSize is the size of ExtraData in bytes.
	ExtraDataSize = 32
)

type BeaconBlockBody struct {
	RawBeaconBlockBody
}

// Empty returns a new BeaconBlockBody with empty fields for the given fork version.
func (b *BeaconBlockBody) Empty(forkVersion uint32) *BeaconBlockBody {
	switch forkVersion {
	case version.Deneb:
		return &BeaconBlockBody{RawBeaconBlockBody: &BeaconBlockBodyDeneb{
			BeaconBlockBodyBase: BeaconBlockBodyBase{},
			ExecutionPayload: &ExecutableDataDeneb{
				LogsBloom: make([]byte, LogsBloomSize),
				ExtraData: make([]byte, ExtraDataSize),
			},
		}}
	default:
		panic("unsupported fork version")
	}
}

// BlockBodyKZGOffset returns the offset of the KZG commitments in the block
// body.
// TODO: I still feel like we need to clean this up somehow.
func BlockBodyKZGOffset(
	slot math.Slot,
	cs common.ChainSpec,
) uint64 {
	switch cs.ActiveForkVersionForSlot(slot) {
	case version.Deneb:
		return KZGMerkleIndexDeneb * cs.MaxBlobCommitmentsPerBlock()
	default:
		panic("unsupported fork version")
	}
}

// BeaconBlockBodyBase represents the base body of a beacon block that is
// shared between all forks.
type BeaconBlockBodyBase struct {
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal crypto.BLSSignature `ssz-size:"96"`
	// Eth1Data is the data from the Eth1 chain.
	Eth1Data *Eth1Data
	// Graffiti is for a fun message or meme.
	Graffiti [32]byte `ssz-size:"32"`
	// Deposits is the list of deposits included in the body.
	Deposits []*Deposit `              ssz-max:"16"`
}

// GetRandaoReveal returns the RandaoReveal of the Body.
func (b *BeaconBlockBodyBase) GetRandaoReveal() crypto.BLSSignature {
	return b.RandaoReveal
}

// SetRandaoReveal sets the RandaoReveal of the Body.
func (b *BeaconBlockBodyBase) SetRandaoReveal(reveal crypto.BLSSignature) {
	b.RandaoReveal = reveal
}

// GetEth1Data returns the Eth1Data of the Body.
func (b *BeaconBlockBodyBase) GetEth1Data() *Eth1Data {
	return b.Eth1Data
}

// SetEth1Data sets the Eth1Data of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetEth1Data(eth1Data *Eth1Data) {
	b.Eth1Data = eth1Data
}

// GetGraffiti returns the Graffiti of the Body.
func (b *BeaconBlockBodyBase) GetGraffiti() bytes.B32 {
	return b.Graffiti
}

// GetDeposits returns the Deposits of the BeaconBlockBodyBase.
func (b *BeaconBlockBodyBase) GetDeposits() []*Deposit {
	return b.Deposits
}

// SetDeposits sets the Deposits of the BeaconBlockBodyBase.
func (b *BeaconBlockBodyBase) SetDeposits(deposits []*Deposit) {
	b.Deposits = deposits
}

// BeaconBlockBodyDeneb represents the body of a beacon block in the Deneb
// chain.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./body.go -objs BeaconBlockBodyDeneb -include ../../../primitives/pkg/crypto,./payload.go,../../../primitives/pkg/eip4844,../../../primitives/pkg/bytes,./eth1data.go,../../../primitives/pkg/math,../../../primitives/pkg/common,./deposit.go,../../../engine-primitives/pkg/engine-primitives/withdrawal.go,./withdrawal_credentials.go,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output body.ssz.go
type BeaconBlockBodyDeneb struct {
	BeaconBlockBodyBase
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutableDataDeneb
	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []eip4844.KZGCommitment `ssz-size:"?,48" ssz-max:"16"`
}

// IsNil checks if the BeaconBlockBodyDeneb is nil.
func (b *BeaconBlockBodyDeneb) IsNil() bool {
	return b == nil
}

// GetExecutionPayload returns the ExecutionPayload of the Body.
func (
	b *BeaconBlockBodyDeneb,
) GetExecutionPayload() *ExecutionPayload {
	return &ExecutionPayload{InnerExecutionPayload: b.ExecutionPayload}
}

// SetExecutionData sets the ExecutionData of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetExecutionData(
	executionData *ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.
		InnerExecutionPayload.(*ExecutableDataDeneb)
	if !ok {
		return errors.New("invalid execution data type")
	}
	return nil
}

// GetBlobKzgCommitments returns the BlobKzgCommitments of the Body.
func (
	b *BeaconBlockBodyDeneb,
) GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash] {
	return b.BlobKzgCommitments
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetBlobKzgCommitments(
	commitments eip4844.KZGCommitments[common.ExecutionHash],
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

	layer[1], err = b.Eth1Data.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	layer[2] = b.GetGraffiti()

	layer[3], err = Deposits(b.GetDeposits()).HashTreeRoot()
	if err != nil {
		return nil, err
	}

	layer[4], err = b.GetExecutionPayload().HashTreeRoot()
	if err != nil {
		return nil, err
	}

	// KZG commitments is not needed
	return layer, nil
}

// Length returns the number of fields in the BeaconBlockBodyDeneb struct.
func (b *BeaconBlockBodyDeneb) Length() uint64 {
	return BodyLengthDeneb
}

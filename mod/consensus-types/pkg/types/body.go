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
	"unsafe"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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

	// ExtraDataSize is the size of ExtraData in bytes.
	ExtraDataSize = 32
)

// Empty returns a new BeaconBlockBody with empty fields
// for the given fork version.
func (b *BeaconBlockBody) Empty(forkVersion uint32) *BeaconBlockBody {
	switch forkVersion {
	case version.Deneb:
		return &BeaconBlockBody{
			ExecutionPayload: &ExecutionPayload{
				ExtraData: make([]byte, ExtraDataSize),
			},
		}
	case version.DenebPlus:
		panic("unsupported fork version")
		// return &BeaconBlockBody{RawBeaconBlockBody:
		// &BeaconBlockBodyDenebPlus{
		// 	BeaconBlockBodyBase: BeaconBlockBodyBase{},
		// 	ExecutionPayload: &ExecutionPayload{
		// 		ExtraData: make([]byte, ExtraDataSize),
		// 	},
		// }}
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

// BeaconBlockBody represents the body of a beacon block in the Deneb
// chain.
type BeaconBlockBody struct {
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal crypto.BLSSignature
	// Eth1Data is the data from the Eth1 chain.
	Eth1Data *Eth1Data
	// Graffiti is for a fun message or meme.
	Graffiti [32]byte
	// Deposits is the list of deposits included in the body.
	Deposits []*Deposit
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutionPayload
	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []eip4844.KZGCommitment
}

// IsNil checks if the BeaconBlockBody is nil.
func (b *BeaconBlockBody) IsNil() bool {
	return b == nil
}

// GetExecutionPayload returns the ExecutionPayload of the Body.
func (
	b *BeaconBlockBody,
) GetExecutionPayload() *ExecutionPayload {
	return b.ExecutionPayload
}

// SetExecutionPayload sets the ExecutionData of the BeaconBlockBody.
func (b *BeaconBlockBody) SetExecutionPayload(
	executionData *ExecutionPayload,
) {
	b.ExecutionPayload = executionData
}

// GetBlobKzgCommitments returns the BlobKzgCommitments of the Body.
func (
	b *BeaconBlockBody,
) GetBlobKzgCommitments() eip4844.KZGCommitments[gethprimitives.ExecutionHash] {
	return b.BlobKzgCommitments
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBody.
func (b *BeaconBlockBody) SetBlobKzgCommitments(
	commitments eip4844.KZGCommitments[gethprimitives.ExecutionHash],
) {
	b.BlobKzgCommitments = commitments
}

// SetEth1Data sets the Eth1Data of the BeaconBlockBody.
func (b *BeaconBlockBody) SetEth1Data(eth1Data *Eth1Data) {
	b.Eth1Data = eth1Data
}

// SetDeposits is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody) GetAttestations() []*AttestationData {
	panic("not implemented")
}

// SetDeposits is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody) SetAttestations(_ []*AttestationData) {
	panic("not implemented")
}

// GetSlashingInfo is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody) GetSlashingInfo() []*SlashingInfo {
	panic("not implemented")
}

// SetSlashingInfo is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody) SetSlashingInfo(_ []*SlashingInfo) {
	panic("not implemented")
}

// GetTopLevelRoots returns the top-level roots of the BeaconBlockBody.
func (b *BeaconBlockBody) GetTopLevelRoots() ([][32]byte, error) {
	var (
		err   error
		layer = make([]common.Root, BodyLengthDeneb)
	)

	layer[0], err = b.GetRandaoReveal().HashTreeRoot()
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
	//#nosec:G103 // Okay to go from common.Root to [32]byte.
	return *(*[][32]byte)(unsafe.Pointer(&layer)), nil
}

// Length returns the number of fields in the BeaconBlockBody struct.
func (b *BeaconBlockBody) Length() uint64 {
	return BodyLengthDeneb
}

// GetRandaoReveal returns the RandaoReveal of the Body.
func (b *BeaconBlockBody) GetRandaoReveal() crypto.BLSSignature {
	return b.RandaoReveal
}

// SetRandaoReveal sets the RandaoReveal of the Body.
func (b *BeaconBlockBody) SetRandaoReveal(reveal crypto.BLSSignature) {
	b.RandaoReveal = reveal
}

// GetEth1Data returns the Eth1Data of the Body.
func (b *BeaconBlockBody) GetEth1Data() *Eth1Data {
	return b.Eth1Data
}

// GetGraffiti returns the Graffiti of the Body.
func (b *BeaconBlockBody) GetGraffiti() common.Bytes32 {
	return b.Graffiti
}

// SetGraffiti sets the Graffiti of the Body.
func (b *BeaconBlockBody) SetGraffiti(graffiti common.Bytes32) {
	b.Graffiti = graffiti
}

// GetDeposits returns the Deposits of the BeaconBlockBody.
func (b *BeaconBlockBody) GetDeposits() []*Deposit {
	return b.Deposits
}

// SetDeposits sets the Deposits of the BeaconBlockBody.
func (b *BeaconBlockBody) SetDeposits(deposits []*Deposit) {
	b.Deposits = deposits
}

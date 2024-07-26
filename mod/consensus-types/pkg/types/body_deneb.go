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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
)

// BeaconBlockBodyDeneb represents the body of a beacon block in the Deneb
// chain.
type BeaconBlockBodyDeneb struct {
	BeaconBlockBodyBase
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutionPayload
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
	return b.ExecutionPayload
}

// SetExecutionPayload sets the ExecutionData of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetExecutionPayload(
	executionData *ExecutionPayload,
) {
	b.ExecutionPayload = executionData
}

// GetBlobKzgCommitments returns the BlobKzgCommitments of the Body.
func (
	b *BeaconBlockBodyDeneb,
) GetBlobKzgCommitments() eip4844.KZGCommitments[gethprimitives.ExecutionHash] {
	return b.BlobKzgCommitments
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetBlobKzgCommitments(
	commitments eip4844.KZGCommitments[gethprimitives.ExecutionHash],
) {
	b.BlobKzgCommitments = commitments
}

// SetEth1Data sets the Eth1Data of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) SetEth1Data(eth1Data *Eth1Data) {
	b.Eth1Data = eth1Data
}

// SetDeposits is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBodyDeneb) GetAttestations() []*AttestationData {
	panic("not implemented")
}

// SetDeposits is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBodyDeneb) SetAttestations(_ []*AttestationData) {
	panic("not implemented")
}

// GetSlashingInfo is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBodyDeneb) GetSlashingInfo() []*SlashingInfo {
	panic("not implemented")
}

// SetSlashingInfo is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBodyDeneb) SetSlashingInfo(_ []*SlashingInfo) {
	panic("not implemented")
}

// GetTopLevelRoots returns the top-level roots of the BeaconBlockBodyDeneb.
func (b *BeaconBlockBodyDeneb) GetTopLevelRoots() ([][32]byte, error) {
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

// Length returns the number of fields in the BeaconBlockBodyDeneb struct.
func (b *BeaconBlockBodyDeneb) Length() uint64 {
	return BodyLengthDeneb
}

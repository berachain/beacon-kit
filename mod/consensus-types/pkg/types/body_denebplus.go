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

// BeaconBlockBodyDenebPlus represents the body of a beacon block in the Deneb
// chain. This is a temporary struct to be used until mainnet.
type BeaconBlockBodyDenebPlus struct {
	BeaconBlockBodyBase
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutionPayload
	// TODO: Put this in BeaconBlockBodyBase for mainnet.
	// Attestations is the list of attestations included in the body.
	Attestations []*AttestationData `ssz-max:"256"`
	// SlashingInfo is the list of slashing info included in the body.
	SlashingInfo []*SlashingInfo `ssz-max:"256"`
	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []eip4844.KZGCommitment `ssz-max:"16"  ssz-size:"?,48"`
}

// IsNil checks if the BeaconBlockBodyDenebPlus is nil.
func (b *BeaconBlockBodyDenebPlus) IsNil() bool {
	return b == nil
}

// GetExecutionPayload returns the ExecutionPayload of the Body.
func (
	b *BeaconBlockBodyDenebPlus,
) GetExecutionPayload() *ExecutionPayload {
	return b.ExecutionPayload
}

// SetExecutionPayload sets the ExecutionData of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) SetExecutionPayload(
	executionData *ExecutionPayload,
) {
	b.ExecutionPayload = executionData
}

// GetBlobKzgCommitments returns the BlobKzgCommitments of the Body.
func (
	b *BeaconBlockBodyDenebPlus,
) GetBlobKzgCommitments() eip4844.KZGCommitments[gethprimitives.ExecutionHash] {
	return b.BlobKzgCommitments
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) SetBlobKzgCommitments(
	commitments eip4844.KZGCommitments[gethprimitives.ExecutionHash],
) {
	b.BlobKzgCommitments = commitments
}

// GetAttestations returns the Attestations of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) GetAttestations() []*AttestationData {
	return b.Attestations
}

// SetAttestations sets the Attestations of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) SetAttestations(
	attestations []*AttestationData,
) {
	b.Attestations = attestations
}

// GetSlashingInfo returns the SlashingInfo of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) GetSlashingInfo() []*SlashingInfo {
	return b.SlashingInfo
}

// SetSlashingInfo sets the SlashingInfo of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) SetSlashingInfo(
	slashingInfo []*SlashingInfo,
) {
	b.SlashingInfo = slashingInfo
}

// SetEth1Data sets the Eth1Data of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) SetEth1Data(eth1Data *Eth1Data) {
	b.Eth1Data = eth1Data
}

// GetTopLevelRoots returns the top-level roots of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) GetTopLevelRoots() ([][32]byte, error) {
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

// Length returns the number of fields in the BeaconBlockBodyDenebPlus struct.
func (b *BeaconBlockBodyDenebPlus) Length() uint64 {
	return BodyLengthDeneb
}

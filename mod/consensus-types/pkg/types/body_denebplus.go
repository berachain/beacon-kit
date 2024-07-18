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

	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
)

// BeaconBlockBodyDenebPlus represents the body of a beacon block in the Deneb
// chain. This is a temporary struct to be used until mainnet.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./body_denebplus.go -objs BeaconBlockBodyDenebPlus -include ./body.go,../../../primitives/pkg/crypto,./payload.go,../../../primitives/pkg/eip4844,../../../primitives/pkg/bytes,./eth1data.go,../../../primitives/pkg/math,../../../primitives/pkg/common,./deposit.go,../../../engine-primitives/pkg/engine-primitives/withdrawal.go,./withdrawal_credentials.go,./attestation_data.go,./slashing_info.go,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output body_denebplus.ssz.go
type BeaconBlockBodyDenebPlus struct {
	BeaconBlockBodyBase
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutableDataDeneb
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
	return &ExecutionPayload{InnerExecutionPayload: b.ExecutionPayload}
}

// SetExecutionData sets the ExecutionData of the BeaconBlockBodyDenebPlus.
func (b *BeaconBlockBodyDenebPlus) SetExecutionData(
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

// SetEth1Data sets the Eth1Data of the BeaconBlockBodyDeneb.
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

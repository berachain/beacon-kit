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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/karalabe/ssz"
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
	// metadata
	//
	// version represents the fork version of the body.
	version uint32

	// contents
	//
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal crypto.BLSSignature
	// Eth1Data is the data from the Eth1 chain.
	Eth1Data *Eth1Data
	// Graffiti is for a fun message or meme.
	Graffiti [32]byte
	// Deposits is the list of deposits included in the body.
	Deposits []*Deposit
}

// Empty returns a new BeaconBlockBody with empty fields
// for the given fork version.
func (b *BeaconBlockBody) Empty(forkVersion uint32) *BeaconBlockBody {
	return &BeaconBlockBody{
		version: forkVersion,
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

// SizeSSZ returns the size of the Eth1Data object in SSZ encoding.
func (b *BeaconBlockBody) SizeSSZ(fixed bool) uint32 {
	size := uint32(0)
	size += uint32(b.RandaoReveal.SizeSSZ())
	size += b.Eth1Data.SizeSSZ()
	size += uint32(b.GetGraffiti().SizeSSZ())
	size += ssz.SizeSliceOfStaticObjects(b.Deposits)
	return size
}

// DefineSSZ defines the SSZ encoding for the BeaconBlockBody object.
func (b *BeaconBlockBody) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &b.RandaoReveal)
	ssz.DefineStaticObject(codec, &b.Eth1Data)
	ssz.DefineStaticBytes(codec, &b.Graffiti)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.Deposits, constants.MaxDepositsPerBlock)
}

// HashTreeRoot computes the SSZ hash tree root of the BeaconBlockBody object.
func (b *BeaconBlockBody) HashTreeRoot() [32]byte {
	return ssz.HashSequential(b)
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

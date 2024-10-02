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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	fastssz "github.com/ferranbt/fastssz"
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

	// ExtraDataSize is the size of ExtraData in bytes.
	ExtraDataSize = 32
)

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
type BeaconBlockBody[LogT interface {
	GetData() []byte
}] struct {
	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal crypto.BLSSignature
	// Eth1Data is the data from the Eth1 chain.
	Eth1Data *Eth1Data
	// Graffiti is for a fun message or meme.
	Graffiti [32]byte
	// Deposits is the list of deposits included in the body.
	Deposits []*Deposit[LogT]
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutionPayload
	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []eip4844.KZGCommitment
}

// Empty returns a new BeaconBlockBody with empty fields
// for the given fork version.
func (b *BeaconBlockBody[LogT]) Empty(
	forkVersion uint32,
) *BeaconBlockBody[LogT] {
	switch forkVersion {
	case version.Deneb:
		return &BeaconBlockBody[LogT]{
			Eth1Data: new(Eth1Data),
			ExecutionPayload: &ExecutionPayload{
				ExtraData: make([]byte, ExtraDataSize),
			},
		}
	default:
		panic("unsupported fork version")
	}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlockBody in SSZ.
func (b *BeaconBlockBody[_]) SizeSSZ(fixed bool) uint32 {
	var size uint32 = 96 + 72 + 32 + 4 + 4 + 4
	if fixed {
		return size
	}

	size += ssz.SizeSliceOfStaticObjects(b.Deposits)
	size += ssz.SizeDynamicObject(b.ExecutionPayload)
	size += ssz.SizeSliceOfStaticBytes(b.BlobKzgCommitments)
	return size
}

// DefineSSZ defines the SSZ serialization of the BeaconBlockBody.
//
//nolint:mnd // TODO: chainspec.
func (b *BeaconBlockBody[_]) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineStaticBytes(codec, &b.RandaoReveal)
	ssz.DefineStaticObject(codec, &b.Eth1Data)
	ssz.DefineStaticBytes(codec, &b.Graffiti)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.Deposits, 16)
	ssz.DefineDynamicObjectOffset(codec, &b.ExecutionPayload)
	ssz.DefineSliceOfStaticBytesOffset(codec, &b.BlobKzgCommitments, 16)

	// Define the dynamic data (fields)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.Deposits, 16)
	ssz.DefineDynamicObjectContent(codec, &b.ExecutionPayload)
	ssz.DefineSliceOfStaticBytesContent(codec, &b.BlobKzgCommitments, 16)
}

// MarshalSSZ serializes the BeaconBlockBody to SSZ-encoded bytes.
func (b *BeaconBlockBody[_]) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ deserializes the BeaconBlockBody from SSZ-encoded bytes.
func (b *BeaconBlockBody[_]) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot returns the SSZ hash tree root of the BeaconBlockBody.
func (b *BeaconBlockBody[_]) HashTreeRoot() common.Root {
	return ssz.HashConcurrent(b)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo serializes the BeaconBlockBody into a writer.
func (b *BeaconBlockBody[_]) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := b.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the BeaconBlockBody object with a hasher.
//
//nolint:mnd // todo fix.
func (b *BeaconBlockBody[_]) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'RandaoReveal'
	hh.PutBytes(b.RandaoReveal[:])

	// Field (1) 'Eth1Data'
	if b.Eth1Data == nil {
		b.Eth1Data = new(Eth1Data)
	}
	if err := b.Eth1Data.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (2) 'Graffiti'
	hh.PutBytes(b.Graffiti[:])

	// Field (3) 'Deposits'
	{
		subIndx := hh.Index()
		num := uint64(len(b.Deposits))
		if num > 16 {
			return fastssz.ErrIncorrectListSize
		}
		for _, elem := range b.Deposits {
			if err := elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
	}

	// Field (4) 'ExecutionPayload'
	if err := b.ExecutionPayload.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (5) 'BlobKzgCommitments'
	{
		if size := len(b.BlobKzgCommitments); size > 16 {
			return fastssz.ErrListTooBigFn(
				"BeaconBlockBody.BlobKzgCommitments",
				size,
				16,
			)
		}
		subIndx := hh.Index()
		for _, i := range b.BlobKzgCommitments {
			hh.PutBytes(i[:])
		}
		numItems := uint64(len(b.BlobKzgCommitments))
		hh.MerkleizeWithMixin(subIndx, numItems, 16)
	}

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BeaconBlockBody object.
func (b *BeaconBlockBody[_]) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}

// IsNil checks if the BeaconBlockBody is nil.
func (b *BeaconBlockBody[_]) IsNil() bool {
	return b == nil
}

// GetExecutionPayload returns the ExecutionPayload of the Body.
func (
	b *BeaconBlockBody[_],
) GetExecutionPayload() *ExecutionPayload {
	return b.ExecutionPayload
}

// SetExecutionPayload sets the ExecutionData of the BeaconBlockBody.
func (b *BeaconBlockBody[_]) SetExecutionPayload(
	executionData *ExecutionPayload,
) {
	b.ExecutionPayload = executionData
}

// GetBlobKzgCommitments returns the BlobKzgCommitments of the Body.
func (
	b *BeaconBlockBody[_],
) GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash] {
	return b.BlobKzgCommitments
}

// SetBlobKzgCommitments sets the BlobKzgCommitments of the
// BeaconBlockBody.
func (b *BeaconBlockBody[_]) SetBlobKzgCommitments(
	commitments eip4844.KZGCommitments[common.ExecutionHash],
) {
	b.BlobKzgCommitments = commitments
}

// SetEth1Data sets the Eth1Data of the BeaconBlockBody.
func (b *BeaconBlockBody[_]) SetEth1Data(eth1Data *Eth1Data) {
	b.Eth1Data = eth1Data
}

// GetAttestations is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody[_]) GetAttestations() []*AttestationData {
	panic("not implemented")
}

// SetAttestations is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody[_]) SetAttestations(_ []*AttestationData) {
	panic("not implemented")
}

// GetSlashingInfo is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody[_]) GetSlashingInfo() []*SlashingInfo {
	panic("not implemented")
}

// SetSlashingInfo is not implemented for BeaconBlockDeneb.
func (b *BeaconBlockBody[_]) SetSlashingInfo(_ []*SlashingInfo) {
	panic("not implemented")
}

// GetTopLevelRoots returns the top-level roots of the BeaconBlockBody.
func (b *BeaconBlockBody[LogT]) GetTopLevelRoots() []common.Root {
	return []common.Root{
		common.Root(b.GetRandaoReveal().HashTreeRoot()),
		b.Eth1Data.HashTreeRoot(),
		common.Root(b.GetGraffiti().HashTreeRoot()),
		Deposits[LogT](b.GetDeposits()).HashTreeRoot(),
		b.GetExecutionPayload().HashTreeRoot(),
		// I think this is a bug.
		common.Root{},
	}
}

// Length returns the number of fields in the BeaconBlockBody struct.
func (b *BeaconBlockBody[_]) Length() uint64 {
	return BodyLengthDeneb
}

// GetRandaoReveal returns the RandaoReveal of the Body.
func (b *BeaconBlockBody[_]) GetRandaoReveal() crypto.BLSSignature {
	return b.RandaoReveal
}

// SetRandaoReveal sets the RandaoReveal of the Body.
func (b *BeaconBlockBody[_]) SetRandaoReveal(reveal crypto.BLSSignature) {
	b.RandaoReveal = reveal
}

// GetEth1Data returns the Eth1Data of the Body.
func (b *BeaconBlockBody[_]) GetEth1Data() *Eth1Data {
	return b.Eth1Data
}

// GetGraffiti returns the Graffiti of the Body.
func (b *BeaconBlockBody[_]) GetGraffiti() common.Bytes32 {
	return b.Graffiti
}

// SetGraffiti sets the Graffiti of the Body.
func (b *BeaconBlockBody[_]) SetGraffiti(graffiti common.Bytes32) {
	b.Graffiti = graffiti
}

// GetDeposits returns the Deposits of the BeaconBlockBody.
func (b *BeaconBlockBody[LogT]) GetDeposits() []*Deposit[LogT] {
	return b.Deposits
}

// SetDeposits sets the Deposits of the BeaconBlockBody.
func (b *BeaconBlockBody[LogT]) SetDeposits(deposits []*Deposit[LogT]) {
	b.Deposits = deposits
}

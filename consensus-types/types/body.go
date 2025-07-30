// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"fmt"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/version"
	ssz "github.com/ferranbt/fastssz"
)

const (
	// BodyLengthDeneb is the number of fields in the BeaconBlockBody struct for Deneb.
	BodyLengthDeneb uint64 = 12

	// BodyLengthElectra is the number of fields in the BeaconBlockBody struct for Electra.
	BodyLengthElectra uint64 = 13

	// KZGPosition is the position of BlobKzgCommitments in the block body.
	KZGPosition uint64 = 11

	// KZGGeneralizedIndex is the index of the KZG commitment root's parent.
	//     (1 << log2ceil(KZGPosition)) | KZGPosition.
	KZGGeneralizedIndex = 27

	// KZGRootIndex is the merkle index of BlobKzgCommitments' root
	// in the merkle tree built from the block body.
	//     2 * KZGGeneralizedIndex.
	KZGRootIndex = KZGGeneralizedIndex * 2

	// KZGInclusionProofDepth is the
	//     Log2Floor(KZGGeneralizedIndex) +
	//     Log2Ceil(MaxBlobCommitmentsPerBlock) + 1
	KZGInclusionProofDepth = 17

	// KZGOffset is the offset of the KZG commitments in the serialized block body.
	KZGOffset = KZGRootIndex * constants.MaxBlobCommitmentsPerBlock
)

// Compile-time assertions to ensure BeaconBlockBody implements necessary interfaces.
var (
	_ constraints.SSZVersionedMarshallableRootable = (*BeaconBlockBody)(nil)
)

// BeaconBlockBody represents the body of a beacon block.
type BeaconBlockBody struct {
	// Must be available within the object to satisfy signature required for SizeSSZ.
	Versionable `json:"-"`

	// RandaoReveal is the reveal of the RANDAO.
	RandaoReveal crypto.BLSSignature
	// Eth1Data is the data from the Eth1 chain.
	Eth1Data *Eth1Data
	// Graffiti is for a fun message or meme.
	Graffiti [32]byte
	// proposerSlashings is unused but left for compatibility.
	proposerSlashings []*ProposerSlashing
	// attesterSlashings is unused but left for compatibility.
	attesterSlashings []*AttesterSlashing
	// attestations is unused but left for compatibility.
	attestations []*Attestation
	// Deposits is the list of deposits included in the body.
	Deposits []*Deposit
	// voluntaryExits is unused but left for compatibility.
	voluntaryExits []*VoluntaryExit
	// syncAggregate is unused but left for compatibility.
	syncAggregate *SyncAggregate
	// ExecutionPayload is the execution payload of the body.
	ExecutionPayload *ExecutionPayload
	// blsToExecutionChanges is unused but left for compatibility.
	blsToExecutionChanges []*BlsToExecutionChange
	// BlobKzgCommitments is the list of KZG commitments for the EIP-4844 blobs.
	BlobKzgCommitments []eip4844.KZGCommitment
	// executionRequests is introduced in Electra. We keep this private so that it must go through Getter/Setter
	// which does a forkVersion check.
	executionRequests *ExecutionRequests
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlockBody in SSZ.
func (b *BeaconBlockBody) SizeSSZ() int {
	var size = 96 + 72 + 32 + 4 + 4 + 4 + 4 + 4 + 160 + 4 + 4 + 4 // sync aggregate is 160 bytes
	includeExecRequest := version.EqualsOrIsAfter(b.GetForkVersion(), version.Electra())
	if includeExecRequest {
		// Add 4 for the offset of dynamic field ExecutionRequests
		size += int(constants.SSZOffsetSize)
	}

	// Dynamic fields
	size += len(b.proposerSlashings) * 16 // UnusedType
	size += len(b.attesterSlashings) * 16 // UnusedType
	size += len(b.attestations) * 16      // UnusedType
	size += len(b.Deposits) * 192         // Deposit
	size += len(b.voluntaryExits) * 16    // UnusedType
	size += b.ExecutionPayload.SizeSSZ()
	size += len(b.blsToExecutionChanges) * 16 // UnusedType
	size += len(b.BlobKzgCommitments) * 48    // KZGCommitment
	if includeExecRequest && b.executionRequests != nil {
		size += b.executionRequests.SizeSSZ()
	}
	return size
}

// MarshalSSZ serializes the BeaconBlockBody to SSZ-encoded bytes.
func (b *BeaconBlockBody) MarshalSSZ() ([]byte, error) {
	err := common.EnforceAllUnused(
		b.GetProposerSlashings(),
		b.GetAttesterSlashings(),
		b.GetAttestations(),
		b.GetVoluntaryExits(),
		b.GetSyncAggregate(),
		b.GetBlsToExecutionChanges(),
	)
	if err != nil {
		return []byte{}, err
	}
	return b.MarshalSSZTo(make([]byte, 0, b.SizeSSZ()))
}

func NewEmptyBeaconBlockBodyWithVersion(version common.Version) *BeaconBlockBody {
	return &BeaconBlockBody{
		Versionable:      NewVersionable(version),
		Eth1Data:         NewEmptyEth1Data(),
		ExecutionPayload: NewEmptyExecutionPayloadWithVersion(version),
		syncAggregate:    &SyncAggregate{},
	}
}

func (b *BeaconBlockBody) ValidateAfterDecodingSSZ() error {
	errUnused := common.EnforceAllUnused(
		b.GetProposerSlashings(),
		b.GetAttesterSlashings(),
		b.GetAttestations(),
		b.GetVoluntaryExits(),
		b.GetSyncAggregate(),
		b.GetBlsToExecutionChanges(),
	)
	return errors.Join(
		b.ExecutionPayload.ValidateAfterDecodingSSZ(),
		errUnused,
	)
}

// HashTreeRoot returns the SSZ hash tree root of the BeaconBlockBody.
func (b *BeaconBlockBody) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)
	if err := b.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */
// NOTE: The fastssz methods below are manually implemented to handle fork-specific logic.

// MarshalSSZTo ssz marshals the BeaconBlockBody object to a target array.
func (b *BeaconBlockBody) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Calculate offsets
	offset := 96 + 72 + 32 + 4 + 4 + 4 + 4 + 4 + 160 + 4 + 4 + 4 // Fixed part size
	includeExecRequest := version.EqualsOrIsAfter(b.GetForkVersion(), version.Electra())
	if includeExecRequest {
		offset += 4 // ExecutionRequests offset
	}

	// Static fields
	dst = append(dst, b.RandaoReveal[:]...)

	// Eth1Data
	eth1Bytes, err := b.Eth1Data.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, eth1Bytes...)

	dst = append(dst, b.Graffiti[:]...)

	// Offsets for dynamic fields
	// ProposerSlashings offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += len(b.proposerSlashings) * 16

	// AttesterSlashings offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += len(b.attesterSlashings) * 16

	// Attestations offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += len(b.attestations) * 16

	// Deposits offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += len(b.Deposits) * 192

	// VoluntaryExits offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += len(b.voluntaryExits) * 16

	// SyncAggregate
	var syncBytes []byte
	if b.syncAggregate == nil {
		// For nil SyncAggregate, use empty SyncAggregate
		emptySyncAgg := &SyncAggregate{}
		syncBytes, err = emptySyncAgg.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	} else {
		syncBytes, err = b.syncAggregate.MarshalSSZ()
		if err != nil {
			return nil, err
		}
	}
	dst = append(dst, syncBytes...)

	// ExecutionPayload offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += b.ExecutionPayload.SizeSSZ()

	// BlsToExecutionChanges offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += len(b.blsToExecutionChanges) * 16

	// BlobKzgCommitments offset
	dst = ssz.MarshalUint32(dst, uint32(offset))
	offset += len(b.BlobKzgCommitments) * 48

	// ExecutionRequests offset (Electra+)
	if includeExecRequest {
		dst = ssz.MarshalUint32(dst, uint32(offset))
	}

	// Dynamic fields
	// ProposerSlashings
	for _, ps := range b.proposerSlashings {
		psBytes, err := ps.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		dst = append(dst, psBytes...)
	}

	// AttesterSlashings
	for _, as := range b.attesterSlashings {
		asBytes, err := as.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		dst = append(dst, asBytes...)
	}

	// Attestations
	for _, att := range b.attestations {
		attBytes, err := att.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		dst = append(dst, attBytes...)
	}

	// Deposits
	for _, dep := range b.Deposits {
		depBytes, err := dep.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		dst = append(dst, depBytes...)
	}

	// VoluntaryExits
	for _, ve := range b.voluntaryExits {
		veBytes, err := ve.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		dst = append(dst, veBytes...)
	}

	// ExecutionPayload
	dst, err = b.ExecutionPayload.MarshalSSZTo(dst)
	if err != nil {
		return nil, err
	}

	// BlsToExecutionChanges
	for _, bec := range b.blsToExecutionChanges {
		becBytes, err := bec.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		dst = append(dst, becBytes...)
	}

	// BlobKzgCommitments
	for _, comm := range b.BlobKzgCommitments {
		dst = append(dst, comm[:]...)
	}

	// ExecutionRequests (Electra+)
	if includeExecRequest && b.executionRequests != nil {
		erBytes, err := b.executionRequests.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		dst = append(dst, erBytes...)
	}

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the BeaconBlockBody object.
func (b *BeaconBlockBody) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 388 { // Minimum size without Electra
		return ssz.ErrSize
	}

	includeExecRequest := version.EqualsOrIsAfter(b.GetForkVersion(), version.Electra())
	if includeExecRequest && len(buf) < 392 { // Minimum size with Electra
		return ssz.ErrSize
	}

	var err error
	size := uint32(len(buf))
	offset := 0

	// Field (0) 'RandaoReveal'
	copy(b.RandaoReveal[:], buf[0:96])
	offset += 96

	// Field (1) 'Eth1Data'
	if b.Eth1Data == nil {
		b.Eth1Data = &Eth1Data{}
	}
	if err = b.Eth1Data.UnmarshalSSZ(buf[offset : offset+72]); err != nil {
		return err
	}
	offset += 72

	// Field (2) 'Graffiti'
	copy(b.Graffiti[:], buf[offset:offset+32])
	offset += 32

	// Read offsets
	proposerSlashingsOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4
	attesterSlashingsOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4
	attestationsOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4
	depositsOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4
	voluntaryExitsOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4

	// Field (8) 'SyncAggregate'
	if b.syncAggregate == nil {
		b.syncAggregate = &SyncAggregate{}
	}
	if err = b.syncAggregate.UnmarshalSSZ(buf[offset : offset+160]); err != nil {
		return err
	}
	offset += 160

	executionPayloadOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4
	blsToExecutionChangesOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4
	blobKzgCommitmentsOffset := ssz.UnmarshallUint32(buf[offset : offset+4])
	offset += 4

	var executionRequestsOffset uint32
	if includeExecRequest {
		executionRequestsOffset = ssz.UnmarshallUint32(buf[offset : offset+4])
		offset += 4
	}

	// Validate offsets
	if proposerSlashingsOffset > size || attesterSlashingsOffset > size ||
		attestationsOffset > size || depositsOffset > size ||
		voluntaryExitsOffset > size || executionPayloadOffset > size ||
		blsToExecutionChangesOffset > size || blobKzgCommitmentsOffset > size ||
		(includeExecRequest && executionRequestsOffset > size) {
		return ssz.ErrInvalidVariableOffset
	}

	// Unmarshal dynamic fields
	// ProposerSlashings
	if proposerSlashingsOffset < attesterSlashingsOffset {
		count := (attesterSlashingsOffset - proposerSlashingsOffset) / 16
		b.proposerSlashings = make([]*ProposerSlashing, count)
		for i := uint32(0); i < count; i++ {
			ps := ProposerSlashing(0)
			b.proposerSlashings[i] = &ps
			if err = b.proposerSlashings[i].UnmarshalSSZ(buf[proposerSlashingsOffset+i*16 : proposerSlashingsOffset+(i+1)*16]); err != nil {
				return err
			}
		}
	}

	// AttesterSlashings
	if attesterSlashingsOffset < attestationsOffset {
		count := (attestationsOffset - attesterSlashingsOffset) / 16
		b.attesterSlashings = make([]*AttesterSlashing, count)
		for i := uint32(0); i < count; i++ {
			as := AttesterSlashing(0)
			b.attesterSlashings[i] = &as
			if err = b.attesterSlashings[i].UnmarshalSSZ(buf[attesterSlashingsOffset+i*16 : attesterSlashingsOffset+(i+1)*16]); err != nil {
				return err
			}
		}
	}

	// Attestations
	if attestationsOffset < depositsOffset {
		count := (depositsOffset - attestationsOffset) / 16
		b.attestations = make([]*Attestation, count)
		for i := uint32(0); i < count; i++ {
			att := Attestation(0)
			b.attestations[i] = &att
			if err = b.attestations[i].UnmarshalSSZ(buf[attestationsOffset+i*16 : attestationsOffset+(i+1)*16]); err != nil {
				return err
			}
		}
	}

	// Deposits
	if depositsOffset < voluntaryExitsOffset {
		count := (voluntaryExitsOffset - depositsOffset) / 192
		b.Deposits = make([]*Deposit, count)
		for i := uint32(0); i < count; i++ {
			b.Deposits[i] = &Deposit{}
			if err = b.Deposits[i].UnmarshalSSZ(buf[depositsOffset+i*192 : depositsOffset+(i+1)*192]); err != nil {
				return err
			}
		}
	}

	// VoluntaryExits
	if voluntaryExitsOffset < executionPayloadOffset {
		count := (executionPayloadOffset - voluntaryExitsOffset) / 16
		b.voluntaryExits = make([]*VoluntaryExit, count)
		for i := uint32(0); i < count; i++ {
			ve := VoluntaryExit(0)
			b.voluntaryExits[i] = &ve
			if err = b.voluntaryExits[i].UnmarshalSSZ(buf[voluntaryExitsOffset+i*16 : voluntaryExitsOffset+(i+1)*16]); err != nil {
				return err
			}
		}
	}

	// ExecutionPayload
	if b.ExecutionPayload == nil {
		b.ExecutionPayload = NewEmptyExecutionPayloadWithVersion(b.GetForkVersion())
	}
	if executionPayloadOffset < blsToExecutionChangesOffset {
		if err = b.ExecutionPayload.UnmarshalSSZ(buf[executionPayloadOffset:blsToExecutionChangesOffset]); err != nil {
			return err
		}
	}

	// BlsToExecutionChanges
	if blsToExecutionChangesOffset < blobKzgCommitmentsOffset {
		count := (blobKzgCommitmentsOffset - blsToExecutionChangesOffset) / 16
		b.blsToExecutionChanges = make([]*BlsToExecutionChange, count)
		for i := uint32(0); i < count; i++ {
			bec := BlsToExecutionChange(0)
			b.blsToExecutionChanges[i] = &bec
			if err = b.blsToExecutionChanges[i].UnmarshalSSZ(buf[blsToExecutionChangesOffset+i*16 : blsToExecutionChangesOffset+(i+1)*16]); err != nil {
				return err
			}
		}
	}

	// BlobKzgCommitments
	var endOffset uint32
	if includeExecRequest {
		endOffset = executionRequestsOffset
	} else {
		endOffset = size
	}

	if blobKzgCommitmentsOffset < endOffset {
		count := (endOffset - blobKzgCommitmentsOffset) / 48
		b.BlobKzgCommitments = make([]eip4844.KZGCommitment, count)
		for i := uint32(0); i < count; i++ {
			copy(b.BlobKzgCommitments[i][:], buf[blobKzgCommitmentsOffset+i*48:blobKzgCommitmentsOffset+(i+1)*48])
		}
	}

	// ExecutionRequests (Electra+)
	if includeExecRequest && executionRequestsOffset < size {
		if b.executionRequests == nil {
			b.executionRequests = &ExecutionRequests{}
		}
		if err = b.executionRequests.UnmarshalSSZ(buf[executionRequestsOffset:]); err != nil {
			return err
		}
	}

	// Initialize nil slices to empty slices for consistency
	if b.proposerSlashings == nil {
		b.proposerSlashings = make([]*ProposerSlashing, 0)
	}
	if b.attesterSlashings == nil {
		b.attesterSlashings = make([]*AttesterSlashing, 0)
	}
	if b.attestations == nil {
		b.attestations = make([]*Attestation, 0)
	}
	if b.voluntaryExits == nil {
		b.voluntaryExits = make([]*VoluntaryExit, 0)
	}
	if b.blsToExecutionChanges == nil {
		b.blsToExecutionChanges = make([]*BlsToExecutionChange, 0)
	}

	return b.ValidateAfterDecodingSSZ()
}

// HashTreeRootWith ssz hashes the BeaconBlockBody object with a hasher.
func (b *BeaconBlockBody) HashTreeRootWith(hh ssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'RandaoReveal'
	hh.PutBytes(b.RandaoReveal[:])

	// Field (1) 'Eth1Data'
	if err := b.Eth1Data.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (2) 'Graffiti'
	hh.PutBytes(b.Graffiti[:])

	// Field (3) 'ProposerSlashings'
	{
		subIndx := hh.Index()
		num := uint64(len(b.proposerSlashings))
		if num > constants.MaxProposerSlashings {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range b.proposerSlashings {
			// ProposerSlashing uses UnusedType which inherits HashTreeRoot from common.UnusedType
			root, err := elem.HashTreeRoot()
			if err != nil {
				return err
			}
			hh.PutBytes(root[:])
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxProposerSlashings)
	}

	// Field (4) 'AttesterSlashings'
	{
		subIndx := hh.Index()
		num := uint64(len(b.attesterSlashings))
		if num > constants.MaxAttesterSlashings {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range b.attesterSlashings {
			// AttesterSlashing uses UnusedType which inherits HashTreeRoot from common.UnusedType
			root, err := elem.HashTreeRoot()
			if err != nil {
				return err
			}
			hh.PutBytes(root[:])
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxAttesterSlashings)
	}

	// Field (5) 'Attestations'
	{
		subIndx := hh.Index()
		num := uint64(len(b.attestations))
		if num > constants.MaxAttestations {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range b.attestations {
			// Attestation uses UnusedType which inherits HashTreeRoot from common.UnusedType
			root, err := elem.HashTreeRoot()
			if err != nil {
				return err
			}
			hh.PutBytes(root[:])
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxAttestations)
	}

	// Field (6) 'Deposits'
	{
		subIndx := hh.Index()
		num := uint64(len(b.Deposits))
		if num > constants.MaxDeposits {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range b.Deposits {
			if err := elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxDeposits)
	}

	// Field (7) 'VoluntaryExits'
	{
		subIndx := hh.Index()
		num := uint64(len(b.voluntaryExits))
		if num > constants.MaxVoluntaryExits {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range b.voluntaryExits {
			// VoluntaryExit uses UnusedType which inherits HashTreeRoot from common.UnusedType
			root, err := elem.HashTreeRoot()
			if err != nil {
				return err
			}
			hh.PutBytes(root[:])
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxVoluntaryExits)
	}

	// Field (8) 'SyncAggregate'
	if b.syncAggregate == nil {
		// For nil SyncAggregate, use empty SyncAggregate
		emptySyncAgg := &SyncAggregate{}
		if err := emptySyncAgg.HashTreeRootWith(hh); err != nil {
			return err
		}
	} else {
		if err := b.syncAggregate.HashTreeRootWith(hh); err != nil {
			return err
		}
	}

	// Field (9) 'ExecutionPayload'
	if err := b.ExecutionPayload.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (10) 'BlsToExecutionChanges'
	{
		subIndx := hh.Index()
		num := uint64(len(b.blsToExecutionChanges))
		if num > constants.MaxBlsToExecutionChanges {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range b.blsToExecutionChanges {
			// BlsToExecutionChange uses UnusedType which inherits HashTreeRoot from common.UnusedType
			root, err := elem.HashTreeRoot()
			if err != nil {
				return err
			}
			hh.PutBytes(root[:])
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxBlsToExecutionChanges)
	}

	// Field (11) 'BlobKzgCommitments'
	{
		subIndx := hh.Index()
		num := uint64(len(b.BlobKzgCommitments))
		if num > 4096 {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range b.BlobKzgCommitments {
			hh.PutBytes(elem[:])
		}
		hh.MerkleizeWithMixin(subIndx, num, 4096)
	}

	// Field (12) 'ExecutionRequests' (Electra+ only)
	if version.EqualsOrIsAfter(b.GetForkVersion(), version.Electra()) {
		if b.executionRequests != nil {
			// ExecutionRequests doesn't have HashTreeRootWith yet
			// Use the HashTreeRoot method from ExecutionRequests
			root, err := b.executionRequests.HashTreeRoot()
			if err != nil {
				return err
			}
			hh.PutBytes(root[:])
		} else {
			// If executionRequests is nil but we're in Electra+, we need to handle this
			// This should not happen in valid blocks, but we handle it gracefully
			hh.PutBytes(make([]byte, 32))
		}
	}

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BeaconBlockBody object.
func (b *BeaconBlockBody) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(b)
}

/* -------------------------------------------------------------------------- */
/*                              Getters/Setters                               */
/* -------------------------------------------------------------------------- */

// GetTopLevelRoots returns the top-level roots of the BeaconBlockBody.
func (b *BeaconBlockBody) GetTopLevelRoots() ([]common.Root, error) {
	var tlrs []common.Root
	var root [32]byte
	var err error

	// RandaoReveal
	tlrs = append(tlrs, common.Root(b.GetRandaoReveal().HashTreeRoot()))

	// Eth1Data
	root, err = b.Eth1Data.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	tlrs = append(tlrs, common.Root(root))

	// Graffiti
	tlrs = append(tlrs, common.Root(b.GetGraffiti().HashTreeRoot()))

	// Collections
	collections := []func() ([32]byte, error){
		b.GetProposerSlashings().HashTreeRoot,
		b.GetAttesterSlashings().HashTreeRoot,
		b.GetAttestations().HashTreeRoot,
		b.GetDeposits().HashTreeRoot,
		b.GetVoluntaryExits().HashTreeRoot,
		b.syncAggregate.HashTreeRoot,
		b.GetExecutionPayload().HashTreeRoot,
		b.GetBlsToExecutionChanges().HashTreeRoot,
	}

	for _, fn := range collections {
		root, err = fn()
		if err != nil {
			return nil, err
		}
		tlrs = append(tlrs, common.Root(root))
	}

	// KzgCommitments intentionally left blank
	tlrs = append(tlrs, common.Root{})

	// ExecutionRequests for Electra+
	if version.EqualsOrIsAfter(b.GetForkVersion(), version.Electra()) {
		er, err := b.GetExecutionRequests()
		if err != nil {
			return nil, err
		}
		root, err = er.HashTreeRoot()
		if err != nil {
			return nil, err
		}
		tlrs = append(tlrs, common.Root(root))
	}

	// Verify length
	if uint64(len(tlrs)) != b.Length() {
		return nil, fmt.Errorf(
			"top-level roots length (%d) does not match expected body length (%d)",
			len(tlrs), b.Length(),
		)
	}

	return tlrs, nil
}

// Length returns the number of fields in the BeaconBlockBody struct
// according to the fork version.
func (b *BeaconBlockBody) Length() uint64 {
	if version.IsBefore(b.GetForkVersion(), version.Electra()) {
		return BodyLengthDeneb
	}
	return BodyLengthElectra
}

func (b *BeaconBlockBody) GetRandaoReveal() crypto.BLSSignature {
	return b.RandaoReveal
}

func (b *BeaconBlockBody) SetRandaoReveal(reveal crypto.BLSSignature) {
	b.RandaoReveal = reveal
}

func (b *BeaconBlockBody) GetEth1Data() *Eth1Data {
	return b.Eth1Data
}

func (b *BeaconBlockBody) SetEth1Data(eth1Data *Eth1Data) {
	b.Eth1Data = eth1Data
}

func (b *BeaconBlockBody) GetGraffiti() common.Bytes32 {
	return b.Graffiti
}

func (b *BeaconBlockBody) SetGraffiti(graffiti common.Bytes32) {
	b.Graffiti = graffiti
}

func (b *BeaconBlockBody) GetProposerSlashings() ProposerSlashings {
	return b.proposerSlashings
}

func (b *BeaconBlockBody) SetProposerSlashings(ps ProposerSlashings) {
	b.proposerSlashings = ps
}

func (b *BeaconBlockBody) GetAttesterSlashings() AttesterSlashings {
	return b.attesterSlashings
}

func (b *BeaconBlockBody) SetAttesterSlashings(ps AttesterSlashings) {
	b.attesterSlashings = ps
}

func (b *BeaconBlockBody) GetVoluntaryExits() VoluntaryExits {
	return b.voluntaryExits
}

func (b *BeaconBlockBody) SetVoluntaryExits(exits VoluntaryExits) {
	b.voluntaryExits = exits
}

func (b *BeaconBlockBody) GetDeposits() Deposits {
	return b.Deposits
}

func (b *BeaconBlockBody) SetDeposits(deposits Deposits) {
	b.Deposits = deposits
}

func (b *BeaconBlockBody) GetAttestations() Attestations {
	return b.attestations
}

func (b *BeaconBlockBody) SetAttestations(attestations Attestations) {
	b.attestations = attestations
}

func (b *BeaconBlockBody) GetSyncAggregate() *SyncAggregate {
	return b.syncAggregate
}

func (b *BeaconBlockBody) SetSyncAggregate(syncAggregate *SyncAggregate) {
	b.syncAggregate = syncAggregate
}

func (b *BeaconBlockBody) GetExecutionPayload() *ExecutionPayload {
	return b.ExecutionPayload
}

func (b *BeaconBlockBody) SetExecutionPayload(executionData *ExecutionPayload) {
	b.ExecutionPayload = executionData
}

func (b *BeaconBlockBody) GetBlsToExecutionChanges() BlsToExecutionChanges {
	return b.blsToExecutionChanges
}

func (b *BeaconBlockBody) SetBlsToExecutionChanges(blsChanges BlsToExecutionChanges) {
	b.blsToExecutionChanges = blsChanges
}

func (b *BeaconBlockBody) GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash] {
	return b.BlobKzgCommitments
}

func (b *BeaconBlockBody) SetBlobKzgCommitments(commitments eip4844.KZGCommitments[common.ExecutionHash]) {
	b.BlobKzgCommitments = commitments
}

func (b *BeaconBlockBody) GetExecutionRequests() (*ExecutionRequests, error) {
	if version.IsBefore(b.GetForkVersion(), version.Electra()) {
		return nil, errors.Wrapf(ErrFieldNotSupportedOnFork, "block version %d", b.GetForkVersion())
	}
	if b.executionRequests == nil {
		return nil, errors.New("retrieved execution requests is nil")
	}
	return b.executionRequests, nil
}

func (b *BeaconBlockBody) SetExecutionRequests(executionRequest *ExecutionRequests) error {
	if executionRequest == nil {
		return errors.New("cannot set execution requests to nil")
	}
	if version.IsBefore(b.GetForkVersion(), version.Electra()) {
		return errors.Wrapf(ErrFieldNotSupportedOnFork, "block version %d", b.GetForkVersion())
	}
	b.executionRequests = executionRequest
	return nil
}

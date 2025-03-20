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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/karalabe/ssz"
)

const (
	// BodyLengthDeneb is the number of fields in the BeaconBlockBodyDeneb
	// struct.
	BodyLengthDeneb uint64 = 12

	// BodyLengthElectra is the number of fields in the BeaconBlockBodyElectra struct.
	// TODO(pectra): Ensure this is propagated where necessary
	BodyLengthElectra uint64 = 13

	// KZGPositionDeneb is the position of BlobKzgCommitments in the block body.
	KZGPositionDeneb = BodyLengthDeneb - 1

	// KZGGeneralizedIndex is the index of the KZG commitment root's parent.
	//     (1 << log2ceil(KZGPositionDeneb)) | KZGPositionDeneb.
	KZGGeneralizedIndex = 27

	// KZGRootIndexDeneb is the merkle index of BlobKzgCommitments' root
	// in the merkle tree built from the block body.
	//     2 * KZGGeneralizedIndex.
	KZGRootIndexDeneb = KZGGeneralizedIndex * 2

	// KZGInclusionProofDepth is the
	//     Log2Floor(KZGGeneralizedIndex) +
	//     Log2Ceil(MaxBlobCommitmentsPerBlock) + 1
	KZGInclusionProofDepth = 17

	// KZGOffsetDeneb is the offset of the KZG commitments in the serialized block body.
	KZGOffsetDeneb = KZGRootIndexDeneb * constants.MaxBlobCommitmentsPerBlock
)

// Compile-time assertions to ensure BeaconBlockBody implements necessary interfaces.
var (
	_ ssz.DynamicObject                            = (*BeaconBlockBody)(nil)
	_ constraints.SSZVersionedMarshallableRootable = (*BeaconBlockBody)(nil)
)

// BeaconBlockBody represents the body of a beacon block.
type BeaconBlockBody struct {
	// Must be available within the object to satisfy signature required for SizeSSZ and DefineSSZ.
	constraints.Versionable `json:"-"`

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
	// executionRequests is introduced in electra. We keep this private so that it must go through Getter/Setter
	// which does a forkVersion check.
	executionRequests *ExecutionRequests
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlockBody in SSZ.
func (b *BeaconBlockBody) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	syncSize := b.syncAggregate.SizeSSZ(siz)
	var size = 96 + 72 + 32 + 4 + 4 + 4 + 4 + 4 + syncSize + 4 + 4 + 4
	if !version.IsBefore(b.GetForkVersion(), version.Electra()) {
		// Add 4 for the offset of dynamic field ExecutionRequests
		size += sszDynamicObjectOffset
	}

	if fixed {
		return size
	}

	size += ssz.SizeSliceOfStaticObjects(siz, b.proposerSlashings)
	size += ssz.SizeSliceOfStaticObjects(siz, b.attesterSlashings)
	size += ssz.SizeSliceOfStaticObjects(siz, b.attestations)
	size += ssz.SizeSliceOfStaticObjects(siz, b.Deposits)
	size += ssz.SizeSliceOfStaticObjects(siz, b.voluntaryExits)
	size += ssz.SizeDynamicObject(siz, b.ExecutionPayload)
	size += ssz.SizeSliceOfStaticObjects(siz, b.blsToExecutionChanges)
	size += ssz.SizeSliceOfStaticBytes(siz, b.BlobKzgCommitments)
	if !version.IsBefore(b.GetForkVersion(), version.Electra()) {
		size += ssz.SizeDynamicObject(siz, b.executionRequests)
	}
	return size
}

// DefineSSZ defines the SSZ serialization of the BeaconBlockBody.
//
//nolint:mnd // TODO: get from accessible chainspec field params
func (b *BeaconBlockBody) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineStaticBytes(codec, &b.RandaoReveal)
	ssz.DefineStaticObject(codec, &b.Eth1Data)
	ssz.DefineStaticBytes(codec, &b.Graffiti)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.proposerSlashings, constants.MaxProposerSlashings)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.attesterSlashings, constants.MaxAttesterSlashings)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.attestations, constants.MaxAttestations)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.Deposits, constants.MaxDeposits)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.voluntaryExits, constants.MaxVoluntaryExits)
	ssz.DefineStaticObject(codec, &b.syncAggregate)
	ssz.DefineDynamicObjectOffset(codec, &b.ExecutionPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.blsToExecutionChanges, constants.MaxBlsToExecutionChanges)
	ssz.DefineSliceOfStaticBytesOffset(codec, &b.BlobKzgCommitments, 4096)
	if !version.IsBefore(b.GetForkVersion(), version.Electra()) {
		ssz.DefineDynamicObjectOffset(codec, &b.executionRequests)
	}

	// Define the dynamic data (fields)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.proposerSlashings, constants.MaxProposerSlashings)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.attesterSlashings, constants.MaxAttesterSlashings)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.attestations, constants.MaxAttestations)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.Deposits, constants.MaxDeposits)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.voluntaryExits, constants.MaxVoluntaryExits)
	ssz.DefineDynamicObjectContent(codec, &b.ExecutionPayload)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.blsToExecutionChanges, constants.MaxBlsToExecutionChanges)
	ssz.DefineSliceOfStaticBytesContent(codec, &b.BlobKzgCommitments, 4096)
	if !version.IsBefore(b.GetForkVersion(), version.Electra()) {
		ssz.DefineDynamicObjectContent(codec, &b.executionRequests)
	}
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
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

func NewEmptyBeaconBlockBodyWithVersion(version common.Version) *BeaconBlockBody {
	return &BeaconBlockBody{
		Versionable:      NewVersionable(version),
		Eth1Data:         NewEmptyEthi1Data(),
		ExecutionPayload: NewEmptyExecutionPayloadWithVersion(version),
		syncAggregate:    &SyncAggregate{},
	}
}

func (b *BeaconBlockBody) EnsureSyntaxFromSSZ() error {
	errUnused := common.EnforceAllUnused(
		b.GetProposerSlashings(),
		b.GetAttesterSlashings(),
		b.GetAttestations(),
		b.GetVoluntaryExits(),
		b.GetSyncAggregate(),
		b.GetBlsToExecutionChanges(),
	)
	return errors.Join(
		b.ExecutionPayload.EnsureSyntaxFromSSZ(),
		errUnused,
	)
}

// HashTreeRoot returns the SSZ hash tree root of the BeaconBlockBody.
func (b *BeaconBlockBody) HashTreeRoot() common.Root {
	return ssz.HashConcurrent(b)
}

// GetTopLevelRoots returns the top-level roots of the BeaconBlockBody.
func (b *BeaconBlockBody) GetTopLevelRoots() []common.Root {
	return []common.Root{
		common.Root(b.GetRandaoReveal().HashTreeRoot()),
		b.Eth1Data.HashTreeRoot(),
		common.Root(b.GetGraffiti().HashTreeRoot()),
		b.GetProposerSlashings().HashTreeRoot(),
		b.GetAttesterSlashings().HashTreeRoot(),
		b.GetAttestations().HashTreeRoot(),
		b.GetDeposits().HashTreeRoot(),
		b.GetVoluntaryExits().HashTreeRoot(),
		b.syncAggregate.HashTreeRoot(),
		b.GetExecutionPayload().HashTreeRoot(),
		b.GetBlsToExecutionChanges().HashTreeRoot(),
		// KzgCommitments intentionally left blank - included separately for inclusion proof
		{},
	}
}

// Length returns the number of fields in the BeaconBlockBody struct.
func (b *BeaconBlockBody) Length() uint64 {
	return BodyLengthDeneb
}

/* -------------------------------------------------------------------------- */
/*                              Getters/Setters                               */
/* -------------------------------------------------------------------------- */

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

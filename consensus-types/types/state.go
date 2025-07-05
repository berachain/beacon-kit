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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	fastssz "github.com/ferranbt/fastssz"
)

// BeaconState represents the entire state of the beacon chain.
type BeaconState struct {
	Versionable `json:"-"`

	// Versioning
	GenesisValidatorsRoot common.Root `json:"genesis_validators_root,omitempty"`
	Slot                  math.Slot   `json:"slot,omitempty"`
	Fork                  *Fork       `json:"fork,omitempty"`

	// History
	LatestBlockHeader *BeaconBlockHeader `json:"latest_block_header,omitempty"`
	BlockRoots        []common.Root      `json:"block_roots,omitempty"`
	StateRoots        []common.Root      `json:"state_roots,omitempty"`

	// Eth1
	Eth1Data                     *Eth1Data               `json:"eth1_data,omitempty"`
	Eth1DepositIndex             uint64                  `json:"eth1_deposit_index,omitempty"`
	LatestExecutionPayloadHeader *ExecutionPayloadHeader `json:"latest_execution_payload_header,omitempty"`

	// Registry
	Validators []*Validator `json:"validators,omitempty"`
	Balances   []uint64     `json:"balances,omitempty"`

	// Randomness
	RandaoMixes []common.Bytes32 `json:"randao_mixes,omitempty"`

	// Withdrawals
	NextWithdrawalIndex          uint64              `json:"next_withdrawal_index,omitempty"`
	NextWithdrawalValidatorIndex math.ValidatorIndex `json:"next_withdrawal_validator_index,omitempty"`

	// Slashing
	Slashings     []math.Gwei `json:"slashings,omitempty"`
	TotalSlashing math.Gwei   `json:"total_slashing,omitempty"`

	// PendingPartialWithdrawals is introduced in electra
	PendingPartialWithdrawals []*PendingPartialWithdrawal `json:"pending_partial_withdrawals,omitempty"`
}

// NewEmptyBeaconStateWithVersion returns a new empty BeaconState with the given fork version.
func NewEmptyBeaconStateWithVersion(version common.Version) *BeaconState {
	return &BeaconState{
		Versionable: NewVersionable(version),
	}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the ssz encoded size in bytes for the BeaconState object.
func (st *BeaconState) SizeSSZ() int {
	/*
		GenesisValidatorsRoot = 32
		Slot = 8
		Fork = 4 + 4 + 8 = 16
		LatestBlockHeader = 8 + 8 + 32 + 32 + 32 = 112
		BlockRoots = 4 (Dynamic field)
		StateRoots = 4 (Dynamic field)
		Eth1Data = 32 + 8 + 32 = 72
		Eth1DepositIndex = 8
		LatestExecutionPayloadHeader = 4 (Dynamic field)
		Validators = 4 (Dynamic field)
		Balances = 4 (Dynamic field)
		RandaoMixes = 4 (Dynamic field)
		NextWithdrawalIndex = 8
		NextWithdrawalValidatorIndex = 8
		Slashings = 4 (Dynamic field)
		TotalSlashing = 8

		// Electra Fork
		PendingPartialWithdrawals = 4 (Dynamic field)
	*/
	var size = 300

	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		// Add 4 for PendingPartialWithdrawals after Electra
		size += 4
	}

	// Dynamic size fields
	size += len(st.BlockRoots) * 32 // Each root is 32 bytes
	size += len(st.StateRoots) * 32 // Each root is 32 bytes
	size += st.LatestExecutionPayloadHeader.SizeSSZ()
	size += len(st.Validators) * 121 // Each validator is 121 bytes
	size += len(st.Balances) * 8     // Each balance is 8 bytes
	size += len(st.RandaoMixes) * 32 // Each mix is 32 bytes
	size += len(st.Slashings) * 8    // Each slashing is 8 bytes
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		size += len(st.PendingPartialWithdrawals) * 24 // Each pending withdrawal is 24 bytes
	}

	return size
}

// MarshalSSZ marshals the BeaconState into SSZ format.
func (st *BeaconState) MarshalSSZ() ([]byte, error) {
	return st.MarshalSSZTo(make([]byte, 0, st.SizeSSZ()))
}

// HashTreeRoot computes the Merkleization of the BeaconState.
func (st *BeaconState) HashTreeRoot() ([32]byte, error) {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	if err := st.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()

}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the BeaconState object with a hasher.
//
//nolint:mnd,funlen,gocognit // todo fix.
func (st *BeaconState) HashTreeRootWith(
	hh fastssz.HashWalker,
) error {
	indx := hh.Index()

	// Field (0) 'GenesisValidatorsRoot'
	hh.PutBytes(st.GenesisValidatorsRoot[:])

	// Field (1) 'Slot'
	hh.PutUint64(uint64(st.Slot))

	// Field (2) 'Fork'
	if st.Fork == nil {
		st.Fork = &Fork{}
	}
	if err := st.Fork.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (3) 'LatestBlockHeader'
	if st.LatestBlockHeader == nil {
		st.LatestBlockHeader = &BeaconBlockHeader{}
	}
	if err := st.LatestBlockHeader.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (4) 'BlockRoots'
	if size := len(st.BlockRoots); size > 8192 {
		return fastssz.ErrListTooBigFn("BeaconState.BlockRoots", size, 8192)
	}
	subIndx := hh.Index()
	for _, i := range st.BlockRoots {
		hh.Append(i[:])
	}
	numItems := uint64(len(st.BlockRoots))
	hh.MerkleizeWithMixin(subIndx, numItems, 8192)

	// Field (5) 'StateRoots'
	if size := len(st.StateRoots); size > 8192 {
		return fastssz.ErrListTooBigFn("BeaconState.StateRoots", size, 8192)
	}
	subIndx = hh.Index()
	for _, i := range st.StateRoots {
		hh.Append(i[:])
	}
	numItems = uint64(len(st.StateRoots))
	hh.MerkleizeWithMixin(subIndx, numItems, 8192)

	// Field (6) 'Eth1Data'
	if st.Eth1Data == nil {
		st.Eth1Data = &Eth1Data{}
	}
	if err := st.Eth1Data.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (7) 'Eth1DepositIndex'
	hh.PutUint64(st.Eth1DepositIndex)

	// Field (8) 'LatestExecutionPayloadHeader'
	if st.LatestExecutionPayloadHeader == nil {
		st.LatestExecutionPayloadHeader = NewEmptyExecutionPayloadHeaderWithVersion(st.GetForkVersion())
	}
	if err := st.LatestExecutionPayloadHeader.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (9) 'Validators'
	subIndx = hh.Index()
	num := uint64(len(st.Validators))
	if num > 1099511627776 {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range st.Validators {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(subIndx, num, 1099511627776)

	// Field (10) 'Balances'
	if size := len(st.Balances); size > 1099511627776 {
		return fastssz.ErrListTooBigFn(
			"BeaconState.Balances",
			size,
			1099511627776,
		)
	}
	subIndx = hh.Index()
	for _, i := range st.Balances {
		hh.AppendUint64(i)
	}
	hh.FillUpTo32()
	numItems = uint64(len(st.Balances))
	hh.MerkleizeWithMixin(
		subIndx,
		numItems,
		fastssz.CalculateLimit(1099511627776, numItems, 8),
	)

	// Field (11) 'RandaoMixes'
	if size := len(st.RandaoMixes); size > 65536 {
		return fastssz.ErrListTooBigFn("BeaconState.RandaoMixes", size, 65536)
	}
	subIndx = hh.Index()
	for _, i := range st.RandaoMixes {
		hh.Append(i[:])
	}
	numItems = uint64(len(st.RandaoMixes))
	hh.MerkleizeWithMixin(subIndx, numItems, 65536)

	// Field (12) 'NextWithdrawalIndex'
	hh.PutUint64(st.NextWithdrawalIndex)

	// Field (13) 'NextWithdrawalValidatorIndex'
	hh.PutUint64(uint64(st.NextWithdrawalValidatorIndex))

	// Field (14) 'Slashings'
	if size := len(st.Slashings); size > 1099511627776 {
		return fastssz.ErrListTooBigFn(
			"BeaconState.Slashings",
			size,
			1099511627776,
		)
	}
	subIndx = hh.Index()
	for _, i := range st.Slashings {
		hh.AppendUint64(uint64(i))
	}
	hh.FillUpTo32()
	numItems = uint64(len(st.Slashings))
	hh.MerkleizeWithMixin(
		subIndx,
		numItems,
		fastssz.CalculateLimit(1099511627776, numItems, 8),
	)

	// Field (15) 'TotalSlashing'
	hh.PutUint64(uint64(st.TotalSlashing))

	// Field (16) 'PendingPartialWithdrawals' post-electra
	if version.EqualsOrIsAfter(st.GetForkVersion(), version.Electra()) {
		subIndx = hh.Index()
		numPPW := uint64(len(st.PendingPartialWithdrawals))
		if numPPW > constants.PendingPartialWithdrawalsLimit {
			return fastssz.ErrIncorrectListSize
		}
		for _, elem := range st.PendingPartialWithdrawals {
			if err := elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, numPPW, constants.PendingPartialWithdrawalsLimit)
	}
	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BeaconState object.
func (st *BeaconState) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(st)
}

// MarshalSSZTo ssz marshals the BeaconState object to a target array.
func (st *BeaconState) MarshalSSZTo(dst []byte) ([]byte, error) {
	// TODO: Implement proper SSZ marshaling for BeaconState
	// For now, return an error to avoid infinite recursion
	return nil, fastssz.ErrSize
}

// UnmarshalSSZ ssz unmarshals the BeaconState object.
func (st *BeaconState) UnmarshalSSZ(buf []byte) error {
	// TODO: Implement full unmarshaling logic
	// This is complex due to fork-specific fields
	return nil
}

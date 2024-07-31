// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// BeaconState represents the entire state of the beacon chain.
//
// TODO: Properly use the generics.
type BeaconState[
	BeaconBlockHeaderT interface {
		ssz.StaticObject
		Empty() BeaconBlockHeaderT
		HashTreeRootWith(fastssz.HashWalker) error
		*B
	},
	Eth1DataPtrT interface {
		ssz.StaticObject
		Empty() Eth1DataPtrT
		HashTreeRootWith(fastssz.HashWalker) error
		*E
	},
	ExecutionPayloadHeaderT interface {
		ssz.DynamicObject
		Empty() ExecutionPayloadHeaderT
		HashTreeRootWith(fastssz.HashWalker) error
		*P
	},
	ForkT interface {
		ssz.StaticObject
		Empty() ForkT
		HashTreeRootWith(fastssz.HashWalker) error
		*F
	},
	ValidatorT interface {
		ssz.StaticObject
		Empty() ValidatorT
		HashTreeRootWith(fastssz.HashWalker) error
		*V
	},
	B, E, P, F, V any,
] struct {
	// Versioning
	GenesisValidatorsRoot common.Root
	Slot                  math.Slot
	Fork                  *Fork

	// History
	LatestBlockHeader BeaconBlockHeaderT
	BlockRoots        []common.Root
	StateRoots        []common.Root

	// Eth1
	Eth1Data                     Eth1DataPtrT
	Eth1DepositIndex             uint64
	LatestExecutionPayloadHeader ExecutionPayloadHeaderT

	// Registry
	Validators []ValidatorT
	Balances   []uint64

	// Randomness
	RandaoMixes []common.Bytes32

	// Withdrawals
	NextWithdrawalIndex          uint64
	NextWithdrawalValidatorIndex math.ValidatorIndex

	// Slashing
	Slashings     []uint64
	TotalSlashing math.Gwei
}

// New creates a new BeaconState.
func (st *BeaconState[
	BeaconBlockHeaderT,
	Eth1DataPtrT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT,
	B, E, P, F, V,
]) New(
	_ uint32,
	genesisValidatorsRoot common.Root,
	slot math.Slot,
	fork ForkT,
	latestBlockHeader BeaconBlockHeaderT,
	blockRoots []common.Root,
	stateRoots []common.Root,
	eth1Data Eth1DataPtrT,
	eth1DepositIndex uint64,
	latestExecutionPayloadHeader ExecutionPayloadHeaderT,
	validators []ValidatorT,
	balances []uint64,
	randaoMixes []common.Bytes32,
	nextWithdrawalIndex uint64,
	nextWithdrawalValidatorIndex math.ValidatorIndex,
	slashings []uint64,
	totalSlashing math.Gwei,
) (*BeaconState[
	BeaconBlockHeaderT,
	Eth1DataPtrT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT,
	B, E, P, F, V,
], error) {
	return &BeaconState[
		BeaconBlockHeaderT,
		Eth1DataPtrT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
		B, E, P, F, V,
	]{
		Slot:                         slot,
		GenesisValidatorsRoot:        genesisValidatorsRoot,
		Fork:                         any(fork).(*Fork),
		LatestBlockHeader:            latestBlockHeader,
		BlockRoots:                   blockRoots,
		StateRoots:                   stateRoots,
		LatestExecutionPayloadHeader: latestExecutionPayloadHeader,
		Eth1Data:                     eth1Data,
		Eth1DepositIndex:             eth1DepositIndex,
		Validators:                   validators,
		Balances:                     balances,
		RandaoMixes:                  randaoMixes,
		NextWithdrawalIndex:          nextWithdrawalIndex,
		NextWithdrawalValidatorIndex: nextWithdrawalValidatorIndex,
		Slashings:                    slashings,
		TotalSlashing:                totalSlashing,
	}, nil
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the ssz encoded size in bytes for the BeaconState object.
func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) SizeSSZ(
	fixed bool,
) uint32 {
	var size uint32 = 300

	if fixed {
		return size
	}

	// Dynamic size fields
	size += ssz.SizeSliceOfStaticBytes(st.BlockRoots)
	size += ssz.SizeSliceOfStaticBytes(st.StateRoots)
	size += ssz.SizeDynamicObject(st.LatestExecutionPayloadHeader)
	size += ssz.SizeSliceOfStaticObjects(st.Validators)
	size += ssz.SizeSliceOfUint64s(st.Balances)
	size += ssz.SizeSliceOfStaticBytes(st.RandaoMixes)
	size += ssz.SizeSliceOfUint64s(st.Slashings)

	return size
}

// DefineSSZ defines the SSZ encoding for the BeaconState object.
//
//nolint:mnd // todo fix.
func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) DefineSSZ(
	codec *ssz.Codec,
) {
	// Versioning
	ssz.DefineStaticBytes(codec, &st.GenesisValidatorsRoot)
	ssz.DefineUint64(codec, &st.Slot)
	ssz.DefineStaticObject(codec, &st.Fork)

	// History
	ssz.DefineStaticObject(codec, &st.LatestBlockHeader)
	ssz.DefineSliceOfStaticBytesOffset(codec, &st.BlockRoots, 8192)
	ssz.DefineSliceOfStaticBytesOffset(codec, &st.StateRoots, 8192)

	// Eth1
	ssz.DefineStaticObject(codec, &st.Eth1Data)
	ssz.DefineUint64(codec, &st.Eth1DepositIndex)
	ssz.DefineDynamicObjectOffset(codec, &st.LatestExecutionPayloadHeader)

	// Registry
	ssz.DefineSliceOfStaticObjectsOffset(codec, &st.Validators, 1099511627776)
	ssz.DefineSliceOfUint64sOffset(codec, &st.Balances, 1099511627776)

	// Randomness
	ssz.DefineSliceOfStaticBytesOffset(codec, &st.RandaoMixes, 65536)

	// Withdrawals
	ssz.DefineUint64(codec, &st.NextWithdrawalIndex)
	ssz.DefineUint64(codec, &st.NextWithdrawalValidatorIndex)

	// // Slashing
	ssz.DefineSliceOfUint64sOffset(codec, &st.Slashings, 1099511627776)
	ssz.DefineUint64(codec, (*uint64)(&st.TotalSlashing))

	// Dynamic content
	ssz.DefineSliceOfStaticBytesContent(codec, &st.BlockRoots, 8192)
	ssz.DefineSliceOfStaticBytesContent(codec, &st.StateRoots, 8192)
	ssz.DefineDynamicObjectContent(codec, &st.LatestExecutionPayloadHeader)
	ssz.DefineSliceOfStaticObjectsContent(codec, &st.Validators, 1099511627776)
	ssz.DefineSliceOfUint64sContent(codec, &st.Balances, 1099511627776)
	ssz.DefineSliceOfStaticBytesContent(codec, &st.RandaoMixes, 65536)
	ssz.DefineSliceOfUint64sContent(codec, &st.Slashings, 1099511627776)
}

// MarshalSSZ marshals the BeaconState into SSZ format.
func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, st.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, st)
}

// UnmarshalSSZ unmarshals the BeaconState from SSZ format.
func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) UnmarshalSSZ(
	buf []byte,
) error {
	return ssz.DecodeFromBytes(buf, st)
}

// HashTreeRoot computes the Merkleization of the BeaconState.
func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(st), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) MarshalSSZTo(
	dst []byte,
) ([]byte, error) {
	bz, err := st.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the BeaconState object with a hasher.
//
//nolint:mnd,funlen,gocognit // todo fix.
func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) HashTreeRootWith(
	hh fastssz.HashWalker,
) error {
	indx := hh.Index()

	// Field (0) 'GenesisValidatorsRoot'
	hh.PutBytes(st.GenesisValidatorsRoot[:])

	// Field (1) 'Slot'
	hh.PutUint64(uint64(st.Slot))

	// Field (2) 'Fork'
	if st.Fork == nil {
		st.Fork = new(Fork)
	}
	if err := st.Fork.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (3) 'LatestBlockHeader'
	if st.LatestBlockHeader == nil {
		st.LatestBlockHeader = st.LatestBlockHeader.Empty()
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
		st.Eth1Data = st.Eth1Data.Empty()
	}
	if err := st.Eth1Data.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (7) 'Eth1DepositIndex'
	hh.PutUint64(st.Eth1DepositIndex)

	// Field (8) 'LatestExecutionPayloadHeader'
	if st.LatestExecutionPayloadHeader == nil {
		st.LatestExecutionPayloadHeader = st.LatestExecutionPayloadHeader.Empty()
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
		hh.AppendUint64(i)
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

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BeaconState object.
func (st *BeaconState[_, _, _, _, _, _, _, _, _, _]) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(st)
}

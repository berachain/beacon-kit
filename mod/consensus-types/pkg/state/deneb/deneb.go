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

package deneb

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// BeaconState represents the state of the Ethereum 2.0 beacon chain.
type BeaconState struct {
	// Versioning
	GenesisValidatorsRoot common.Root
	Slot                  math.Slot
	Fork                  *types.Fork

	// History
	LatestBlockHeader *types.BeaconBlockHeader
	BlockRoots        []common.Root
	StateRoots        []common.Root

	// Eth1
	Eth1Data                     *types.Eth1Data
	Eth1DepositIndex             uint64
	LatestExecutionPayloadHeader *types.ExecutionPayloadHeader

	// Registry
	Validators []*types.Validator
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

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the ssz encoded size in bytes for the BeaconState object
func (b *BeaconState) SizeSSZ(fixed bool) uint32 {
	var size uint32 = 300

	if fixed {
		return size
	}

	// Dynamic size fields
	size += ssz.SizeSliceOfStaticBytes(b.BlockRoots)
	size += ssz.SizeSliceOfStaticBytes(b.StateRoots)
	size += ssz.SizeDynamicObject(b.LatestExecutionPayloadHeader)
	size += ssz.SizeSliceOfStaticObjects(b.Validators)
	size += ssz.SizeSliceOfUint64s(b.Balances)
	size += ssz.SizeSliceOfStaticBytes(b.RandaoMixes)
	size += ssz.SizeSliceOfUint64s(b.Slashings)

	return size
}

// DefineSSZ defines the SSZ encoding for the BeaconState object.
func (b *BeaconState) DefineSSZ(codec *ssz.Codec) {
	// Versioning
	ssz.DefineStaticBytes(codec, &b.GenesisValidatorsRoot)
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineStaticObject(codec, &b.Fork)

	// History
	ssz.DefineStaticObject(codec, &b.LatestBlockHeader)
	ssz.DefineSliceOfStaticBytesOffset(codec, &b.BlockRoots, 8192)
	ssz.DefineSliceOfStaticBytesOffset(codec, &b.StateRoots, 8192)

	// Eth1
	ssz.DefineStaticObject(codec, &b.Eth1Data)
	ssz.DefineUint64(codec, &b.Eth1DepositIndex)
	ssz.DefineDynamicObjectOffset(codec, &b.LatestExecutionPayloadHeader)

	// Registry
	ssz.DefineSliceOfStaticObjectsOffset(codec, &b.Validators, 1099511627776)
	ssz.DefineSliceOfUint64sOffset(codec, &b.Balances, 1099511627776)

	// Randomness
	ssz.DefineSliceOfStaticBytesOffset(codec, &b.RandaoMixes, 65536)

	// Withdrawals
	ssz.DefineUint64(codec, &b.NextWithdrawalIndex)
	ssz.DefineUint64(codec, &b.NextWithdrawalValidatorIndex)

	// // Slashing
	ssz.DefineSliceOfUint64sOffset(codec, &b.Slashings, 1099511627776)
	ssz.DefineUint64(codec, (*uint64)(&b.TotalSlashing))

	// Dynamic content
	ssz.DefineSliceOfStaticBytesContent(codec, &b.BlockRoots, 8192)
	ssz.DefineSliceOfStaticBytesContent(codec, &b.StateRoots, 8192)
	ssz.DefineDynamicObjectContent(codec, &b.LatestExecutionPayloadHeader)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.Validators, 1099511627776)
	ssz.DefineSliceOfUint64sContent(codec, &b.Balances, 1099511627776)
	ssz.DefineSliceOfStaticBytesContent(codec, &b.RandaoMixes, 65536)
	ssz.DefineSliceOfUint64sContent(codec, &b.Slashings, 1099511627776)
}

func (b *BeaconState) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, b)
}

func (b *BeaconState) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

func (b *BeaconState) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(b), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

func (b *BeaconState) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := b.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the BeaconState object with a hasher
func (b *BeaconState) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'GenesisValidatorsRoot'
	hh.PutBytes(b.GenesisValidatorsRoot[:])

	// Field (1) 'Slot'
	hh.PutUint64(uint64(b.Slot))

	// Field (2) 'Fork'
	if b.Fork == nil {
		b.Fork = new(types.Fork)
	}
	if err := b.Fork.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (3) 'LatestBlockHeader'
	if b.LatestBlockHeader == nil {
		b.LatestBlockHeader = new(types.BeaconBlockHeader)
	}
	if err := b.LatestBlockHeader.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (4) 'BlockRoots'
	if size := len(b.BlockRoots); size > 8192 {
		return fastssz.ErrListTooBigFn("BeaconState.BlockRoots", size, 8192)
	}
	subIndx := hh.Index()
	for _, i := range b.BlockRoots {
		hh.Append(i[:])
	}
	numItems := uint64(len(b.BlockRoots))
	hh.MerkleizeWithMixin(subIndx, numItems, 8192)

	// Field (5) 'StateRoots'
	if size := len(b.StateRoots); size > 8192 {
		return fastssz.ErrListTooBigFn("BeaconState.StateRoots", size, 8192)
	}
	subIndx = hh.Index()
	for _, i := range b.StateRoots {
		hh.Append(i[:])
	}
	numItems = uint64(len(b.StateRoots))
	hh.MerkleizeWithMixin(subIndx, numItems, 8192)

	// Field (6) 'Eth1Data'
	if b.Eth1Data == nil {
		b.Eth1Data = new(types.Eth1Data)
	}
	if err := b.Eth1Data.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (7) 'Eth1DepositIndex'
	hh.PutUint64(b.Eth1DepositIndex)

	// Field (8) 'LatestExecutionPayloadHeader'
	if err := b.LatestExecutionPayloadHeader.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (9) 'Validators'
	subIndx = hh.Index()
	num := uint64(len(b.Validators))
	if num > 1099511627776 {
		return fastssz.ErrIncorrectListSize
	}
	for _, elem := range b.Validators {
		if err := elem.HashTreeRootWith(hh); err != nil {
			return err
		}
	}
	hh.MerkleizeWithMixin(subIndx, num, 1099511627776)

	// Field (10) 'Balances'
	if size := len(b.Balances); size > 1099511627776 {
		return fastssz.ErrListTooBigFn(
			"BeaconState.Balances",
			size,
			1099511627776,
		)
	}
	subIndx = hh.Index()
	for _, i := range b.Balances {
		hh.AppendUint64(i)
	}
	hh.FillUpTo32()
	numItems = uint64(len(b.Balances))
	hh.MerkleizeWithMixin(
		subIndx,
		numItems,
		fastssz.CalculateLimit(1099511627776, numItems, 8),
	)

	// Field (11) 'RandaoMixes'
	if size := len(b.RandaoMixes); size > 65536 {
		return fastssz.ErrListTooBigFn("BeaconState.RandaoMixes", size, 65536)
	}
	subIndx = hh.Index()
	for _, i := range b.RandaoMixes {
		hh.Append(i[:])
	}
	numItems = uint64(len(b.RandaoMixes))
	hh.MerkleizeWithMixin(subIndx, numItems, 65536)

	// Field (12) 'NextWithdrawalIndex'
	hh.PutUint64(b.NextWithdrawalIndex)

	// Field (13) 'NextWithdrawalValidatorIndex'
	hh.PutUint64(uint64(b.NextWithdrawalValidatorIndex))

	// Field (14) 'Slashings'
	if size := len(b.Slashings); size > 1099511627776 {
		return fastssz.ErrListTooBigFn(
			"BeaconState.Slashings",
			size,
			1099511627776,
		)
	}
	subIndx = hh.Index()
	for _, i := range b.Slashings {
		hh.AppendUint64(i)
	}
	hh.FillUpTo32()
	numItems = uint64(len(b.Slashings))
	hh.MerkleizeWithMixin(
		subIndx,
		numItems,
		fastssz.CalculateLimit(1099511627776, numItems, 8),
	)

	// Field (15) 'TotalSlashing'
	hh.PutUint64(uint64(b.TotalSlashing))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BeaconState object
func (b *BeaconState) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}

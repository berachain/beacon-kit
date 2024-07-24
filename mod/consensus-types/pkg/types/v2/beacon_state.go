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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/karalabe/ssz"
)

// Cached static size computed on package init.
var staticSizeCacheBeaconState = 32 + 8 + (*Fork)(nil).SizeSSZ() + (*BeaconBlockHeader)(nil).SizeSSZ() + 4 + 4 + (*Eth1Data)(nil).SizeSSZ() + 8 + 4 + 4 + 4 + 4 + 8 + 8 + 4 + 8

type BeaconState struct {
	// Versioning
	//
	//nolint:lll
	GenesisValidatorsRoot common.Root
	Slot                  math.Slot
	Fork                  *Fork

	// History
	LatestBlockHeader *BeaconBlockHeader
	BlockRoots        []common.Root
	StateRoots        []common.Root

	// Eth1
	Eth1Data                     *Eth1Data
	Eth1DepositIndex             math.U64
	LatestExecutionPayloadHeader *ExecutionPayloadHeader

	// Registry
	Validators []*Validator
	Balances   []uint64

	// Randomness
	RandaoMixes []common.Bytes32

	// Withdrawals
	NextWithdrawalIndex          math.U64
	NextWithdrawalValidatorIndex math.ValidatorIndex

	// Slashing
	Slashings     []uint64
	TotalSlashing math.Gwei
}

// SizeSSZ returns either the static size of the object if fixed == true, or
// the total size otherwise.
func (obj *BeaconState) SizeSSZ(fixed bool) uint32 {
	var size = uint32(staticSizeCacheBeaconState)
	if fixed {
		return size
	}
	size += ssz.SizeSliceOfStaticBytes(obj.BlockRoots)
	size += ssz.SizeSliceOfStaticBytes(obj.StateRoots)
	size += ssz.SizeDynamicObject(obj.LatestExecutionPayloadHeader)
	size += ssz.SizeSliceOfStaticObjects(obj.Validators)
	size += ssz.SizeSliceOfUint64s(obj.Balances)
	size += ssz.SizeSliceOfStaticBytes(obj.RandaoMixes)
	size += ssz.SizeSliceOfUint64s(obj.Slashings)

	return size
}

// DefineSSZ defines how an object is encoded/decoded.
func (obj *BeaconState) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineStaticBytes(codec, &obj.GenesisValidatorsRoot)                    // Field  ( 0) -        GenesisValidatorsRoot - 32 bytes
	ssz.DefineUint64(codec, &obj.Slot)                                          // Field  ( 1) -                         Slot -  8 bytes
	ssz.DefineStaticObject(codec, &obj.Fork)                                    // Field  ( 2) -                         Fork -  ? bytes (Fork)
	ssz.DefineStaticObject(codec, &obj.LatestBlockHeader)                       // Field  ( 3) -            LatestBlockHeader -  ? bytes (BeaconBlockHeader)
	ssz.DefineSliceOfStaticBytesOffset(codec, &obj.BlockRoots, 8192)            // Offset ( 4) -                   BlockRoots -  4 bytes
	ssz.DefineSliceOfStaticBytesOffset(codec, &obj.StateRoots, 8192)            // Offset ( 5) -                   StateRoots -  4 bytes
	ssz.DefineStaticObject(codec, &obj.Eth1Data)                                // Field  ( 6) -                     Eth1Data -  ? bytes (Eth1Data)
	ssz.DefineUint64(codec, &obj.Eth1DepositIndex)                              // Field  ( 7) -             Eth1DepositIndex -  8 bytes
	ssz.DefineDynamicObjectOffset(codec, &obj.LatestExecutionPayloadHeader)     // Offset ( 8) - LatestExecutionPayloadHeader -  4 bytes
	ssz.DefineSliceOfStaticObjectsOffset(codec, &obj.Validators, 1099511627776) // Offset ( 9) -                   Validators -  4 bytes
	ssz.DefineSliceOfUint64sOffset(codec, &obj.Balances, 1099511627776)         // Offset (10) -                     Balances -  4 bytes
	ssz.DefineSliceOfStaticBytesOffset(codec, &obj.RandaoMixes, 65536)          // Offset (11) -                  RandaoMixes -  4 bytes
	ssz.DefineUint64(codec, &obj.NextWithdrawalIndex)                           // Field  (12) -          NextWithdrawalIndex -  8 bytes
	ssz.DefineUint64(codec, &obj.NextWithdrawalValidatorIndex)                  // Field  (13) - NextWithdrawalValidatorIndex -  8 bytes
	ssz.DefineSliceOfUint64sOffset(codec, &obj.Slashings, 1099511627776)        // Offset (14) -                    Slashings -  4 bytes
	ssz.DefineUint64(codec, &obj.TotalSlashing)                                 // Field  (15) -                TotalSlashing -  8 bytes

	// Define the dynamic data (fields)
	ssz.DefineSliceOfStaticBytesContent(codec, &obj.BlockRoots, 8192)            // Field  ( 4) -                   BlockRoots - ? bytes
	ssz.DefineSliceOfStaticBytesContent(codec, &obj.StateRoots, 8192)            // Field  ( 5) -                   StateRoots - ? bytes
	ssz.DefineDynamicObjectContent(codec, &obj.LatestExecutionPayloadHeader)     // Field  ( 8) - LatestExecutionPayloadHeader - ? bytes
	ssz.DefineSliceOfStaticObjectsContent(codec, &obj.Validators, 1099511627776) // Field  ( 9) -                   Validators - ? bytes
	ssz.DefineSliceOfUint64sContent(codec, &obj.Balances, 1099511627776)         // Field  (10) -                     Balances - ? bytes
	ssz.DefineSliceOfStaticBytesContent(codec, &obj.RandaoMixes, 65536)          // Field  (11) -                  RandaoMixes - ? bytes
	ssz.DefineSliceOfUint64sContent(codec, &obj.Slashings, 1099511627776)        // Field  (14) -                    Slashings - ? bytes
}

func (s *BeaconState) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(s), nil
}

func (s *BeaconState) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, s)
}

func (s *BeaconState) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, s)
}

func (s *BeaconState) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, s.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, s)
}

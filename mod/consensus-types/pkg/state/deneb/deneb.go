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
	"github.com/karalabe/ssz"
)

type BeaconState struct {
	// Versioning
	//
	//nolint:lll
	GenesisValidatorsRoot common.Root `json:"genesisValidatorsRoot" ssz-size:"32"`
	Slot                  math.Slot   `json:"slot"`
	Fork                  *types.Fork `json:"fork"`

	// History
	LatestBlockHeader *types.BeaconBlockHeader `json:"latestBlockHeader"`
	BlockRoots        []common.Root            `json:"blockRoots" ssz-size:"8192"`
	StateRoots        []common.Root            `json:"stateRoots" ssz-size:"8192"`

	// Eth1
	Eth1Data                     *types.Eth1Data               `json:"eth1Data"`
	Eth1DepositIndex             uint64                        `json:"eth1DepositIndex"`
	LatestExecutionPayloadHeader *types.ExecutionPayloadHeader `json:"latestExecutionPayloadHeader"`

	// Registry
	Validators []*types.Validator `ssz-max:"1099511627776"`
	Balances   []uint64           `ssz-max:"1099511627776"`

	// Randomness
	RandaoMixes []common.Bytes32 `json:"randaoMixes" ssz-size:"65536"`

	// Withdrawals
	NextWithdrawalIndex          uint64              `json:"nextWithdrawalIndex"`
	NextWithdrawalValidatorIndex math.ValidatorIndex `json:"nextWithdrawalValidatorIndex"`

	// Slashing
	Slashings     []uint64  `json:"slashings" ssz-size:"8192"`
	TotalSlashing math.Gwei `json:"totalSlashing"`
}

// SizeSSZ returns the ssz encoded size in bytes for the BeaconState object
func (b *BeaconState) SizeSSZ(fixed bool) uint32 {
	var size uint32 = 300

	if fixed {
		return size
	}

	// // Dynamic size fields
	size += ssz.SizeSliceOfStaticBytes(b.BlockRoots)
	size += ssz.SizeSliceOfStaticBytes(b.StateRoots)
	size += ssz.SizeDynamicObject(b.LatestExecutionPayloadHeader)
	size += ssz.SizeSliceOfStaticObjects(b.Validators)
	size += ssz.SizeSliceOfUint64s(b.Balances)
	// size += ssz.SizeSliceOfStaticBytes(b.RandaoMixes)
	// size += ssz.SizeSliceOfUint64s(b.Slashings)

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
	// ssz.DefineSliceOfStaticBytesOffset(codec, &b.RandaoMixes, 65536)

	// Withdrawals
	ssz.DefineUint64(codec, &b.NextWithdrawalIndex)
	ssz.DefineUint64(codec, &b.NextWithdrawalValidatorIndex)

	// // Slashing
	// ssz.DefineSliceOfUint64sOffset(codec, &b.Slashings, 8192)
	ssz.DefineUint64(codec, (*uint64)(&b.TotalSlashing))

	// Dynamic content
	ssz.DefineSliceOfStaticBytesContent(codec, &b.BlockRoots, 8192)
	ssz.DefineSliceOfStaticBytesContent(codec, &b.StateRoots, 8192)
	ssz.DefineDynamicObjectContent(codec, &b.LatestExecutionPayloadHeader)
	ssz.DefineSliceOfStaticObjectsContent(codec, &b.Validators, 1099511627776)
	ssz.DefineSliceOfUint64sContent(codec, &b.Balances, 1099511627776)
	// ssz.DefineSliceOfStaticBytesContent(codec, &b.RandaoMixes, 65536)
	// ssz.DefineSliceOfUint64sContent(codec, &b.Slashings, 1099511627776)
}

func (b *BeaconState) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, b)
}

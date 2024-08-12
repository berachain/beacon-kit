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

package block

import (
	sdkcollections "cosmossdk.io/collections"
	sdkindexes "cosmossdk.io/collections/indexes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
)

const (
	blockRootsIndexName       = "block_roots"
	executionNumbersIndexName = "execution_numbers"
	stateRootsIndexName       = "state_roots"
)

type indexes[BeaconBlockT BeaconBlock[BeaconBlockT]] struct {
	BlockRoots       *sdkindexes.Unique[[]byte, math.Slot, BeaconBlockT]
	ExecutionNumbers *sdkindexes.Unique[math.U64, math.Slot, BeaconBlockT]
	StateRoots       *sdkindexes.Unique[[]byte, math.Slot, BeaconBlockT]
}

// IndexesList returns a list of all indexes associated with the
// validatorsIndex.
func (i indexes[BeaconBlockT]) IndexesList() []sdkcollections.Index[
	math.Slot, BeaconBlockT,
] {
	return []sdkcollections.Index[math.Slot, BeaconBlockT]{
		i.BlockRoots,
		i.ExecutionNumbers,
		i.StateRoots,
	}
}

func newIndexes[BeaconBlockT BeaconBlock[BeaconBlockT]](
	sb *sdkcollections.SchemaBuilder,
) indexes[BeaconBlockT] {
	return indexes[BeaconBlockT]{
		BlockRoots: sdkindexes.NewUnique(
			sb,
			sdkcollections.NewPrefix(blockRootsIndexName),
			blockRootsIndexName,
			sdkcollections.BytesKey,
			encoding.U64Key,
			func(_ math.Slot, blk BeaconBlockT) ([]byte, error) {
				root := blk.HashTreeRoot()
				return root[:], nil
			},
		),
		ExecutionNumbers: sdkindexes.NewUnique(
			sb,
			sdkcollections.NewPrefix(executionNumbersIndexName),
			executionNumbersIndexName,
			encoding.U64Key,
			encoding.U64Key,
			func(_ math.Slot, blk BeaconBlockT) (math.U64, error) {
				return blk.GetExecutionNumber(), nil
			},
		),
		StateRoots: sdkindexes.NewUnique(
			sb,
			sdkcollections.NewPrefix(stateRootsIndexName),
			stateRootsIndexName,
			sdkcollections.BytesKey,
			encoding.U64Key,
			func(_ math.Slot, blk BeaconBlockT) ([]byte, error) {
				root := blk.GetStateRoot()
				return root[:], nil
			},
		),
	}
}

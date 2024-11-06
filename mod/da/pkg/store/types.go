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

package store

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlock is an interface for beacon blocks.
type BeaconBlock interface {
	GetSlot() math.U64
}

// BlockEvent is an interface for block events.
type BlockEvent[BeaconBlockT BeaconBlock] interface {
	Data() BeaconBlockT
}

// IndexDB is a database that allows prefixing by index.
type IndexDB interface {
	Has(index uint64, key []byte) (bool, error)
	Set(index uint64, key []byte, value []byte) error

	// Prune returns error if start > end
	Prune(start uint64, end uint64) error
}

// BeaconBlockBody is the body of a beacon block.
type BeaconBlockBody interface {
	// GetBlobKzgCommitments returns the KZG commitments for the blob.
	GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
}

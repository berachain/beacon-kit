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
	"errors"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche"
)

var (
	// errNilBeaconBlock is an error for when
	// the beacon block in an abci request is nil.
	errNilBeaconBlock = errors.New("nil beacon block")

	// errNilBlock is an error for when the abci request
	// is nil.
	errNilBlock = errors.New("nil abci request")
)

func (b *StatelessBlock) GetBeaconBlock(
	forkVersion uint32,
) (miniavalanche.BeaconBlockT, error) {
	var blk miniavalanche.BeaconBlockT
	if b == nil {
		return blk, errNilBlock
	}
	if b.BlkContent.BeaconBlockByte == nil {
		return blk, errNilBeaconBlock
	}
	return blk.NewFromSSZ(b.BlkContent.BeaconBlockByte, forkVersion)
}

func (b *StatelessBlock) GetBlobSidecars() (
	miniavalanche.BlobSidecarsT,
	error,
) {
	var sidecars miniavalanche.BlobSidecarsT
	if b == nil {
		return sidecars, errNilBlock
	}
	if b.BlkContent.BlobsBytes == nil {
		return sidecars, errNilBeaconBlock
	}
	sidecars = sidecars.Empty()
	return sidecars, sidecars.UnmarshalSSZ(b.BlkContent.BlobsBytes)
}

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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// BeaconBlockDenebPlus represents a block in the beacon chain during
// the DenebPlus fork.
type BeaconBlockDenebPlus struct {
	// BeaconBlockHeaderBase is the base of the BeaconBlockDenebPlus.
	BeaconBlockHeaderBase
	// Body is the body of the BeaconBlockDenebPlus, containing the block's
	// operations.
	Body *BeaconBlockBodyDenebPlus
}

// Version identifies the version of the BeaconBlockDenebPlus.
func (b *BeaconBlockDenebPlus) Version() uint32 {
	return version.DenebPlus
}

// IsNil checks if the BeaconBlockDenebPlus instance is nil.
func (b *BeaconBlockDenebPlus) IsNil() bool {
	return b == nil
}

// SetStateRoot sets the state root of the BeaconBlockDenebPlus.
func (b *BeaconBlockDenebPlus) SetStateRoot(root common.Root) {
	b.StateRoot = root
}

// GetBody retrieves the body of the BeaconBlockDenebPlus.
func (b *BeaconBlockDenebPlus) GetBody() *BeaconBlockBody {
	return &BeaconBlockBody{RawBeaconBlockBody: b.Body}
}

// GetHeader builds a BeaconBlockHeader from the BeaconBlockDenebPlus.
func (b BeaconBlockDenebPlus) GetHeader() *BeaconBlockHeader {
	bodyRoot, err := b.GetBody().HashTreeRoot()
	if err != nil {
		return nil
	}

	return &BeaconBlockHeader{
		BeaconBlockHeaderBase: BeaconBlockHeaderBase{
			Slot:            b.Slot,
			ProposerIndex:   b.ProposerIndex,
			ParentBlockRoot: b.ParentBlockRoot,
			StateRoot:       b.StateRoot,
		},
		BodyRoot: bodyRoot,
	}
}

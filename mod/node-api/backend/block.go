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

package backend

import (
	"fmt"

	types "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BlockHeader returns the block header at the given slot.
func (b Backend[
	_, _, _, BeaconBlockHeaderT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
	_,
]) BlockHeaderAtSlot(slot math.Slot) (BeaconBlockHeaderT, error) {
	var blockHeader BeaconBlockHeaderT

	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return blockHeader, err
	}

	blockHeader, err = st.GetLatestBlockHeader()
	return blockHeader, err
}

// GetBlockRoot returns the root of the block at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlockRootAtSlot(slot math.Slot) (common.Root, error) {
	st, slot, err := b.stateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}

	// As calculated by the beacon chain. Ideally, this logic
	// should be abstracted by the beacon chain.
	return st.GetBlockRootAtIndex(slot.Unwrap() % b.cs.SlotsPerHistoricalRoot())
}

// TODO: Implement this.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlockRewardsAtSlot(math.Slot) (*types.BlockRewardsData, error) {
	return &types.BlockRewardsData{
		ProposerIndex:     1,
		Total:             1,
		Attestations:      1,
		SyncAggregate:     1,
		ProposerSlashings: 1,
		AttesterSlashings: 1,
	}, nil
}

func (b Backend[
	_, _, _, BeaconBlockHeaderT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
	_,
]) BlockSignatureAtSlot(slot math.Slot) (bytes.B48, error) {
	// Get the block header
	header, err := b.BlockHeaderAtSlot(slot)
	if err != nil {
		return bytes.B48{}, err
	}

	// Extract the signature from the header
	// This will depend on how your BeaconBlockHeaderT type stores the signature
	signature, err := b.extractSignatureFromHeader(header)
	if err != nil {
		return bytes.B48{}, err
	}

	return signature, nil
}

func (b Backend[
	_, _, _, BeaconBlockHeaderT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
	_,
]) extractSignatureFromHeader(header BeaconBlockHeaderT) (bytes.B48, error) {
	// Implementation depends on how signatures are stored in your header structure
	// This is a placeholder implementation
	rawSignature := header.GetSignatures()
	if len(rawSignature) != 48 {
		return bytes.B48{}, fmt.Errorf("invalid signature length: expected 48, got %d", len(rawSignature))
	}

	//var signature bytes.B48
	//copy(signature[:], rawSignature)
	return rawSignature, nil
}

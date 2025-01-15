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

package blockchain

import (
	"github.com/berachain/beacon-kit/chain-spec/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
)

func (s *Service[
	_, _, ConsensusBlockT, _, _, _,
]) processPruning(beaconBlk *ctypes.BeaconBlock) error {
	// prune availability store
	start, end := availabilityPruneRangeFn(
		beaconBlk.GetSlot().Unwrap(), s.chainSpec)
	err := s.storageBackend.AvailabilityStore().Prune(start, end)
	if err != nil {
		return err
	}

	// prune deposit store
	start, end = depositPruneRangeFn(
		beaconBlk.GetBody().GetDeposits(), s.chainSpec)
	err = s.storageBackend.DepositStore().Prune(start, end)

	if err != nil {
		return err
	}

	return nil
}

func depositPruneRangeFn([]*ctypes.Deposit, chain.ChainSpec) (uint64, uint64) {
	// The whole deposit list is validated in consensus and its Merkle root is part of
	// Beacon State. Therefore every node must keep the full deposit list and deposits
	// pruning must be turned off.
	return 0, 0
}

//nolint:unparam // this is ok
func availabilityPruneRangeFn(
	slot uint64, cs chain.ChainSpec) (uint64, uint64) {
	window := cs.MinEpochsForBlobsSidecarsRequest() * cs.SlotsPerEpoch()
	if slot < window {
		return 0, 0
	}

	return 0, slot - window
}

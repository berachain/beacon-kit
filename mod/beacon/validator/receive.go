// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package validator

import (
	"context"
)

// VerifyIncomingBlobs receives blobs from the network and processes them.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) VerifyIncomingBlobs(
	_ context.Context,
	blk BeaconBlockT,
	blobs BlobSidecarsT,
) error {
	if blk.IsNil() {
		s.logger.Error(
			"aborting blob verification on nil block ⛔️ ",
		)
		return ErrNilBlk
	}

	s.logger.Info(
		"received incoming blob sidecars 🚔 ",
		"state_root", blk.GetStateRoot(),
	)

	if err := s.verifyBlobProofs(blk.GetSlot(), blobs); err != nil {
		s.logger.Error(
			"rejecting incoming blob sidecars ❌ ",
			"error", err,
		)
		return err
	}

	s.logger.Info(
		"blob sidecars verification succeeded - accepting incoming blobs 💦 ",
		"num_blobs", blobs.Len(),
	)
	return nil
}

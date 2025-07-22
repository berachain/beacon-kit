// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package state

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
)

// ParentProposerPubkey returns the parent proposer pubkey for the given timestamp.
// It must return nil if we are before Electra1.
//
//nolint:nilnil // TODO: consider addressing this
func (s *StateDB) ParentProposerPubkey(timestamp math.U64) (*crypto.BLSPubkey, error) {
	if version.IsBefore(s.cs.ActiveForkVersionForTimestamp(timestamp), version.Electra1()) {
		return nil, nil
	}

	latestBlockHeader, err := s.GetLatestBlockHeader()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving latest block header: %w", err)
	}
	prevProposer, err := s.ValidatorByIndex(latestBlockHeader.GetProposerIndex())
	if err != nil {
		return nil, fmt.Errorf("failed retrieving prev proposer: %w", err)
	}
	p := prevProposer.GetPubkey()
	return &p, nil
}

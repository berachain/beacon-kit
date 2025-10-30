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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blobreactor_test

import (
	"errors"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/da/blobreactor"
	"github.com/cometbft/cometbft/p2p"
	"github.com/stretchr/testify/require"
)

func TestReputationManager_BanAndExpiration(t *testing.T) {
	t.Parallel()
	logger := log.NewTestLogger(t)
	config := blobreactor.DefaultReputationConfig()
	config.BanPeriod = 200 * time.Millisecond
	rm := blobreactor.NewReputationManager(logger, config)

	peerID := p2p.ID("test-peer")

	// Record violations until banned
	violationsNeeded := (config.MaxReputationScore-config.DisconnectThreshold)/config.BadBehaviorPenalty + 1
	for range violationsNeeded {
		rm.RecordBadBehavior(peerID, errors.New("test violation"))
	}

	// Peer should be banned
	require.False(t, rm.ShouldAcceptPeer(peerID))

	// Peer should be accepted with reset score after ban expires
	require.Eventually(t, func() bool {
		return rm.ShouldAcceptPeer(peerID)
	}, config.BanPeriod+50*time.Millisecond, 10*time.Millisecond)
}

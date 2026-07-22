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
//

package cometbft

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/log/phuslu"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
)

func TestHaltPointReached(t *testing.T) {
	t.Parallel()

	blockTime := time.Unix(1700000000, 0)

	tests := []struct {
		name       string
		haltHeight uint64
		haltTime   uint64
		height     int64
		blockTime  time.Time
		want       bool
	}{
		{name: "disabled", height: 100, blockTime: blockTime, want: false},
		{name: "below halt height", haltHeight: 100, height: 99, blockTime: blockTime, want: false},
		{name: "at halt height", haltHeight: 100, height: 100, blockTime: blockTime, want: true},
		{name: "past halt height", haltHeight: 100, height: 101, blockTime: blockTime, want: true},
		{name: "before halt time", haltTime: 1700000001, height: 100, blockTime: blockTime, want: false},
		{name: "at halt time", haltTime: 1700000000, height: 100, blockTime: blockTime, want: true},
		{name: "past halt time", haltTime: 1699999999, height: 100, blockTime: blockTime, want: true},
		{name: "zero block time does not halt", haltTime: 1700000000, height: 100, blockTime: time.Time{}, want: false},
		{name: "halt time reached before halt height", haltHeight: 200, haltTime: 1700000000, height: 100, blockTime: blockTime, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, haltPointReached(tt.haltHeight, tt.haltTime, tt.height, tt.blockTime))
		})
	}
}

// TestEnsureNotHalted pins ensureNotHalted to the last finalized block. FinalizeBlock at height N runs it before
// updating the finalized fields, so state held at N-1 means "may block N be finalized".
func TestEnsureNotHalted(t *testing.T) {
	t.Parallel()

	s := &Service{}
	require.NoError(t, s.ensureNotHalted(), "halt disabled")

	s.haltHeight = 10
	s.finalizedHeight = 9
	require.NoError(t, s.ensureNotHalted(), "halt block itself must finalize")

	s.finalizedHeight = 10
	require.Error(t, s.ensureNotHalted(), "block past halt height must be refused")

	s = &Service{haltTime: 1700000000}
	require.NoError(t, s.ensureNotHalted(), "unseeded finalized time must not refuse")

	s.finalizedTime = time.Unix(1700000000, 0)
	require.Error(t, s.ensureNotHalted(), "block after halt time must be refused")
}

// haltedService returns a Service whose halt point has been reached.
func haltedService(ctx context.Context) *Service {
	return &Service{
		logger:          phuslu.NewLogger(io.Discard, nil),
		ctx:             ctx,
		haltHeight:      10,
		finalizedHeight: 10,
	}
}

// TestFinalizeBlockParksAtHaltPoint pins the FinalizeBlock halt gate. A block past the halt point must park
// until shutdown rather than return an error, which CometBFT escalates to a CONSENSUS FAILURE panic.
func TestFinalizeBlockParksAtHaltPoint(t *testing.T) {
	t.Parallel() // safe, no other test reads haltShutdownSlack

	restore := haltShutdownSlack
	haltShutdownSlack = 50 * time.Millisecond
	t.Cleanup(func() { haltShutdownSlack = restore })

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	s := haltedService(ctx)

	done := make(chan error, 1)
	go func() {
		_, err := s.FinalizeBlock(ctx, &cmtabci.FinalizeBlockRequest{Height: 11})
		done <- err
	}()

	select {
	case err := <-done:
		t.Fatalf("FinalizeBlock returned instead of parking until shutdown: %v", err)
	case <-time.After(250 * time.Millisecond):
	}

	cancel()
	select {
	case err := <-done:
		require.ErrorContains(t, err, "halt point")
	case <-time.After(5 * time.Second):
		t.Fatal("FinalizeBlock did not return after context cancel plus slack")
	}
}

// TestProposalHandlersGateAtHaltPoint pins the non-panicking proposal gates,
// which keep a halting validator from helping decide a post-halt block.
func TestProposalHandlersGateAtHaltPoint(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	s := haltedService(ctx)

	prepResp, err := s.PrepareProposal(ctx, &cmtabci.PrepareProposalRequest{Height: 11})
	require.NoError(t, err)
	require.Empty(t, prepResp.Txs, "halted node must propose nothing")

	procResp, err := s.ProcessProposal(ctx, &cmtabci.ProcessProposalRequest{Height: 11})
	require.NoError(t, err)
	require.Equal(t, cmtabci.PROCESS_PROPOSAL_STATUS_REJECT, procResp.Status)
}

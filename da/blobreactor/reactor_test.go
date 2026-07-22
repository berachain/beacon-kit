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

//nolint:paralleltest // Tests cannot run in parallel due to race condition in CometBFT's p2p.MakeConnectedSwitches
package blobreactor_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/blobreactor"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/mock"
	"github.com/stretchr/testify/require"
)

const maxBlobsPerBlock = 6

type stubBlobStore struct {
	blobs map[uint64][][]byte // slot -> raw sidecars
	delay time.Duration       // optional delay to simulate slow responses
}

func newStubBlobStore() *stubBlobStore {
	return &stubBlobStore{blobs: make(map[uint64][][]byte)}
}

func (m *stubBlobStore) GetByIndex(index uint64) ([][]byte, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.blobs[index], nil
}

func (m *stubBlobStore) setBlobs(t *testing.T, slot uint64, sidecars datypes.BlobSidecars) {
	t.Helper()
	sidecarBzs := make([][]byte, len(sidecars))
	for i, sidecar := range sidecars {
		data, err := sidecar.MarshalSSZ()
		require.NoError(t, err)
		sidecarBzs[i] = data
	}
	m.blobs[slot] = sidecarBzs
}

// noopVerifier accepts any sidecars; content verification is exercised
// elsewhere (the reactor treats it as a black box).
type noopVerifier struct{}

func (noopVerifier) VerifySidecars(
	context.Context,
	datypes.BlobSidecars,
	*ctypes.BeaconBlockHeader,
	eip4844.KZGCommitments[common.ExecutionHash],
) error {
	return nil
}

type noOpTelemetrySink struct{}

func (noOpTelemetrySink) IncrementCounter(string, ...string)        {}
func (noOpTelemetrySink) SetGauge(string, int64, ...string)         {}
func (noOpTelemetrySink) MeasureSince(string, time.Time, ...string) {}

// createTestSidecars builds structurally valid sidecars all bound to the same
// header at the given slot.
func createTestSidecars(t *testing.T, slot math.Slot, count int) datypes.BlobSidecars {
	t.Helper()
	sidecars := make(datypes.BlobSidecars, count)
	for i := range count {
		sidecars[i] = &datypes.BlobSidecar{
			Index: uint64(i), //#nosec:G115 // test
			SignedBeaconBlockHeader: &ctypes.SignedBeaconBlockHeader{
				Header: &ctypes.BeaconBlockHeader{Slot: slot},
			},
			InclusionProof: make([]common.Root, ctypes.KZGInclusionProofDepth),
		}
	}
	return sidecars
}

func blockRootOf(sidecars datypes.BlobSidecars) common.Root {
	return sidecars[0].GetBeaconBlockHeader().HashTreeRoot()
}

func newTestReactor(t *testing.T, store blobreactor.BlobStore, cfg blobreactor.Config) *blobreactor.BlobReactor {
	t.Helper()
	return blobreactor.NewBlobReactor(store, noopVerifier{}, log.NewTestLogger(t), cfg, maxBlobsPerBlock, noOpTelemetrySink{})
}

func testConfig() blobreactor.Config {
	return blobreactor.Config{
		RequestTimeout: 2 * time.Second,
		FetchTimeout:   4 * time.Second,
	}
}

func makeConnectedReactors(
	t *testing.T, n int, stores []*stubBlobStore,
) ([]*blobreactor.BlobReactor, []*p2p.Switch) {
	t.Helper()

	p2pConfig := config.DefaultP2PConfig()
	p2pConfig.ListenAddress = "tcp://127.0.0.1:0"
	p2pConfig.RootDir = t.TempDir()

	reactors := make([]*blobreactor.BlobReactor, n)
	for i := range n {
		reactors[i] = newTestReactor(t, stores[i], testConfig())
	}

	initSwitch := func(i int, sw *p2p.Switch) *p2p.Switch {
		sw.AddReactor(blobreactor.ReactorName, reactors[i])
		return sw
	}
	switches := p2p.MakeConnectedSwitches(p2pConfig, n, initSwitch, p2p.Connect2Switches)
	return reactors, switches
}

func stopSwitches(switches []*p2p.Switch) {
	for _, s := range switches {
		_ = s.Stop()
	}
}

// requireCount returns a verifier enforcing the expected sidecar count, the
// same guarantee real callers provide (count vs block commitments).
func requireCount(expected int) func(datypes.BlobSidecars) error {
	return func(sidecars datypes.BlobSidecars) error {
		if len(sidecars) != expected {
			return fmt.Errorf("expected %d sidecars, got %d", expected, len(sidecars))
		}
		return nil
	}
}

// Short data must never read as success: a peer with nothing (empty response, a miss) and a peer with fewer
// sidecars than the block's commitments (count mismatch, a verification failure) must both count as peer
// failures. This is the exact failure mode that allowed silent DA gaps in the earlier design.
func TestBlobReactor_ByRootShortDataIsFailure(t *testing.T) {
	slot := math.Slot(99)
	sidecars := createTestSidecars(t, slot, 2)

	partialStore := newStubBlobStore()
	partialStore.setBlobs(t, slot.Unwrap(), sidecars[:1]) // one of the block's two sidecars

	reactors, switches := makeConnectedReactors(
		t, 3, []*stubBlobStore{newStubBlobStore(), newStubBlobStore(), partialStore})
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	got, err := reactors[0].RequestSidecarsByRoot(
		t.Context(), slot, blockRootOf(sidecars), requireCount(2), // block commits to 2 blobs
	)
	require.Error(t, err)
	require.ErrorIs(t, err, blobreactor.ErrAllPeersFailed)
	require.Nil(t, got)
}

// A by-root response whose sidecars belong to a different block root must be
// rejected even if it verifies structurally.
func TestBlobReactor_ByRootWrongRootIsFailure(t *testing.T) {
	slot := math.Slot(77)
	sidecars := createTestSidecars(t, slot, 1)

	servingStore := newStubBlobStore()
	servingStore.setBlobs(t, slot.Unwrap(), sidecars)

	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), servingStore})
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	wrongRoot := common.Root{0xde, 0xad, 0xbe, 0xef}
	got, err := reactors[0].RequestSidecarsByRoot(t.Context(), slot, wrongRoot, requireCount(1))
	require.Error(t, err)
	require.Nil(t, got)
}

// The happy path with peer retry: the first peer misses, the second serves the full set, which comes back in
// index order.
func TestBlobReactor_ByRootPeerRetry(t *testing.T) {
	slot := math.Slot(123)
	sidecars := createTestSidecars(t, slot, 2)

	validStore := newStubBlobStore()
	validStore.setBlobs(t, slot.Unwrap(), sidecars)

	reactors, switches := makeConnectedReactors(
		t, 3, []*stubBlobStore{newStubBlobStore(), newStubBlobStore(), validStore})
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	got, err := reactors[0].RequestSidecarsByRoot(t.Context(), slot, blockRootOf(sidecars), requireCount(2))
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, uint64(0), got[0].GetIndex())
	require.Equal(t, uint64(1), got[1].GetIndex())
}

func TestBlobReactor_ByRootTimeout(t *testing.T) {
	slot := math.Slot(789)
	sidecars := createTestSidecars(t, slot, 1)

	slowStore := newStubBlobStore()
	slowStore.delay = 500 * time.Millisecond
	slowStore.setBlobs(t, slot.Unwrap(), sidecars)

	p2pConfig := config.DefaultP2PConfig()
	p2pConfig.ListenAddress = "tcp://127.0.0.1:0"
	p2pConfig.RootDir = t.TempDir()

	stores := []*stubBlobStore{newStubBlobStore(), slowStore}
	reactors := make([]*blobreactor.BlobReactor, 2)
	for i := range 2 {
		reactors[i] = newTestReactor(t, stores[i], blobreactor.Config{
			RequestTimeout: 100 * time.Millisecond,
			FetchTimeout:   200 * time.Millisecond,
		})
	}
	initSwitch := func(i int, sw *p2p.Switch) *p2p.Switch {
		sw.AddReactor(blobreactor.ReactorName, reactors[i])
		return sw
	}
	switches := p2p.MakeConnectedSwitches(p2pConfig, 2, initSwitch, p2p.Connect2Switches)
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(slot + 10)
	}

	_, err := reactors[0].RequestSidecarsByRoot(t.Context(), slot, blockRootOf(sidecars), requireCount(1))
	require.Error(t, err)
}

func TestBlobReactor_NoPeers(t *testing.T) {
	reactor := newTestReactor(t, newStubBlobStore(), testConfig())
	_, err := reactor.RequestSidecarsByRoot(t.Context(), 1, common.Root{}, requireCount(1))
	require.ErrorIs(t, err, blobreactor.ErrNoPeersAvailable)
}

func TestBlobReactor_ByRangeBasic(t *testing.T) {
	var (
		slotA     = math.Slot(10)
		slotB     = math.Slot(12)
		sidecarsA = createTestSidecars(t, slotA, 2)
		sidecarsB = createTestSidecars(t, slotB, 1)
	)

	servingStore := newStubBlobStore()
	servingStore.setBlobs(t, slotA.Unwrap(), sidecarsA)
	servingStore.setBlobs(t, slotB.Unwrap(), sidecarsB)

	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), servingStore})
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(20)
	}

	expected := map[math.Slot]int{slotA: 2, slotB: 1}
	verify := func(slot math.Slot, sidecars datypes.BlobSidecars) error {
		want, ok := expected[slot]
		if !ok {
			return blobreactor.ErrSlotNotRequested
		}
		if len(sidecars) != want {
			return fmt.Errorf("slot %d: expected %d sidecars, got %d", slot.Unwrap(), want, len(sidecars))
		}
		return nil
	}

	verified, err := reactors[0].RequestSidecarsByRange(t.Context(), slotA, 5, verify)
	require.NoError(t, err)
	require.Len(t, verified, 2)
	require.Len(t, verified[slotA], 2)
	require.Len(t, verified[slotB], 1)
}

// Slots returned by the peer that the requester does not need are skipped
// without failing the response.
func TestBlobReactor_ByRangeSkipsUnrequestedSlots(t *testing.T) {
	var (
		slotWanted   = math.Slot(10)
		slotUnwanted = math.Slot(11)
	)

	servingStore := newStubBlobStore()
	servingStore.setBlobs(t, slotWanted.Unwrap(), createTestSidecars(t, slotWanted, 1))
	servingStore.setBlobs(t, slotUnwanted.Unwrap(), createTestSidecars(t, slotUnwanted, 1))

	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), servingStore})
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(20)
	}

	verify := func(slot math.Slot, sidecars datypes.BlobSidecars) error {
		if slot != slotWanted {
			return blobreactor.ErrSlotNotRequested
		}
		if len(sidecars) != 1 {
			return errors.New("unexpected count")
		}
		return nil
	}

	verified, err := reactors[0].RequestSidecarsByRange(t.Context(), slotWanted, 4, verify)
	require.NoError(t, err)
	require.Len(t, verified, 1)
	require.Contains(t, verified, slotWanted)
}

// A peer whose response contains an invalid slot poisons the whole response.
func TestBlobReactor_ByRangeBadSlotPoisonsResponse(t *testing.T) {
	slot := math.Slot(10)

	servingStore := newStubBlobStore()
	servingStore.setBlobs(t, slot.Unwrap(), createTestSidecars(t, slot, 1))

	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), servingStore})
	defer stopSwitches(switches)

	for _, r := range reactors {
		r.SetHeadSlot(20)
	}

	verify := func(math.Slot, datypes.BlobSidecars) error {
		return errors.New("verification failed")
	}

	_, err := reactors[0].RequestSidecarsByRange(t.Context(), slot, 2, verify)
	require.ErrorIs(t, err, blobreactor.ErrAllPeersFailed)
}

// BroadcastSidecars must make the sidecars available to the proposer's own
// push cache and deliver them to peers (which self-verify and cache them).
func TestBlobReactor_PushDelivery(t *testing.T) {
	slot := math.Slot(42)
	sidecars := createTestSidecars(t, slot, 2)
	root := blockRootOf(sidecars)

	reactors, switches := makeConnectedReactors(
		t, 3, []*stubBlobStore{newStubBlobStore(), newStubBlobStore(), newStubBlobStore()})
	defer stopSwitches(switches)

	raw, err := sidecars.MarshalSSZ()
	require.NoError(t, err)

	// Reactor 0 is the proposer.
	require.NoError(t, reactors[0].BroadcastSidecars(raw))

	// The proposer immediately finds its own sidecars.
	require.NotNil(t, reactors[0].GetPushedSidecars(root))

	// Peers receive, self-verify and cache the push.
	require.Eventually(t, func() bool {
		return reactors[1].GetPushedSidecars(root) != nil &&
			reactors[2].GetPushedSidecars(root) != nil
	}, 5*time.Second, 50*time.Millisecond)

	// And a validator that missed the push could fetch by root from anyone,
	// even though nothing was persisted to any store yet.
	got, err := reactors[1].RequestSidecarsByRoot(t.Context(), slot, root, requireCount(2))
	require.NoError(t, err)
	require.Len(t, got, 2)
}

// Broadcasting an empty sidecar list is a no-op.
func TestBlobReactor_PushEmptyIsNoop(t *testing.T) {
	reactor := newTestReactor(t, newStubBlobStore(), testConfig())
	empty := datypes.BlobSidecars{}
	raw, err := empty.MarshalSSZ()
	require.NoError(t, err)
	require.NoError(t, reactor.BroadcastSidecars(raw))
}

// A push whose claimed root does not match its sidecars must not be cached.
func TestBlobReactor_PushRootMismatchRejected(t *testing.T) {
	slot := math.Slot(5)
	sidecars := createTestSidecars(t, slot, 1)

	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), newStubBlobStore()})
	defer stopSwitches(switches)

	raw, err := sidecars.MarshalSSZ()
	require.NoError(t, err)

	// Send a push with a bogus root directly over the wire.
	push := &blobreactor.SidecarsPush{BlockRoot: common.Root{0xbb}, SidecarData: raw}
	pushBz, err := push.MarshalSSZ()
	require.NoError(t, err)

	sw := switches[0]
	peers := sw.Peers().Copy()
	require.Len(t, peers, 1)
	sent := peers[0].Send(p2p.Envelope{
		ChannelID: blobreactor.BlobChannel,
		Message:   blobreactor.NewBlobMessageForTest(blobreactor.MessageTypePush, pushBz),
	})
	require.True(t, sent)

	// The receiver must never cache it, neither under the bogus root nor the
	// real one.
	time.Sleep(500 * time.Millisecond)
	require.Nil(t, reactors[1].GetPushedSidecars(common.Root{0xbb}))
	require.Nil(t, reactors[1].GetPushedSidecars(blockRootOf(sidecars)))
}

// A received Have suppresses the full-payload push toward its sender: the announcing peer is recorded as
// holding the root and forwardPush skips it.
func TestBlobReactor_HaveSuppressesPush(t *testing.T) {
	slot := math.Slot(7)
	sidecars := createTestSidecars(t, slot, 1)
	root := blockRootOf(sidecars)

	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), newStubBlobStore()})
	defer stopSwitches(switches)

	// Reactor 1 announces it holds the root.
	have := &blobreactor.SidecarsHave{BlockRoot: root}
	haveBz, err := have.MarshalSSZ()
	require.NoError(t, err)
	peers := switches[1].Peers().Copy()
	require.Len(t, peers, 1)
	require.True(t, peers[0].Send(p2p.Envelope{
		ChannelID: blobreactor.BlobChannel,
		Message:   blobreactor.NewBlobMessageForTest(blobreactor.MessageTypeHave, haveBz),
	}))

	node1ID := switches[1].NodeInfo().ID()
	require.Eventually(t, func() bool {
		return reactors[0].PeerKnowsRoot(node1ID, root)
	}, 5*time.Second, 20*time.Millisecond)

	// Reactor 0 then broadcasts the payload; the announcing peer must not receive it.
	raw, err := sidecars.MarshalSSZ()
	require.NoError(t, err)
	require.NoError(t, reactors[0].BroadcastSidecars(raw))

	time.Sleep(500 * time.Millisecond)
	require.Nil(t, reactors[1].GetPushedSidecars(root))
}

// Accepting a push makes the receiver announce the root to its other peers, so their re-pushes toward it are
// suppressed without a payload exchange.
func TestBlobReactor_PushTriggersHaveAnnouncement(t *testing.T) {
	slot := math.Slot(9)
	sidecars := createTestSidecars(t, slot, 1)
	root := blockRootOf(sidecars)

	reactors, switches := makeConnectedReactors(
		t, 3, []*stubBlobStore{newStubBlobStore(), newStubBlobStore(), newStubBlobStore()})
	defer stopSwitches(switches)

	raw, err := sidecars.MarshalSSZ()
	require.NoError(t, err)
	push := &blobreactor.SidecarsPush{BlockRoot: root, SidecarData: raw}
	pushBz, err := push.MarshalSSZ()
	require.NoError(t, err)

	// Send the payload from node 0 to node 1 only, so node 2 can learn about it only from node 1's announcement.
	node1ID := switches[1].NodeInfo().ID()
	target := switches[0].Peers().Get(node1ID)
	require.NotNil(t, target)
	require.True(t, target.Send(p2p.Envelope{
		ChannelID: blobreactor.BlobChannel,
		Message:   blobreactor.NewBlobMessageForTest(blobreactor.MessageTypePush, pushBz),
	}))

	// Node 2 learns that node 1 holds the root without ever receiving the payload itself.
	require.Eventually(t, func() bool {
		return reactors[2].PeerKnowsRoot(node1ID, root)
	}, 5*time.Second, 20*time.Millisecond)
	require.Nil(t, reactors[2].GetPushedSidecars(root))
}

// Pushes are only cached for slots just ahead of the receiver's head; stale and far-future fabrications are
// ignored before any expensive verification.
func TestBlobReactor_PushSlotWindowEnforced(t *testing.T) {
	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), newStubBlobStore()})
	defer stopSwitches(switches)

	reactors[1].SetHeadSlot(100)

	broadcast := func(slot math.Slot, count int) common.Root {
		sidecars := createTestSidecars(t, slot, count)
		raw, err := sidecars.MarshalSSZ()
		require.NoError(t, err)
		require.NoError(t, reactors[0].BroadcastSidecars(raw))
		return blockRootOf(sidecars)
	}

	staleRoot := broadcast(100, 1)  // == head: already finalized here, useless
	futureRoot := broadcast(300, 1) // far beyond the tip window
	tipRoot := broadcast(101, 1)    // head+1: a real proposal's slot

	require.Eventually(t, func() bool {
		return reactors[1].GetPushedSidecars(tipRoot) != nil
	}, 5*time.Second, 50*time.Millisecond)
	require.Nil(t, reactors[1].GetPushedSidecars(staleRoot))
	require.Nil(t, reactors[1].GetPushedSidecars(futureRoot))
}

// A single peer can only occupy a bounded number of push-cache slots with sets that are not yet bound to a
// verified proposal, so junk pushes cannot flush the cache.
func TestBlobReactor_PushPerPeerUnverifiedCap(t *testing.T) {
	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), newStubBlobStore()})
	defer stopSwitches(switches)

	reactors[1].SetHeadSlot(100)

	roots := make([]common.Root, 3)
	for i := range roots {
		sidecars := createTestSidecars(t, math.Slot(101+i), 1)
		raw, err := sidecars.MarshalSSZ()
		require.NoError(t, err)
		require.NoError(t, reactors[0].BroadcastSidecars(raw))
		roots[i] = blockRootOf(sidecars)
	}

	// The third push evicts the same peer's oldest unverified entry.
	require.Eventually(t, func() bool {
		return reactors[1].GetPushedSidecars(roots[2]) != nil
	}, 5*time.Second, 50*time.Millisecond)
	require.NotNil(t, reactors[1].GetPushedSidecars(roots[1]))
	require.Nil(t, reactors[1].GetPushedSidecars(roots[0]))
}

// A permuted-but-valid push must be cached in index order: consumers build positional blob bundles from the
// cached set, and a permuted bundle poisons re-proposals.
func TestBlobReactor_PushStoredSorted(t *testing.T) {
	reactors, switches := makeConnectedReactors(t, 2, []*stubBlobStore{newStubBlobStore(), newStubBlobStore()})
	defer stopSwitches(switches)

	sidecars := createTestSidecars(t, 42, 2)
	root := blockRootOf(sidecars)
	permuted := datypes.BlobSidecars{sidecars[1], sidecars[0]}
	raw, err := permuted.MarshalSSZ()
	require.NoError(t, err)
	require.NoError(t, reactors[0].BroadcastSidecars(raw))

	require.Eventually(t, func() bool {
		return reactors[1].GetPushedSidecars(root) != nil
	}, 5*time.Second, 50*time.Millisecond)

	got := reactors[1].GetPushedSidecars(root)
	require.Equal(t, uint64(0), got[0].GetIndex())
	require.Equal(t, uint64(1), got[1].GetIndex())
}

// A cached push that later fails verification against the real proposal must be evictable, so it stops
// shadowing the honest data and is no longer served. Guards against push-cache poisoning where an attacker
// races a same-root copy to arrive first.
func TestBlobReactor_DiscardEvictsPoisonedPush(t *testing.T) {
	slot := math.Slot(42)
	sidecars := createTestSidecars(t, slot, 2)
	root := blockRootOf(sidecars)

	reactor := newTestReactor(t, newStubBlobStore(), testConfig())

	raw, err := sidecars.MarshalSSZ()
	require.NoError(t, err)
	require.NoError(t, reactor.BroadcastSidecars(raw))
	require.NotNil(t, reactor.GetPushedSidecars(root))

	reactor.DiscardPushedSidecars(root)
	require.Nil(t, reactor.GetPushedSidecars(root), "discarded push must no longer be served")

	// Discarding a root that is not cached is a no-op.
	reactor.DiscardPushedSidecars(common.Root{0xff})
}

// NotifySidecarsObtained must replace any existing (possibly poisoned) entry for the root with the verified
// data, not merely flag the stale entry as verified.
func TestBlobReactor_NotifyReplacesStaleEntry(t *testing.T) {
	slot := math.Slot(7)
	good := createTestSidecars(t, slot, 2)
	root := blockRootOf(good)

	reactor := newTestReactor(t, newStubBlobStore(), testConfig())

	// Seed the cache with a different set under the same root (stand-in for a poisoned push).
	stale := createTestSidecars(t, slot, 1)
	staleRaw, err := stale.MarshalSSZ()
	require.NoError(t, err)
	require.NoError(t, reactor.BroadcastSidecars(staleRaw))

	reactor.NotifySidecarsObtained(root, good)
	require.Len(t, reactor.GetPushedSidecars(root), 2, "the verified set must replace the stale one")
}

// The serve lane must hand out every queued by-root task before any by-range task, regardless of arrival
// order, since tip-critical requests may not queue behind bulk sync serving.
func TestBlobReactor_ServeLanePrefersByRoot(t *testing.T) {
	reactor := newTestReactor(t, newStubBlobStore(), testConfig())
	src := mock.NewPeer(nil)

	for range 3 {
		reactor.EnqueueByRangeForTest(src, &blobreactor.SidecarsByRangeRequest{})
		reactor.EnqueueByRootForTest(src, &blobreactor.SidecarsByRootRequest{})
	}

	served := make([]string, 0, 6)
	for range 6 {
		lane, ok := reactor.RunOneServeTaskForTest()
		require.True(t, ok)
		served = append(served, lane)
	}
	require.Equal(t,
		[]string{"by_root", "by_root", "by_root", "by_range", "by_range", "by_range"},
		served)
}

// A saturated lane drops the newest messages instead of blocking the peer's receive loop, and the drop is
// bounded exactly by the queue capacity.
func TestBlobReactor_QueueFullDropsNewest(t *testing.T) {
	reactor := newTestReactor(t, newStubBlobStore(), testConfig())
	src := mock.NewPeer(nil)

	// Workers are not running, so the queue fills deterministically.
	for range blobreactor.PushQueueCapacityForTest + 5 {
		reactor.EnqueuePushForTest(src, &blobreactor.SidecarsPush{})
	}
	require.Equal(t, blobreactor.PushQueueCapacityForTest, reactor.PushQueueLenForTest(),
		"overflow beyond the queue capacity must be dropped, not queued")
}

// blockingBlobStore blocks GetByIndex until released, to hold a serve handler in flight.
type blockingBlobStore struct {
	started chan struct{}
	release chan struct{}
}

func (b *blockingBlobStore) GetByIndex(uint64) ([][]byte, error) {
	close(b.started)
	<-b.release
	return nil, nil
}

// Stop must wait for in-flight handlers but not for queued-but-unstarted work.
func TestBlobReactor_StopWaitsForInFlight(t *testing.T) {
	store := &blockingBlobStore{started: make(chan struct{}), release: make(chan struct{})}
	reactor := newTestReactor(t, store, testConfig())
	require.NoError(t, reactor.Start())

	reactor.EnqueueByRootForTest(mock.NewPeer(nil), &blobreactor.SidecarsByRootRequest{})
	<-store.started

	stopDone := make(chan struct{})
	var stopErr error
	go func() {
		stopErr = reactor.Stop()
		close(stopDone)
	}()

	select {
	case <-stopDone:
		t.Fatal("Stop returned while a handler was still in flight")
	case <-time.After(100 * time.Millisecond):
	}

	close(store.release)
	select {
	case <-stopDone:
	case <-time.After(5 * time.Second):
		t.Fatal("Stop did not return after the in-flight handler finished")
	}
	require.NoError(t, stopErr)
}

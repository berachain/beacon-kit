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

package blockchain_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/blockchain/mocks"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/filedb"
	testutils "github.com/berachain/beacon-kit/testing/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// The fakes below stand in for the blob distribution components so the FinalizeSidecars decision tree can be
// exercised in isolation: which tier ran, in which order, and what ended up persisted or queued.

type fakeBlobProcessor struct {
	persistCalls int
	verifyErr    error
}

func (f *fakeBlobProcessor) ProcessSidecars(avs *dastore.Store, sidecars datypes.BlobSidecars) error {
	f.persistCalls++
	return avs.Persist(sidecars)
}

func (f *fakeBlobProcessor) VerifySidecars(
	context.Context, datypes.BlobSidecars, *ctypes.BeaconBlockHeader, eip4844.KZGCommitments[common.ExecutionHash],
) error {
	return f.verifyErr
}

type fakeBlobRequester struct {
	pushed    datypes.BlobSidecars
	byRoot    datypes.BlobSidecars
	byRootErr error
	// calls records tier invocations across the fakes ("by_root" here, "el" on the reconstructor).
	calls *[]string
}

func (f *fakeBlobRequester) GetPushedSidecars(common.Root) datypes.BlobSidecars       { return f.pushed }
func (f *fakeBlobRequester) DiscardPushedSidecars(common.Root)                        {}
func (f *fakeBlobRequester) NotifySidecarsObtained(common.Root, datypes.BlobSidecars) {}
func (f *fakeBlobRequester) SetHeadSlot(math.Slot)                                    {}
func (f *fakeBlobRequester) FetchTimeout() time.Duration                              { return time.Second }

func (f *fakeBlobRequester) RequestSidecarsByRoot(
	_ context.Context, _ math.Slot, _ common.Root, verify func(datypes.BlobSidecars) error,
) (datypes.BlobSidecars, error) {
	*f.calls = append(*f.calls, "by_root")
	if f.byRootErr != nil {
		return nil, f.byRootErr
	}
	if err := verify(f.byRoot); err != nil {
		return nil, err
	}
	return f.byRoot, nil
}

func (f *fakeBlobRequester) RequestSidecarsByRange(
	context.Context, math.Slot, uint64, func(math.Slot, datypes.BlobSidecars) error,
) (map[math.Slot]datypes.BlobSidecars, error) {
	panic("by-range must not be used by the tip path")
}

type fakeBlobReconstructor struct {
	sidecars datypes.BlobSidecars
	err      error
	calls    *[]string
}

func (f *fakeBlobReconstructor) ReconstructSidecars(
	context.Context, *ctypes.SignedBeaconBlock,
) (datypes.BlobSidecars, error) {
	*f.calls = append(*f.calls, "el")
	return f.sidecars, f.err
}

type fakeBlobFetcher struct {
	queued []*ctypes.SignedBeaconBlock
}

func (f *fakeBlobFetcher) Start(context.Context) {}
func (f *fakeBlobFetcher) Stop()                 {}
func (f *fakeBlobFetcher) SetHeadSlot(math.Slot) {}
func (f *fakeBlobFetcher) PendingRequests() int  { return len(f.queued) }
func (f *fakeBlobFetcher) QueueBlobRequest(blk *ctypes.SignedBeaconBlock) error {
	f.queued = append(f.queued, blk)
	return nil
}

// finalizeSidecarsEnv wires a blockchain Service whose blob components are all fakes except a real
// availability store, so persistence and IsDataAvailable behave exactly as in production.
type finalizeSidecarsEnv struct {
	svc           *blockchain.Service
	store         *dastore.Store
	processor     *fakeBlobProcessor
	requester     *fakeBlobRequester
	reconstructor *fakeBlobReconstructor
	fetcher       *fakeBlobFetcher
	calls         []string
	cs            chain.Spec
}

func newFinalizeSidecarsEnv(t *testing.T) *finalizeSidecarsEnv {
	t.Helper()
	cs, err := chain.NewSpec(spec.DevnetChainSpecData())
	require.NoError(t, err)
	require.Equal(t, int64(2), cs.BlobConsensusEnableHeight())

	logger := log.NewNopLogger()
	store := dastore.New(
		filedb.NewRangeDB(
			filedb.NewDB(filedb.WithRootDirectory(t.TempDir()),
				filedb.WithFileExtension("ssz"),
				filedb.WithDirectoryPermissions(os.ModePerm),
				filedb.WithLogger(logger),
			),
		),
		logger,
	)

	storageBackend := &mocks.StorageBackend{}
	storageBackend.On("AvailabilityStore").Return(store)

	env := &finalizeSidecarsEnv{store: store, cs: cs}
	env.processor = &fakeBlobProcessor{}
	env.requester = &fakeBlobRequester{calls: &env.calls}
	env.reconstructor = &fakeBlobReconstructor{calls: &env.calls}
	env.fetcher = &fakeBlobFetcher{}
	env.svc = blockchain.NewService(
		storageBackend,
		env.processor,
		env.requester,
		env.reconstructor,
		env.fetcher,
		nil, // deposit contract unused
		logger,
		cs,
		nil, // execution engine unused
		nil, // local builder unused
		nil, // state processor unused
		metrics.NewNoOpTelemetrySink(),
	)
	return env
}

// makeBlobBlock builds a signed block committing to one blob (slot 10 from the generator) and its matching
// sidecar set.
func makeBlobBlock(t *testing.T, cs chain.Spec) (*ctypes.SignedBeaconBlock, datypes.BlobSidecars) {
	t.Helper()
	blk := testutils.GenerateValidBeaconBlock(t, cs.ActiveForkVersionForTimestamp(1))
	signedBlk := &ctypes.SignedBeaconBlock{BeaconBlock: blk, Signature: crypto.BLSSignature{}}

	commitments := blk.GetBody().GetBlobKzgCommitments()
	require.Len(t, commitments, 1)

	header := ctypes.NewSignedBeaconBlockHeader(blk.GetHeader(), crypto.BLSSignature{})
	sidecars := datypes.BlobSidecars{
		datypes.BuildBlobSidecar(
			0, header, &eip4844.Blob{}, commitments[0], eip4844.KZGProof{},
			make([]common.Root, ctypes.KZGInclusionProofDepth),
		),
	}
	return signedBlk, sidecars
}

func testCtx(t *testing.T) sdk.Context {
	t.Helper()
	return sdk.Context{}.WithContext(t.Context())
}

// Outside the DA window nothing is enforced, fetched, or queued.
func TestFinalizeSidecars_OutsideDAWindowIsNoop(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, _ := makeBlobBlock(t, env.cs)

	outside := int64((env.cs.MinEpochsForBlobsSidecarsRequest().Unwrap() + 2) * env.cs.SlotsPerEpoch())
	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), outside, signedBlk, nil))
	require.Zero(t, env.processor.persistCalls)
	require.Empty(t, env.fetcher.queued)
	require.Empty(t, env.calls)
}

// Sidecars handed in by the caller (legacy tx, or verified during ProcessProposal) are persisted with a strict
// availability check and nothing is fetched or queued.
func TestFinalizeSidecars_ProvidedSidecarsArePersisted(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, sidecars := makeBlobBlock(t, env.cs)
	blk := signedBlk.GetBeaconBlock()

	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), int64(blk.GetSlot().Unwrap()), signedBlk, sidecars))
	require.Equal(t, 1, env.processor.persistCalls)
	require.True(t, env.store.IsDataAvailable(t.Context(), blk.GetSlot(), blk.GetBody()))
	require.Empty(t, env.fetcher.queued)
	require.Empty(t, env.calls)
}

// A block with no blob commitments requires nothing at all.
func TestFinalizeSidecars_NoCommitmentsIsNoop(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, _ := makeBlobBlock(t, env.cs)
	signedBlk.GetBeaconBlock().GetBody().SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash]{})

	slot := int64(signedBlk.GetBeaconBlock().GetSlot().Unwrap())
	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), slot, signedBlk, nil))
	require.Zero(t, env.processor.persistCalls)
	require.Empty(t, env.fetcher.queued)
	require.Empty(t, env.calls)
}

// Sidecars already in the store (crash-restart replay, or a prior fetch) short-circuit before any retrieval.
func TestFinalizeSidecars_AlreadyAvailableShortCircuits(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, sidecars := makeBlobBlock(t, env.cs)
	require.NoError(t, env.store.Persist(sidecars))

	slot := int64(signedBlk.GetBeaconBlock().GetSlot().Unwrap())
	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), slot, signedBlk, nil))
	require.Zero(t, env.processor.persistCalls)
	require.Empty(t, env.fetcher.queued)
	require.Empty(t, env.calls)
}

// At the tip, a missing set is fetched synchronously (here via the push cache) and persisted.
func TestFinalizeSidecars_TipFetchSuccess(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, sidecars := makeBlobBlock(t, env.cs)
	env.requester.pushed = sidecars
	blk := signedBlk.GetBeaconBlock()

	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), int64(blk.GetSlot().Unwrap()), signedBlk, nil))
	require.Equal(t, 1, env.processor.persistCalls)
	require.True(t, env.store.IsDataAvailable(t.Context(), blk.GetSlot(), blk.GetBody()))
	require.Empty(t, env.fetcher.queued)
}

// When every tier fails at the tip, the block still finalizes and a background fetch is queued. The tiers must
// have been tried in order: EL first, then by-root.
func TestFinalizeSidecars_TipFetchFailureQueues(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, _ := makeBlobBlock(t, env.cs)
	env.reconstructor.err = errors.New("blobs not in EL pool")
	env.requester.byRootErr = errors.New("all peers failed")

	slot := int64(signedBlk.GetBeaconBlock().GetSlot().Unwrap())
	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), slot, signedBlk, nil))
	require.Equal(t, []string{"el", "by_root"}, env.calls)
	require.Len(t, env.fetcher.queued, 1)
	require.Zero(t, env.processor.persistCalls)
}

// During block sync the tip fetch is skipped entirely; the request goes straight to the background queue.
func TestFinalizeSidecars_SyncingQueuesWithoutFetching(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, _ := makeBlobBlock(t, env.cs)

	syncingTo := int64(signedBlk.GetBeaconBlock().GetSlot().Unwrap()) + 100
	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), syncingTo, signedBlk, nil))
	require.Empty(t, env.calls, "no synchronous retrieval during sync")
	require.Len(t, env.fetcher.queued, 1)
}

// Below the enable height the legacy strictness is preserved: a blob block finalizing without its sidecars is
// an error, never a queued fetch.
func TestFinalizeSidecars_LegacyMissingSidecarsIsFatal(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, _ := makeBlobBlock(t, env.cs)
	signedBlk.GetBeaconBlock().Slot = 1 // below the devnet enable height (2)

	err := env.svc.FinalizeSidecars(testCtx(t), 1, signedBlk, nil)
	require.ErrorIs(t, err, blockchain.ErrDataNotAvailable)
	require.Empty(t, env.fetcher.queued)
	require.Empty(t, env.calls)
}

// Stored sidecars matching the block's commitments but bound to a different header (a different proposal at the
// same slot) must not short-circuit; the node falls through to retrieval and the stale set is replaced.
func TestFinalizeSidecars_StaleStoredSidecarsAreNotReused(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, sidecars := makeBlobBlock(t, env.cs)
	blk := signedBlk.GetBeaconBlock()
	commitments := blk.GetBody().GetBlobKzgCommitments()

	staleHdr := *blk.GetHeader()
	staleHdr.ProposerIndex++
	stale := datypes.BlobSidecars{
		datypes.BuildBlobSidecar(
			0, ctypes.NewSignedBeaconBlockHeader(&staleHdr, crypto.BLSSignature{}),
			&eip4844.Blob{}, commitments[0], eip4844.KZGProof{},
			make([]common.Root, ctypes.KZGInclusionProofDepth),
		),
	}
	require.NoError(t, env.store.Persist(stale))

	env.requester.pushed = sidecars
	require.NoError(t, env.svc.FinalizeSidecars(testCtx(t), int64(blk.GetSlot().Unwrap()), signedBlk, nil))
	require.Equal(t, 1, env.processor.persistCalls, "the stale set must be replaced, not reused")

	got, err := env.store.GetBlobSidecars(blk.GetSlot())
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, blk.GetHeader().HashTreeRoot(), got[0].GetBeaconBlockHeader().HashTreeRoot())
}

// Once blob consensus is enabled at the orphaned height, startup must keep sidecars persisted by an interrupted
// finalization: replayed blocks no longer carry them, so they may be the only local copy.
func TestPruneOrphanedBlobs_KeepsSidecarsWhenBlobConsensusEnabled(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	signedBlk, sidecars := makeBlobBlock(t, env.cs)
	require.NoError(t, env.store.Persist(sidecars))

	// Crash between FinalizeBlock and Commit: comet's last committed height is one below the persisted slot.
	slot := signedBlk.GetBeaconBlock().GetSlot()
	require.NoError(t, env.svc.PruneOrphanedBlobs(int64(slot.Unwrap())-1))

	got, err := env.store.GetBlobSidecars(slot)
	require.NoError(t, err)
	require.NotEmpty(t, got, "the only local copy must survive startup")
}

// Below the enable height blobs still replay from the consensus tx, so the legacy orphan cleanup is preserved.
func TestPruneOrphanedBlobs_LegacyDeletesOrphanedSlot(t *testing.T) {
	t.Parallel()
	env := newFinalizeSidecarsEnv(t)
	_, sidecars := makeBlobBlock(t, env.cs)
	for _, sc := range sidecars {
		sc.GetBeaconBlockHeader().SetSlot(1) // below the devnet enable height (2)
	}
	require.NoError(t, env.store.Persist(sidecars))

	require.NoError(t, env.svc.PruneOrphanedBlobs(0))

	got, err := env.store.GetBlobSidecars(1)
	require.NoError(t, err)
	require.Empty(t, got)
}

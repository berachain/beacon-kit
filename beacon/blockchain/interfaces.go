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

package blockchain

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/blobconsensus"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/delay"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/block"
	"github.com/berachain/beacon-kit/storage/deposit"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine interface {
	// NotifyNewPayload notifies the execution client of new payload.
	NotifyNewPayload(
		ctx context.Context,
		req ctypes.NewPayloadRequest,
		retryOnSyncingStatus bool,
	) error
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *ctypes.ForkchoiceUpdateRequest,
	) (*engineprimitives.PayloadID, error)
}

// LocalBuilder is the interface for the builder service.
type LocalBuilder interface {
	// Enabled returns true if the local builder is enabled.
	Enabled() bool
	// RequestPayloadAsync requests a new payload for the given slot.
	RequestPayloadAsync(
		ctx context.Context,
		r *builder.RequestPayloadData,
	) (*engineprimitives.PayloadID, common.Version, error)
	CacheLatestVerifiedPayload(
		latestEnvelopeSlot math.Slot,
		latestEnvelope ctypes.BuiltExecutionPayloadEnv,
	)
}

// StateProcessor defines the interface for processing various state transitions
// in the beacon chain.
type StateProcessor interface {
	// InitializeBeaconStateFromEth1 initializes the premined beacon
	// state from the eth1 deposits.
	InitializeBeaconStateFromEth1(
		*statedb.StateDB,
		ctypes.Deposits,
		*ctypes.ExecutionPayloadHeader,
		common.Version,
	) (transition.ValidatorUpdates, error)
	// ProcessFork prepares the state for the fork version at the given timestamp.
	ProcessFork(
		st *statedb.StateDB, timestamp math.U64, logUpgrade bool,
	) error
	// ProcessSlots processes the state transition for a range of slots.
	ProcessSlots(
		*statedb.StateDB, math.Slot,
	) (transition.ValidatorUpdates, error)
	// Transition processes the state transition for a given block.
	Transition(
		core.ReadOnlyContext,
		*statedb.StateDB,
		*ctypes.BeaconBlock,
	) (transition.ValidatorUpdates, error)
	GetSignatureVerifierFn(*statedb.StateDB) (
		func(
			blk *ctypes.BeaconBlock,
			signature crypto.BLSSignature) error,
		error,
	)
}

// StorageBackend defines an interface for accessing various storage components
// required by the beacon node.
type StorageBackend interface {
	// AvailabilityStore returns the availability store for the given context.
	AvailabilityStore() *dastore.Store
	// StateFromContext retrieves the beacon state from the given context.
	StateFromContext(context.Context) *statedb.StateDB
	// DepositStore retrieves the deposit store.
	DepositStore() deposit.StoreManager
	// BlockStore retrieves the block store.
	BlockStore() *block.KVStore[*ctypes.BeaconBlock]
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// IncrementCounter increments the counter identified by
	// the provided key.
	IncrementCounter(key string, args ...string)

	// SetGauge sets a gauge metric to the specified value.
	SetGauge(key string, value int64, args ...string)

	// MeasureSince measures the time since the provided start time,
	// identified by the provided keys.
	MeasureSince(key string, start time.Time, args ...string)
}

//nolint:revive // its ok
type BlockchainI interface {
	ProcessGenesisData(
		context.Context,
		[]byte,
	) (transition.ValidatorUpdates, error)
	ParseBeaconBlock(req encoding.ABCIRequest) (
		*ctypes.SignedBeaconBlock,
		datypes.BlobSidecars,
		error,
	)
	// ProcessProposal verifies a proposal and also returns the block's verified blob sidecars, so the caller
	// can cache them for FinalizeBlock.
	ProcessProposal(
		sdk.Context,
		*cmtabci.ProcessProposalRequest,
		[]byte, // this node address
	) (transition.ValidatorUpdates, datypes.BlobSidecars, error)
	FinalizeSidecars(
		ctx sdk.Context,
		syncingToHeight int64,
		signedBlk *ctypes.SignedBeaconBlock,
		blobs datypes.BlobSidecars,
	) error
	FinalizeBlock(
		sdk.Context,
		*cmtabci.FinalizeBlockRequest,
	) (transition.ValidatorUpdates, error)
	PostFinalizeBlockOps(
		sdk.Context,
		*ctypes.BeaconBlock,
	) error
	PruneOrphanedBlobs(lastBlockHeight int64) error
	// PendingBlobRequests reports queued background blob fetches; a node with pending in-window fetches must not report itself as synced.
	PendingBlobRequests() int
}

// BlobRequester is the blob reactor surface the blockchain service uses to obtain sidecars that no longer ride inside CometBFT
// blocks. Satisfied by *blobreactor.BlobReactor.
type BlobRequester interface {
	// GetPushedSidecars returns sidecars pushed for the given block root, or nil.
	GetPushedSidecars(root common.Root) datypes.BlobSidecars
	// DiscardPushedSidecars evicts a cached push that failed verification against the actual proposal, so it
	// no longer shadows an honest re-push or gets served to peers.
	DiscardPushedSidecars(root common.Root)
	// NotifySidecarsObtained caches fully verified sidecars for serving and re-gossips them to peers that are not known to have them.
	NotifySidecarsObtained(root common.Root, sidecars datypes.BlobSidecars)
	// RequestSidecarsByRoot fetches one block's verified sidecars from peers.
	RequestSidecarsByRoot(
		ctx context.Context,
		slot math.Slot,
		blockRoot common.Root,
		verify func(datypes.BlobSidecars) error,
	) (datypes.BlobSidecars, error)
	// RequestSidecarsByRange fetches verified sidecars for a slot range from peers; the verify callback is invoked once per returned slot.
	RequestSidecarsByRange(
		ctx context.Context,
		startSlot math.Slot,
		count uint64,
		verify func(math.Slot, datypes.BlobSidecars) error,
	) (map[math.Slot]datypes.BlobSidecars, error)
	// SetHeadSlot updates the requester's view of the chain head.
	SetHeadSlot(slot math.Slot)
	// FetchTimeout is the overall deadline for retrieving one block's sidecars at the tip.
	FetchTimeout() time.Duration
}

// BlobReconstructor rebuilds a block's sidecars from blobs fetched off the local execution client. Satisfied by *blob.Reconstructor.
type BlobReconstructor interface {
	ReconstructSidecars(ctx context.Context, signedBlk *ctypes.SignedBeaconBlock) (datypes.BlobSidecars, error)
}

// BlobFetcher asynchronously fetches blob sidecars of finalized blocks.
type BlobFetcher interface {
	// Start begins the background blob fetching process.
	Start(ctx context.Context)
	// Stop gracefully shuts down the blob fetcher.
	Stop()
	// QueueBlobRequest queues a fetch of the given block's sidecars.
	QueueBlobRequest(signedBlk *ctypes.SignedBeaconBlock) error
	// SetHeadSlot updates the head slot for DA-window expiry and serving.
	SetHeadSlot(slot math.Slot)
	// PendingRequests returns the number of queued fetch requests.
	PendingRequests() int
}

// BlobProcessor is the interface for the blobs processor.
type BlobProcessor interface {
	// ProcessSidecars processes the blobs and ensures they match the local
	// state.
	ProcessSidecars(
		avs *dastore.Store,
		sidecars datypes.BlobSidecars,
	) error
	// VerifySidecars verifies the blobs and ensures they match the local state.
	VerifySidecars(
		ctx context.Context,
		sidecars datypes.BlobSidecars,
		blkHeader *ctypes.BeaconBlockHeader,
		kzgCommitments eip4844.KZGCommitments[common.ExecutionHash],
	) error
}

type PruningChainSpec interface {
	MinEpochsForBlobsSidecarsRequest() math.Epoch
	SlotsPerEpoch() uint64
}

type ServiceChainSpec interface {
	PruningChainSpec
	chain.BlobSpec
	chain.ForkSpec
	chain.ForkVersionSpec
	delay.ConfigGetter
	blobconsensus.ConfigGetter

	EpochsPerHistoricalVector() uint64
	SlotToEpoch(slot math.Slot) math.Epoch
	Eth1FollowDistance() uint64
	MaxDepositsPerBlock() uint64
}

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

package blobreactor

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/p2p"
	"golang.org/x/time/rate"
)

const (
	// BlobChannel is the custom channel ID for blob distribution.
	BlobChannel = byte(0x70)

	// ReactorName is the registered name for the blob reactor in CometBFT's switch.
	ReactorName = "BLOBREACTOR"

	defaultPriority           = 5
	defaultSendQueueCapacity  = 20
	defaultRecvBufferCapacity = 4 << 20
	// defaultRecvMessageCapacity bounds a single inbound message; must fit a full by-range response (responseByteBudget) plus envelope
	// overhead.
	defaultRecvMessageCapacity = 4 << 20

	// responseByteBudget bounds the total sidecar bytes in one response.
	responseByteBudget = 2 << 20

	// gossipWorkerCount is the maximum number of inbound pushes verified concurrently. Each verification is
	// a CPU-bound KZG batch check taking several milliseconds, so this caps how many cores inbound gossip
	// can occupy. Honest traffic is about one distinct push per block, everything above that is burst
	// headroom.
	gossipWorkerCount = 4

	// serveWorkerCount is the maximum number of peer requests answered concurrently. Serving is disk-bound
	// blob store reading, so this caps concurrent reads and bounds response-assembly memory at roughly
	// serveWorkerCount times responseByteBudget.
	serveWorkerCount = 6

	// The lane queues are small burst buffers in front of the workers. Overflow is dropped (see enqueue).
	pushQueueCapacity    = 16
	byRootQueueCapacity  = 16
	byRangeQueueCapacity = 8

	// knownRootsPerPeer bounds the per-peer memory of roots the peer is known to have.
	knownRootsPerPeer = 64

	// pushCacheSize bounds the number of recent blocks whose pushed sidecars are held in memory for the tip of the chain.
	pushCacheSize = 16

	// Per-peer inbound rate limits. Requests cover by-root and by-range; gossip covers pushes and haves.
	requestsPerSecond = 10
	requestsBurst     = 20
	gossipPerSecond   = 8
	gossipBurst       = 16

	// Peer scoring. Scores order peer selection for fetches; they are not a ban mechanism.
	scoreSuccess     = 1
	scoreFailure     = -1
	scoreJunk        = -5
	scoreMin         = -100
	scoreMax         = 100
	scoreSpoofedResp = -20
)

// peerState tracks everything the reactor knows about one connected peer.
type peerState struct {
	peerID p2p.ID
	// score orders peer selection; junk responses drive it down.
	score int
	// knownRoots tracks block roots this peer is known to hold (it pushed or announced them, or we pushed to it).
	knownRoots *rootSet
	// reqLimiter rate-limits inbound by-root/by-range requests.
	reqLimiter *rate.Limiter
	// gossipLimiter rate-limits inbound pushes and haves.
	gossipLimiter *rate.Limiter
}

// rootSet is a fixed-capacity set of roots with FIFO eviction.
type rootSet struct {
	order []common.Root
	set   map[common.Root]struct{}
	next  int
}

func newRootSet(capacity int) *rootSet {
	return &rootSet{
		order: make([]common.Root, capacity),
		set:   make(map[common.Root]struct{}, capacity),
	}
}

func (rs *rootSet) add(root common.Root) {
	if _, ok := rs.set[root]; ok {
		return
	}
	delete(rs.set, rs.order[rs.next])
	rs.order[rs.next] = root
	rs.set[root] = struct{}{}
	rs.next = (rs.next + 1) % len(rs.order)
}

func (rs *rootSet) has(root common.Root) bool {
	_, ok := rs.set[root]
	return ok
}

// pendingRequest correlates an in-flight request with the peer it was sent to.
type pendingRequest struct {
	peerID p2p.ID
	ch     chan *SidecarsResponse
}

// BlobReactor handles p2p blob sidecar distribution for BeaconKit. It implements the CometBFT p2p.Reactor
// interface, and its channel rides the p2p connections CometBFT already maintains.
//
// At the tip of the chain, delivery is a three-lane system, fastest lane first:
//
//  1. Push (the normal case). When the proposer hands its block to CometBFT it simultaneously pushes the
//     sidecars to its directly connected peers (BroadcastSidecars), and every node re-forwards them once after
//     verifying them against the proposal. The sidecars flood the network in parallel with the block proposal,
//     so validators usually have them in their local push cache by the time they vote.
//
//  2. Own execution client (first fallback). A validator that missed the push asks its own EL for the blobs
//     (engine_getBlobsV2) and rebuilds the sidecars locally, usually a hit since blob transactions travel the
//     EL mempool with their blobs attached.
//
//  3. Ask a peer (last resort). A by-root request (RequestSidecarsByRoot), preferring peers that announced
//     having the data; the proposer serves these straight from its push cache.
//
// Whichever lane delivered, the node re-pushes to peers that still lack the sidecars (NotifySidecarsObtained),
// so the network converges even if the proposer's own pushes reached only a few nodes. Nodes that hold a set
// also announce it cheaply (Have), so peers about to re-push the full payload skip them and flood duplication
// stays low. A validator that gets nothing in time votes against the proposal, and the round retries with
// larger timeouts. Nodes catching up after downtime use a fourth bulk lane, by-range requests for whole
// slot windows (RequestSidecarsByRange), served from the availability store every node keeps for the DA
// window.
//
// The reactor itself is just the transport. It moves sidecars between peers, caches recent pushes, and rate
// limits and scores peers. Deciding which lane to try, and when, is the caller's job in beacon/blockchain.
//
// Inbound messages are never handled on CometBFT's per-peer receive loop, since blocking it would stall all
// traffic from that peer, consensus included. Receive only decodes and rate-limits, then queues each message
// as a plain-data task (pushTask, byRootTask, byRangeTask) for one of two worker lanes started in OnStart.
// The gossip lane verifies pushes, which is CPU-bound work, and the serve lane answers requests, which is
// disk-bound work, so one kind of load cannot starve the other. Queues are bounded and overflow is dropped
// (see enqueue).
//
// Peers are never trusted. Sidecars carry inclusion proofs against their header and KZG proofs for their
// blobs, so every set can be verified on arrival, whether it came in as a push or as a response, and a peer
// that sends junk simply loses score while the next one is tried.
type BlobReactor struct {
	service.BaseService
	sw *p2p.Switch

	blobStore     BlobStore
	verifier      SidecarVerifier
	logger        log.Logger
	config        Config
	metrics       *blobReactorMetrics
	maxBlobsPerTx uint64

	// stateMu protects peers and headSlot.
	stateMu  sync.RWMutex
	peers    map[p2p.ID]*peerState
	headSlot math.Slot

	// pushMu protects the push cache of recent tip sidecars.
	pushMu    sync.RWMutex
	pushCache *pushCache

	// pendingMu protects in-flight request correlation state.
	pendingMu sync.Mutex
	pending   map[uint64]pendingRequest

	// Inbound work scheduling. Fixed workers started in OnStart drain the bounded lane queues (see the lane
	// constants). Closing quit stops the workers, and workersWg waits for in-flight handlers only. Queued but
	// unstarted tasks are abandoned.
	pushQueue    chan pushTask
	byRootQueue  chan byRootTask
	byRangeQueue chan byRangeTask
	quit         chan struct{}
	workersWg    sync.WaitGroup

	// stopped is set when OnStop begins. It short-circuits Receive and enqueueing.
	stopped atomic.Bool
}

// NewBlobReactor creates a new blob reactor.
func NewBlobReactor(
	blobStore BlobStore,
	verifier SidecarVerifier,
	logger log.Logger,
	cfg Config,
	maxBlobsPerBlock uint64,
	sink TelemetrySink,
) *BlobReactor {
	br := &BlobReactor{
		blobStore:     blobStore,
		verifier:      verifier,
		logger:        logger,
		config:        cfg,
		metrics:       newBlobReactorMetrics(sink),
		maxBlobsPerTx: maxBlobsPerBlock,
		peers:         make(map[p2p.ID]*peerState),
		pushCache:     newPushCache(pushCacheSize),
		pending:       make(map[uint64]pendingRequest),
		pushQueue:     make(chan pushTask, pushQueueCapacity),
		byRootQueue:   make(chan byRootTask, byRootQueueCapacity),
		byRangeQueue:  make(chan byRangeTask, byRangeQueueCapacity),
		quit:          make(chan struct{}),
	}
	br.BaseService = *service.NewBaseService(nil, ReactorName, br)
	return br
}

// SetSwitch allows setting a switch. Part of p2p.Reactor.
func (br *BlobReactor) SetSwitch(sw *p2p.Switch) {
	br.sw = sw
}

// GetChannels returns the channel descriptors. Part of p2p.Reactor.
func (br *BlobReactor) GetChannels() []*p2p.ChannelDescriptor {
	return []*p2p.ChannelDescriptor{
		{
			ID:                  BlobChannel,
			Priority:            defaultPriority,
			SendQueueCapacity:   defaultSendQueueCapacity,
			RecvBufferCapacity:  defaultRecvBufferCapacity,
			RecvMessageCapacity: defaultRecvMessageCapacity,
			MessageType:         &BlobMessage{},
		},
	}
}

// InitPeer is called by the switch before the peer is started. Part of p2p.Reactor.
func (br *BlobReactor) InitPeer(peer p2p.Peer) p2p.Peer {
	br.AddPeer(peer)
	return peer
}

// AddPeer is called by the switch after the peer is added and successfully started. Part of p2p.Reactor.
func (br *BlobReactor) AddPeer(peer p2p.Peer) {
	br.stateMu.Lock()
	if _, ok := br.peers[peer.ID()]; !ok {
		br.peers[peer.ID()] = &peerState{
			peerID:        peer.ID(),
			knownRoots:    newRootSet(knownRootsPerPeer),
			reqLimiter:    rate.NewLimiter(requestsPerSecond, requestsBurst),
			gossipLimiter: rate.NewLimiter(gossipPerSecond, gossipBurst),
		}
	}
	peerCount := len(br.peers)
	br.stateMu.Unlock()

	br.metrics.setPeerCount(peerCount)
	br.logger.Info("BlobReactor added peer", "peer", peer.ID())
}

// RemovePeer is called by the switch when the peer is stopped. Part of p2p.Reactor.
func (br *BlobReactor) RemovePeer(peer p2p.Peer, reason any) {
	br.stateMu.Lock()
	delete(br.peers, peer.ID())
	peerCount := len(br.peers)
	br.stateMu.Unlock()

	br.metrics.setPeerCount(peerCount)
	br.logger.Info("BlobReactor removed peer", "peer", peer.ID(), "reason", reason)
}

// SetHeadSlot updates the reactor's view of the current chain head. It bounds by-range serving and the push acceptance window.
func (br *BlobReactor) SetHeadSlot(slot math.Slot) {
	br.stateMu.Lock()
	br.headSlot = slot
	br.stateMu.Unlock()
}

// FetchTimeout returns the configured overall deadline for retrieving one block's sidecars at the tip of the chain.
func (br *BlobReactor) FetchTimeout() time.Duration {
	return br.config.FetchTimeout
}

func (br *BlobReactor) getHeadSlot() math.Slot {
	br.stateMu.RLock()
	defer br.stateMu.RUnlock()
	return br.headSlot
}

// adjustScore moves a peer's selection score by delta, clamped.
func (br *BlobReactor) adjustScore(peerID p2p.ID, delta int) {
	br.stateMu.Lock()
	defer br.stateMu.Unlock()
	ps, ok := br.peers[peerID]
	if !ok {
		return
	}
	ps.score += delta
	if ps.score > scoreMax {
		ps.score = scoreMax
	}
	if ps.score < scoreMin {
		ps.score = scoreMin
	}
}

// markPeerHasRoot records that the given peer is known to hold the sidecars for root (it pushed or announced them, or we pushed to it).
func (br *BlobReactor) markPeerHasRoot(peerID p2p.ID, root common.Root) {
	br.stateMu.Lock()
	defer br.stateMu.Unlock()
	if ps, ok := br.peers[peerID]; ok {
		ps.knownRoots.add(root)
	}
}

// pushTask is a queued inbound push, sidecars a peer sent us unsolicited. The gossip lane verifies them
// (KZG, CPU-bound) via handlePush before caching them.
type pushTask struct {
	src  p2p.Peer
	push *SidecarsPush
}

// byRootTask is a queued request for one block's sidecars, typically from a validator that needs them to
// vote within the current consensus round. The serve lane always answers these before any byRangeTask.
type byRootTask struct {
	src p2p.Peer
	req *SidecarsByRootRequest
}

// byRangeTask is a queued bulk request for a slot range from a peer catching up after downtime. The serve
// lane answers these with whatever capacity by-root traffic leaves free, since the sender retries on a slow
// cadence anyway.
type byRangeTask struct {
	src p2p.Peer
	req *SidecarsByRangeRequest
}

// enqueue puts a task on a lane queue, or drops it when the queue is full or the reactor is stopping.
// Dropping is safe because no message is a single point of delivery. A dropped push is recoverable through
// by-root or the EL, and a dropped request is retried by its sender against another peer. Dropping the
// newest rather than evicting the oldest is deliberate, a flood can only ever shed its own tail and never
// evict an already-queued honest message.
func enqueue[T any](br *BlobReactor, queue chan T, task T, peerID p2p.ID, taskType string) {
	if br.stopped.Load() {
		return
	}
	select {
	case queue <- task:
	default:
		br.logger.Warn("BlobReactor queue full, dropping message", "peer", peerID, "task_type", taskType)
		br.metrics.observeQueueFull(taskType)
	}
}

// gossipWorker drains push handling, the CPU-bound lane (KZG self-verification of pushed sidecars).
func (br *BlobReactor) gossipWorker() {
	defer br.workersWg.Done()
	for {
		// Exit promptly once stopping, even with work still queued.
		select {
		case <-br.quit:
			return
		default:
		}
		select {
		case task := <-br.pushQueue:
			br.handlePush(task.src, task.push)
		case <-br.quit:
			return
		}
	}
}

// serveWorker drains request serving, the disk-bound lane, via runOneServeTask.
func (br *BlobReactor) serveWorker() {
	defer br.workersWg.Done()
	for {
		if _, ok := br.runOneServeTask(); !ok {
			return
		}
	}
}

// runOneServeTask picks and runs the next serve-lane task, preferring tip-critical by-root requests over
// bulk by-range ones. It reports which lane it served, or ok=false once the reactor is stopping. The
// priority holds at task boundaries only. An in-flight by-range serve is never preempted, so worst-case
// by-root queueing latency is one serve duration. Sustained by-root load can starve by-range entirely,
// which is the intended failure mode, since syncers retry on their own cadence while tip traffic must fit
// a consensus round.
func (br *BlobReactor) runOneServeTask() (string, bool) {
	select {
	case <-br.quit:
		return "", false
	default:
	}
	select {
	case task := <-br.byRootQueue:
		br.handleByRootRequest(task.src, task.req)
		return "by_root", true
	default:
	}
	select {
	case task := <-br.byRootQueue:
		br.handleByRootRequest(task.src, task.req)
		return "by_root", true
	case task := <-br.byRangeQueue:
		br.handleByRangeRequest(task.src, task.req)
		return "by_range", true
	case <-br.quit:
		return "", false
	}
}

// Receive is called by the switch for every envelope on the blob channel. Part of p2p.Reactor.
//
//nolint:gocognit,funlen // flat message dispatch
func (br *BlobReactor) Receive(envelope p2p.Envelope) {
	if br.stopped.Load() {
		return
	}

	blobMsg, ok := envelope.Message.(*BlobMessage)
	if !ok {
		br.logger.Error("BlobReactor received non-blob message",
			"peer", envelope.Src.ID(), "type", fmt.Sprintf("%T", envelope.Message))
		return
	}
	if len(blobMsg.Data) < 1 {
		br.adjustScore(envelope.Src.ID(), scoreJunk)
		return
	}

	msgType := MessageType(blobMsg.Data[0])
	msgData := blobMsg.Data[1:]
	src := envelope.Src
	br.metrics.observeMessageReceived(msgType.String())

	switch msgType {
	case MessageTypePush:
		// Rate-limit before decoding. A push is up to ~832 KiB and the limiter needs only the peer ID, so an
		// over-limit peer must not cost an unmarshal.
		if !br.allowGossip(src.ID()) {
			br.metrics.observeRateLimited("push")
			return
		}
		var push SidecarsPush
		if err := push.UnmarshalSSZ(msgData); err != nil {
			br.logger.Warn("Failed to unmarshal sidecars push", "error", err, "peer", src.ID())
			br.adjustScore(src.ID(), scoreJunk)
			return
		}
		enqueue(br, br.pushQueue, pushTask{src: src, push: &push}, src.ID(), "push")

	case MessageTypeHave:
		var have SidecarsHave
		if err := have.UnmarshalSSZ(msgData); err != nil {
			br.adjustScore(src.ID(), scoreJunk)
			return
		}
		if !br.allowGossip(src.ID()) {
			br.metrics.observeRateLimited("have")
			return
		}
		// Cheap enough to handle inline.
		br.markPeerHasRoot(src.ID(), have.BlockRoot)

	case MessageTypeByRootRequest:
		var req SidecarsByRootRequest
		if err := req.UnmarshalSSZ(msgData); err != nil {
			br.adjustScore(src.ID(), scoreJunk)
			return
		}
		if !br.allowRequest(src.ID()) {
			br.metrics.observeRateLimited("by_root_request")
			return
		}
		enqueue(br, br.byRootQueue, byRootTask{src: src, req: &req}, src.ID(), "by_root_request")

	case MessageTypeByRangeRequest:
		var req SidecarsByRangeRequest
		if err := req.UnmarshalSSZ(msgData); err != nil {
			br.adjustScore(src.ID(), scoreJunk)
			return
		}
		if !br.allowRequest(src.ID()) {
			br.metrics.observeRateLimited("by_range_request")
			return
		}
		enqueue(br, br.byRangeQueue, byRangeTask{src: src, req: &req}, src.ID(), "by_range_request")

	case MessageTypeResponse:
		// Only decode responses from peers we have an in-flight request to; anything else is unsolicited and
		// not worth the decode. No score penalty: a response arriving after our timeout deregistered the
		// request is indistinguishable from junk, and punishing it would hurt honest-but-slow peers.
		if !br.hasPendingFrom(src.ID()) {
			br.metrics.observeUnsolicitedResponse()
			return
		}
		var resp SidecarsResponse
		if err := resp.UnmarshalSSZ(msgData); err != nil {
			br.logger.Warn("Failed to unmarshal sidecars response", "error", err, "peer", src.ID())
			br.adjustScore(src.ID(), scoreJunk)
			return
		}
		// Correlation is cheap; the waiting requester does the heavy lifting.
		br.handleResponse(src, &resp)

	default:
		br.logger.Warn("BlobReactor received unknown message type", "type", uint8(msgType), "peer", src.ID())
		br.adjustScore(src.ID(), scoreJunk)
	}
}

func (br *BlobReactor) allowRequest(peerID p2p.ID) bool {
	br.stateMu.RLock()
	ps := br.peers[peerID]
	br.stateMu.RUnlock()
	if ps == nil {
		return false
	}
	return ps.reqLimiter.Allow()
}

func (br *BlobReactor) allowGossip(peerID p2p.ID) bool {
	br.stateMu.RLock()
	ps := br.peers[peerID]
	br.stateMu.RUnlock()
	if ps == nil {
		return false
	}
	return ps.gossipLimiter.Allow()
}

// sendToPeer marshals and sends a typed message to a peer, blocking until it is queued (or the peer's send queue rejects it).
func (br *BlobReactor) sendToPeer(peer p2p.Peer, msgType MessageType, msg sszMarshaler) bool {
	bz, err := msg.MarshalSSZ()
	if err != nil {
		br.logger.Error("Failed to marshal blob message", "type", msgType.String(), "error", err)
		return false
	}
	return peer.Send(p2p.Envelope{ChannelID: BlobChannel, Message: newBlobMessage(msgType, bz)})
}

// trySendToPeer is like sendToPeer but never blocks; used for gossip and request responses, where dropping under backpressure is
// acceptable.
func (br *BlobReactor) trySendToPeer(peer p2p.Peer, msgType MessageType, msg sszMarshaler) bool {
	bz, err := msg.MarshalSSZ()
	if err != nil {
		br.logger.Error("Failed to marshal blob message", "type", msgType.String(), "error", err)
		return false
	}
	return peer.TrySend(p2p.Envelope{ChannelID: BlobChannel, Message: newBlobMessage(msgType, bz)})
}

type sszMarshaler interface {
	MarshalSSZ() ([]byte, error)
}

// getPeer returns the live switch peer for an ID, or nil.
func (br *BlobReactor) getPeer(peerID p2p.ID) p2p.Peer {
	if br.sw == nil {
		return nil
	}
	peer := br.sw.Peers().Get(peerID)
	if peer == nil || !peer.IsRunning() {
		return nil
	}
	return peer
}

// OnStart implements service.Service. It launches the worker lanes. The switch starts the reactor before any
// peer can deliver messages, so nothing is enqueued before workers exist.
func (br *BlobReactor) OnStart() error {
	br.logger.Info("Starting BlobReactor")
	br.workersWg.Add(gossipWorkerCount + serveWorkerCount)
	for range gossipWorkerCount {
		go br.gossipWorker()
	}
	for range serveWorkerCount {
		go br.serveWorker()
	}
	return nil
}

// OnStop implements service.Service. It waits for in-flight handlers. Queued but unstarted tasks are dropped.
func (br *BlobReactor) OnStop() {
	br.logger.Info("Stopping BlobReactor")
	br.stopped.Store(true)
	close(br.quit)
	br.workersWg.Wait()
	br.logger.Info("BlobReactor stopped, in-flight work completed")
}

// decodeSidecarsChunk decodes one slot's SSZ-encoded BlobSidecars and applies structural sanity checks shared by every inbound path:
// non-empty, bounded count, all sidecars bound to the same header, and unique in-bounds indices.
func (br *BlobReactor) decodeSidecarsChunk(chunk []byte) (datypes.BlobSidecars, error) {
	var sidecars datypes.BlobSidecars
	if err := sszUnmarshalSidecars(chunk, &sidecars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sidecars: %w", err)
	}
	if len(sidecars) == 0 {
		return nil, errors.New("empty sidecars chunk")
	}
	if uint64(len(sidecars)) > br.maxBlobsPerTx {
		return nil, fmt.Errorf("too many sidecars: %d (max %d)", len(sidecars), br.maxBlobsPerTx)
	}
	if err := sidecars.ValidateBlockRoots(); err != nil {
		return nil, fmt.Errorf("sidecars bound to different headers: %w", err)
	}
	seen := make(map[uint64]struct{}, len(sidecars))
	for _, sc := range sidecars {
		idx := sc.GetIndex()
		if idx >= uint64(len(sidecars)) {
			return nil, fmt.Errorf("sidecar index %d out of bounds (%d sidecars)", idx, len(sidecars))
		}
		if _, dup := seen[idx]; dup {
			return nil, fmt.Errorf("duplicate sidecar index %d", idx)
		}
		seen[idx] = struct{}{}
	}
	// Return the set in index order regardless of wire order. Consumers cache blob bundles positionally, so a
	// permuted-but-valid set from a racing peer must never survive to poison a re-proposal.
	sortSidecarsByIndex(sidecars)
	return sidecars, nil
}

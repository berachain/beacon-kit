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

package blobreactor

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/p2p"
)

const (
	// BlobChannel is our custom channel ID for blob requests/responses
	BlobChannel = byte(0x70)

	// ReactorName is the registered name for the blob reactor in CometBFT's switch
	ReactorName = "BLOBREACTOR"

	defaultSleepDuration       = 100 * time.Millisecond
	defaultPriority            = 5
	defaultSendQueueCapacity   = 100
	defaultRecvBufferCapacity  = 1024 * 1024
	defaultRecvMessageCapacity = 1024 * 1024

	defaultMaxRequestWorkers = 10

	maxBlobsPerBlock = 6
)

// blobRequestError wraps an error with a status for metrics tracking.
type blobRequestError struct {
	err    error
	status string
}

func (e *blobRequestError) Error() string {
	return e.err.Error()
}

func (e *blobRequestError) Unwrap() error {
	return e.err
}

func newBlobRequestError(err error, status string) error {
	return &blobRequestError{err: err, status: status}
}

// BlobReactor handles P2P blob distribution for BeaconKit.
// It implements the CometBFT Reactor interface.
type BlobReactor struct {
	service.BaseService
	sw *p2p.Switch

	blobStore BlobStore  // Storage backend for checking which blobs exist locally
	logger    log.Logger // Logger for the reactor
	config    Config     // Config for the reactor
	metrics   *Metrics

	// Track peers and our head slot
	stateMu  sync.RWMutex // Protects peers and headSlot
	peers    map[p2p.ID]struct{}
	headSlot math.Slot // Our own head slot (updated by blockchain service)
	nodeKey  string    // Our nodeKey (identity)

	// Concurrent request/response handling with per-request channels
	responseMu    sync.RWMutex
	responseChans map[uint64]chan *BlobResponse // requestID -> response channel

	// Worker pool for controlled concurrency
	requestWorkers chan struct{}  // semaphore for limiting concurrent request handlers
	workersWg      sync.WaitGroup // tracks active worker goroutines

	// Request ID counter
	nextRequestID atomic.Uint64 // atomic counter for generating unique request IDs

	// Shutdown flag to prevent new workers during stop
	stopped atomic.Bool // set to true when OnStop begins
}

// NewBlobReactor creates a new blob reactor with storage backend
func NewBlobReactor(blobStore BlobStore, logger log.Logger, cfg Config, metrics *Metrics) *BlobReactor {
	br := &BlobReactor{
		peers:          make(map[p2p.ID]struct{}),
		blobStore:      blobStore,
		logger:         logger,
		config:         cfg,
		metrics:        metrics,
		responseChans:  make(map[uint64]chan *BlobResponse),
		requestWorkers: make(chan struct{}, defaultMaxRequestWorkers),
	}
	br.BaseService = *service.NewBaseService(nil, ReactorName, br)
	return br
}

// SetSwitch allows setting a switch.
func (br *BlobReactor) SetSwitch(sw *p2p.Switch) {
	br.logger.Info("BlobReactor SetSwitch called", "switch", sw)
	br.sw = sw
}

// GetChannels returns the list of MConnection.ChannelDescriptor.
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

// InitPeer is called by the switch before the peer is started. Use it to
// initialize data for the peer (e.g. peer state).
func (br *BlobReactor) InitPeer(peer p2p.Peer) p2p.Peer {
	br.AddPeer(peer)
	return peer
}

// AddPeer is called by the switch after the peer is added and successfully started.
func (br *BlobReactor) AddPeer(peer p2p.Peer) {
	br.stateMu.Lock()
	br.peers[peer.ID()] = struct{}{}
	br.stateMu.Unlock()

	br.logger.Info("Added peer", "peer", peer.ID())
}

// RemovePeer is called by the switch when the peer is stopped (due to error or other reason).
func (br *BlobReactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	br.stateMu.Lock()
	delete(br.peers, peer.ID())
	br.stateMu.Unlock()

	br.logger.Info("Removed peer", "peer", peer.ID(), "reason", reason)
}

// spawnWorker attempts to spawn a worker goroutine to handle the given task.
// Returns true if worker was spawned, false if pool is full or reactor is stopped.
func (br *BlobReactor) spawnWorker(task func(), peerID p2p.ID, taskType string) {
	select {
	case br.requestWorkers <- struct{}{}:
		// Double-check stopped flag after acquiring worker slot to prevent race
		if br.stopped.Load() {
			<-br.requestWorkers // Release slot
			br.logger.Debug("Dropping message, reactor stopped during worker acquisition", "peer", peerID, "task_type", taskType)
			return
		}
		br.workersWg.Add(1)
		go func() {
			defer func() {
				<-br.requestWorkers
				br.workersWg.Done()
			}()
			task()
		}()
	default:
		br.logger.Warn("Worker pool full, dropping message", "peer", peerID, "task_type", taskType)
		br.metrics.observeWorkerPoolFull(taskType)
	}
}

// Receive is called by the switch when an envelope is received from any connected
// peer on any of the channels registered by the reactor
func (br *BlobReactor) Receive(envelope p2p.Envelope) {
	// Ignore messages if reactor is stopped
	if br.stopped.Load() {
		return
	}

	br.logger.Info("Received message on BlobChannel",
		"peer", envelope.Src.ID(),
		"channel", envelope.ChannelID,
		"peer_is_running", envelope.Src.IsRunning())

	// Get the message from the envelope
	blobMsg, ok := envelope.Message.(*BlobMessage)
	if !ok {
		br.logger.Error("Failed to cast message to BlobMessage", "peer", envelope.Src.ID(), "type", fmt.Sprintf("%T", envelope.Message))
		return
	}

	// Validate message has minimum length
	if len(blobMsg.Data) < 1 {
		br.logger.Error("Received message too short", "size", len(blobMsg.Data), "peer", envelope.Src.ID())
		return
	}

	msgType := MessageType(blobMsg.Data[0])
	msgData := blobMsg.Data[1:]

	br.logger.Info("Processing message", "type", msgType, "data_size", len(msgData), "peer", envelope.Src.ID())

	switch msgType {
	case MessageTypeRequest:
		var req BlobRequest
		if err := req.UnmarshalSSZ(msgData); err != nil {
			br.logger.Error("Failed to unmarshal BlobRequest", "error", err, "peer", envelope.Src.ID())
			return
		}
		br.logger.Info("Received blob request", "slot", req.Slot.Unwrap(), "request_id", req.RequestID, "peer", envelope.Src.ID())

		handleRequest := func() {
			br.handleBlobRequest(envelope.Src, &req)
		}
		br.spawnWorker(handleRequest, envelope.Src.ID(), "request")

	case MessageTypeResponse:
		var resp BlobResponse
		if err := resp.UnmarshalSSZ(msgData); err != nil {
			br.logger.Error("Failed to unmarshal BlobResponse", "error", err, "peer", envelope.Src.ID(), "data_size", len(msgData))
			return
		}
		br.logger.Info("Received blob response",
			"slot", resp.Slot.Unwrap(),
			"request_id", resp.RequestID,
			"peer", envelope.Src.ID(),
			"sidecar_data_size", len(resp.SidecarData))

		handleResponse := func() {
			br.handleBlobResponse(envelope.Src, &resp)
		}
		br.spawnWorker(handleResponse, envelope.Src.ID(), "response")

	default:
		br.logger.Warn("Received unknown message type", "type", msgType, "peer", envelope.Src.ID())
	}
}

func (br *BlobReactor) SetNodeKey(nodeKey string) {
	br.nodeKey = nodeKey
}

// SetHeadSlot updates the reactor's view of the current blockchain head slot.
func (br *BlobReactor) SetHeadSlot(slot math.Slot) {
	br.stateMu.Lock()
	br.headSlot = slot
	br.stateMu.Unlock()
}

// handleBlobRequest processes incoming blob requests and sends back blobs
func (br *BlobReactor) handleBlobRequest(peer p2p.Peer, req *BlobRequest) {
	br.logger.Info("Received blob request", "slot", req.Slot.Unwrap(), "request_id", req.RequestID, "peer", peer.ID())

	// Get our current head slot to include in response
	br.stateMu.RLock()
	headSlot := br.headSlot
	br.stateMu.RUnlock()

	// Fetch blobs from storage - if not found or error, sidecarBzs will be nil
	sidecarBzs, err := br.blobStore.GetByIndex(req.Slot.Unwrap())
	if err != nil {
		br.logger.Error("Failed to fetch blobs from storage", "slot", req.Slot.Unwrap(), "request_id", req.RequestID, "error", err)
	}
	resp := &BlobResponse{
		Slot:        req.Slot,
		RequestID:   req.RequestID,
		HeadSlot:    headSlot,
		SidecarData: encodeBlobSidecarsSSZ(sidecarBzs),
	}

	respBytes, err := resp.MarshalSSZ()
	if err != nil {
		br.logger.Error("Failed to marshal response", "slot", req.Slot.Unwrap(), "request_id", req.RequestID, "error", err)
		return
	}

	// Prepend message type
	msgData := append([]byte{byte(MessageTypeResponse)}, respBytes...)

	// Send response back to peer
	if !peer.Send(p2p.Envelope{ChannelID: BlobChannel, Message: NewBlobMessage(msgData)}) {
		br.logger.Warn("Failed to send blob response",
			"peer", peer.ID(),
			"slot", req.Slot,
			"request_id", req.RequestID,
			"data_size", len(msgData))
		// If sending response failed, the caller will timeout and try another peer
		return
	}

	br.logger.Info("Sent blob response", "slot", req.Slot.Unwrap(), "request_id", req.RequestID, "peer", peer.ID(), "data_size", len(msgData))
}

// handleBlobResponse processes incoming blob responses
func (br *BlobReactor) handleBlobResponse(peer p2p.Peer, resp *BlobResponse) {
	br.logger.Info("Received blob response",
		"slot", resp.Slot.Unwrap(),
		"request_id", resp.RequestID,
		"peer", peer.ID(),
		"data_size", len(resp.SidecarData),
		"peer_head", resp.HeadSlot)

	// Look up and remove the response channel for this request ID
	br.responseMu.Lock()
	respChan, exists := br.responseChans[resp.RequestID]
	if exists {
		delete(br.responseChans, resp.RequestID)
	}
	br.responseMu.Unlock()

	if !exists {
		br.logger.Info("No waiting channel for response (request may have timed out)",
			"request_id", resp.RequestID,
			"slot", resp.Slot.Unwrap())
		return
	}

	// Try to deliver the response
	select {
	case respChan <- resp:
		br.logger.Info("Delivered response to waiting request", "request_id", resp.RequestID, "slot", resp.Slot.Unwrap())
	default:
		br.logger.Warn("Response channel full, dropping response", "request_id", resp.RequestID, "slot", resp.Slot.Unwrap())
	}
}

// RequestBlobs fetches all blobs for a given slot from peers.
// Returns all blob sidecars for the slot, or an error if none could be retrieved.
// The context controls cancellation and timeout for the entire operation.
func (br *BlobReactor) RequestBlobs(
	ctx context.Context,
	slot math.Slot,
	verifier func(datypes.BlobSidecars) error) ([]*datypes.BlobSidecar, error) {
	br.logger.Info("RequestBlobs called", "slot", slot.Unwrap())

	// Check if we have any peers at all
	br.stateMu.RLock()
	peerCount := len(br.peers)
	br.stateMu.RUnlock()

	br.metrics.setPeerPoolSize(peerCount, peerCount)

	if peerCount == 0 {
		br.logger.Error("No peers available for blob request", "slot", slot.Unwrap())
		return nil, ErrNoPeersAvailable
	}

	// Track which peers we've already tried
	triedPeers := make(map[p2p.ID]bool)

	start := time.Now()

	// Continue trying while we have untried peers
	for {
		// Check context before trying next peer
		select {
		case <-ctx.Done():
			br.logger.Warn("Request cancelled before all peers tried", "slot", slot.Unwrap(), "peers_tried", len(triedPeers))
			br.metrics.recordOverallRequestComplete(statusTimeout, start)
			return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
		default:
		}

		// Select next untried peer
		peerID := br.selectUntriedPeer(triedPeers)
		if peerID == "" {
			// No more peers to try
			break
		}

		// Mark this peer as tried
		triedPeers[peerID] = true

		// Update available peers metric
		br.metrics.setPeerPoolSize(peerCount-len(triedPeers), peerCount)

		// Try to request blobs from this peer
		sidecars, err := br.requestBlobsFromPeer(ctx, peerID, slot)
		if err != nil {
			// Record per-peer failure with status
			status := statusInvalidResponse
			var reqErr *blobRequestError
			if errors.As(err, &reqErr) {
				status = reqErr.status
			}
			br.metrics.recordPeerAttempt(status)
			br.logger.Warn("Failed to get blobs from peer", "peer", peerID, "error", err)
			continue
		}

		// Sort sidecars by index to ensure correct order
		sort.Slice(sidecars, func(i, j int) bool { return sidecars[i].GetIndex() < sidecars[j].GetIndex() })

		// Verify the blobs before returning
		if verifyErr := verifier(sidecars); verifyErr != nil {
			br.metrics.recordPeerAttempt(statusVerifyFailed)
			br.logger.Warn("Blob verification failed, trying next peer",
				"slot", slot.Unwrap(),
				"count", len(sidecars),
				"peer", peerID,
				"error", verifyErr)
			continue
		}

		// Success - record both per-peer and overall metrics
		br.metrics.recordPeerAttempt(statusSuccess)
		br.metrics.recordOverallRequestComplete(statusSuccess, start)
		br.logger.Info("Successfully retrieved and verified blobs", "slot", slot.Unwrap(), "peer", peerID, "count", len(sidecars))
		return sidecars, nil
	}

	br.logger.Error("Failed to retrieve blobs from all peers", "slot", slot.Unwrap(), "peers_tried", len(triedPeers))
	br.metrics.recordOverallRequestComplete(statusAllPeersFailed, start)
	return nil, ErrAllPeersFailed
}

// selectUntriedPeer returns a random untried peer, or empty string if all peers have been tried.
func (br *BlobReactor) selectUntriedPeer(triedPeers map[p2p.ID]bool) p2p.ID {
	br.stateMu.RLock()
	defer br.stateMu.RUnlock()

	// Build list of untried peers
	var untried []p2p.ID
	for peerID := range br.peers {
		if !triedPeers[peerID] {
			untried = append(untried, peerID)
		}
	}

	if len(untried) == 0 {
		return "" // All peers tried or no peers available
	}

	// Return random untried peer to distribute load
	return untried[rand.Intn(len(untried))] // #nosec G404 // weak rng is acceptable for peer selection
}

// requestBlobsFromPeer sends a blob request to a specific peer and waits for response.
func (br *BlobReactor) requestBlobsFromPeer(ctx context.Context, peerID p2p.ID, slot math.Slot) (datypes.BlobSidecars, error) {
	peer := br.sw.Peers().Get(peerID)
	if peer == nil {
		return nil, newBlobRequestError(fmt.Errorf("peer %s not found", peerID), statusPeerNotFound)
	}

	if !peer.IsRunning() {
		return nil, newBlobRequestError(fmt.Errorf("peer %s not running", peerID), statusPeerNotFound)
	}

	// Generate unique request ID
	requestID := br.nextRequestID.Add(1)

	req := &BlobRequest{
		Slot:      slot,
		RequestID: requestID,
	}

	var err error
	reqBytes, err := req.MarshalSSZ()
	if err != nil {
		return nil, newBlobRequestError(fmt.Errorf("failed to marshal request from peer %s: %w", peerID, err), statusMarshalFailed)
	}

	// Create a dedicated response channel for this request
	respChan := make(chan *BlobResponse, 1)

	// Register the response channel and track active requests
	br.responseMu.Lock()
	br.responseChans[requestID] = respChan
	activeRequests := len(br.responseChans)
	br.responseMu.Unlock()
	br.metrics.setActiveRequests(activeRequests)

	cleanup := func() {
		br.logger.Debug("Cleaning up response channel", "request_id", requestID)
		br.responseMu.Lock()
		delete(br.responseChans, requestID)
		br.metrics.setActiveRequests(len(br.responseChans))
		br.responseMu.Unlock()
	}
	defer cleanup()

	msgData := append([]byte{byte(MessageTypeRequest)}, reqBytes...)
	if !peer.Send(p2p.Envelope{ChannelID: BlobChannel, Message: NewBlobMessage(msgData)}) {
		return nil, newBlobRequestError(fmt.Errorf("failed to send blob request to peer %s", peerID), statusSendFailed)
	}

	br.logger.Info("Sent blob request, waiting for response", "slot", slot.Unwrap(), "peer", peerID, "request_id", requestID)

	// Wait for response with timeout, respecting parent context
	timeoutCtx, cancel := context.WithTimeout(ctx, br.config.RequestTimeout)
	defer cancel()

	select {
	case resp := <-respChan:
		br.logger.Info("Received response", "slot", resp.Slot.Unwrap(), "peer", peerID, "data_size", len(resp.SidecarData))

		if resp.Slot != slot {
			err = fmt.Errorf("peer %s returned wrong slot: expected %d, got %d", peerID, slot.Unwrap(), resp.Slot.Unwrap())
			return nil, newBlobRequestError(err, statusInvalidResponse)
		}

		if resp.HeadSlot < resp.Slot {
			err = fmt.Errorf("peer %s head (%d) not at requested slot (%d)", peerID, resp.HeadSlot.Unwrap(), resp.Slot.Unwrap())
			return nil, newBlobRequestError(err, statusInvalidResponse)
		}

		if len(resp.SidecarData) > defaultRecvMessageCapacity {
			err = fmt.Errorf("peer %s sent oversized response: %d bytes (max %d)", peerID, len(resp.SidecarData), defaultRecvMessageCapacity)
			return nil, newBlobRequestError(err, statusInvalidResponse)
		}

		var sidecars datypes.BlobSidecars
		if len(resp.SidecarData) > 0 {
			if err = ssz.Unmarshal(resp.SidecarData, &sidecars); err != nil {
				err = fmt.Errorf("failed to unmarshal sidecars from peer %s: %w", peerID, err)
				return nil, newBlobRequestError(err, statusInvalidResponse)
			}
		}

		if len(sidecars) > maxBlobsPerBlock {
			err = fmt.Errorf("peer %s sent too many blobs: %d (max %d)", peerID, len(sidecars), maxBlobsPerBlock)
			return nil, newBlobRequestError(err, statusInvalidResponse)
		}

		return sidecars, nil

	case <-timeoutCtx.Done():
		if ctx.Err() != nil {
			return nil, newBlobRequestError(fmt.Errorf("request cancelled from peer %s: %w", peerID, ctx.Err()), statusTimeout)
		}
		err = fmt.Errorf("request timed out from peer %s after %v", peerID, br.config.RequestTimeout)
		return nil, newBlobRequestError(err, statusTimeout)
	}
}

func (br *BlobReactor) OnStart() error {
	br.logger.Info("Starting BlobReactor", "node_key", br.nodeKey)
	return nil
}

func (br *BlobReactor) OnStop() {
	br.logger.Info("Stopping BlobReactor", "node_key", br.nodeKey)

	// Set stop flag to prevent new workers from being spawned
	// This must happen before waiting for existing workers
	br.stopped.Store(true)

	// Wait for all worker goroutines to complete
	br.workersWg.Wait()

	br.logger.Info("BlobReactor stopped, all workers completed")
}

// encodeBlobSidecarsSSZ takes multiple SSZ-encoded BlobSidecar bytes and combines them
// into a single SSZ-encoded BlobSidecars (slice) format.
// The encoding is: 4-byte offset (always 4) + concatenated sidecars.
//
//nolint:mnd // ok for now
func encodeBlobSidecarsSSZ(sidecarBzs [][]byte) []byte {
	totalSize := 4
	for _, data := range sidecarBzs {
		totalSize += len(data)
	}

	result := make([]byte, totalSize)

	// Write offset (4) in little-endian - data starts after the offset
	binary.LittleEndian.PutUint32(result[0:4], 4)

	// Concatenate all sidecars after the offset (if any)
	pos := 4
	for _, data := range sidecarBzs {
		pos += copy(result[pos:], data)
	}

	return result
}

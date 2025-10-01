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

	defaultSleepDuration       = 100 * time.Millisecond
	defaultPriority            = 5
	defaultSendQueueCapacity   = 100
	defaultRecvBufferCapacity  = 1024 * 1024
	defaultRecvMessageCapacity = 1024 * 1024

	defaultMaxRequestWorkers = 10
)

// BlobReactor handles P2P blob distribution for BeaconKit.
// It implements the CometBFT Reactor interface.
type BlobReactor struct {
	service.BaseService
	sw *p2p.Switch

	blobStore BlobStore  // Storage backend for checking which blobs exist locally
	logger    log.Logger // Logger for the reactor
	config    Config     // Config for the reactor

	// Track peers and our head slot
	stateMu  sync.RWMutex // Protects peers and headSlot
	peers    map[p2p.ID]struct{}
	headSlot math.Slot // Our own head slot (updated by blockchain service)
	nodeKey  string    // Our nodeKey (identity)

	// Concurrent request/response handling with per-request channels
	responseMu    sync.RWMutex
	responseChans map[uint64]chan *BlobResponse // requestID -> response channel

	// Worker pool for controlled concurrency
	requestWorkers chan struct{} // semaphore for limiting concurrent request handlers

	// Request ID counter
	nextRequestID atomic.Uint64 // atomic counter for generating unique request IDs
}

// NewBlobReactor creates a new blob reactor with storage backend
func NewBlobReactor(blobStore BlobStore, logger log.Logger, cfg Config) *BlobReactor {
	br := &BlobReactor{
		peers:          make(map[p2p.ID]struct{}),
		blobStore:      blobStore,
		logger:         logger,
		config:         cfg,
		responseChans:  make(map[uint64]chan *BlobResponse),
		requestWorkers: make(chan struct{}, defaultMaxRequestWorkers),
	}
	br.BaseService = *service.NewBaseService(nil, "BlobReactor", br)
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

// Receive is called by the switch when an envelope is received from any connected
// peer on any of the channels registered by the reactor
func (br *BlobReactor) Receive(envelope p2p.Envelope) {
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

		select {
		case br.requestWorkers <- struct{}{}:
			go func() {
				defer func() { <-br.requestWorkers }()
				br.handleBlobRequest(envelope.Src, &req)
			}()
		default:
			br.logger.Warn("Worker pool full, dropping blob request",
				"slot", req.Slot,
				"request_id", req.RequestID,
				"peer", envelope.Src.ID())
		}

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

		select {
		case br.requestWorkers <- struct{}{}:
			go func() {
				defer func() { <-br.requestWorkers }()
				br.handleBlobResponse(envelope.Src, &resp)
			}()
		default:
			br.logger.Warn("Worker pool full, dropping response",
				"slot", resp.Slot.Unwrap(),
				"request_id", resp.RequestID,
				"peer", envelope.Src.ID())
		}

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

// HeadSlot returns the current blockchain head slot.
func (br *BlobReactor) HeadSlot() math.Slot {
	br.stateMu.RLock()
	defer br.stateMu.RUnlock()
	return br.headSlot
}

// handleBlobRequest processes incoming blob requests and sends back blobs
func (br *BlobReactor) handleBlobRequest(peer p2p.Peer, req *BlobRequest) {
	br.logger.Info("Received blob request", "slot", req.Slot.Unwrap(), "request_id", req.RequestID, "peer", peer.ID())

	// Get our current head slot to include in response
	br.stateMu.RLock()
	headSlot := br.headSlot
	br.stateMu.RUnlock()

	var errorMsg string
	var sidecarBzs [][]byte

	// TESTING: Simulate failure when requesting blobs that are divisible by 1000. Make them fail for 10000 slots
	// so that they will eventually succeed (to test retries being successful).
	if req.Slot.Unwrap()%1000 == 0 && headSlot.Unwrap() > req.Slot.Unwrap()+10000 {
		br.logger.Warn("TESTING: Simulating blob request failure", "slot", req.Slot.Unwrap())
		errorMsg = "simulated failure for testing"
	} else {
		var err error
		sidecarBzs, err = br.blobStore.GetByIndex(req.Slot.Unwrap())
		if err != nil {
			br.logger.Error("Failed to fetch blobs from storage", "slot", req.Slot.Unwrap(), "request_id", req.RequestID, "error", err)
			errorMsg = err.Error()
		}
	}

	resp := &BlobResponse{
		Slot:        req.Slot,
		RequestID:   req.RequestID,
		Error:       errorMsg,
		SidecarData: EncodeBlobSidecarsSSZ(sidecarBzs),
		HeadSlot:    headSlot,
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
			"error_msg", errorMsg,
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

	// Look up the response channel for this request ID
	br.responseMu.RLock()
	respChan, exists := br.responseChans[resp.RequestID]
	br.responseMu.RUnlock()

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
func (br *BlobReactor) RequestBlobs(
	slot math.Slot,
	expectedBlobs int,
	verifier func(datypes.BlobSidecars) error) ([]*datypes.BlobSidecar, error) {
	br.logger.Info("RequestBlobs called", "slot", slot.Unwrap())

	// Check if we have any peers at all
	br.stateMu.RLock()
	peerCount := len(br.peers)
	br.stateMu.RUnlock()

	if peerCount == 0 {
		br.logger.Error("No peers available for blob request", "slot", slot.Unwrap())
		return nil, ErrNoPeersAvailable
	}

	// Track which peers we've already tried
	triedPeers := make(map[p2p.ID]bool)

	// Continue trying while we have untried peers
	for {
		// Select next untried peer
		peerID := br.selectUntriedPeer(triedPeers)
		if peerID == "" {
			// No more peers to try
			break
		}

		// Mark this peer as tried
		triedPeers[peerID] = true

		// Try to request blobs from this peer
		sidecars, err := br.requestBlobsFromPeer(peerID, slot)
		if err != nil {
			br.logger.Warn("Failed to get blobs from peer", "peer", peerID, "error", err)
			continue
		}

		if len(sidecars) != expectedBlobs {
			br.logger.Warn("Received unexpected number of blob sidecars from peer",
				"peer", peerID,
				"slot", slot,
				"expected", expectedBlobs,
				"actual", len(sidecars))
			continue
		}

		// Sort sidecars by index to ensure correct order
		sort.Slice(sidecars, func(i, j int) bool { return sidecars[i].GetIndex() < sidecars[j].GetIndex() })

		// Verify the blobs before returning
		if verifyErr := verifier(sidecars); verifyErr != nil {
			br.logger.Warn("Blob verification failed, trying next peer",
				"slot", slot.Unwrap(),
				"count", len(sidecars),
				"peer", peerID,
				"error", verifyErr)
			continue
		}

		br.logger.Info("Successfully retrieved and verified blobs", "slot", slot.Unwrap(), "peer", peerID, "count", len(sidecars))
		return sidecars, nil
	}

	br.logger.Error("Failed to retrieve blobs from all peers", "slot", slot.Unwrap(), "peers_tried", len(triedPeers))
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
func (br *BlobReactor) requestBlobsFromPeer(peerID p2p.ID, slot math.Slot) (datypes.BlobSidecars, error) {
	peer := br.sw.Peers().Get(peerID)
	if peer == nil {
		return nil, fmt.Errorf("peer %s not found", peerID)
	}

	if !peer.IsRunning() {
		return nil, fmt.Errorf("peer %s not running", peerID)
	}

	// Generate unique request ID
	requestID := br.nextRequestID.Add(1)

	req := &BlobRequest{
		Slot:      slot,
		RequestID: requestID,
	}

	reqBytes, err := req.MarshalSSZ()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request from peer %s: %w", peerID, err)
	}

	// Create a dedicated response channel for this request
	respChan := make(chan *BlobResponse, 1)

	// Register the response channel
	br.responseMu.Lock()
	br.responseChans[requestID] = respChan
	br.responseMu.Unlock()

	cleanup := func() {
		br.logger.Debug("Cleaning up response channel", "request_id", requestID)
		br.responseMu.Lock()
		delete(br.responseChans, requestID)
		br.responseMu.Unlock()
	}
	defer cleanup()

	msgData := append([]byte{byte(MessageTypeRequest)}, reqBytes...)
	if !peer.Send(p2p.Envelope{ChannelID: BlobChannel, Message: NewBlobMessage(msgData)}) {
		return nil, fmt.Errorf("failed to send blob request to peer %s", peerID)
	}

	br.logger.Info("Sent blob request, waiting for response", "slot", slot.Unwrap(), "peer", peerID, "request_id", requestID)

	// Wait for response with timeout
	ctx, cancel := context.WithTimeout(context.Background(), br.config.RequestTimeout)
	defer cancel()

	select {
	case resp := <-respChan:
		br.logger.Info("Received response",
			"slot", resp.Slot.Unwrap(),
			"peer", peerID,
			"data_size", len(resp.SidecarData),
			"error", resp.Error)

		// Check if peer reported an error
		if resp.Error != "" {
			return nil, fmt.Errorf("peer %s reported error: %s", peerID, resp.Error)
		}

		if resp.HeadSlot < resp.Slot {
			return nil, fmt.Errorf("peer %s head (%d) not at requested slot (%d)", peerID, resp.HeadSlot.Unwrap(), resp.Slot.Unwrap())
		}

		var sidecars datypes.BlobSidecars
		if len(resp.SidecarData) > 0 {
			if err = ssz.Unmarshal(resp.SidecarData, &sidecars); err != nil {
				return nil, fmt.Errorf("failed to unmarshal sidecars from peer %s: %w", peerID, err)
			}
		}

		return sidecars, nil

	case <-ctx.Done():
		return nil, fmt.Errorf("request timed out from peer %s after %v", peerID, br.config.RequestTimeout)
	}
}

func (br *BlobReactor) OnStart() error {
	br.logger.Info("Starting BlobReactor", "node_key", br.nodeKey)
	return nil
}

func (br *BlobReactor) OnStop() {
	br.logger.Info("Stopping BlobReactor", "node_key", br.nodeKey)
}

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

// BlobReactor handles P2P blob distribution for BeaconKit
type BlobReactor struct {
	blobStore BlobStore  // Storage backend for checking which blobs exist locally
	logger    log.Logger // Logger for the reactor
	config    Config     // Config for the reactor

	service.BaseService             // Embedding BaseService to manage lifecycle
	sw                  *p2p.Switch // The switch is set by the p2p layer

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

func (br *BlobReactor) SetNodeKey(nodeKey string) {
	br.nodeKey = nodeKey
}

// SetHeadSlot updates the reactor's view of the current blockchain head slot.
// Called by the blockchain service after processing each block.
func (br *BlobReactor) SetHeadSlot(slot uint64) {
	br.stateMu.Lock()
	br.headSlot = math.Slot(slot)
	br.stateMu.Unlock()
}

func (br *BlobReactor) GetChannels() []*p2p.ChannelDescriptor {
	br.logger.Info("BlobReactor GetChannels called", "channel_id", fmt.Sprintf("0x%02X", BlobChannel))
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

// SetSwitch implements Reactor by setting the switch
func (br *BlobReactor) SetSwitch(sw *p2p.Switch) {
	br.logger.Info("BlobReactor SetSwitch called", "switch", sw)
	br.sw = sw
}

// InitPeer is called when a peer is initialized
func (br *BlobReactor) InitPeer(peer p2p.Peer) p2p.Peer {
	br.AddPeer(peer)
	return peer
}

// AddPeer is called when a peer is added
func (br *BlobReactor) AddPeer(peer p2p.Peer) {
	br.stateMu.Lock()
	br.peers[peer.ID()] = struct{}{}
	br.stateMu.Unlock()

	br.logger.Info("Added peer", "peer", peer.ID())
}

// RemovePeer is called when a peer is removed
func (br *BlobReactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	br.stateMu.Lock()
	delete(br.peers, peer.ID())
	br.stateMu.Unlock()

	br.logger.Info("Removed peer", "peer", peer.ID(), "reason", reason)
}

// Receive handles incoming messages
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
		br.logger.Info("Received blob request", "slot", req.Slot, "request_id", req.RequestID, "peer", envelope.Src.ID())

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
			"slot", resp.Slot,
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
				"slot", resp.Slot,
				"request_id", resp.RequestID,
				"peer", envelope.Src.ID())
		}

	default:
		br.logger.Warn("Received unknown message type", "type", msgType, "peer", envelope.Src.ID())
	}
}

// handleBlobRequest processes incoming blob requests and sends back blobs
func (br *BlobReactor) handleBlobRequest(peer p2p.Peer, req *BlobRequest) {
	br.logger.Info("Received blob request", "slot", req.Slot, "request_id", req.RequestID, "peer", peer.ID())

	// Get our current head slot to include in response
	br.stateMu.RLock()
	headSlot := br.headSlot
	br.stateMu.RUnlock()

	var errorMsg string
	sidecarBzs, err := br.blobStore.GetByIndex(req.Slot.Unwrap())
	if err != nil {
		br.logger.Error("Failed to fetch blobs from storage", "slot", req.Slot, "request_id", req.RequestID, "error", err)
		errorMsg = err.Error()
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
		br.logger.Error("Failed to marshal response", "slot", req.Slot, "request_id", req.RequestID, "error", err)
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

	br.logger.Info("Sent blob response", "slot", req.Slot, "request_id", req.RequestID, "peer", peer.ID(), "data_size", len(msgData))
}

// handleBlobResponse processes incoming blob responses
func (br *BlobReactor) handleBlobResponse(peer p2p.Peer, resp *BlobResponse) {
	br.logger.Info("Received blob response",
		"slot", resp.Slot,
		"request_id", resp.RequestID,
		"peer", peer.ID(),
		"data_size", len(resp.SidecarData), "peer_head", resp.HeadSlot)

	// Look up the response channel for this request ID
	br.responseMu.RLock()
	respChan, exists := br.responseChans[resp.RequestID]
	br.responseMu.RUnlock()

	if !exists {
		br.logger.Info("No waiting channel for response (request may have timed out)",
			"request_id", resp.RequestID,
			"slot", resp.Slot)
		return
	}

	// Try to deliver the response
	select {
	case respChan <- resp:
		br.logger.Info("Delivered response to waiting request", "request_id", resp.RequestID, "slot", resp.Slot)
	default:
		br.logger.Warn("Response channel full, dropping response", "request_id", resp.RequestID, "slot", resp.Slot)
	}
}

// RequestBlobs fetches all blobs for a given slot from peers.
// Returns all blob sidecars for the slot, or an error if none could be retrieved.
//
//nolint:funlen,gocognit,maintidx // ok for now
func (br *BlobReactor) RequestBlobs(slot uint64, verifier func(datypes.BlobSidecars) error) ([]*datypes.BlobSidecar, error) {
	br.logger.Info("RequestBlobs called", "slot", slot)

	br.stateMu.RLock()
	ourHead := br.headSlot
	peers := make([]p2p.ID, 0, len(br.peers))
	for peerID := range br.peers {
		peers = append(peers, peerID)
	}
	br.stateMu.RUnlock()

	br.logger.Info("Current state", "our_head", ourHead, "requested_slot", slot, "num_peers", len(peers))

	if len(peers) == 0 {
		br.logger.Error("No peers available for blob request", "slot", slot)
		return nil, ErrNoPeersAvailable
	}

	// Randomize peer order to distribute load
	rand.Shuffle(len(peers), func(i, j int) { peers[i], peers[j] = peers[j], peers[i] })

	br.logger.Info("Found peers for blob request", "slot", slot, "num_peers", len(peers), "peers", peers)

	// Track which peers we've already tried
	triedPeers := make(map[p2p.ID]bool)

	// Continue trying while we have untried peers
	for {
		// Check if we have any untried peers left
		hasUntriedPeer := false
		for _, peerID := range peers {
			if !triedPeers[peerID] {
				hasUntriedPeer = true
				break
			}
		}

		if !hasUntriedPeer {
			// All current peers have been tried, check if there are new peers
			br.stateMu.RLock()
			newPeersFound := false
			for peerID := range br.peers {
				if !triedPeers[peerID] {
					// Found a peer we haven't tried yet
					newPeersFound = true
					// Rebuild peer list with current peers
					peers = make([]p2p.ID, 0, len(br.peers))
					for p := range br.peers {
						peers = append(peers, p)
					}
					break
				}
			}
			br.stateMu.RUnlock()

			if !newPeersFound {
				// No new peers, we've tried everyone
				break
			}

			// Randomize the new peer list
			rand.Shuffle(len(peers), func(i, j int) { peers[i], peers[j] = peers[j], peers[i] })
			br.logger.Info("Refreshed peer list after exhausting attempts", "new_total", len(peers), "tried_so_far", len(triedPeers))
		}

		// Find next untried peer
		var peerID p2p.ID
		found := false
		for _, p := range peers {
			if !triedPeers[p] {
				peerID = p
				found = true
				break
			}
		}

		if !found {
			continue // This shouldn't happen but safety check
		}

		// Mark this peer as tried
		triedPeers[peerID] = true
		peer := br.sw.Peers().Get(peerID)
		if peer == nil {
			br.logger.Info("Peer no longer connected, skipping", "peer", peerID)
			continue
		}

		if !peer.IsRunning() {
			br.logger.Info("Peer not running, skipping", "peer", peerID)
			continue
		}

		// Generate unique request ID
		requestID := br.nextRequestID.Add(1)

		req := &BlobRequest{
			Slot:      math.Slot(slot),
			RequestID: requestID,
		}

		reqBytes, err := req.MarshalSSZ()
		if err != nil {
			br.logger.Error("Failed to marshal request", "error", err)
			continue
		}

		// Create a dedicated response channel for this request
		respChan := make(chan *BlobResponse, 1)

		// Register the response channel
		br.responseMu.Lock()
		br.responseChans[requestID] = respChan
		br.responseMu.Unlock()

		cleanup := func() {
			br.logger.Info("Cleaning up response channel", "request_id", requestID)
			br.responseMu.Lock()
			delete(br.responseChans, requestID)
			br.responseMu.Unlock()
			br.logger.Info("Cleaned up response channel", "request_id", requestID)
		}

		msgData := append([]byte{byte(MessageTypeRequest)}, reqBytes...)
		if !peer.Send(p2p.Envelope{ChannelID: BlobChannel, Message: NewBlobMessage(msgData)}) {
			br.logger.Error("Failed to send blob request to peer", "peer", peerID, "slot", slot)
			cleanup()
			continue
		}

		br.logger.Info("Sent blob request, waiting for response", "slot", slot, "peer", peerID, "request_id", requestID)

		// Wait for response with timeout
		ctx, cancel := context.WithTimeout(context.Background(), br.config.RequestTimeout)
		br.logger.Info("Starting wait for response", "request_id", requestID, "timeout_ms", br.config.RequestTimeout.Milliseconds())

		select {
		case resp := <-respChan:
			cancel() // Cancel context immediately on response
			br.logger.Info("Received response", "slot", resp.Slot, "data_size", len(resp.SidecarData), "error", resp.Error)

			// Check if peer reported an error
			if resp.Error != "" {
				br.logger.Warn("Peer reported error fetching blobs", "slot", slot, "peer", peerID, "error", resp.Error)
				cleanup()
				continue
			}

			if resp.HeadSlot < resp.Slot {
				br.logger.Info(
					"Peer head was not at requested slot, trying next peer",
					"slot", slot,
					"peer_head", resp.HeadSlot,
					"peer", peerID)
				cleanup()
				continue
			}

			var sidecars datypes.BlobSidecars
			if len(resp.SidecarData) > 0 {
				if err = ssz.Unmarshal(resp.SidecarData, &sidecars); err != nil {
					br.logger.Error("Failed to unmarshal sidecars from response", "error", err, "peer", peerID)
					cleanup()
					continue
				}
			}

			indices := make([]uint64, len(sidecars))
			for i, sc := range sidecars {
				indices[i] = sc.GetIndex()
			}
			br.logger.Info("Blob indices from peer", "slot", slot, "indices", indices, "count", len(sidecars))

			// If peer returned no blobs, try next peer
			if len(sidecars) == 0 {
				br.logger.Warn("Peer returned no blobs despite having the slot, trying next peer",
					"slot", slot,
					"peer", peerID,
					"peer_head", resp.HeadSlot)
				cleanup()
				continue
			}

			// IMPORTANT: Sort sidecars by index to ensure correct order. The storage backend may return them in arbitrary order
			sort.Slice(sidecars, func(i, j int) bool { return sidecars[i].GetIndex() < sidecars[j].GetIndex() })

			// Warn if indices are not sequential starting from 0 (this should fail later in verification but this
			// is a good debugging warning)
			for i, sc := range sidecars {
				// #nosec G115
				if sc.GetIndex() != uint64(i) {
					br.logger.Warn("Non-sequential blob indices detected",
						"slot", slot,
						"expected_index", i,
						"actual_index",
						sc.GetIndex(), "peer", peerID)
				}
			}

			// Verify the blobs before returning
			if verifyErr := verifier(sidecars); verifyErr != nil {
				br.logger.Warn("Blob verification failed, trying next peer", "slot", slot, "peer", peerID, "error", verifyErr)
				cleanup()
				continue
			}

			br.logger.Info("Successfully retrieved and verified blobs", "slot", slot, "peer", peerID, "count", len(sidecars))
			cleanup()
			return sidecars, nil

		case <-ctx.Done():
			cancel() // Cancel context on timeout
			br.logger.Warn("Request timed out, trying next peer",
				"slot", slot,
				"peer", peerID,
				"request_id", requestID,
				"timeout", br.config.RequestTimeout)
			cleanup()
			continue
		}
	}

	br.logger.Error("Failed to retrieve blobs from all peers", "slot", slot, "peers_tried", len(triedPeers))
	return nil, ErrAllPeersFailed
}

func (br *BlobReactor) OnStart() error {
	br.logger.Info("Starting BlobReactor", "node_key", br.nodeKey)
	return nil
}

func (br *BlobReactor) OnStop() {
	br.logger.Info("Stopping BlobReactor", "node_key", br.nodeKey)
}

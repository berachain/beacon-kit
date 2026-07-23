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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/cometbft/cometbft/p2p"
)

// NewBlobMessageForTest exposes the internal envelope constructor to tests
// that need to craft raw wire messages.
func NewBlobMessageForTest(msgType MessageType, msgBz []byte) *BlobMessage {
	return newBlobMessage(msgType, msgBz)
}

// PeerKnowsRoot reports whether the reactor has recorded that the given peer holds the sidecars for root.
func (br *BlobReactor) PeerKnowsRoot(peerID p2p.ID, root common.Root) bool {
	br.stateMu.RLock()
	defer br.stateMu.RUnlock()
	ps, ok := br.peers[peerID]
	return ok && ps.knownRoots.has(root)
}

// PushQueueCapacityForTest exposes the push lane capacity so scheduler tests can assert the drop boundary
// exactly.
const PushQueueCapacityForTest = pushQueueCapacity

// EnqueuePushForTest, EnqueueByRootForTest and EnqueueByRangeForTest expose the lane queues so scheduler
// behavior (bounded drop, serve priority, shutdown draining) can be pinned without crafting wire messages.
func (br *BlobReactor) EnqueuePushForTest(src p2p.Peer, push *SidecarsPush) {
	enqueue(br, br.pushQueue, pushTask{src: src, push: push}, "test", "push")
}

func (br *BlobReactor) EnqueueByRootForTest(src p2p.Peer, req *SidecarsByRootRequest) {
	enqueue(br, br.byRootQueue, byRootTask{src: src, req: req}, "test", "by_root_request")
}

func (br *BlobReactor) EnqueueByRangeForTest(src p2p.Peer, req *SidecarsByRangeRequest) {
	enqueue(br, br.byRangeQueue, byRangeTask{src: src, req: req}, "test", "by_range_request")
}

// PushQueueLenForTest reports the number of queued push tasks.
func (br *BlobReactor) PushQueueLenForTest() int { return len(br.pushQueue) }

// RunOneServeTaskForTest runs a single serve-lane scheduling decision synchronously, returning the lane label
// ("by_root" or "by_range") that was served, or ok=false once the reactor is stopping. Call only with tasks
// already queued, it blocks otherwise.
func (br *BlobReactor) RunOneServeTaskForTest() (string, bool) {
	return br.runOneServeTask()
}

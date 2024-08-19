// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package block

// storeMetrics is a struct that contains metrics for the block store.
type storeMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newStoreMetrics creates a new storeMetrics.
func newStoreMetrics(sink TelemetrySink) *storeMetrics {
	return &storeMetrics{
		sink: sink,
	}
}

// markPruneBlockFailure marks a block prune failure.
func (sm *storeMetrics) markPruneBlockFailure(err error) {
	var labels []string
	if err != nil {
		labels = append(labels, "reason", "error")
	} else {
		labels = append(labels, "reason", "panic")
	}

	sm.sink.IncrementCounter(
		"beacon_kit.block-store.prune_block_failure", labels...,
	)
}

// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package filedb

import "context"

// RangeDBPruner prunes old indexes from the range DB.
type RangeDBPruner struct {
	db             *RangeDB
	notifyNewIndex chan uint64
	pruneWindow    uint64
}

func NewRangeDBPruner(
	db *RangeDB,
	pruneWindow uint64,
	notifyNewIndex chan uint64,
) *RangeDBPruner {
	return &RangeDBPruner{
		db:             db,
		notifyNewIndex: notifyNewIndex,
		pruneWindow:    pruneWindow,
	}
}

// Start starts a goroutine that listens for new indices to prune.
func (p *RangeDBPruner) Start(ctx context.Context) {
	go p.runLoop(ctx)
}

// NotifyNewIndex notifies the pruner of a new index.
func (p *RangeDBPruner) runLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case index := <-p.notifyNewIndex:
			// Do nothing until the prune window is reached.
			if p.pruneWindow > index {
				continue
			}

			// Prune all indexes below the most recently seen
			// index minus the prune window.
			if err := p.db.DeleteRange(0, index-p.pruneWindow); err != nil {
				return
			}
		}
	}
}

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

package cometbft

import (
	"testing"
	"time"
)

func TestFinalizeBlock_NextBlockDelay(t *testing.T) {
	tests := []struct {
		name            string
		initialTime     time.Time
		curBlockTime    time.Time
		curBlockHeight  int64
		targetBlockTime time.Duration
		nextBlockDelay  time.Duration
	}{
		{
			name:            "block is produced in 1s, but target is 2s => wait 1s",
			initialTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			curBlockTime:    time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
			curBlockHeight:  1,
			targetBlockTime: 2 * time.Second,
			nextBlockDelay:  time.Second,
		},
		{
			name:            "block is produced 1s past deadline => no delay",
			initialTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			curBlockTime:    time.Date(2024, 1, 1, 0, 0, 5, 0, time.UTC),
			curBlockHeight:  2,
			targetBlockTime: 2 * time.Second,
			nextBlockDelay:  1 * time.Microsecond, // no delay
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nextBlockDelay(tt.initialTime, tt.curBlockTime, tt.curBlockHeight, tt.targetBlockTime)
			if got != tt.nextBlockDelay {
				t.Errorf("nextBlockDelay() = %v, want %v", got, tt.nextBlockDelay)
			}
		})
	}
}

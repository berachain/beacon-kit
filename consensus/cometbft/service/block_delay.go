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

package cometbft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// maxDelayBetweenBlocks is the maximum delay between two consecutive blocks.
// If the last block time minus the previous block time is greater than
// maxDelayBetweenBlocks, then we reset `FinalizeBlockResponse.NextBlockDelay`
// to default.
//
// This is needed because the network may stall for a long time and we don't
// want to rush in new blocks as the network resumes its operation.
const maxDelayBetweenBlocks = 30 * time.Minute

type blockDelay struct {
	// - genesis time if height is 0 OR
	// - block time of the last block if height > 0 and time between blocks is
	// greater than maxDelayBetweenBlocks.
	InitialTime time.Time

	// - InitChainRequest.InitialHeight OR
	// - last block height if time between blocks is greater than maxDelayBetweenBlocks.
	InitialHeight int64

	// PreviousBlockTime is the time of the previous block.
	PreviousBlockTime time.Time
}

// CONTRACT: called only once upon genesis during or after InitChain.
func blockDelayUponGenesis(genesisTime time.Time, initialHeight int64) *blockDelay {
	return &blockDelay{
		InitialTime:       genesisTime,
		InitialHeight:     initialHeight,
		PreviousBlockTime: genesisTime,
	}
}

// Panics if it fails to read from the buffer.
func blockDelayFromBytes(
	bz []byte,
) *blockDelay {
	buf := bytes.NewReader(bz)
	var initialTime, prevBlockTime int64
	var initialHeight int64

	err := binary.Read(buf, binary.LittleEndian, &initialTime)
	if err != nil {
		panic(fmt.Sprintf("failed to read InitialTime: %v", err))
	}
	err = binary.Read(buf, binary.LittleEndian, &initialHeight)
	if err != nil {
		panic(fmt.Sprintf("failed to read InitialHeight: %v", err))
	}
	err = binary.Read(buf, binary.LittleEndian, &prevBlockTime)
	if err != nil {
		panic(fmt.Sprintf("failed to read PreviousBlockTime: %v", err))
	}

	return &blockDelay{
		InitialTime:       time.Unix(0, initialTime),
		InitialHeight:     initialHeight,
		PreviousBlockTime: time.Unix(0, prevBlockTime),
	}
}

// Next returns the duration to wait before proposing the next block.
func (d *blockDelay) Next(curBlockTime time.Time, curBlockHeight int64,
	targetBlockTime time.Duration) time.Duration {
	// Until `timeout_commit` is removed from the CometBFT config,
	// `FinalizeBlockResponse.NextBlockDelay` can't be exactly 0. If it's set to
	// 0, then `timeout_commit` from the config will be used, which is not what
	// we want since we're trying to control the block time.
	const noDelay = 1 * time.Microsecond

	// Reset the initial time and height if the time between blocks is greater
	// than maxDelayBetweenBlocks.
	if curBlockTime.Sub(d.PreviousBlockTime) > maxDelayBetweenBlocks {
		d.InitialTime = curBlockTime
		d.InitialHeight = curBlockHeight - 1
	}

	t := d.InitialTime.Add(targetBlockTime * time.Duration(curBlockHeight-d.InitialHeight))
	if curBlockTime.Before(t) {
		return t.Sub(curBlockTime)
	}

	// Update the previous block time.
	d.PreviousBlockTime = curBlockTime

	return noDelay
}

// ToBytes converts the blockDelay to bytes.
//
// Panics if it fails to write to the buffer.
func (d *blockDelay) ToBytes() []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, d.InitialTime.UnixNano())
	if err != nil {
		panic(fmt.Sprintf("failed to write InitialTime: %v", err))
	}
	err = binary.Write(buf, binary.LittleEndian, d.InitialHeight)
	if err != nil {
		panic(fmt.Sprintf("failed to write InitialHeight: %v", err))
	}
	err = binary.Write(buf, binary.LittleEndian, d.PreviousBlockTime.UnixNano())
	if err != nil {
		panic(fmt.Sprintf("failed to write PreviousBlockTime: %v", err))
	}

	return buf.Bytes()
}

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

const (
	// MaxDelayBetweenBlocks is the maximum delay between two consecutive blocks.
	// If the last block time minus the previous block time is greater than
	// MaxDelayBetweenBlocks, then we reset `FinalizeBlockResponse.NextBlockDelay`
	// to default.
	//
	// This is needed because the network may stall for a long time and we don't
	// want to rush in new blocks as the network resumes its operation.
	MaxDelayBetweenBlocks = 5 * time.Minute

	// TargetBlockTime is the desired block time.
	//
	// Note that it CAN'T be lower than the minimal (floor) block time in the
	// network, which is comprised of the time to a) propose a new block b)
	// gather 2/3+ prevotes c) gather 2/3+ precommits.
	TargetBlockTime = 2 * time.Second
)

// BlockDelayFromBytesError is returned when the block delay can't be read from
// the buffer.
type BlockDelayFromBytesError struct {
	field string
	err   error
}

func (e *BlockDelayFromBytesError) Error() string {
	return fmt.Sprintf("%s: %v", e.field, e.err)
}

func (e *BlockDelayFromBytesError) Unwrap() error {
	return e.err
}

// BlockDelay is used to calculate the delay before proposing the next block.
//
// The general formula is:
//
//	NextBlockTime = InitialTime + TargetBlockTime * (CurrentHeight - InitialHeight)
type BlockDelay struct {
	// InitialTime is a checkpoint in time from which we start calculating the
	// next block time.
	InitialTime time.Time

	// InitialHeight is a checkpoint in blockchain's height from which we start
	// calculating the next block time.
	InitialHeight int64

	// PreviousBlockTime is the time of the previous block.
	PreviousBlockTime time.Time
}

// BlockDelayFromBytes converts the bytes to a blockDelay.
//
// Expected format:
//
//	InitialTime (int64) | InitialHeight (int64) | PreviousBlockTime (int64)
//	(little endian)
//
// It returns ErrBlockDelayFromBytes if it fails to read from the buffer.
func BlockDelayFromBytes(
	bz []byte,
) (*BlockDelay, error) {
	buf := bytes.NewReader(bz)
	var initialTime, prevBlockTime int64
	var initialHeight int64

	err := binary.Read(buf, binary.LittleEndian, &initialTime)
	if err != nil {
		return nil, &BlockDelayFromBytesError{
			field: "InitialTime",
			err:   err,
		}
	}
	err = binary.Read(buf, binary.LittleEndian, &initialHeight)
	if err != nil {
		return nil, &BlockDelayFromBytesError{
			field: "InitialHeight",
			err:   err,
		}
	}
	err = binary.Read(buf, binary.LittleEndian, &prevBlockTime)
	if err != nil {
		return nil, &BlockDelayFromBytesError{
			field: "PreviousBlockTime",
			err:   err,
		}
	}

	return &BlockDelay{
		InitialTime:       time.Unix(0, initialTime),
		InitialHeight:     initialHeight,
		PreviousBlockTime: time.Unix(0, prevBlockTime),
	}, nil
}

// Next returns the duration to wait before proposing the next block.
func (d *BlockDelay) Next(curBlockTime time.Time, curBlockHeight int64) time.Duration {
	// Until `timeout_commit` is removed from the CometBFT config,
	// `FinalizeBlockResponse.NextBlockDelay` can't be exactly 0. If it's set to
	// 0, then `timeout_commit` from the config will be used, which is not what
	// we want since we're trying to control the block time.
	const noDelay = 1 * time.Microsecond

	// Reset the initial time and height if the time between blocks is greater
	// than MaxDelayBetweenBlocks. This makes the current time and height the
	// initial one as if the upgrade happened just now.
	if curBlockTime.Sub(d.PreviousBlockTime) > MaxDelayBetweenBlocks {
		d.InitialTime = curBlockTime
		d.InitialHeight = curBlockHeight
		d.PreviousBlockTime = curBlockTime
		return TargetBlockTime
	}

	d.PreviousBlockTime = curBlockTime

	t := d.InitialTime.Add(TargetBlockTime * time.Duration(curBlockHeight-d.InitialHeight))
	if curBlockTime.Before(t) {
		return t.Sub(curBlockTime)
	}
	return noDelay
}

// ToBytes converts the blockDelay to bytes.
//
// Format:
//
//	InitialTime (int64) | InitialHeight (int64) | PreviousBlockTime (int64)
//	(little endian)
//
// Panics if it fails to write to the buffer.
func (d *BlockDelay) ToBytes() []byte {
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

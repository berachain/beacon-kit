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

package cometbft_test

import (
	"testing"
	"time"

	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockDelayFromBytes(t *testing.T) {
	t.Parallel()

	d1 := &cometbft.BlockDelay{
		InitialTime:       time.Now().Add(-10 * time.Minute),
		InitialHeight:     5,
		PreviousBlockTime: time.Now().Add(-5 * time.Minute),
	}

	b := d1.ToBytes()
	d2, err := cometbft.BlockDelayFromBytes(b)
	require.NoError(t, err)

	assert.Equal(t, d1.InitialTime.Unix(), d2.InitialTime.Unix())
	assert.Equal(t, d1.InitialHeight, d2.InitialHeight)
	assert.Equal(t, d1.PreviousBlockTime.Unix(), d2.PreviousBlockTime.Unix())
}

func TestBlockDelayNext_NoDelay(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()
	initialHeight := int64(1)
	d := &cometbft.BlockDelay{
		InitialTime:       genesisTime,
		InitialHeight:     initialHeight,
		PreviousBlockTime: genesisTime,
	}

	curBlockTime := genesisTime.Add(10 * time.Second)
	curBlockHeight := int64(2)

	delay := d.Next(curBlockTime, curBlockHeight)

	assert.Equal(t, 1*time.Microsecond, delay)

	// InitialTime/Height are not updated, PreviousBlockTime is updated
	assert.Equal(t, genesisTime, d.InitialTime)
	assert.Equal(t, initialHeight, d.InitialHeight)
	assert.Equal(t, curBlockTime, d.PreviousBlockTime)
}

func TestBlockDelayNext_WithDelay(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()
	initialHeight := int64(1)
	d := &cometbft.BlockDelay{
		InitialTime:       genesisTime,
		InitialHeight:     initialHeight,
		PreviousBlockTime: genesisTime,
	}

	curBlockTime := genesisTime.Add(2 * time.Second)
	curBlockHeight := int64(3)

	delay := d.Next(curBlockTime, curBlockHeight)

	assert.Equal(t, 2*time.Second, delay)

	// InitialTime/Height are not updated, PreviousBlockTime is updated
	assert.Equal(t, genesisTime, d.InitialTime)
	assert.Equal(t, initialHeight, d.InitialHeight)
	assert.Equal(t, curBlockTime, d.PreviousBlockTime)
}

func TestBlockDelayNext_ResetOnStall(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()
	initialHeight := int64(1)
	d := &cometbft.BlockDelay{
		InitialTime:       genesisTime,
		InitialHeight:     initialHeight,
		PreviousBlockTime: genesisTime,
	}

	curBlockTime := genesisTime.Add(cometbft.MaxDelayBetweenBlocks + 1*time.Minute)
	curBlockHeight := int64(10)

	delay := d.Next(curBlockTime, curBlockHeight)

	assert.Equal(t, d.InitialTime, curBlockTime)
	assert.Equal(t, d.InitialHeight, curBlockHeight)
	assert.Equal(t, cometbft.TargetBlockTime, delay)
}

func TestBlockDelayNext_Catchup(t *testing.T) {
	t.Parallel()

	genesisTime := time.Now()
	initialHeight := int64(1)
	d := &cometbft.BlockDelay{
		InitialTime:       genesisTime,
		InitialHeight:     initialHeight,
		PreviousBlockTime: genesisTime,
	}

	curBlockTime := genesisTime.Add(2 * time.Second)
	curBlockHeight := int64(3)

	delay := d.Next(curBlockTime, curBlockHeight)

	assert.Equal(t, 2*time.Second, delay)

	// Ideal
	// height 4 -> 3*2 = 6
	// height 5 -> 4*2 = 8
	// height 6 -> 5*2 = 10
	// height 7 -> 6*2 = 12

	curBlockHeight++
	curBlockTime = genesisTime.Add(8 * time.Second)
	// T(h4) = 8s
	// delay = 1us
	// T(h5 = 9s)
	// delay = 1us
	// T(h6) = 10s
	// delay = 1us
	// T(h7) = 11s
	// delay = 1
	delay = d.Next(curBlockTime, curBlockHeight)
	assert.Equal(t, 1*time.Microsecond, delay)

	for curBlockHeight++; curBlockHeight < 7; curBlockHeight++ {
		curBlockTime = curBlockTime.Add(1 * time.Second)
		delay = d.Next(curBlockTime, curBlockHeight)
		assert.Equal(t, 1*time.Microsecond, delay)
	}
	curBlockTime = curBlockTime.Add(1 * time.Second)
	delay = d.Next(curBlockTime, curBlockHeight)

	assert.Equal(t, 1*time.Second, delay)
}

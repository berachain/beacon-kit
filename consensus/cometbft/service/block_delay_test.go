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
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlockDelayUponGenesis(t *testing.T) {
	genesisTime := time.Now()
	initialHeight := int64(1)

	d := blockDelayUponGenesis(genesisTime, initialHeight)

	assert.Equal(t, genesisTime, d.InitialTime)
	assert.Equal(t, initialHeight, d.InitialHeight)
	assert.Equal(t, genesisTime, d.PreviousBlockTime)
}

func TestBlockDelayFromBytes(t *testing.T) {
	d1 := &blockDelay{
		InitialTime:       time.Now().Add(-10 * time.Minute),
		InitialHeight:     5,
		PreviousBlockTime: time.Now().Add(-5 * time.Minute),
	}

	b := d1.ToBytes()
	d2 := blockDelayFromBytes(b)

	assert.Equal(t, d1.InitialTime.Unix(), d2.InitialTime.Unix())
	assert.Equal(t, d1.InitialHeight, d2.InitialHeight)
	assert.Equal(t, d1.PreviousBlockTime.Unix(), d2.PreviousBlockTime.Unix())
}

func TestBlockDelayNext_NoDelay(t *testing.T) {
	genesisTime := time.Now()
	initialHeight := int64(1)
	d := blockDelayUponGenesis(genesisTime, initialHeight)

	curBlockTime := genesisTime.Add(10 * time.Second)
	curBlockHeight := int64(2)
	targetBlockTime := 5 * time.Second
	delay := d.Next(curBlockTime, curBlockHeight, targetBlockTime)

	assert.Equal(t, 1*time.Microsecond, delay)
}

func TestBlockDelayNext_WithDelay(t *testing.T) {
	genesisTime := time.Now()
	initialHeight := int64(1)
	d := blockDelayUponGenesis(genesisTime, initialHeight)

	curBlockTime := genesisTime.Add(8 * time.Second)
	curBlockHeight := int64(3)
	targetBlockTime := 5 * time.Second
	delay := d.Next(curBlockTime, curBlockHeight, targetBlockTime)

	assert.Equal(t, delay, 2*time.Second)
}

func TestBlockDelayNext_ResetOnStall(t *testing.T) {
	genesisTime := time.Now()
	initialHeight := int64(1)
	d := blockDelayUponGenesis(genesisTime, initialHeight)

	curBlockTime := genesisTime.Add(maxDelayBetweenBlocks + 1*time.Minute)
	curBlockHeight := int64(10)
	targetBlockTime := 5 * time.Second

	delay := d.Next(curBlockTime, curBlockHeight, targetBlockTime)

	assert.Equal(t, d.InitialTime, curBlockTime)
	assert.Equal(t, d.InitialHeight, curBlockHeight-1)
	assert.Equal(t, delay, targetBlockTime)
}

func TestBlockDelaySerialization(t *testing.T) {
	d := &blockDelay{
		InitialTime:       time.Now(),
		InitialHeight:     10,
		PreviousBlockTime: time.Now().Add(-1 * time.Minute),
	}

	b := d.ToBytes()
	var d2 blockDelay
	err := json.Unmarshal(b, &d2)
	assert.NoError(t, err)
	assert.Equal(t, d.InitialTime.Unix(), d2.InitialTime.Unix())
	assert.Equal(t, d.InitialHeight, d2.InitialHeight)
	assert.Equal(t, d.PreviousBlockTime.Unix(), d2.PreviousBlockTime.Unix())
}

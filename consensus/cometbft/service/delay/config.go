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

package delay

import (
	"math"
	"time"
)

const (
	// maxDelayBetweenBlocks is the maximum delay between two consecutive blocks.
	// If the last block time minus the previous block time is greater than
	// maxDelayBetweenBlocks, then we reset `FinalizeBlockResponse.NextBlockDelay`
	// to default.
	//
	// This is needed because the network may stall for a long time and we don't
	// want to rush in new blocks as the network resumes its operation.
	maxDelayBetweenBlocks = 5 * time.Minute

	// targetBlockTime is the desired block time.
	//
	// Note that it CAN'T be lower than the minimal (floor) block time in the
	// network, which is comprised of the time to a) propose a new block b)
	// gather 2/3+ prevotes c) gather 2/3+ precommits.
	targetBlockTime = 2 * time.Second

	// Delay to use before the upgrade to SBT.
	constBlockDelay = 500 * time.Millisecond

	// Height to enable SBT. Changes will be applied in the next block after this
	sbtEnableHeight = math.MaxInt64

	// Height at which consensus params are upgraded to use SBT
	sbtConsensusParamUpdate = 0

	// Until `timeout_commit` is removed from the CometBFT config,
	// `FinalizeBlockResponse.NextBlockDelay` can't be exactly 0. If it's set to
	// 0, then `timeout_commit` from the config will be used, which is not what
	// we want since we're trying to control the block time.
	noDelay = 1 * time.Microsecond
)

type ConfigGetter interface {
	SbtMaxBlockDelay() time.Duration
	SbtTargetBlockTime() time.Duration
	SbtConstBlockDelay() time.Duration
	SbtConsensusUpdateHeight() int64
	SbtConsensusEnableHeight() int64
}

type Config struct {
	MaxBlockDelay   time.Duration `mapstructure:"max-block-delay"`
	TargetBlockTime time.Duration `mapstructure:"target-block-time"`
	ConstBlockDelay time.Duration `mapstructure:"const-block-delay"`

	// ConsensusParamUpdate decides the height where we set the ConsensusEnableHeight
	// In turn ConsensusEnableHeight, initializes the BlockDelay struct which will then
	// be used to calculate block delay in the next height.
	ConsensusUpdateHeight int64 `mapstructure:"consensus-update-height"`
	ConsensusEnableHeight int64 `mapstructure:"consensus-enable-height"`
}

func DefaultConfig() Config {
	return Config{
		MaxBlockDelay:         maxDelayBetweenBlocks,
		TargetBlockTime:       targetBlockTime,
		ConstBlockDelay:       constBlockDelay,
		ConsensusUpdateHeight: sbtConsensusParamUpdate,
		ConsensusEnableHeight: sbtEnableHeight,
	}
}

func (c Config) SbtMaxBlockDelay() time.Duration {
	return c.MaxBlockDelay
}
func (c Config) SbtTargetBlockTime() time.Duration {
	return c.TargetBlockTime
}
func (c Config) SbtConstBlockDelay() time.Duration {
	return c.ConstBlockDelay
}
func (c Config) SbtConsensusUpdateHeight() int64 {
	return c.ConsensusUpdateHeight
}
func (c Config) SbtConsensusEnableHeight() int64 {
	return c.ConsensusEnableHeight
}

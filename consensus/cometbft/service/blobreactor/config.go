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

import "math"

const (
	// Height at which consensus params are upgraded to use blobreactor
	blobConsensusUpdateHeight int64 = 0

	// Height to enable blobreactor
	blobConsensusEnableHeight int64 = math.MaxInt64

	// maxBytes is the maximum size of blob data in bytes for BlobReactor consensus params (e.g., 800KB).
	blobMaxBytes int64 = 800 * 1024 // 800KB
)

// ConfigGetter provides read access to BlobReactor configuration.
type ConfigGetter interface {
	// BlobConsensusUpdateHeight returns the height at which BlobReactor consensus params are updated.
	// This is when the parameters are set but not yet active.
	BlobConsensusUpdateHeight() int64
	// BlobConsensusEnableHeight returns the height when P2P blob distribution via BlobReactor is enabled.
	// A value of 0 means the BlobReactor is disabled.
	BlobConsensusEnableHeight() int64
	// BlobMaxBytes returns the maximum size of blob data in bytes for BlobReactor consensus params.
	BlobMaxBytes() int64
}

// Config contains configuration for the P2P BlobReactor component.
type Config struct {
	// ConsensusUpdateHeight is the height at which BlobReactor consensus params are updated.
	// This is when the parameters are set but not yet active.
	ConsensusUpdateHeight int64 `mapstructure:"consensus-update-height"`
	// ConsensusEnableHeight is the height when P2P blob distribution via BlobReactor is enabled.
	// A value of 0 means the BlobReactor is disabled.
	ConsensusEnableHeight int64 `mapstructure:"consensus-enable-height"`
	// MaxBytes is the maximum size of blob data in bytes for BlobReactor consensus params (e.g., 800KB).
	MaxBytes int64 `mapstructure:"max-bytes"`
}

// BlobConsensusEnableHeight returns the height when blob processing is enabled.
func (c Config) BlobConsensusEnableHeight() int64 {
	return c.ConsensusEnableHeight
}

// BlobMaxBytes returns the maximum blob size in bytes.
func (c Config) BlobMaxBytes() int64 {
	return c.MaxBytes
}

func DefaultConfig() Config {
	return Config{
		ConsensusUpdateHeight: blobConsensusUpdateHeight,
		ConsensusEnableHeight: blobConsensusEnableHeight,
		MaxBytes:              blobMaxBytes,
	}
}

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

// Package blobconsensus holds the chain-spec parameter gating the transition
// from carrying blob sidecars as the second consensus transaction to
// distributing them outside CometBFT via the blob reactor.
//
// Unlike stable block time (SBT), this transition does not change any CometBFT
// consensus parameters, so a single enable height is sufficient. Below the
// enable height every proposal carries two txs (block, sidecars); at and above
// it proposals carry a single tx (block) and sidecars travel on the blob
// reactor p2p channel.
package blobconsensus

// disabledEnableHeight disables the transition; blobs keep riding as the second consensus tx.
const disabledEnableHeight int64 = 0

// ConfigGetter exposes the blob consensus activation parameters.
type ConfigGetter interface {
	// BlobConsensusEnableHeight returns the height at which blob sidecars stop being carried as a consensus tx. A value of 0 means disabled.
	BlobConsensusEnableHeight() int64
	// IsBlobConsensusEnabled returns true if, at the given height, blob sidecars are distributed via the blob reactor instead of a consensus
	// tx.
	IsBlobConsensusEnabled(height int64) bool
}

// Config is the blob consensus section of the chain spec.
type Config struct {
	// EnableHeight is the height at which blob sidecars stop being carried as the second consensus tx. A value of 0 means disabled.
	EnableHeight int64 `mapstructure:"enable-height"`
}

func DefaultConfig() Config {
	return Config{
		EnableHeight: disabledEnableHeight,
	}
}

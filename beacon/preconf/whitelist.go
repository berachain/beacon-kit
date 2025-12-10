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

package preconf

import (
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/crypto"
)

// Whitelist defines the interface for checking if a validator is whitelisted
// for preconfirmation support.
type Whitelist interface {
	// IsWhitelisted returns true if the given public key is in the whitelist.
	IsWhitelisted(pubkey crypto.BLSPubkey) bool
}

// whitelist is the concrete implementation of Whitelist.
type whitelist struct {
	validators map[crypto.BLSPubkey]struct{}
}

// NewWhitelist creates a new Whitelist from a slice of public keys.
func NewWhitelist(pubkeys []crypto.BLSPubkey, logger *phuslu.Logger) Whitelist {
	validators := make(map[crypto.BLSPubkey]struct{}, len(pubkeys))
	for _, pk := range pubkeys {
		validators[pk] = struct{}{}
	}

	if logger != nil {
		logger.Info("Preconf whitelist loaded", "count", len(validators))
	}

	return &whitelist{
		validators: validators,
	}
}

// IsWhitelisted returns true if the given public key is in the whitelist.
func (w *whitelist) IsWhitelisted(pubkey crypto.BLSPubkey) bool {
	if w == nil || w.validators == nil {
		return false
	}
	_, ok := w.validators[pubkey]
	return ok
}

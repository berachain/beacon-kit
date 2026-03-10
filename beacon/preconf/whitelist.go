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
	"sync/atomic"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/crypto"
)

// Whitelist defines the interface for checking if a validator is whitelisted
// for preconfirmation support.
type Whitelist interface {
	// IsWhitelisted returns true if the given public key is in the whitelist.
	IsWhitelisted(pubkey crypto.BLSPubkey) bool
}

type validatorSet map[crypto.BLSPubkey]struct{}

// reloadableWhitelist is the concrete implementation of Whitelist with hot-reload support.
// Allows reloading the whitelist from source at runtime.
type reloadableWhitelist struct {
	path    string
	current atomic.Pointer[validatorSet]
}

// Len returns the number of validators in the current whitelist.
// Returns 0 if the whitelist is not loaded.
func (r *reloadableWhitelist) Len() int {
	w := r.current.Load()
	if w == nil {
		return 0
	}
	return len(*w)
}

// NewWhitelist creates a new reloadableWhitelist with the given path
// and logger, also loading the initial whitelist from the path.
//
//nolint:revive // reloadableWhitelist is intentionally unexported; callers use the Whitelist interface
func NewWhitelist(path string) (*reloadableWhitelist, error) {
	pubkeys, err := LoadWhitelist(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load preconf whitelist from: %s", path)
	}

	validators := pubkeysToSet(pubkeys)

	r := &reloadableWhitelist{path: path}
	r.current.Store(validators)

	return r, nil
}

// IsWhitelisted checks if the given public key is in the current whitelist.
func (r *reloadableWhitelist) IsWhitelisted(pubkey crypto.BLSPubkey) bool {
	w := r.current.Load()
	if w == nil {
		return false
	}
	_, ok := (*w)[pubkey]
	return ok
}

// Reload reloads the whitelist from the stored file path. If loading fails
// or the resulting set is empty, it returns an error and keeps the existing whitelist.
func (r *reloadableWhitelist) Reload() error {
	pubkeys, err := LoadWhitelist(r.path)
	if err != nil {
		return errors.Wrapf(err, "failed to load preconf whitelist from: %s", r.path)
	}
	if len(pubkeys) == 0 {
		return errors.New("reloaded preconf whitelist is empty")
	}
	r.current.Store(pubkeysToSet(pubkeys))

	return nil
}

func pubkeysToSet(pubkeys []crypto.BLSPubkey) *validatorSet {
	m := make(validatorSet, len(pubkeys))
	for _, pk := range pubkeys {
		m[pk] = struct{}{}
	}
	return &m
}

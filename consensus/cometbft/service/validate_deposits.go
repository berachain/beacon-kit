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
	"encoding/hex"
	"fmt"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
)

// validateDeposits validates the provided deposits.
func validateDeposits(deposits []types.Deposit) error {
	if len(deposits) == 0 {
		return errors.New("at least one deposit is required")
	}

	seenPubkeys := make(map[string]bool)

	for i, deposit := range deposits {
		depositIndex := deposit.GetIndex()
		//#nosec:G701 // realistically fine in practice.
		// Validate index matches position
		if depositIndex.Unwrap() != uint64(i) {
			return fmt.Errorf(
				"deposit index %d does not match position %d",
				depositIndex,
				i,
			)
		}

		if isZeroBytes(deposit.Pubkey[:]) {
			return fmt.Errorf("deposit %d has a zeroed public key", i)
		}
		// Check for duplicate pubkeys
		pubkeyHex := hex.EncodeToString(deposit.Pubkey[:])
		if seenPubkeys[pubkeyHex] {
			return fmt.Errorf("duplicate pubkey found in deposit %d", i)
		}
		seenPubkeys[pubkeyHex] = true

		if isZeroBytes(deposit.Credentials[:]) {
			return fmt.Errorf(
				"deposit %d has zeroed withdrawal credentials",
				i,
			)
		}

		if deposit.Amount == 0 {
			return fmt.Errorf("deposit %d has zero amount", i)
		}

		if isZeroBytes(deposit.Signature[:]) {
			return fmt.Errorf("deposit %d has a zeroed signature", i)
		}
	}

	return nil
}

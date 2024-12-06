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
	"strings"

	"github.com/berachain/beacon-kit/primitives/common"
)

const expectedHexLength = 8

// isValidForkVersion returns true if the provided fork version is valid.
// A valid fork version must:
// - Start with "0x"
// - Be followed by exactly 8 hexadecimal characters.
func isValidForkVersion(forkVersion common.Version) bool {
	forkVersionStr := forkVersion.String()
	if !strings.HasPrefix(forkVersionStr, "0x") {
		return false
	}

	// Remove "0x" prefix and verify remaining characters
	hexPart := strings.TrimPrefix(forkVersionStr, "0x")

	// Should have exactly 8 characters after 0x prefix
	if len(hexPart) != expectedHexLength {
		return false
	}

	// Verify it's a valid hex number
	_, err := hex.DecodeString(hexPart)
	return err == nil
}

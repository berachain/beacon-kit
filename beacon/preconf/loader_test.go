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

package preconf_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/stretchr/testify/require"
)

func TestLoadWhitelist_ValidFile(t *testing.T) {
	t.Parallel()

	// Create temp file with valid JSON whitelist
	content := `[
		"0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a",
		"0xa572cbea904d67468808c8eb50a9450c9721db309128012543902d0ac358a62ae28f75bb8f1c7c42c39a8c5529bf0f4e"
	]`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "whitelist.json")
	err := os.WriteFile(tmpFile, []byte(content), 0o644)
	require.NoError(t, err)

	pubkeys, err := preconf.LoadWhitelist(tmpFile)
	require.NoError(t, err)
	require.Len(t, pubkeys, 2)

	// Verify loaded pubkeys work with whitelist
	w := preconf.NewWhitelist(pubkeys, nil)
	require.True(t, w.IsWhitelisted(pubkeys[0]))
	require.True(t, w.IsWhitelisted(pubkeys[1]))

	// Verify non-whitelisted key returns false
	notWhitelisted := crypto.BLSPubkey{}
	require.False(t, w.IsWhitelisted(notWhitelisted))
}

func TestLoadWhitelist_InvalidPubkey(t *testing.T) {
	t.Parallel()

	content := `["0xinvalid"]`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "whitelist.json")
	err := os.WriteFile(tmpFile, []byte(content), 0o644)
	require.NoError(t, err)

	_, err = preconf.LoadWhitelist(tmpFile)
	require.Error(t, err)
}

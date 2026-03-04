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
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/beacon/preconf"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/stretchr/testify/require"
)

func TestWhitelist_Reload_UpdatesKeys(t *testing.T) {
	t.Parallel()

	pkA, err := parser.ConvertPubkey(pubkeyAHex)
	require.NoError(t, err)
	pkB, err := parser.ConvertPubkey(pubkeyBHex)
	require.NoError(t, err)

	tmpFile := filepath.Join(t.TempDir(), "whitelist.json")
	writeWhitelistFile(t, tmpFile, pubkeyAHex)

	wl, err := preconf.NewWhitelist(tmpFile)
	require.NoError(t, err)
	require.Equal(t, 1, wl.Len())
	require.True(t, wl.IsWhitelisted(pkA))
	require.False(t, wl.IsWhitelisted(pkB))

	// Add key B to the file and reload.
	writeWhitelistFile(t, tmpFile, pubkeyAHex, pubkeyBHex)
	require.NoError(t, wl.Reload())
	require.Equal(t, 2, wl.Len())
	require.True(t, wl.IsWhitelisted(pkA))
	require.True(t, wl.IsWhitelisted(pkB))

	// Remove key A from the file and reload.
	writeWhitelistFile(t, tmpFile, pubkeyBHex)
	require.NoError(t, wl.Reload())
	require.Equal(t, 1, wl.Len())
	require.False(t, wl.IsWhitelisted(pkA))
	require.True(t, wl.IsWhitelisted(pkB))
}

func TestWhitelist_Reload_KeepsExistingOnError(t *testing.T) {
	t.Parallel()

	pkA, err := parser.ConvertPubkey(pubkeyAHex)
	require.NoError(t, err)

	tmpFile := filepath.Join(t.TempDir(), "whitelist.json")
	writeWhitelistFile(t, tmpFile, pubkeyAHex)

	wl, err := preconf.NewWhitelist(tmpFile)
	require.NoError(t, err)
	require.True(t, wl.IsWhitelisted(pkA))

	// Corrupt the file - reload should fail but preserve the existing set.
	err = os.WriteFile(tmpFile, []byte("not valid json"), 0o644)
	require.NoError(t, err)

	require.Error(t, wl.Reload())
	require.Equal(t, 1, wl.Len())
	require.True(t, wl.IsWhitelisted(pkA))

	// write valid file with 0 whitelisted validators - reload should fail because whitelist cannot be empty.
	writeWhitelistFile(t, tmpFile)
	require.Error(t, wl.Reload())
	require.Equal(t, 1, wl.Len())
	require.True(t, wl.IsWhitelisted(pkA))
}

func writeWhitelistFile(t *testing.T, path string, hexKeys ...string) {
	t.Helper()
	content, err := json.Marshal(hexKeys)
	require.NoError(t, err)
	err = os.WriteFile(path, content, 0o644)
	require.NoError(t, err)
}

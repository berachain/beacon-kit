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

package node_test

import (
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/node-api/handlers/node"
	"github.com/berachain/beacon-kit/node-api/handlers/node/mocks"
	"github.com/berachain/beacon-kit/node-api/handlers/node/types"
	"github.com/stretchr/testify/require"
)

func TestNodeVersion(t *testing.T) {
	t.Parallel()

	// setup test
	backend := mocks.NewBackend(t)
	h := node.NewHandler(backend)

	var (
		appName = "testing-beacond"
		version = "x.x.x-rc-14-gebeee824a"
		os      = "theOs"
		arch    = "theArch"
	)
	backend.EXPECT().GetVersionData().Return(appName, version, os, arch).Once()

	// test
	res, err := h.Version(nil) // input does not matter here

	// checks
	require.NoError(t, err)
	require.NotNil(t, res)
	require.IsType(t, types.DataResponse{}, res)
	dr, _ := res.(types.DataResponse)
	require.IsType(t, &types.VersionData{}, dr.Data)
	data, _ := dr.Data.(*types.VersionData)

	require.True(t, strings.Contains(data.Version, appName))
	require.True(t, strings.Contains(data.Version, version))
	require.True(t, strings.Contains(data.Version, os))
	require.True(t, strings.Contains(data.Version, arch))
}

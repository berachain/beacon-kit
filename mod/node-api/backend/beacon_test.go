// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package backend_test

import (
	"context"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/backend/mocks"
	response "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/stretchr/testify/require"
)

func TestGetGenesisDetails(t *testing.T) {
	sdb := &mocks.StateDB{}
	b := backend.NewMockBackend()
	expected := &response.GenesisData{
		GenesisTime:           0,
		GenesisValidatorsRoot: primitives.Root{0x01},
		GenesisForkVersion:    primitives.Version{0x01},
	}
	sdb.EXPECT().GetGenesisDetails().Return(expected, nil)
	actual, err := b.GetGenesis(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestGetBlockHeader(t *testing.T) {
	bdb := &mocks.BlockDB{}
	b := backend.NewMockBackend()
	expected := &response.BlockHeaderData{
		Root:      primitives.Root{0x01},
		Canonical: true,
		Header: response.MessageResponse{
			Message: types.NewBeaconBlockHeader(
				0,
				0,
				primitives.Root{0x01},
				primitives.Root{0x01},
				primitives.Root{0x01},
			),
		},
		Signature: crypto.BLSSignature{0x01},
	}
	bdb.EXPECT().GetBlockHeader().Return(expected, nil)
	actual, err := b.GetBlockHeader(context.Background(), "0")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

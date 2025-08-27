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

package beacon_test

import (
	"errors"
	"testing"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	beaconecho "github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetGenesis(t *testing.T) {
	t.Parallel()

	var (
		testGenesisRoot        = common.Root{0x10, 0x20, 0x30}
		testGenesisForkVersion = common.Version{0xff, 0x11, 0x22, 0x33}
		testGenesisTime        = math.U64(123456789)

		errTestError = errors.New("test error when missing item")
	)

	testCases := []struct {
		name                string
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "sunny path",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GenesisValidatorsRoot().Return(testGenesisRoot, nil).Once()
				b.EXPECT().GenesisForkVersion().Return(testGenesisForkVersion, nil).Once()
				b.EXPECT().GenesisTime().Return(testGenesisTime, nil).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenesisResponse{}, res)
				gr, _ := res.(beacontypes.GenesisResponse)

				require.Equal(t, testGenesisRoot, gr.Data.GenesisValidatorsRoot)
				require.Equal(t, testGenesisForkVersion.String(), gr.Data.GenesisForkVersion)
				require.Equal(t, testGenesisTime.Base10(), gr.Data.GenesisTime)
			},
		},
		{
			name: "genesis not ready",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GenesisValidatorsRoot().Return(common.Root{}, nil).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.ErrorIs(t, err, types.ErrNotFound)
				require.Contains(t, err.Error(), "Chain genesis info is not yet known")
				require.Nil(t, res)
			},
		},
		{
			name: "failed loading validator root",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GenesisValidatorsRoot().Return(common.Root{}, errTestError).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.ErrorIs(t, err, errTestError)
				require.Nil(t, res)
			},
		},
		{
			name: "failed loading fork version",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GenesisValidatorsRoot().Return(testGenesisRoot, nil).Once()
				b.EXPECT().GenesisForkVersion().Return(common.Version{}, errTestError).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.ErrorIs(t, err, errTestError)
				require.Nil(t, res)
			},
		},
		{
			name: "failed loading genesis time",
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GenesisValidatorsRoot().Return(testGenesisRoot, nil).Once()
				b.EXPECT().GenesisForkVersion().Return(testGenesisForkVersion, nil).Once()
				b.EXPECT().GenesisTime().Return(0, errTestError).Once()
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.ErrorIs(t, err, errTestError)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// setup test
			backend := mocks.NewBackend(t)
			h := beacon.NewHandler(backend)
			h.SetLogger(noop.NewLogger[log.Logger]())
			e := echo.New()
			e.Validator = &beaconecho.CustomValidator{
				Validator: beaconecho.ConstructValidator(),
			}

			// set expectations
			tc.setMockExpectations(backend)

			// test
			res, err := h.GetGenesis(nil) // input not relevant for this API

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

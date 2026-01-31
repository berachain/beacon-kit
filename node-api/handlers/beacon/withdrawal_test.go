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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/mapping"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/middleware"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPendingPartialWithdrawals(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	preElectraFork := &ctypes.Fork{
		PreviousVersion: version.Capella(),
		CurrentVersion:  version.Deneb(),
		Epoch:           math.Epoch(200),
	}
	electraFork := &ctypes.Fork{
		PreviousVersion: version.Deneb(),
		CurrentVersion:  version.Electra(),
		Epoch:           math.Epoch(200),
	}
	testPendingPartialWithdrawals := []*ctypes.PendingPartialWithdrawal{
		{
			ValidatorIndex:    0,
			Amount:            math.Gwei(1_000_000),
			WithdrawableEpoch: math.Epoch(1),
		},
		{
			ValidatorIndex:    1,
			Amount:            math.Gwei(999_999_999),
			WithdrawableEpoch: math.Epoch(2),
		},
	}

	testCases := []struct {
		name                string
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "post-Electra withdrawal requests present",
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)

				require.NoError(t, st.SetFork(electraFork))
				require.NoError(t, st.SetPendingPartialWithdrawals(testPendingPartialWithdrawals))

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.PendingPartialWithdrawalsResponse{}, res)
				resp, _ := res.(beacontypes.PendingPartialWithdrawalsResponse)

				require.Equal(t, version.Name(electraFork.CurrentVersion), resp.Version)

				require.IsType(t, []*beacontypes.PendingPartialWithdrawalData{}, resp.GenericResponse.Data)
				data, _ := resp.GenericResponse.Data.([]*beacontypes.PendingPartialWithdrawalData)

				require.Len(t, data, len(testPendingPartialWithdrawals))
				for i, w := range testPendingPartialWithdrawals {
					require.Equal(t, math.ValidatorIndex(data[i].ValidatorIndex), w.ValidatorIndex)
					require.Equal(t, math.Gwei(data[i].Amount), w.Amount)
					require.Equal(t, math.Epoch(data[i].WithdrawalEpoch), w.WithdrawableEpoch)
				}
			},
		},
		{
			name: "post-Electra withdrawal no requests",
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)

				require.NoError(t, st.SetFork(electraFork))
				require.NoError(t, st.SetPendingPartialWithdrawals(nil))

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.PendingPartialWithdrawalsResponse{}, res)
				resp, _ := res.(beacontypes.PendingPartialWithdrawalsResponse)

				require.Equal(t, version.Name(electraFork.CurrentVersion), resp.Version)

				require.IsType(t, []*beacontypes.PendingPartialWithdrawalData{}, resp.GenericResponse.Data)
				data, _ := resp.GenericResponse.Data.([]*beacontypes.PendingPartialWithdrawalData)

				require.Empty(t, data)
			},
		},
		{
			name: "pre-Electra - error",
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)

				require.NoError(t, st.SetFork(preElectraFork))

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.ErrorIs(t, err, middleware.ErrInvalidRequest)
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup test
			backend := mocks.NewBackend(t)
			h := beacon.NewHandler(backend, cs, noop.NewLogger[log.Logger]())
			e := echo.New()
			e.Validator = &middleware.CustomValidator{
				Validator: middleware.ConstructValidator(),
			}

			// set expectations
			tc.setMockExpectations(backend)

			// create input
			input := beacontypes.GetPendingPartialWithdrawalsRequest{
				StateIDRequest: handlertypes.StateIDRequest{
					StateID: mapping.StateIDGenesis,
				},
			}
			inputBytes, err := json.Marshal(input) //nolint:musttag //  TODO:fix
			require.NoError(t, err)
			body := strings.NewReader(string(inputBytes))
			req := httptest.NewRequest(http.MethodGet, "/", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // otherwise code=415, message=Unsupported Media Type
			c := e.NewContext(req, httptest.NewRecorder())

			// test
			res, err := h.GetPendingPartialWithdrawals(c)

			// check
			tc.check(t, res, err)
		})
	}
}

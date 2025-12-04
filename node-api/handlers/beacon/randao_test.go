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
	"strconv"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/node-api/middleware"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetRandao(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	var (
		testSlot  = math.Slot(1234)
		testEpoch = cs.SlotToEpoch(testSlot)
		index     = testEpoch.Unwrap() % cs.EpochsPerHistoricalVector()
		testMix   = common.Bytes32{0x01, 0x02, 0x03}
	)

	testCases := []struct {
		name                string
		inputs              func() beacontypes.GetRandaoRequest
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "randao request - specified epoch",
			inputs: func() beacontypes.GetRandaoRequest {
				stateID := utils.StateIDHead
				epoch := strconv.Itoa(int(testEpoch.Unwrap()))
				return beacontypes.GetRandaoRequest{
					StateIDRequest: handlertypes.StateIDRequest{
						StateID: stateID,
					},
					EpochOptionalRequest: beacontypes.EpochOptionalRequest{
						Epoch: epoch,
					},
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)

				require.NoError(t, st.UpdateRandaoMixAtIndex(index, testMix))
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, testSlot, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, common.Bytes32{}, gr.Data)
				data, _ := gr.Data.(common.Bytes32)
				require.Equal(t, testMix, data)
			},
		},
		{
			name: "randao request - unspecified epoch",
			inputs: func() beacontypes.GetRandaoRequest {
				stateID := utils.StateIDHead
				return beacontypes.GetRandaoRequest{
					StateIDRequest: handlertypes.StateIDRequest{
						StateID: stateID,
					},
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)

				require.NoError(t, st.UpdateRandaoMixAtIndex(index, testMix))
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, testSlot, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, common.Bytes32{}, gr.Data)
				data, _ := gr.Data.(common.Bytes32)
				require.Equal(t, testMix, data)
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
			input := tc.inputs()
			inputBytes, err := json.Marshal(input) //nolint:musttag //  TODO:fix
			require.NoError(t, err)
			body := strings.NewReader(string(inputBytes))
			req := httptest.NewRequest(http.MethodGet, "/", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // otherwise code=415, message=Unsupported Media Type
			c := e.NewContext(req, httptest.NewRecorder())

			// test
			res, err := h.GetRandao(c)

			// check
			tc.check(t, res, err)
		})
	}
}

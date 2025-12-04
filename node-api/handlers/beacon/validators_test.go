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
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetValidator(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	// Create some input validators and store them to a readonly state
	stateValidators := createStateValidators(cs)

	testCases := []struct {
		name                string
		inputs              func() beacontypes.GetStateValidatorRequest
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "get existing validator",
			inputs: func() beacontypes.GetStateValidatorRequest {
				return beacontypes.GetStateValidatorRequest{
					StateIDRequest: handlertypes.StateIDRequest{
						StateID: utils.StateIDHead,
					},
					ValidatorID: stateValidators[0].Validator.PublicKey,
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, &beacontypes.ValidatorData{}, gr.Data)
				data, _ := gr.Data.(*beacontypes.ValidatorData)

				require.Equal(t, stateValidators[0], data)
			},
		},
		{
			name: "get unknown validator - 1",
			inputs: func() beacontypes.GetStateValidatorRequest {
				unknownValIdx := strconv.Itoa(2025)
				return beacontypes.GetStateValidatorRequest{
					StateIDRequest: handlertypes.StateIDRequest{
						StateID: utils.StateIDHead,
					},
					ValidatorID: unknownValIdx,
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.ErrorIs(t, err, handlertypes.ErrNotFound)
				require.Nil(t, res)
			},
		},
		{
			name: "get unknown validator - 2",
			inputs: func() beacontypes.GetStateValidatorRequest {
				unknownValPk := bytes.B48{0xff, 0xff}
				return beacontypes.GetStateValidatorRequest{
					StateIDRequest: handlertypes.StateIDRequest{
						StateID: utils.StateIDHead,
					},
					ValidatorID: unknownValPk.String(),
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()
				require.ErrorIs(t, err, handlertypes.ErrNotFound)
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

			// create API inputs
			input := tc.inputs()
			inputBytes, err := json.Marshal(input) //nolint:musttag //  TODO:fix
			require.NoError(t, err)
			body := strings.NewReader(string(inputBytes))
			req := httptest.NewRequest(http.MethodGet, "/", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // otherwise code=415, message=Unsupported Media Type
			c := e.NewContext(req, httptest.NewRecorder())

			// set expectations
			tc.setMockExpectations(backend)

			// test
			res, err := h.GetStateValidator(c)

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

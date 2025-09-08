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
	"strconv"
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/middleware"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetStateValidatorBalances(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	// Create some input validators and store them to a readonly state
	stateValidators := createStateValidators(cs)

	testCases := []struct {
		name                string
		inputs              func() (math.Slot, []string)
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res []*beacontypes.ValidatorBalanceData, err error)
	}{
		{
			name: "all validators",
			inputs: func() (math.Slot, []string) {
				return 0, nil
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAtSlot(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorBalanceData, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)

				expectedBalances := make([]*beacontypes.ValidatorBalanceData, 0, len(stateValidators))
				for _, val := range stateValidators {
					expectedBalances = append(expectedBalances,
						&beacontypes.ValidatorBalanceData{
							Index:   val.Index,
							Balance: val.Balance,
						},
					)
				}

				require.Equal(t, expectedBalances, res)
			},
		},
		{
			name: "single validators",
			inputs: func() (math.Slot, []string) {
				return 0, []string{
					stateValidators[0].Validator.PublicKey,
					stateValidators[1].Validator.PublicKey,
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAtSlot(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorBalanceData, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)

				expectedBalances := []*beacontypes.ValidatorBalanceData{
					{
						Index:   stateValidators[0].Index,
						Balance: stateValidators[0].Balance,
					},
					{
						Index:   stateValidators[1].Index,
						Balance: stateValidators[1].Balance,
					},
				}
				require.Equal(t, expectedBalances, res)
			},
		},
		{
			name: "mixed know and unknow validators",
			inputs: func() (math.Slot, []string) {
				unknownValPk := bytes.B48{0xff, 0xff}
				unknownValIdx := strconv.Itoa(2025)
				return 0, []string{
					unknownValIdx,
					stateValidators[0].Validator.PublicKey,
					unknownValPk.String(),
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				st := makeTestState(t, cs)
				addTestValidators(t, stateValidators, st)

				// slot is not really tested here, we just return zero
				b.EXPECT().StateAtSlot(mock.Anything).Return(st, math.Slot(0), nil)
			},
			check: func(t *testing.T, res []*beacontypes.ValidatorBalanceData, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)

				expectedBalances := []*beacontypes.ValidatorBalanceData{
					{
						Index:   stateValidators[0].Index,
						Balance: stateValidators[0].Balance,
					},
				}
				require.Equal(t, expectedBalances, res)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// setup test
			backend := mocks.NewBackend(t)
			h := beacon.NewHandler(backend, cs, noop.NewLogger[log.Logger]())
			e := echo.New()
			e.Validator = &middleware.CustomValidator{
				Validator: middleware.ConstructValidator(),
			}

			// set expectations
			tc.setMockExpectations(backend)

			// test
			res, err := h.GetValidatorBalance(tc.inputs())

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

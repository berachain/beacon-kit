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
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/mapping"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/middleware"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetBlobSidecars(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	testSidecars := datypes.BlobSidecars{
		{
			Index: 25,
			SignedBeaconBlockHeader: &ctypes.SignedBeaconBlockHeader{
				Header: &ctypes.BeaconBlockHeader{
					Slot:            math.Slot(987),
					ProposerIndex:   math.ValidatorIndex(2),
					ParentBlockRoot: common.Root{0x1, 0x2},
					StateRoot:       common.Root{0x3, 0x4},
					BodyRoot:        common.Root{0x5, 0x6},
				},
			},
		},
	}

	testCases := []struct {
		name                string
		inputs              func() beacontypes.GetBlobSidecarsRequest
		setMockExpectations func(*testing.T, *mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "success",
			inputs: func() beacontypes.GetBlobSidecarsRequest {
				return beacontypes.GetBlobSidecarsRequest{
					BlockIDRequest: handlertypes.BlockIDRequest{
						BlockID: mapping.StateIDHead,
					},
					Indices: nil,
				}
			},
			setMockExpectations: func(t *testing.T, b *mocks.Backend) {
				t.Helper()

				b.EXPECT().GetSyncData().Return(int64(1234), int64(1234))
				b.EXPECT().GetBlobSidecarsAtSlot(mock.Anything).Return(testSidecars, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.SidecarsResponse{}, res)
				sr, _ := res.(beacontypes.SidecarsResponse)

				require.Len(t, sr.Data, 1)
				require.Equal(t, beacontypes.SidecarFromConsensus(testSidecars[0]), sr.Data[0])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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
			tc.setMockExpectations(t, backend)

			// test
			res, err := h.GetBlobSidecars(c)

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

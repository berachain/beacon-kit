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
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	beaconecho "github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

//nolint:maintidx // multiple test cases
func TestGetBlockHeaders(t *testing.T) {
	t.Parallel()

	// testHeaders to build test cases on top of
	testParentHeader := &types.BeaconBlockHeader{
		Slot:            math.Slot(10),
		ProposerIndex:   math.ValidatorIndex(1234),
		ParentBlockRoot: common.Root{'p', 'a', 'r', 'e', 'n', 't', 'b', 'l', 'o', 'c', 'k', 'r', 'o', 'o', 't'},
		StateRoot:       common.Root{'p', 'a', 'r', 'e', 'n', 't', 's', 't', 'a', 't', 'e', 'r', 'o', 'o', 't'},
		BodyRoot:        common.Root{'p', 'a', 'r', 'e', 'n', 't', 'r', 'o', 'o', 't'},
	}
	testHeader := &types.BeaconBlockHeader{
		Slot:            testParentHeader.Slot + 1,
		ProposerIndex:   math.ValidatorIndex(5678),
		ParentBlockRoot: testParentHeader.BodyRoot,
		StateRoot:       common.Root{'d', 'u', 'm', 'm', 'y', 's', 't', 'a', 't', 'e', 'r', 'o', 'o', 't'},
		BodyRoot:        common.Root{'d', 'u', 'm', 'm', 'y', 'r', 'o', 'o', 't'},
	}
	wrongSlot := testHeader.Slot + 1234
	errTestHeaderNotFound := errors.New("test header not found error")

	testCases := []struct {
		name                string
		inputs              func() beacontypes.GetBlockHeadersRequest
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "GetBlockHeaders - success - no query params",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{},
					ParentRoot:  "",
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				// math.Slot(0) is a special parametere as it indicates head should be picked.
				// For these tests we return testHeader as head
				b.EXPECT().BlockHeaderAtSlot(math.Slot(0)).Return(testHeader, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, []beacontypes.BlockHeaderResponse{}, gr.Data)
				data, _ := gr.Data.([]beacontypes.BlockHeaderResponse)
				require.Len(t, data, 1)

				require.Equal(t, testHeader.BodyRoot, data[0].Root)
				expectedHeader := &beacontypes.BeaconBlockHeader{
					Slot:          testHeader.Slot.Base10(),
					ProposerIndex: testHeader.ProposerIndex.Base10(),
					ParentRoot:    testHeader.ParentBlockRoot.Hex(),
					StateRoot:     testHeader.StateRoot.Hex(),
					BodyRoot:      testHeader.BodyRoot.Hex(),
				}
				require.Equal(t, expectedHeader, data[0].Header.Message)
			},
		},
		{
			name: "GetBlockHeaders - success - slot only",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{
						Slot: testHeader.Slot.Base10(),
					},
					ParentRoot: "",
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().BlockHeaderAtSlot(testHeader.Slot).Return(testHeader, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, []beacontypes.BlockHeaderResponse{}, gr.Data)
				data, _ := gr.Data.([]beacontypes.BlockHeaderResponse)
				require.Len(t, data, 1)

				require.Equal(t, testHeader.BodyRoot, data[0].Root)
				expectedHeader := &beacontypes.BeaconBlockHeader{
					Slot:          testHeader.Slot.Base10(),
					ProposerIndex: testHeader.ProposerIndex.Base10(),
					ParentRoot:    testHeader.ParentBlockRoot.Hex(),
					StateRoot:     testHeader.StateRoot.Hex(),
					BodyRoot:      testHeader.BodyRoot.Hex(),
				}
				require.Equal(t, expectedHeader, data[0].Header.Message)
			},
		},
		{
			name: "GetBlockHeaders - failure - invalid slot",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{
						Slot: "AAAA",
					},
					ParentRoot: "",
				}
			},
			setMockExpectations: func(*mocks.Backend) {
				// nothing to set here, slot is invalid
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, handlertypes.ErrInvalidRequest)
			},
		},
		{
			name: "GetBlockHeaders - failure - unindexed slot",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{
						Slot: testHeader.Slot.Base10(),
					},
					ParentRoot: "",
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().BlockHeaderAtSlot(testHeader.Slot).Return(nil, errTestHeaderNotFound)
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				// Implicitly ensuring that 404 error code is returned
				// (see responseFromError implementation)
				require.ErrorIs(t, err, handlertypes.ErrNotFound)
			},
		},
		{
			name: "GetBlockHeaders - success - parent root only",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{},
					ParentRoot:  testHeader.ParentBlockRoot.Hex(),
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(testParentHeader.Slot, nil)
				b.EXPECT().BlockHeaderAtSlot(testHeader.Slot).Return(testHeader, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, []beacontypes.BlockHeaderResponse{}, gr.Data)
				data, _ := gr.Data.([]beacontypes.BlockHeaderResponse)
				require.Len(t, data, 1)

				require.Equal(t, testHeader.BodyRoot, data[0].Root)
				expectedHeader := &beacontypes.BeaconBlockHeader{
					Slot:          testHeader.Slot.Base10(),
					ProposerIndex: testHeader.ProposerIndex.Base10(),
					ParentRoot:    testHeader.ParentBlockRoot.Hex(),
					StateRoot:     testHeader.StateRoot.Hex(),
					BodyRoot:      testHeader.BodyRoot.Hex(),
				}
				require.Equal(t, expectedHeader, data[0].Header.Message)
			},
		},
		{
			name: "GetBlockHeaders - failure - invalid parent root",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{},
					ParentRoot:  "AN-INVALID-HEX",
				}
			},
			setMockExpectations: func(*mocks.Backend) {
				// nothing to set here, slot is invalid
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, handlertypes.ErrInvalidRequest)
			},
		},
		{
			name: "GetBlockHeaders - failure - unindexed parent block",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{},
					ParentRoot:  testHeader.ParentBlockRoot.Hex(),
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				// Assume parent header is not known
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(0, errTestHeaderNotFound)
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				// Implicitly ensuring that 404 error code is returned
				// (see responseFromError implementation)
				require.ErrorIs(t, err, handlertypes.ErrNotFound)
			},
		},
		{
			name: "GetBlockHeaders - success - slot and parent root",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{
						Slot: testHeader.Slot.Base10(),
					},
					ParentRoot: testHeader.ParentBlockRoot.Hex(),
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(testParentHeader.Slot, nil)
				b.EXPECT().BlockHeaderAtSlot(testHeader.Slot).Return(testHeader, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, []beacontypes.BlockHeaderResponse{}, gr.Data)
				data, _ := gr.Data.([]beacontypes.BlockHeaderResponse)
				require.Len(t, data, 1)

				require.Equal(t, testHeader.BodyRoot, data[0].Root)
				expectedHeader := &beacontypes.BeaconBlockHeader{
					Slot:          testHeader.Slot.Base10(),
					ProposerIndex: testHeader.ProposerIndex.Base10(),
					ParentRoot:    testHeader.ParentBlockRoot.Hex(),
					StateRoot:     testHeader.StateRoot.Hex(),
					BodyRoot:      testHeader.BodyRoot.Hex(),
				}
				require.Equal(t, expectedHeader, data[0].Header.Message)
			},
		},
		{
			name: "GetBlockHeaders - failure - mismatched slot and parent root",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{
						Slot: wrongSlot.Base10(),
					},
					ParentRoot: testHeader.ParentBlockRoot.Hex(),
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(testParentHeader.Slot, nil)
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, beacon.ErrMismatchedSlotAndParentBlock)
			},
		},
	}

	for _, tc := range testCases {
		// capture range variable
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
			res, err := h.GetBlockHeaders(c)

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

func TestGetBlockHeaderByID(t *testing.T) {
	t.Parallel()

	// a test testHeader to build test cases on top of
	testHeader := &types.BeaconBlockHeader{
		Slot:            math.Slot(1234),
		ProposerIndex:   math.ValidatorIndex(5678),
		ParentBlockRoot: common.Root{'p', 'a', 'r', 'e', 'n', 't', 'b', 'l', 'o', 'c', 'k', 'r', 'o', 'o', 't'},
		StateRoot:       common.Root{'d', 'u', 'm', 'm', 'y', 's', 't', 'a', 't', 'e', 'r', 'o', 'o', 't'},
		BodyRoot:        common.Root{'d', 'u', 'm', 'm', 'y', 'r', 'o', 'o', 't'},
	}
	errTestHeaderNotFound := errors.New("test header not found error")

	testCases := []struct {
		name                string
		inputs              func() beacontypes.GetBlockHeaderRequest
		setMockExpectations func(*mocks.Backend)
		check               func(t *testing.T, res any, err error)
	}{
		{
			name: "GetBlockHeaderByID - success - id is slot",
			inputs: func() beacontypes.GetBlockHeaderRequest {
				return beacontypes.GetBlockHeaderRequest{
					BlockIDRequest: handlertypes.BlockIDRequest{
						BlockID: testHeader.Slot.Base10(),
					},
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().BlockHeaderAtSlot(testHeader.Slot).Return(testHeader, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, &beacontypes.BlockHeaderResponse{}, gr.Data)
				data, _ := gr.Data.(*beacontypes.BlockHeaderResponse)

				require.Equal(t, testHeader.BodyRoot, data.Root)
				expectedHeader := &beacontypes.BeaconBlockHeader{
					Slot:          testHeader.Slot.Base10(),
					ProposerIndex: testHeader.ProposerIndex.Base10(),
					ParentRoot:    testHeader.ParentBlockRoot.Hex(),
					StateRoot:     testHeader.StateRoot.Hex(),
					BodyRoot:      testHeader.BodyRoot.Hex(),
				}
				require.Equal(t, expectedHeader, data.Header.Message)
			},
		},
		{
			name: "GetBlockHeaderByID - failure - invalid slot",
			inputs: func() beacontypes.GetBlockHeaderRequest {
				return beacontypes.GetBlockHeaderRequest{
					BlockIDRequest: handlertypes.BlockIDRequest{
						BlockID: "zzzz",
					},
				}
			},
			setMockExpectations: func(*mocks.Backend) {
				// nothing to set here, slot is invalid
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, handlertypes.ErrInvalidRequest)
			},
		},
		{
			name: "GetBlockHeaderByID - failure - unindexed slot",
			inputs: func() beacontypes.GetBlockHeaderRequest {
				return beacontypes.GetBlockHeaderRequest{
					BlockIDRequest: handlertypes.BlockIDRequest{
						BlockID: testHeader.Slot.Base10(),
					},
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().BlockHeaderAtSlot(testHeader.Slot).Return(nil, errTestHeaderNotFound)
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()

				// Implicitly ensuring that 404 error code is returned
				// (see responseFromError implementation)
				require.ErrorIs(t, err, handlertypes.ErrNotFound)
			},
		},
		{
			name: "GetBlockHeaderByID - success - id is block root",
			inputs: func() beacontypes.GetBlockHeaderRequest {
				return beacontypes.GetBlockHeaderRequest{
					BlockIDRequest: handlertypes.BlockIDRequest{
						BlockID: testHeader.BodyRoot.Hex(),
					},
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				b.EXPECT().GetSlotByBlockRoot(testHeader.BodyRoot).Return(testHeader.Slot, nil)
				b.EXPECT().BlockHeaderAtSlot(testHeader.Slot).Return(testHeader, nil)
			},
			check: func(t *testing.T, res any, err error) {
				t.Helper()

				require.NoError(t, err)
				require.NotNil(t, res)
				require.IsType(t, beacontypes.GenericResponse{}, res)
				gr, _ := res.(beacontypes.GenericResponse)
				require.IsType(t, &beacontypes.BlockHeaderResponse{}, gr.Data)
				data, _ := gr.Data.(*beacontypes.BlockHeaderResponse)

				require.Equal(t, testHeader.BodyRoot, data.Root)
				expectedHeader := &beacontypes.BeaconBlockHeader{
					Slot:          testHeader.Slot.Base10(),
					ProposerIndex: testHeader.ProposerIndex.Base10(),
					ParentRoot:    testHeader.ParentBlockRoot.Hex(),
					StateRoot:     testHeader.StateRoot.Hex(),
					BodyRoot:      testHeader.BodyRoot.Hex(),
				}
				require.Equal(t, expectedHeader, data.Header.Message)
			},
		},
		{
			name: "GetBlockHeaderByID - failure - invalid block root",
			inputs: func() beacontypes.GetBlockHeaderRequest {
				return beacontypes.GetBlockHeaderRequest{
					BlockIDRequest: handlertypes.BlockIDRequest{
						BlockID: "AN-INVALID-HEX",
					},
				}
			},
			setMockExpectations: func(*mocks.Backend) {
				// nothing to set here, slot is invalid
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, handlertypes.ErrInvalidRequest)
			},
		},
		{
			name: "GetBlockHeaderByID - failure - unindexed block",
			inputs: func() beacontypes.GetBlockHeaderRequest {
				return beacontypes.GetBlockHeaderRequest{
					BlockIDRequest: handlertypes.BlockIDRequest{
						BlockID: testHeader.BodyRoot.Hex(),
					},
				}
			},
			setMockExpectations: func(b *mocks.Backend) {
				// Assume parent header is not known
				b.EXPECT().GetSlotByBlockRoot(testHeader.BodyRoot).Return(0, errTestHeaderNotFound)
			},
			check: func(t *testing.T, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, errTestHeaderNotFound)
			},
		},
	}

	for _, tc := range testCases {
		// capture range variable
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
			res, err := h.GetBlockHeaderByID(c)

			// finally do checks
			tc.check(t, res, err)
		})
	}
}

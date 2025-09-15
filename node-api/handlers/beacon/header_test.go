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

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/middleware"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//nolint:maintidx // multiple test cases
func TestGetBlockHeaders(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	// testHeaders to build test cases on top of
	testParentHeader := &ctypes.BeaconBlockHeader{
		Slot:            math.Slot(10),
		ProposerIndex:   math.ValidatorIndex(1234),
		ParentBlockRoot: common.Root{'p', 'a', 'r', 'e', 'n', 't', 'b', 'l', 'o', 'c', 'k', 'r', 'o', 'o', 't'},
		StateRoot:       common.Root{'p', 'a', 'r', 'e', 'n', 't', 's', 't', 'a', 't', 'e', 'r', 'o', 'o', 't'},
		BodyRoot:        common.Root{'p', 'a', 'r', 'e', 'n', 't', 'r', 'o', 'o', 't'},
	}
	testHeader := &ctypes.BeaconBlockHeader{
		Slot:            testParentHeader.Slot + 1,
		ProposerIndex:   math.ValidatorIndex(5678),
		ParentBlockRoot: testParentHeader.BodyRoot,
		StateRoot:       common.Root{}, // set in test cases
		BodyRoot:        common.Root{'d', 'u', 'm', 'm', 'y', 'r', 'o', 'o', 't'},
	}
	wrongSlot := testHeader.Slot + 1234
	errTestHeaderNotFound := errors.New("test header not found error")

	testCases := []struct {
		name                string
		inputs              func() beacontypes.GetBlockHeadersRequest
		setMockExpectations func(*testing.T, *mocks.Backend) common.Root
		check               func(t *testing.T, expectedStateRoot common.Root, res any, err error)
	}{
		{
			name: "GetBlockHeaders - success - no query params",
			inputs: func() beacontypes.GetBlockHeadersRequest {
				return beacontypes.GetBlockHeadersRequest{
					SlotRequest: beacontypes.SlotRequest{},
					ParentRoot:  "",
				}
			},
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()

				st := makeTestState(t, cs)
				stateRoot := testDummyState(t, cs, st, testHeader)
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
				return stateRoot
			},
			check: func(t *testing.T, expectedStateRoot common.Root, res any, err error) {
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
					StateRoot:     expectedStateRoot.Hex(),
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
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()

				st := makeTestState(t, cs)
				stateRoot := testDummyState(t, cs, st, testHeader)
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
				return stateRoot
			},
			check: func(t *testing.T, expectedStateRoot common.Root, res any, err error) {
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
					StateRoot:     expectedStateRoot.Hex(),
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
			setMockExpectations: func(*testing.T, *mocks.Backend) common.Root {
				// nothing to set here, slot is invalid
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
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
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(nil, math.Slot(0), errTestHeaderNotFound)
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
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
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()

				st := makeTestState(t, cs)
				stateRoot := testDummyState(t, cs, st, testHeader)
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(testParentHeader.Slot, nil)
				return stateRoot
			},
			check: func(t *testing.T, expectedStateRoot common.Root, res any, err error) {
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
					StateRoot:     expectedStateRoot.Hex(),
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
			setMockExpectations: func(*testing.T, *mocks.Backend) common.Root {
				// nothing to set here, slot is invalid
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
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
			setMockExpectations: func(_ *testing.T, b *mocks.Backend) common.Root {
				// Assume parent header is not known
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(0, errTestHeaderNotFound)
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
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
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()

				st := makeTestState(t, cs)
				stateRoot := testDummyState(t, cs, st, testHeader)
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(testParentHeader.Slot, nil)
				return stateRoot
			},
			check: func(t *testing.T, expectedStateRoot common.Root, res any, err error) {
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
					StateRoot:     expectedStateRoot.Hex(),
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
			setMockExpectations: func(_ *testing.T, b *mocks.Backend) common.Root {
				b.EXPECT().GetSlotByBlockRoot(testParentHeader.BodyRoot).Return(testParentHeader.Slot, nil)
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, beacon.ErrMismatchedSlotAndParentBlock)
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

			// create API inputs
			input := tc.inputs()
			inputBytes, err := json.Marshal(input) //nolint:musttag //  TODO:fix
			require.NoError(t, err)
			body := strings.NewReader(string(inputBytes))
			req := httptest.NewRequest(http.MethodGet, "/", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // otherwise code=415, message=Unsupported Media Type
			c := e.NewContext(req, httptest.NewRecorder())

			// set expectations
			expectedStateRoot := tc.setMockExpectations(t, backend)

			// test
			res, err := h.GetBlockHeaders(c)

			// finally do checks
			tc.check(t, expectedStateRoot, res, err)
		})
	}
}

func TestGetBlockHeaderByID(t *testing.T) {
	t.Parallel()

	cs, errSpec := spec.MainnetChainSpec()
	require.NoError(t, errSpec)

	// a test testHeader to build test cases on top of
	testHeader := &ctypes.BeaconBlockHeader{
		Slot:            math.Slot(1234),
		ProposerIndex:   math.ValidatorIndex(5678),
		ParentBlockRoot: common.Root{'p', 'a', 'r', 'e', 'n', 't', 'b', 'l', 'o', 'c', 'k', 'r', 'o', 'o', 't'},
		StateRoot:       common.Root{}, // set in test cases
		BodyRoot:        common.Root{'d', 'u', 'm', 'm', 'y', 'r', 'o', 'o', 't'},
	}
	errTestHeaderNotFound := errors.New("test header not found error")

	testCases := []struct {
		name                string
		inputs              func() beacontypes.GetBlockHeaderRequest
		setMockExpectations func(*testing.T, *mocks.Backend) common.Root
		check               func(t *testing.T, expectedStateRoot common.Root, res any, err error)
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
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()

				st := makeTestState(t, cs)
				stateRoot := testDummyState(t, cs, st, testHeader)
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
				return stateRoot
			},
			check: func(t *testing.T, expectedStateRoot common.Root, res any, err error) {
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
					StateRoot:     expectedStateRoot.Hex(),
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
			setMockExpectations: func(*testing.T, *mocks.Backend) common.Root {
				// nothing to set here, slot is invalid
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
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
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(nil, math.Slot(0), errTestHeaderNotFound)
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
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
			setMockExpectations: func(t *testing.T, b *mocks.Backend) common.Root {
				t.Helper()

				st := makeTestState(t, cs)
				stateRoot := testDummyState(t, cs, st, testHeader)
				b.EXPECT().StateAndSlotFromHeight(mock.Anything).Return(st, math.Slot(0), nil)
				b.EXPECT().GetSlotByBlockRoot(testHeader.BodyRoot).Return(testHeader.Slot, nil)
				return stateRoot
			},
			check: func(t *testing.T, expectedStateRoot common.Root, res any, err error) {
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
					StateRoot:     expectedStateRoot.Hex(),
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
			setMockExpectations: func(*testing.T, *mocks.Backend) common.Root {
				// nothing to set here, slot is invalid
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
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
			setMockExpectations: func(_ *testing.T, b *mocks.Backend) common.Root {
				// Assume parent header is not known
				b.EXPECT().GetSlotByBlockRoot(testHeader.BodyRoot).Return(0, errTestHeaderNotFound)
				return common.Root{}
			},
			check: func(t *testing.T, _ common.Root, _ any, err error) {
				t.Helper()
				require.ErrorIs(t, err, errTestHeaderNotFound)
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

			// create API inputs
			input := tc.inputs()
			inputBytes, err := json.Marshal(input) //nolint:musttag //  TODO:fix
			require.NoError(t, err)
			body := strings.NewReader(string(inputBytes))
			req := httptest.NewRequest(http.MethodGet, "/", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // otherwise code=415, message=Unsupported Media Type
			c := e.NewContext(req, httptest.NewRecorder())

			// set expectations
			expectedStateRoot := tc.setMockExpectations(t, backend)

			// test
			res, err := h.GetBlockHeaderByID(c)

			// finally do checks
			tc.check(t, expectedStateRoot, res, err)
		})
	}
}

// TestDummyState sets a few dummy state elements, enough to
// be able to call HashTreeRoot over the state.
func testDummyState(
	t *testing.T,
	cs chain.Spec,
	st *statedb.StateDB,
	testHeader *ctypes.BeaconBlockHeader,
) common.Root {
	t.Helper()

	dummyFork := &ctypes.Fork{
		PreviousVersion: version.Deneb(),
		CurrentVersion:  version.Electra(),
		Epoch:           math.Epoch(200),
	}
	require.NoError(t, st.SetSlot(1234))
	require.NoError(t, st.SetFork(dummyFork))
	require.NoError(t, st.SetGenesisValidatorsRoot(common.Root{0x01}))
	require.NoError(t, st.SetLatestBlockHeader(testHeader))
	for i := range cs.HistoricalRootsLimit() {
		require.NoError(t, st.UpdateBlockRootAtIndex(i, common.Root{}))
		require.NoError(t, st.UpdateStateRootAtIndex(i, common.Root{}))
	}
	require.NoError(t, st.SetLatestExecutionPayloadHeader(&ctypes.ExecutionPayloadHeader{
		Versionable: ctypes.NewVersionable(dummyFork.CurrentVersion),
	}))
	require.NoError(t, st.SetEth1Data(&ctypes.Eth1Data{}))
	require.NoError(t, st.SetEth1DepositIndex(0))
	require.NoError(t, st.AddValidator(&ctypes.Validator{}))
	require.NoError(t, st.SetBalance(math.ValidatorIndex(0), math.Gwei(0)))
	for i := range cs.EpochsPerHistoricalVector() {
		require.NoError(t, st.UpdateRandaoMixAtIndex(i, common.Bytes32{}))
	}
	require.NoError(t, st.SetNextWithdrawalIndex(0))
	require.NoError(t, st.SetNextWithdrawalValidatorIndex(math.ValidatorIndex(0)))
	require.NoError(t, st.SetPendingPartialWithdrawals([]*ctypes.PendingPartialWithdrawal{}))

	return st.HashTreeRoot()
}

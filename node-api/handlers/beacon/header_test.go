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

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	beaconecho "github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/mocks"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetBlockHeaders(t *testing.T) {
	t.Parallel()

	// setup test
	backend := mocks.NewBackend(t)
	h := beacon.NewHandler(backend)
	h.SetLogger(noop.NewLogger[log.Logger]())

	// define input data
	inputSlot := math.Slot(10)
	input := beacontypes.GetBlockHeadersRequest{
		SlotRequest: beacontypes.SlotRequest{
			Slot: inputSlot.Base10(),
		},
		ParentRoot: "",
	}

	header := &types.BeaconBlockHeader{
		Slot:            inputSlot,
		ProposerIndex:   math.ValidatorIndex(1234),
		ParentBlockRoot: common.Root{'d', 'u', 'm', 'm', 'y', ' ', 'b', 'l', 'o', 'c', 'k', 'r', 'o', 'o', 't'},
		StateRoot:       common.Root{'d', 'u', 'm', 'm', 'y', ' ', 's', 't', 'a', 't', 'e', 'r', 'o', 'o', 't'},
		BodyRoot:        common.Root{'d', 'u', 'm', 'm', 'y', ' ', 'r', 'o', 'o', 't'},
	}

	// create API inputs
	e := echo.New()
	e.Validator = &beaconecho.CustomValidator{
		Validator: beaconecho.ConstructValidator(),
	}

	inputBytes, err := json.Marshal(input) //nolint:musttag //  TODO:fix
	require.NoError(t, err)
	body := strings.NewReader(string(inputBytes))
	req := httptest.NewRequest(http.MethodGet, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // otherwise code=415, message=Unsupported Media Type
	c := e.NewContext(req, httptest.NewRecorder())

	// setup mocks
	backend.EXPECT().BlockHeaderAtSlot(inputSlot).Return(header, nil)

	// test
	res, err := h.GetBlockHeaders(c)

	// checks
	require.NoError(t, err)
	require.NotNil(t, res)
	require.IsType(t, beacontypes.GenericResponse{}, res)
	gr, _ := res.(beacontypes.GenericResponse)
	require.IsType(t, &beacontypes.BlockHeaderResponse{}, gr.Data)
	data, _ := gr.Data.(*beacontypes.BlockHeaderResponse)

	require.Equal(t, header.BodyRoot, data.Root)
	expectedHeader := &beacontypes.BeaconBlockHeader{
		Slot:          header.Slot.Base10(),
		ProposerIndex: header.ProposerIndex.Base10(),
		ParentRoot:    header.ParentBlockRoot.Hex(),
		StateRoot:     header.StateRoot.Hex(),
		BodyRoot:      header.BodyRoot.Hex(),
	}
	require.Equal(t, expectedHeader, data.Header.Message)
}

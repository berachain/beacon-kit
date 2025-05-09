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
//

package cometbft

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NOTE: Partially copied from: https://github.com/cosmos/cosmos-sdk/blob/960d44842b9e313cbe762068a67a894ac82060ab/baseapp/abci.go#L1098
func (s *Service) handleQueryStore(path []string, req *abci.QueryRequest) abci.QueryResponse {
	queryable, ok := s.sm.GetCommitMultiStore().(storetypes.Queryable)
	if !ok {
		return queryResult(
			errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "multi-store does not support queries"),
		)
	}

	req.Path = "/" + strings.Join(path[1:], "/") // req.Path == "/beacon"

	if req.Height <= 1 && req.Prove {
		return queryResult(
			errorsmod.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			))
	}

	sdkReq := storetypes.RequestQuery(*req)
	resp, err := queryable.Query(&sdkReq)
	if err != nil {
		return queryResult(err)
	}
	resp.Height = req.Height

	return abci.QueryResponse(*resp)
}

// NOTE: Copied from here: https://github.com/cosmos/cosmos-sdk/blob/960d44842b9e313cbe762068a67a894ac82060ab/baseapp/errors.go#L37-L46.
// This was made public in v0.53, under the sdkerrors module.
// NOTE: the debug parameter has been removed since Service does not expose this functionality.
//
// queryResult returns a ResponseQuery from an error. It will try to parse ABCI
// info from the error.
func queryResult(err error) abci.QueryResponse {
	space, code, log := errorsmod.ABCIInfo(err, false)
	return abci.QueryResponse{
		Codespace: space,
		Code:      code,
		Log:       log,
	}
}

// NOTE: Copied from here: https://github.com/cosmos/cosmos-sdk/blob/960d44842b9e313cbe762068a67a894ac82060ab/baseapp/abci.go#L1153-L1165
//
// splitABCIQueryPath splits a string path using the delimiter '/'.
//
// e.g. "this/is/funny" becomes []string{"this", "is", "funny"}
func splitABCIQueryPath(requestPath string) []string {
	path := strings.Split(requestPath, "/")

	// first element is empty string
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	return path
}

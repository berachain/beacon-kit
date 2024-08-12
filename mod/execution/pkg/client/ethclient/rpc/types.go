// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package rpc

import (
	"fmt"

	json "github.com/goccy/go-json"
)

// Request represents an Ethereum JSON-RPC request.
type Request struct {
	// ID is the request ID.
	ID int `json:"id"`
	// JSONRPC is the JSON-RPC version.
	JSONRPC string `json:"jsonrpc"`
	// Method is the RPC method to be called.
	Method string `json:"method"`
	// Params are the parameters for the RPC method.
	Params any `json:"params"`
}

// Response represents an Ethereum JSON-RPC response.
type Response struct {
	// ID is the request ID.
	ID int `json:"id"`
	// JSONRPC is the JSON-RPC version.
	JSONRPC string `json:"jsonrpc"`
	// Result is the raw JSON-RPC response result.
	Result json.RawMessage `json:"result"`
	// Error is the JSON-RPC error, if any.
	Error *Error `json:"error"`
}

// Error represents an Ethereum JSON-RPC error.
type Error struct {
	// Code is the error code.
	Code int `json:"code"`
	// Message is the error message.
	Message string `json:"message"`
}

// Error returns a formatted error string.
func (err Error) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

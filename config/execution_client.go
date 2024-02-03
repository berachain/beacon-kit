// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package config

// Client is the configuration struct for the execution client.
type ExecutionClient struct {
	// RPCDialURL is the HTTP url of the execution client JSON-RPC endpoint.
	RPCDialURL string
	// RPCTimeout is the RPC timeout for execution client requests.
	RPCTimeout uint64
	// RPCRetries is the number of retries before shutting down consensus client.
	RPCRetries uint64
	// JWTSecretPath is the path to the JWT secret.
	JWTSecretPath string
	// RequiredChainID is the chain id that the consensus client must be connected to.
	RequiredChainID uint64
}

// DefaultExecutionClientConfig returns the default configuration for the execution client.
func DefaultExecutionClientConfig() ExecutionClient {
	return ExecutionClient{
		RPCDialURL:      "http://localhost:8551",
		RPCTimeout:      5, //nolint:gomnd // default config.
		RPCRetries:      3, //nolint:gomnd // default config.
		JWTSecretPath:   "./app/jwt.hex",
		RequiredChainID: 7, //nolint:gomnd // default config.
	}
}

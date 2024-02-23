// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package ethclient

// BeaconKitSupportedCapabilities returns the full list of capabilities
// of the beacon kit client.
func BeaconKitSupportedCapabilities() []string {
	return []string{
		NewPayloadMethodV2,
		NewPayloadMethodV3,
		ForkchoiceUpdatedMethodV2,
		ForkchoiceUpdatedMethodV3,
		GetPayloadMethodV2,
		GetPayloadMethodV3,
		GetClientVersionV1,
	}
}

// Constants for JSON-RPC method names.
const (
	// NewPayloadMethodV2 for creating a new payload in Capella.
	NewPayloadMethodV2 = "engine_newPayloadV2"
	// NewPayloadMethodV3 for creating a new payload in Deneb.
	NewPayloadMethodV3 = "engine_newPayloadV3"
	// ForkchoiceUpdatedMethodV2 for updating fork choice in Capella.
	ForkchoiceUpdatedMethodV2 = "engine_forkchoiceUpdatedV2"
	// ForkchoiceUpdatedMethodV3 for updating fork choice in Deneb.
	ForkchoiceUpdatedMethodV3 = "engine_forkchoiceUpdatedV3"
	// GetPayloadMethodV2 for retrieving a payload in Capella.
	GetPayloadMethodV2 = "engine_getPayloadV2"
	// GetPayloadMethodV3 for retrieving a payload in Deneb.
	GetPayloadMethodV3 = "engine_getPayloadV3"
	// BlockByHashMethod for retrieving a block by its hash.
	BlockByHashMethod = "eth_getBlockByHash"
	// BlockByNumberMethod for retrieving a block by its number.
	BlockByNumberMethod = "eth_getBlockByNumber"
	// ExchangeCapabilities for exchanging capabilities with the peer.
	ExchangeCapabilities = "engine_exchangeCapabilities"
	// GetClientVersionV1 for retrieving the capabilities of the peer.
	GetClientVersionV1 = "engine_getClientVersionV1"
)

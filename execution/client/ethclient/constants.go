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

package ethclient

// BeaconKitSupportedCapabilities returns the full list of capabilities
// of the beacon kit client.
func BeaconKitSupportedCapabilities() []string {
	return []string{
		NewPayloadMethodV3,
		NewPayloadMethodV4,
		NewPayloadMethodV4P11,
		ForkchoiceUpdatedMethodV3,
		ForkchoiceUpdatedMethodV3P11,
		GetPayloadMethodV3,
		GetPayloadMethodV4,
		GetPayloadMethodV4P11,
		GetClientVersionV1,
	}
}

// Constants for JSON-RPC method names.
const (
	// NewPayloadMethodV3 for creating a new payload in Deneb.
	NewPayloadMethodV3 = "engine_newPayloadV3"
	// NewPayloadMethodV4 for creating a new payload in Electra.
	NewPayloadMethodV4 = "engine_newPayloadV4"
	// NewPayloadMethodV4P11 for creating a new payload for Electra1 (Pectra11).
	NewPayloadMethodV4P11 = "engine_newPayloadV4P11"
	// ForkchoiceUpdatedMethodV3 for updating fork choice in Deneb.
	ForkchoiceUpdatedMethodV3 = "engine_forkchoiceUpdatedV3"
	// ForkchoiceUpdatedMethodV4P11 for updating fork choice in Electra1 (Pectra11).
	ForkchoiceUpdatedMethodV3P11 = "engine_forkchoiceUpdatedV3P11"
	// GetPayloadMethodV3 for retrieving a payload in Deneb.
	GetPayloadMethodV3 = "engine_getPayloadV3"
	// GetPayloadMethodV4 for retrieving a payload in Electra.
	GetPayloadMethodV4 = "engine_getPayloadV4"
	// GetPayloadMethodV4P11 for retrieving a payload in Electra1 (Pectra11).
	GetPayloadMethodV4P11 = "engine_getPayloadV4P11"
	// BlockByHashMethod for retrieving a block by its hash.
	BlockByHashMethod = "eth_getBlockByHash"
	// BlockByNumberMethod for retrieving a block by its number.
	BlockByNumberMethod = "eth_getBlockByNumber"
	// ExchangeCapabilities for exchanging capabilities with the peer.
	ExchangeCapabilities = "engine_exchangeCapabilities"
	// GetClientVersionV1 for retrieving the capabilities of the peer.
	GetClientVersionV1 = "engine_getClientVersionV1"
)

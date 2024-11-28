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

package suite

import (
	"time"

	"github.com/ethereum/go-ethereum/params"
)

// Ether represents the number of wei in one ether, used for Ethereum
// transactions.
const (
	Ether   = params.Ether
	OneGwei = uint64(params.GWei) // 1 Gwei = 1e9 wei
	TenGwei = 10 * OneGwei        // 10 Gwei = 1e10 wei
)

// EtherTransferGasLimit specifies the gas limit for a standard Ethereum
// transfer.
// This is the amount of gas required to perform a basic ether transfer.
const (
	EtherTransferGasLimit uint64 = 21000 // Standard gas limit for ether transfer
)

// DefaultE2ETestTimeout defines the default timeout duration for end-to-end
// tests. This is used to specify how long to wait for a test before considering
// it failed.
const (
	DefaultE2ETestTimeout = 60 * 10 * time.Second // timeout for E2E tests
)

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

package constants

const (
	// LogsBloomLength the length of a LogsBloom in bytes.
	LogsBloomLength = 256

	// ExtraDataLength is the length of the extra data in bytes.
	ExtraDataLength = 32

	// DepositContractDepth is the depth of the deposit contract merkle tree.
	DepositContractDepth uint64 = 32

	// MaxDeposits is the maximum number of deposits supported by the
	// deposit tree (2**32). This is different from the enforced
	// MaxDepositsPerBlock.
	MaxDeposits uint64 = 1 << DepositContractDepth

	// MaxBytesPerTx is the maximum number of bytes per transaction.
	MaxBytesPerTx uint64 = 1073741824
)

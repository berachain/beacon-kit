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

package state

import (
	"math"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
)

const (
	// EVMInflationWithdrawalIndex is the fixed withdrawal index to be used for
	// the EVM inflation withdrawal in every block. It should remain unused
	// during processing.
	EVMInflationWithdrawalIndex = math.MaxUint64

	// EVMInflationWithdrawalValidatorIndex is the fixed validator index to be
	// used for the EVM inflation withdrawal in every block. It should remain
	// unused during processing.
	EVMInflationWithdrawalValidatorIndex = math.MaxUint64
)

// Boonet special case for emergency minting of EVM tokens. TODO: remove with
// other special cases.
const (
	// EVMMintingSlot is the slot at which we force a single withdrawal to
	// mint EVMMintingAmount EVM tokens to EVMMintingAddress. No other
	// withdrawals are inserted at this slot.
	EVMMintingSlot uint64 = spec.BoonetFork1Height

	// EVMMintingAddress is the address at which we mint EVM tokens to.
	EVMMintingAddress = "0x8a73D1380345942F1cb32541F1b19C40D8e6C94B"

	// EVMMintingAmount is the amount of EVM tokens to mint.
	EVMMintingAmount uint64 = 530000000000000000
)

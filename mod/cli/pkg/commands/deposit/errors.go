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

package deposit

import "errors"

var (
	// ErrValidatorPrivateKeyRequired is returned when the validator private key
	// is required but not provided.
	ErrValidatorPrivateKeyRequired = errors.New(
		"validator private key required",
	)
	// ErrDepositTransactionFailed is returned when the deposit transaction
	// fails.
	ErrDepositTransactionFailed = errors.New(
		"deposit transaction failed",
	)
	// ErrPrivateKeyRequired is returned when the broadcast flag is set but a
	// private key is not provided.
	ErrPrivateKeyRequired = errors.New(
		"private key required",
	)

	// ErrDepositReceiptEmpty is returned when the deposit receipt is nil.
	ErrDepositReceiptEmpty = errors.New(
		"deposit receipt is nil")
)

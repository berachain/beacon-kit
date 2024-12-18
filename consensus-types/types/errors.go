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

package types

import "github.com/berachain/beacon-kit/errors"

var (
	// ErrDepositMessage is an error for when the deposit signature doesn't
	// match.
	ErrDepositMessage = errors.New("invalid deposit message")

	// ErrInvalidWithdrawalCredentials is an error for when the.
	ErrInvalidWithdrawalCredentials = errors.New(
		"invalid withdrawal credentials",
	)

	// ErrForkVersionNotSupported is an error for when the fork
	// version is not supported.
	ErrForkVersionNotSupported = errors.New("fork version not supported")

	// ErrInclusionProofDepthExceeded is an error for when the
	// KZG_COMMITMENT_INCLUSION_PROOF_DEPTH calculation overflows.
	ErrInclusionProofDepthExceeded = errors.New("inclusion proof depth exceeded")

	// ErrNilPayloadHeader is an error for when the payload header is nil.
	ErrNilPayloadHeader = errors.New("nil payload header")
)

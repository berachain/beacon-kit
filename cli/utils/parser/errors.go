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

package parser

import "errors"

var (
	// ErrInvalidPubKeyLength is returned when the public key is invalid.
	ErrInvalidPubKeyLength = errors.New(
		"invalid public key length",
	)

	// ErrInvalidWithdrawalCredentialsLength is returned when the withdrawal
	// credentials are invalid.
	ErrInvalidWithdrawalCredentialsLength = errors.New(
		"invalid withdrawal credentials length",
	)

	// ErrInvalidAmount is returned when the deposit amount is invalid.
	ErrInvalidAmount = errors.New(
		"invalid amount",
	)

	// ErrInvalidSignatureLength is returned when the signature is invalid.
	ErrInvalidSignatureLength = errors.New(
		"invalid signature length",
	)

	// ErrInvalidVersionLength is returned when the deposit version is invalid.
	ErrInvalidVersionLength = errors.New(
		"invalid version",
	)

	// ErrInvalidRootLength is returned when the deposit root is invalid.
	ErrInvalidRootLength = errors.New(
		"invalid root length",
	)

	// ErrInvalid0xPrefixedHexString is returned when the input string is not
	// a valid 0x prefixed hex string.
	ErrInvalid0xPrefixedHexString = errors.New(
		"invalid 0x prefixed hex string",
	)
)

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

package hex

import (
	"errors"
)

var (
	ErrEmptyString        = errors.New("empty hex string")
	ErrMissingPrefix      = errors.New("hex string without 0x prefix")
	ErrOddLength          = errors.New("hex string of odd length")
	ErrNonQuotedString    = errors.New("non-quoted hex string")
	ErrInvalidString      = errors.New("invalid hex string")
	ErrLeadingZero        = errors.New("hex number with leading zero digits")
	ErrEmptyNumber        = errors.New("hex string \"0x\"")
	ErrUint64Range        = errors.New("hex number > 64 bits")
	ErrBig256Range        = errors.New("hex number > 256 bits")
	ErrInvalidBigWordSize = errors.New("weird big.Word size")
)

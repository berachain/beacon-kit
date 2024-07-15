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

package blockchain

import "github.com/berachain/beacon-kit/mod/errors"

var (
	// ErrInvalidPayload indicates that the payload of a beacon block is
	// invalid.
	ErrInvalidPayload = errors.New("invalid payload")
	// ErrNoPayloadInBeaconBlock indicates that a beacon block was expected to
	// have a payload, but none was found.
	ErrNoPayloadInBeaconBlock = errors.New("no payload in beacon block")
	// ErrNilBlkBody is an error for when the block body is nil.
	ErrNilBlkBody = errors.New("nil block body")
	// ErrNilBlk is an error for when the beacon block is nil.
	ErrNilBlk = errors.New("nil beacon block")
	// ErrDataNotAvailable indicates that the required data is not available.
	ErrDataNotAvailable = errors.New("data not available")
)

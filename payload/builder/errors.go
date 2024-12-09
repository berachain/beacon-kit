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

package builder

import "github.com/berachain/beacon-kit/errors"

var (
	// ErrPayloadBuilderDisabled is returned when the payload builder is
	// disabled.
	ErrPayloadBuilderDisabled = errors.New("payload builder is disabled")

	// ErrNilPayloadID is returned when a nil payload ID is received.
	ErrNilPayloadID = errors.New("received nil payload ID")

	// ErrPayloadIDNotFound is returned when a payload ID is not found in the
	// cache.
	ErrPayloadIDNotFound = errors.New("unable to find payload ID in cache")

	// ErrNilPayloadEnvelope is returned when a nil payload envelope is
	// received.
	ErrNilPayloadEnvelope = errors.New("received nil payload envelope")
)

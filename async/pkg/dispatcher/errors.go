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

package dispatcher

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/pkg/async"
)

//nolint:gochecknoglobals // errors
var (
	// ErrNotFound is returned when a broker is not found for an eventID.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when a broker is already registered for an
	// eventID.
	ErrAlreadyExists = errors.New("already exists")
	// errBrokerNotFound is a helper function to wrap the ErrNotFound error
	// with the eventID.
	errBrokerNotFound = func(eventID async.EventID) error {
		return errors.Wrapf(
			ErrNotFound,
			"publisher not found for eventID: %s",
			eventID,
		)
	}
	// errBrokerAlreadyExists is a helper function to wrap the ErrAlreadyExists
	// error with the eventID.
	errBrokerAlreadyExists = func(eventID async.EventID) error {
		return errors.Wrapf(
			ErrAlreadyExists,
			"publisher already exists for eventID: %s",
			eventID,
		)
	}
)

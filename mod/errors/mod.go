// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package errors

import (
	"github.com/cockroachdb/errors"
)

// TODO: eventually swap out via build flags if we believe there is value
// to doing so.
//
//nolint:gochecknoglobals // used an alias.
var (
	New   = errors.New
	Newf  = errors.Newf
	Wrap  = errors.Wrap
	Wrapf = errors.Wrapf
	Is    = errors.Is
	IsAny = errors.IsAny
	As    = errors.As
	Join  = errors.Join
)

// DetailedError is a custom error type that includes a message and a flag
// indicating if the error is fatal.
type DetailedError struct {
	// Msg is the error message.
	error
	// fatal indicates if the error is fatal.
	fatal bool
}

// Error returns the error message.
func WrapNonFatal(err error) error {
	return &DetailedError{
		error: err,
		fatal: false,
	}
}

// WrapFatal creates a new DetailedError with the
// provided message and fatal flag.
func WrapFatal(err error) error {
	return &DetailedError{
		error: err,
		fatal: true,
	}
}

// IsFatal checks if the provided error is a
// DetailedError and if it is fatal.
func IsFatal(err error) bool {
	// If the error is nil, obviouisly it is not fatal.
	if err == nil {
		return false
	}

	// Otherwise check for our custom error.
	var customErr *DetailedError
	if errors.As(err, &customErr) {
		if customErr == nil {
			return false
		}

		// If the underlying error is nil, we
		// return false.
		if customErr.error == nil {
			return false
		}

		// Otherwise check the custom fatal field.
		return customErr.fatal
	}

	// All other errors are fatal.
	return true
}

// JoinFatal checks if any of the provided errors is a
// DetailedError and if it is fatal.
func JoinFatal(errs ...error) error {
	fatal := false
	for _, err := range errs {
		if IsFatal(err) {
			fatal = true
			break
		}
	}
	retErr := errors.Join(errs...)
	if fatal {
		return WrapFatal(retErr)
	}
	return WrapNonFatal(retErr)
}

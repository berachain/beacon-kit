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

package errors

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

// TODO: eventually swap out via build flags if we believe there is value
// to doing so.
//
//nolint:gochecknoglobals // used an alias.
var (
	New   = errors.New
	Wrap  = errors.Wrap
	Wrapf = errors.Wrapf
	Is    = errors.Is
	As    = errors.As
	Join  = stderrors.Join
)

// IsAny checks if the provided error is any of the provided errors.
func IsAny(err error, errs ...error) bool {
	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}

// DetailedError is a custom error type that includes a message and a flag
// indicating if the error is fatal.
type DetailedError struct {
	// Msg is the error message.
	error
	// fatal indicates if the error is fatal.
	fatal bool
}

// Error returns a string representation.
func (e *DetailedError) Error() string {
	return e.error.Error()
}

// Unwrap returns the wrapped error.
func (e *DetailedError) Unwrap() error {
	return e.error
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
	// If the error is nil, obviously it is not fatal.
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

	// Only declared fatal errors are fatal.
	return false
}

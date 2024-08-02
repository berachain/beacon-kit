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

package depinject

import (
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

var (
	ErrTargetMustBePointer = errors.New("target must be a pointer")
)

func ProvideError(err error, fn any) error {
	// If we still don't have a function, return a more descriptive error
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return errors.Errorf(
			"Error in %s: fn must be a function, got %T",
			"provide",
			fn,
		)
	}
	funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	// Get the return types of the function
	funcType := reflect.TypeOf(fn)
	if funcType.NumOut() < 1 {
		return errors.Wrapf(
			err,
			"Error in %s: provider %s must return at least one value",
			"provide",
			funcName,
		)
	}
	returnType := funcType.Out(0).String()
	return errors.Wrapf(
		err,
		"Error in %s:\n\n"+
			"Can't resolve type %s\n\n"+
			"from provider %s\n\n"+
			"in function %T\n\n"+
			"error from dig:",
		"provide",
		returnType,
		funcName,
		fn,
	)
}

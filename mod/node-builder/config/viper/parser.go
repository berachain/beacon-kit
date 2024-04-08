// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package viper

import (
	"net/url"
	"reflect"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
)

// StringToSliceHookFunc returns a DecodeHookFunc that converts
// string to a `primitives.ExecutionAddresses` by parsing the string.
func StringToExecutionAddressFunc() mapstructure.DecodeHookFunc {
	return StringTo(
		func(s string) (primitives.ExecutionAddress, error) {
			return common.HexToAddress(s), nil
		},
	)
}

// StringToDialURLFunc returns a DecodeHookFunc that converts
// string to *url.URL by parsing the string.
func StringToDialURLFunc() mapstructure.DecodeHookFunc {
	return StringTo(
		func(s string) (*url.URL, error) {
			url, err := url.Parse(s)
			if err != nil {
				return nil, err
			}
			return url, nil
		},
	)
}

// string to *jwt.Secret by reading the file at the given path.
func StringTo[T any](
	constructor func(string) (T, error),
) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		var retType T
		if t != reflect.TypeOf(retType) {
			return data, nil
		}

		// Convert it by parsing
		return constructor(data.(string))
	}
}

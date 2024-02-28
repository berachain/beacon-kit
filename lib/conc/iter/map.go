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

package iter

import (
	"reflect"

	"github.com/sourcegraph/conc/iter"
)

// MapErrNoNils is a wrapper around MapErr
// that filters out nil elements from the result.
func MapErrNoNils[T, R any](input []T, f func(*T) (R, error)) ([]R, error) {
	elems, err := iter.Mapper[T, R]{}.MapErr(input, f)
	if err != nil {
		return nil, err
	}

	nonNilElems := make([]R, 0, len(elems))
	for _, elem := range elems {
		if !IsNil(elem) {
			nonNilElems = append(nonNilElems, elem)
		}
	}
	return nonNilElems, nil
}

func IsNil[T any](elem T) bool {
	isPtr := reflect.TypeOf(elem).Kind() == reflect.Pointer
	if !isPtr {
		return false
	}
	elemVal := reflect.ValueOf(elem).Interface()
	nilVal := reflect.Zero(reflect.TypeOf(elem)).Interface()
	return elemVal == nilVal
}

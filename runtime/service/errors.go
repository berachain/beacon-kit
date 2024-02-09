// SPDX-License-Identifier: MIT
//
// # Copyright (c) 2023 Berachain Foundation
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
//
//nolint:gochecknoglobals // this file contains functions for use as errors.
package service

import "fmt"

var (
	// errServiceAlreadyExists defines an error for when a service already exists.
	errServiceAlreadyExists = func(serviceName string) error {
		return fmt.Errorf("service already exists: %v", serviceName)
	}

	// errInputIsNotPointer defines an error for when the input must be of pointer type.
	errInputIsNotPointer = func(valueType interface{}) error {
		return fmt.Errorf(
			"input must be of pointer type, received value type instead: %T", valueType,
		)
	}

	// errUnknownService defines an error for when an unknown service is encountered.
	errUnknownService = func(serviceType interface{}) error {
		return fmt.Errorf("unknown service: %T", serviceType)
	}
)

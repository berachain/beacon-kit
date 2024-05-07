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

package noop

// Logger is a logger that performs no operations. It can be used in
// environments where logging should be disabled. It implements the Logger
// interface with no-op methods.
type Logger[KeyValT any] struct{}

// NewLogger creates a blank no-op logger.
func NewLogger() *Logger[any] {
	return &Logger[any]{}
}

// Info logs an informational message with associated key-value pairs. This
// method does nothing.
func (n *Logger[KeyValT]) Info(string, ...KeyValT) {
	// No operation
}

// Warn logs a warning message with associated key-value pairs. This method does
// nothing.
func (n *Logger[KeyValT]) Warn(string, ...KeyValT) {
	// No operation
}

// Error logs an error message with associated key-value pairs. This method does
// nothing.
func (n *Logger[KeyValT]) Error(string, ...KeyValT) {
	// No operation
}

// Debug logs a debug message with associated key-value pairs. This method does
// nothing.
func (n *Logger[KeyValT]) Debug(string, ...KeyValT) {
	// No operation
}

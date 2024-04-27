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

package log

// Logger is extremely similar to the Cosmos-SDK Logger interface.
type Logger[KeyValT any] interface {
	// Info takes a message and a set of key/value pairs and logs with level
	// INFO.
	// The key of the tuple must be a string.
	Info(msg string, keyVals ...KeyValT)

	// Warn takes a message and a set of key/value pairs and logs with level
	// WARN.
	// The key of the tuple must be a string.
	Warn(msg string, keyVals ...KeyValT)

	// Error takes a message and a set of key/value pairs and logs with level
	// ERR.
	// The key of the tuple must be a string.
	Error(msg string, keyVals ...KeyValT)

	// Debug takes a message and a set of key/value pairs and logs with level
	// DEBUG.
	// The key of the tuple must be a string.
	Debug(msg string, keyVals ...KeyValT)
}

// AdvancedLogger is extremely similar to the Cosmos-SDK Logger interface,
// however we introduce a generic to allow for more flexibility in
// the underlying logger implementation.
type AdvancedLogger[KeyValT, LoggerT any] interface {
	Logger[KeyValT]
	// With returns a new wrapped logger with additional context provided by a
	// set.
	With(keyVals ...KeyValT) LoggerT

	// Impl returns the underlying logger implementation.
	// It is used to access the full functionalities of the underlying logger.
	// Advanced users can type cast the returned value to the actual logger.
	Impl() any
}

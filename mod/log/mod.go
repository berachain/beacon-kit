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

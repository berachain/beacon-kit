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

// Logger represents a basic logger that is extremely similar to the Cosmos-SDK
// Logger interface.
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

// ConfigurableLogger extends the basic logger with the ability to configure
// the logger with a config.
type ConfigurableLogger[
	ConfigurableLoggerT, KeyValT any, ConfigT any,
] interface {
	Logger[KeyValT]
	WithConfig(config ConfigT) ConfigurableLoggerT
}

// ColorLogger extends the basic logger with the ability to configure the
// logger with key and key value colors.
type ColorLogger[KeyValT any] interface {
	Logger[KeyValT]
	// AddKeyColor sets the log color for a key.
	AddKeyColor(key any, color Color)
	// AddKeyValColor sets the log color for a key and its value.
	AddKeyValColor(key any, val any, color Color)
}

// AdvancedLogger extends the color logger with the ability to wrap the logger
// with additional context and to access the underlying logger implementation.
type AdvancedLogger[KeyValT any] interface {
	ColorLogger[KeyValT]
	// With returns a new wrapped logger with additional context provided by a
	// set.
	With(keyVals ...KeyValT) AdvancedLogger[KeyValT]
	// Impl returns the underlying logger implementation.
	// It is used to access the full functionalities of the underlying logger.
	// Advanced users can type cast the returned value to the actual logger.
	Impl() any
}

// Color is a string that holds the hex color code for the color.
type Color string

// Raw returns the raw color code.
func (c Color) Raw() string {
	return string(c)
}

const (
	// colours.
	Reset   Color = "\x1b[0m"
	Black   Color = "\x1b[30m"
	Red     Color = "\x1b[31m"
	Green   Color = "\x1b[32m"
	Yellow  Color = "\x1b[33m"
	Blue    Color = "\x1b[34m"
	Magenta Color = "\x1b[35m"
	Cyan    Color = "\x1b[36m"
	White   Color = "\x1b[37m"

	Gray          Color = "\x1b[90m"
	BrightRed     Color = "\x1b[91m"
	BrightGreen   Color = "\x1b[92m"
	BrightYellow  Color = "\x1b[93m"
	BrightBlue    Color = "\x1b[94m"
	BrightMagenta Color = "\x1b[95m"
	BrightCyan    Color = "\x1b[96m"
	BrightWhite   Color = "\x1b[97m"

	BrightBackgroundWhite Color = "\x1b[107m"
	BrightBackgroundBlue  Color = "\x1b[104m"
)

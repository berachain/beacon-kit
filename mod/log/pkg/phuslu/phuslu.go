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

package phuslu

import (
	"io"

	"github.com/phuslu/log"
)

// Logger is a wrapper around phuslogger.
type Logger[ImplT any] struct {
	// logger is the underlying logger implementation.
	logger *log.Logger
	// context is a map of key-value pairs that are added to every log entry.
	context log.Fields
}

// NewLogger creates a new logger with the given log level, ConsoleWriter, and
// default configuration.
func NewLogger[ImplT any](
	level string, out io.Writer,
) *Logger[ImplT] {
	cfg := DefaultConfig()
	logger := &log.Logger{
		Level:      log.ParseLevel(level),
		TimeFormat: cfg.TimeFormat,
		Writer: &log.ConsoleWriter{
			Writer:    out,
			Formatter: (NewFormatter().Format),
		},
	}
	return &Logger[ImplT]{
		logger:  logger,
		context: make(log.Fields),
	}
}

// Info logs a message at level Info.
func (l *Logger[ImplT]) Info(msg string, keyVals ...any) {
	l.msgWithContext(msg, l.logger.Info(), keyVals...)
}

// Warn logs a message at level Warn.
func (l *Logger[ImplT]) Warn(msg string, keyVals ...any) {
	l.msgWithContext(msg, l.logger.Warn(), keyVals...)
}

// Error logs a message at level Error.
func (l *Logger[ImplT]) Error(msg string, keyVals ...any) {
	l.msgWithContext(msg, l.logger.Error(), keyVals...)
}

// Debug logs a message at level Debug.
func (l *Logger[ImplT]) Debug(msg string, keyVals ...any) {
	// In a special case for debug, we check to see if the
	// logger level is set to debug before logging the message.
	// We don't do this in other log levels since they are more common
	// and we would add the overhead of the if check, when the happy
	// path in their case, is to print out the line.
	if l.logger.Level > log.DebugLevel {
		return
	}
	l.msgWithContext(msg, l.logger.Debug(), keyVals...)
}

// Impl returns the underlying logger implementation.
func (l *Logger[ImplT]) Impl() any {
	return l.logger
}

// With returns a new wrapped logger with additional context provided by a set.
func (l Logger[ImplT]) With(keyVals ...any) ImplT {
	newLogger := l

	// Perform a deep copy of the map with preallocated size.
	newLogger.context = make(log.Fields, len(l.context)+len(keyVals)/2)
	for k, v := range l.context {
		newLogger.context[k] = v
	}

	// Add the new context to the existing context.
	for i := 0; i < len(keyVals); i += 2 {
		key, ok := keyVals[i].(string)
		if !ok {
			continue
		}
		newLogger.context[key] = keyVals[i+1]
	}

	return any(&newLogger).(ImplT)
}

// msgWithContext logs a message with keyVals and current context.
func (l *Logger[ImplT]) msgWithContext(
	msg string, e *log.Entry, keyVals ...any,
) {
	e.Fields(l.context).KeysAndValues(keyVals...).Msg(msg)
}

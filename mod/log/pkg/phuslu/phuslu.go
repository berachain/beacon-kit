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
	// out is the writer to write logs to.
	out io.Writer
	// formatter is the formatter to use for the logger.
	formatter *Formatter
	// message is the function to use for logging messages.
	message func(msg string, ctx log.Fields, e *log.Entry, keyVals ...any)
}

// NewLogger creates a new logger with the given log level, ConsoleWriter, and
// default configuration.
func NewLogger[ImplT any](
	level string, out io.Writer,
) *Logger[ImplT] {
	formatter := NewFormatter()
	pLogger := &log.Logger{
		Level: log.ParseLevel(level),
	}
	logger := &Logger[ImplT]{
		logger:    pLogger,
		context:   make(log.Fields),
		out:       out,
		formatter: formatter,
	}
	// we can remove this and assume valid config is provided before any
	// messaging occurs
	// logger.message = logger.msgWithContext
	return logger
}

// Info logs a message at level Info.
func (l *Logger[ImplT]) Info(msg string, keyVals ...any) {
	l.message(msg, l.context, l.logger.Info(), keyVals...)
}

// Warn logs a message at level Warn.
func (l *Logger[ImplT]) Warn(msg string, keyVals ...any) {
	l.message(msg, l.context, l.logger.Warn(), keyVals...)
}

// Error logs a message at level Error.
func (l *Logger[ImplT]) Error(msg string, keyVals ...any) {
	l.message(msg, l.context, l.logger.Error(), keyVals...)
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
	l.message(msg, l.context, l.logger.Debug(), keyVals...)
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

// SetLevel sets the log level of the logger.
func (l *Logger[ImplT]) SetLevel(level string) {
	l.logger.Level = log.ParseLevel(level)
}

// we need this to set the config post-creation of the logger.
// This is so cooked but necessary because there is no way to pass a populated
// config to the logger at creation time, because the dependent viper instance
// is not yet populated.
func (l *Logger[ImplT]) WithConfig(cfg Config) *Logger[ImplT] {
	l.withTimeFormat(cfg.TimeFormat)
	l.withStyle(cfg.Style)
	l.withVerbosity(cfg.Verbose)
	return l
}

/* -------------------------------------------------------------------------- */
/*                                 messaging                                  */
/* -------------------------------------------------------------------------- */

// msgWithContext logs a message with keyVals and current context.
func msgWithContext(
	msg string, ctx log.Fields, e *log.Entry, keyVals ...any,
) {
	e.Fields(ctx).KeysAndValues(keyVals...).Msg(msg)
}

// msgWithoutContext logs a message with keyVals and without context.
func msgWithoutContext(
	msg string, _ log.Fields, e *log.Entry, keyVals ...any,
) {
	e.KeysAndValues(keyVals...).Msg(msg)
}

/* -------------------------------------------------------------------------- */
/*                             configuration                                  */
/* -------------------------------------------------------------------------- */

// sets the style of the logger.
func (l *Logger[Impl]) withStyle(style string) {
	if style == StylePretty {
		l.useConsoleWriter()
	} else if style == StyleJSON {
		l.useJSONWriter()
	}
}

// withVerbosity sets the verbosity of the logger.
func (l *Logger[Impl]) withVerbosity(verbose bool) {
	if verbose {
		l.message = msgWithContext
	} else {
		l.message = msgWithoutContext
	}
}

// useConsoleWriter sets the logger to use a console writer.
func (l *Logger[ImplT]) useConsoleWriter() {
	l.setWriter(&log.ConsoleWriter{
		Writer:    l.out,
		Formatter: l.formatter.Format,
	})
}

// useJSONWriter sets the logger to use a IOWriter wrapper.
func (l *Logger[ImplT]) useJSONWriter() {
	l.setWriter(log.IOWriter{Writer: l.out})
}

// setWriter sets the writer of the logger.
func (l *Logger[ImplT]) setWriter(writer log.Writer) {
	l.logger.Writer = writer
}

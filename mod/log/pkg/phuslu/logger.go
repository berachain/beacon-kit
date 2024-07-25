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
}

// NewLogger initializes a new wrapped phuslogger with the provided config.
func NewLogger[ImplT any](
	out io.Writer,
	cfg *Config,
) *Logger[ImplT] {
	logger := &Logger[ImplT]{
		logger:    &log.Logger{},
		context:   make(log.Fields),
		out:       out,
		formatter: NewFormatter(),
	}
	logger.WithConfig(*cfg)
	return logger
}

// Info logs a message at level Info.
func (l *Logger[ImplT]) Info(msg string, keyVals ...any) {
	if l.logger.Level > log.InfoLevel {
		return
	}
	l.msgWithContext(msg, l.logger.Info(), keyVals...)
}

// Warn logs a message at level Warn.
func (l *Logger[ImplT]) Warn(msg string, keyVals ...any) {
	if l.logger.Level > log.WarnLevel {
		return
	}
	l.msgWithContext(msg, l.logger.Warn(), keyVals...)
}

// Error logs a message at level Error.
func (l *Logger[ImplT]) Error(msg string, keyVals ...any) {
	if l.logger.Level > log.ErrorLevel {
		return
	}
	l.msgWithContext(msg, l.logger.Error(), keyVals...)
}

// Debug logs a message at level Debug.
func (l *Logger[ImplT]) Debug(msg string, keyVals ...any) {
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

// Writer returns the io.Writer of the logger.
func (l *Logger[ImplT]) Writer() io.Writer {
	return l.out
}

// msgWithContext logs a message with keyVals and current context.
func (l *Logger[Impl]) msgWithContext(
	msg string, e *log.Entry, keyVals ...any,
) {
	e.Fields(l.context).KeysAndValues(keyVals...).Msg(msg)
}

/* -------------------------------------------------------------------------- */
/*                                   config                                   */
/* -------------------------------------------------------------------------- */

// Temporary workaround to allow dynamic configuration post-logger creation.
// This is necessary due to dependencies on runtime-populated configurations.
func (l *Logger[ImplT]) WithConfig(cfg Config) *Logger[ImplT] {
	l.withTimeFormat(cfg.TimeFormat)
	l.withStyle(cfg.Style)
	l.withLogLevel(cfg.LogLevel)
	return l
}

// AddKeyColor applies a color to log entries based on their keys.
func (l *Logger[ImplT]) AddKeyColor(key any, color string) {
	l.formatter.AddKeyColor(key.(string), color)
}

// AddKeyValColor applies specific colors to log entries based on their keys and
// values.
func (l *Logger[ImplT]) AddKeyValColor(key any, val any, color string) {
	l.formatter.AddKeyValColor(key.(string), val.(string), color)
}

// sets the style of the logger.
func (l *Logger[Impl]) withStyle(style string) {
	if style == StylePretty {
		l.useConsoleWriter()
	} else if style == StyleJSON {
		l.useJSONWriter()
	}
}

// SetLevel sets the log level of the logger.
func (l *Logger[ImplT]) withLogLevel(level string) {
	l.logger.Level = log.ParseLevel(level)
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

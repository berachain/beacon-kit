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

type Logger[KeyValT, ImplT any] struct {
	logger *log.Logger
}

func NewLogger[KeyValT, ImplT any](level string, out io.Writer) *Logger[KeyValT, ImplT] {
	logger := &log.Logger{
		Level:      log.ParseLevel(level),
		TimeFormat: "15:04:05",
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
			Writer:         out,
		},
	}
	return &Logger[KeyValT, ImplT]{
		logger: logger,
	}
}

func (l *Logger[KeyValT, ImplT]) Info(msg string, keyVals ...KeyValT) {
	l.logger.Info().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Warn(msg string, keyVals ...KeyValT) {
	l.logger.Warn().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Error(msg string, keyVals ...KeyValT) {
	l.logger.Error().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Debug(msg string, keyVals ...KeyValT) {
	l.logger.Debug().Msg(msg)
}

func (l *Logger[KeyValT, ImplT]) Impl() any {
	return l.logger
}

func (l *Logger[KeyValT, ImplT]) With(keyVals ...KeyValT) ImplT {
	newLogger := *l
	return any(&newLogger).(ImplT)
}

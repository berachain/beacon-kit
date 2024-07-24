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

package echo

import (
	"fmt"
	"io"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/labstack/echo/v4"
	gommonlog "github.com/labstack/gommon/log"
)

// assert that logger implements echo.Logger
var _ echo.Logger = &logger{}

// logger is an adapter that allows a log.Logger to be used as an echo.Logger.
type logger struct {
	log.Logger[any]
	echoLogger echo.Logger
}

func NewLogger(l log.Logger[any], el echo.Logger) echo.Logger {
	return &logger{
		Logger:     l,
		echoLogger: el,
	}
}

// Output returns the io.Writer to which log output is written.
func (l *logger) Output() io.Writer {
	// TODO: don't do a type assertion
	return l.echoLogger.Output()
}

func (l *logger) Debug(i ...interface{}) {
	l.Logger.Debug(fmt.Sprint(i...))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, args...))
}

func (l *logger) Debugj(j gommonlog.JSON) {
	// Empty body
}

func (l *logger) Info(i ...interface{}) {
	l.Logger.Info(fmt.Sprint(i...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, args...))
}

func (l *logger) Infoj(j gommonlog.JSON) {
	// Empty body
}

func (l *logger) Warn(i ...interface{}) {
	l.Logger.Warn(fmt.Sprint(i...))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, args...))
}

func (l *logger) Warnj(j gommonlog.JSON) {
	// Empty body
}

func (l *logger) Error(i ...interface{}) {
	l.Logger.Error(fmt.Sprint(i...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func (l *logger) Errorj(j gommonlog.JSON) {
	// Empty body
}

func (l *logger) Fatal(i ...interface{}) {
	l.Logger.Error(fmt.Sprint(i...))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func (l *logger) Fatalj(j gommonlog.JSON) {
	// Empty body
}

func (l *logger) Panic(i ...interface{}) {
	l.Logger.Error(fmt.Sprint(i...))
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func (l *logger) Panicj(j gommonlog.JSON) {
	// Empty body
}

/* -------------------------------------------------------------------------- */
/*                                   noops                                    */
/* -------------------------------------------------------------------------- */

func (l *logger) SetOutput(w io.Writer) {
	// noop
}

func (l *logger) Prefix() string {
	return ""
}

func (l *logger) SetPrefix(p string) {
	// noop
}

func (l *logger) Level() gommonlog.Lvl {
	return 0
}

func (l *logger) SetLevel(v gommonlog.Lvl) {
	// noop
}

func (l *logger) SetHeader(h string) {
	// noop
}

func (l *logger) Print(i ...interface{}) {
	// noop
}

func (l *logger) Printf(format string, args ...interface{}) {
	// noop
}

func (l *logger) Printj(j gommonlog.JSON) {
	// noop
}

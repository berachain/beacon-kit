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

import "time"

// colours.
const (
	reset      = "\x1b[0m"
	black      = "\x1b[30m"
	red        = "\x1b[31m"
	green      = "\x1b[32m"
	yellow     = "\x1b[33m"
	blue       = "\x1b[34m"
	magenta    = "\x1b[35m"
	cyan       = "\x1b[36m"
	white      = "\x1b[37m"
	gray       = "\x1b[90m"
	lightWhite = "\x1b[97m"
)

// Config is a structure that defines the configuration for the logger.
type Config struct {
	// TimeFormat is a string that defines the format of the time in
	// the logger.
	TimeFormat string `mapstructure:"trace_color"`

	// Colours for the different log levels.
	TraceColor   string `mapstructure:"trace_color"`
	DebugColor   string `mapstructure:"debug_color"`
	InfoColor    string `mapstructure:"info_color"`
	WarnColor    string `mapstructure:"warn_color"`
	ErrorColor   string `mapstructure:"error_color"`
	FatalColor   string `mapstructure:"fatal_color"`
	PanicColor   string `mapstructure:"panic_color"`
	DefaultColor string `mapstructure:"default_color"`

	// Labels for the different log levels.
	TraceLabel   string `mapstructure:"trace_label"`
	DebugLabel   string `mapstructure:"debug_label"`
	InfoLabel    string `mapstructure:"info_label"`
	WarnLabel    string `mapstructure:"warn_label"`
	ErrorLabel   string `mapstructure:"error_label"`
	FatalLabel   string `mapstructure:"fatal_label"`
	PanicLabel   string `mapstructure:"panic_label"`
	DefaultLabel string `mapstructure:"default_label"`
}

// DefaultConfig is a function that returns a new Config with default values.
func DefaultConfig() Config {
	return Config{
		TimeFormat:   time.RFC3339,
		TraceColor:   magenta,
		DebugColor:   yellow,
		InfoColor:    green,
		WarnColor:    yellow,
		ErrorColor:   red,
		FatalColor:   red,
		PanicColor:   red,
		DefaultColor: gray,
		TraceLabel:   "TRCE",
		DebugLabel:   "DBUG",
		InfoLabel:    "INFO",
		WarnLabel:    "WARN",
		ErrorLabel:   "ERRR",
		FatalLabel:   "FATAL",
		PanicLabel:   "PANIC",
		DefaultLabel: " ???",
	}
}

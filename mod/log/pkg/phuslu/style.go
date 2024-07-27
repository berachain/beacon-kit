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

import "github.com/berachain/beacon-kit/mod/log"

type Color = log.Color

const (
	reset                 = log.Reset
	black                 = log.Black
	red                   = log.Red
	green                 = log.Green
	yellow                = log.Yellow
	blue                  = log.Blue
	magenta               = log.Magenta
	cyan                  = log.Cyan
	white                 = log.White
	gray                  = log.Gray
	lightWhite            = log.BrightWhite
	brightBackgroundWhite = log.BrightBackgroundWhite

	// log levels.
	traceColor   = magenta
	debugColor   = yellow
	infoColor    = green
	warnColor    = yellow
	errorColor   = red
	fatalColor   = red
	panicColor   = red
	defaultColor = gray
	apiColor     = blue
	traceLabel   = "TRCE"
	debugLabel   = "DBUG"
	infoLabel    = "INFO"
	warnLabel    = "WARN"
	errorLabel   = "ERRR"
	fatalLabel   = "FATAL"
	panicLabel   = "PANIC"
	defaultLabel = "???"
	apiLabel     = "API"

	// output styles flags.
	StylePretty = "pretty"
	StyleJSON   = "json"
)

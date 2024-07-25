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

// Color is a string that holds the hex color code for the color.
type Color string

// Raw returns the raw color code.
func (c Color) Raw() string {
	return string(c)
}

// String returns the human-readable string representation of the color.
func (c Color) String() string {
	switch c {
	case reset:
		return "reset"
	case black:
		return "black"
	case red:
		return "red"
	case green:
		return "green"
	case yellow:
		return "yellow"
	case blue:
		return "blue"
	case magenta:
		return "magenta"
	case cyan:
		return "cyan"
	case white:
		return "white"
	case gray:
		return "gray"
	case lightWhite:
		return "lightWhite"
	}
	return "unknown"
}

// ToColor converts a human-readable string to a Color.
func ToColor(color string) Color {
	switch color {
	case "black":
		return black
	case "red":
		return red
	case "green":
		return green
	case "yellow":
		return yellow
	case "blue":
		return blue
	case "magenta":
		return magenta
	case "cyan":
		return cyan
	case "white":
		return white
	case "gray":
		return gray
	case "lightWhite":
		return lightWhite
	}
	return reset
}

const (
	// colours.
	reset      Color = "\x1b[0m"
	black      Color = "\x1b[30m"
	red        Color = "\x1b[31m"
	green      Color = "\x1b[32m"
	yellow     Color = "\x1b[33m"
	blue       Color = "\x1b[34m"
	magenta    Color = "\x1b[35m"
	cyan       Color = "\x1b[36m"
	white      Color = "\x1b[37m"
	gray       Color = "\x1b[90m"
	lightWhite Color = "\x1b[97m"

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

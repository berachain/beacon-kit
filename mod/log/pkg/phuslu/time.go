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
	"time"
)

const (
	// time formats.
	rfc3339     = "RFC3339"
	rfc3339Nano = "RFC3339Nano"
	rfc1123     = "RFC1123"
	rfc1123Z    = "RFC1123Z"
	kitchen     = "Kitchen"
	layout      = "Layout"
	ansic       = "ANSIC"
	unixDate    = "UnixDate"
	rubyDate    = "RubyDate"
	timeOnly    = "TimeOnly"
	dateOnly    = "DateOnly"
)

// withTimeFormat sets the time format for the logger.
func (l *Logger) withTimeFormat(formatStr string) {
	l.logger.TimeFormat = parseFormat(formatStr)
}

// parseFormat parses the time format string.
func parseFormat(formatStr string) string {
	switch formatStr {
	case rfc3339:
		return time.RFC3339
	case rfc3339Nano:
		return time.RFC3339Nano
	case rfc1123:
		return time.RFC1123
	case rfc1123Z:
		return time.RFC1123Z
	case kitchen:
		return time.Kitchen
	case layout:
		return time.Layout
	case ansic:
		return time.ANSIC
	case unixDate:
		return time.UnixDate
	case rubyDate:
		return time.RubyDate
	case timeOnly:
		return time.TimeOnly
	case dateOnly:
		return time.DateOnly
	default:
		return time.RFC3339
	}
}

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
// LICENSOR AS EXPRESSLY REQUIred BY THIS LICENSE).
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

// Formatter is a custom formatter for log messages.
type Formatter struct{}

// NewFormatter creates a new Formatter with default settings.
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Format formats the log message.
func (f *Formatter) Format(
	out io.Writer,
	args *log.FormatterArgs,
) (int, error) {
	buffer, ok := byteBufferPool.Get().(*byteBuffer)
	if !ok {
		panic("failed to get byte buffer from pool")
	}

	buffer.Reset()
	defer byteBufferPool.Put(buffer)

	var color, level string
	switch args.Level {
	case "trace":
		color, level = magenta, "TRCE"
	case "debug":
		color, level = yellow, "DBUG"
	case "info":
		color, level = green, "INFO"
	case "warn":
		color, level = yellow, "WARN"
	case "error":
		color, level = red, "ERRR"
	case "fatal":
		color, level = red, "FATAL"
	case "panic":
		color, level = red, "PANIC"
	default:
		color, level = gray, " ???"
	}

	f.printWithColor(args, buffer, color, level)
	f.ensureLineBreak(buffer)

	if args.Stack != "" {
		buffer.Bytes = append(buffer.Bytes, args.Stack...)
		if args.Stack[len(args.Stack)-1] != '\n' {
			buffer.Bytes = append(buffer.Bytes, '\n')
		}
	}

	return out.Write(buffer.Bytes)
}

// printWithColor prints the log message with color.
func (f *Formatter) printWithColor(
	args *log.FormatterArgs,
	b *byteBuffer,
	color, level string,
) {
	f.formatHeader(args, b, color, level)

	b.Bytes = append(b.Bytes, ' ')
	b.Bytes = append(b.Bytes, args.Message...)
	for _, kv := range args.KeyValues {
		b.Bytes = append(b.Bytes, ' ')
		if kv.Key == "error" || kv.Key == "err" {
			b.Bytes = append(b.Bytes, red...)
		}
		b.Bytes = append(b.Bytes, kv.Key...)
		b.Bytes = append(b.Bytes, '=')
		b.Bytes = append(b.Bytes, kv.Value...)
		b.Bytes = append(b.Bytes, reset...)
	}
}

// formatHeader formats the header of the log message.
func (f *Formatter) formatHeader(
	args *log.FormatterArgs,
	b *byteBuffer,
	color, level string,
) {
	b.Bytes = append(b.Bytes, gray...)
	b.Bytes = append(b.Bytes, args.Time...)
	b.Bytes = append(b.Bytes, ' ')
	b.Bytes = append(b.Bytes, color...)
	b.Bytes = append(b.Bytes, level...)
	if args.Caller != "" {
		b.Bytes = append(b.Bytes, args.Goid...)
		b.Bytes = append(b.Bytes, ' ')
		b.Bytes = append(b.Bytes, args.Caller...)
		b.Bytes = append(b.Bytes, ' ')
	}
	b.Bytes = append(b.Bytes, reset...)
}

// ensureLineBreak ensures the log message ends with a line break.
func (f *Formatter) ensureLineBreak(b *byteBuffer) {
	if b.Bytes == nil {
		b.Bytes = make([]byte, 0)
	}
	length := len(b.Bytes)
	if length == 0 || b.Bytes[length-1] != '\n' {
		b.Bytes = append(b.Bytes, '\n')
	}
}

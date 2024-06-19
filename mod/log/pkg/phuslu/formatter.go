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
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/phuslu/log"
)

// colours.
const (
	reset   = "\x1b[0m"
	black   = "\x1b[30m"
	red     = "\x1b[31m"
	green   = "\x1b[32m"
	yellow  = "\x1b[33m"
	blue    = "\x1b[34m"
	magenta = "\x1b[35m"
	cyan    = "\x1b[36m"
	white   = "\x1b[37m"
	gray    = "\x1b[90m"
)

// byteBuffer is a byte buffer.
type byteBuffer struct {
	Bytes []byte
}

// Write writes to the byte buffer.
func (b *byteBuffer) Write(bytes []byte) (int, error) {
	b.Bytes = append(b.Bytes, bytes...)
	return len(bytes), nil
}

// byteBufferPool is a pool of byte buffers.
//
//nolint:gochecknoglobals // buffer pool
var byteBufferPool = sync.Pool{
	New: func() any {
		return new(byteBuffer)
	},
}

func resetBuffer(b *byteBuffer) {
	if b.Bytes != nil {
		b.Bytes = b.Bytes[:0]
	} else {
		b.Bytes = make([]byte, 0)
	}
}

// customFormatter is a custom Formatter to pass into the ConsoleWriter.
func customFormatter(out io.Writer, args *log.FormatterArgs) (int, error) {
	buffer, ok := byteBufferPool.Get().(*byteBuffer)
	if !ok {
		panic("failed to get byte buffer from pool")
	}
	resetBuffer(buffer)
	defer byteBufferPool.Put(buffer)

	// TODO: pull out of config
	var color, level string
	switch args.Level {
	case "trace":
		color, level = magenta, "TRACE"
	case "debug":
		color, level = yellow, "DEBUG"
	case "info":
		color, level = green, " BET"
	case "warn":
		color, level = yellow, "WARN"
	case "error":
		color, level = red, "ERROR"
	case "fatal":
		color, level = red, " FTL"
	case "panic":
		color, level = red, " PNC"
	default:
		color, level = gray, " ???"
	}

	// TODO: pull out of config
	colorOutput := true
	quoteString := true
	endWithMessage := false

	// pretty console writer
	if colorOutput {
		printWithColor(args, buffer, color, level, quoteString, endWithMessage)
	} else {
		printWithoutColor(args, buffer, level, quoteString, endWithMessage)
	}

	// add line break if needed
	ensureLineBreak(buffer)

	// stack
	if args.Stack != "" {
		buffer.Bytes = append(buffer.Bytes, args.Stack...)
		if args.Stack[len(args.Stack)-1] != '\n' {
			buffer.Bytes = append(buffer.Bytes, '\n')
		}
	}

	return out.Write(buffer.Bytes)
}

// printWithColor prints the log message with color.
func printWithColor(args *log.FormatterArgs,
	b *byteBuffer,
	color, level string,
	quoteString, endWithMessage bool) {
	// header
	formatHeader(args, b, true, color, level)
	if !endWithMessage {
		fmt.Fprintf(b, " %s", args.Message)
	}
	// key and values
	for _, kv := range args.KeyValues {
		if quoteString && kv.ValueType == 's' {
			kv.Value = strconv.Quote(kv.Value)
		}
		if kv.Key == "error" {
			fmt.Fprintf(b, " %s%s=%s%s", red, kv.Key, kv.Value, reset)
		} else {
			fmt.Fprintf(b, " %s%s=%s%s%s", cyan, kv.Key, gray, kv.Value, reset)
		}
	}
	// message
	if endWithMessage {
		fmt.Fprintf(b, "%s %s", reset, args.Message)
	}
}

func formatHeader(args *log.FormatterArgs, b *byteBuffer, colorEnabled bool,
	color, level string) {
	headerColor, resetColor := "", ""
	if colorEnabled {
		headerColor, resetColor = color, reset
	}
	fmt.Fprintf(b, "%s%s%s %s%s%s ", gray, args.Time, resetColor, headerColor,
		level, resetColor)
	if args.Caller != "" {
		fmt.Fprintf(b, "%s %s %s💦%s", args.Goid, args.Caller, cyan, resetColor)
	} else {
		fmt.Fprintf(b, "%s💦%s", cyan, resetColor)
	}
}

func ensureLineBreak(b *byteBuffer) {
	if b.Bytes == nil {
		b.Bytes = make([]byte, 0)
	}
	if len(b.Bytes) == 0 || b.Bytes[len(b.Bytes)-1] != '\n' {
		b.Bytes = append(b.Bytes, '\n')
	}
}

// printWithoutColor prints the log message without color.
func printWithoutColor(args *log.FormatterArgs,
	b *byteBuffer,
	level string,
	quoteString, endWithMessage bool) {
	// header
	fmt.Fprintf(b, "%s %s ", args.Time, level)
	if args.Caller != "" {
		fmt.Fprintf(b, "%s %s 💦", args.Goid, args.Caller)
	} else {
		fmt.Fprint(b, "💦")
	}
	if !endWithMessage {
		fmt.Fprintf(b, " %s", args.Message)
	}
	// key and values
	for _, kv := range args.KeyValues {
		if quoteString && kv.ValueType == 's' {
			b.Bytes = append(b.Bytes, ' ')
			b.Bytes = append(b.Bytes, kv.Key...)
			b.Bytes = append(b.Bytes, '=')
			b.Bytes = strconv.AppendQuote(b.Bytes, kv.Value)
		} else {
			fmt.Fprintf(b, " %s=%s", kv.Key, kv.Value)
		}
	}
	// message
	if endWithMessage {
		fmt.Fprintf(b, " %s", args.Message)
	}
}

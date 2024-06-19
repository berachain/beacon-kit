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
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/phuslu/log"
)

// byteBuffer is a byte buffer
type byteBuffer struct {
	Bytes []byte
}

// Write writes to the byte buffer
func (b *byteBuffer) Write(bytes []byte) (int, error) {
	b.Bytes = append(b.Bytes, bytes...)
	return len(bytes), nil
}

// byteBufferPool is a pool of byte buffers
var byteBufferPool = sync.Pool{
	New: func() any {
		return new(byteBuffer)
	},
}

// customFormatter is a custom Formatter to pass into the ConsoleWriter
func customFormatter(out io.Writer, args *log.FormatterArgs) (
	n int, err error) {
	b := byteBufferPool.Get().(*byteBuffer)
	b.Bytes = b.Bytes[:0]
	defer byteBufferPool.Put(b)

	const (
		Reset   = "\x1b[0m"
		Black   = "\x1b[30m"
		Red     = "\x1b[31m"
		Green   = "\x1b[32m"
		Yellow  = "\x1b[33m"
		Blue    = "\x1b[34m"
		Magenta = "\x1b[35m"
		Cyan    = "\x1b[36m"
		White   = "\x1b[37m"
		Gray    = "\x1b[90m"
	)

	// TODO: pull out of config
	var color, level string
	switch args.Level {
	case "trace":
		color, level = Magenta, "TRACE"
	case "debug":
		color, level = Yellow, "DEBUG"
	case "info":
		color, level = Green, " BET"
	case "warn":
		color, level = Yellow, "WARN"
	case "error":
		color, level = Red, "ERROR"
	case "fatal":
		color, level = Red, " FTL"
	case "panic":
		color, level = Red, " PNC"
	default:
		color, level = Gray, " ???"
	}

	// TODO: pull out of config
	colorOutput := true
	quoteString := true
	endWithMessage := false

	// pretty console writer
	if colorOutput {
		printWithColor(args, b, color, level, quoteString, endWithMessage)
	} else {
		printWithoutColor(args, b, level, quoteString, endWithMessage)
	}

	// add line break if needed
	if b.Bytes[len(b.Bytes)-1] != '\n' {
		b.Bytes = append(b.Bytes, '\n')
	}

	// stack
	if args.Stack != "" {
		b.Bytes = append(b.Bytes, args.Stack...)
		if args.Stack[len(args.Stack)-1] != '\n' {
			b.Bytes = append(b.Bytes, '\n')
		}
	}

	return out.Write(b.Bytes)
}

// printWithColor prints the log message with color
func printWithColor(args *log.FormatterArgs,
	b *byteBuffer,
	color, level string,
	quoteString, endWithMessage bool) {
	const (
		Reset   = "\x1b[0m"
		Black   = "\x1b[30m"
		Red     = "\x1b[31m"
		Green   = "\x1b[32m"
		Yellow  = "\x1b[33m"
		Blue    = "\x1b[34m"
		Magenta = "\x1b[35m"
		Cyan    = "\x1b[36m"
		White   = "\x1b[37m"
		Gray    = "\x1b[90m"
	)
	// header
	fmt.Fprintf(b, "%s%s%s %s%s%s ", Gray, args.Time, Reset, color, level,
		Reset)
	if args.Caller != "" {
		fmt.Fprintf(b, "%s %s %s💦%s", args.Goid, args.Caller, Cyan, Reset)
	} else {
		fmt.Fprintf(b, "%s💦%s", Cyan, Reset)
	}
	if !endWithMessage {
		fmt.Fprintf(b, " %s", args.Message)
	}
	// key and values
	for _, kv := range args.KeyValues {
		if quoteString && kv.ValueType == 's' {
			kv.Value = strconv.Quote(kv.Value)
		}
		if kv.Key == "error" {
			fmt.Fprintf(b, " %s%s=%s%s", Red, kv.Key, kv.Value, Reset)
		} else {
			fmt.Fprintf(b, " %s%s=%s%s%s", Cyan, kv.Key, Gray, kv.Value, Reset)
		}
	}
	// message
	if endWithMessage {
		fmt.Fprintf(b, "%s %s", Reset, args.Message)
	}
}

// printWithoutColor prints the log message without color
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

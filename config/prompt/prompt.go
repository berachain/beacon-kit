// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// DefaultPrompt prompts the user for any text and performs no validation.
// If nothing is entered it returns the default.
func DefaultPrompt(
	cmd *cobra.Command, promptText, defaultValue string,
) (string, error) {
	var (
		response string
		au       = aurora.NewAurora(true)
		scanner  = bufio.NewScanner(os.Stdin)
	)

	if defaultValue != "" {
		cmd.Print(fmt.Sprintf(
			"%s (%s: %s):\n", promptText, au.BrightGreen("default"), defaultValue))
	} else {
		cmd.Print(fmt.Sprintf("%s:\n", promptText))
	}

	if ok := scanner.Scan(); ok {
		item := scanner.Text()
		response = strings.TrimRight(item, "\r\n")
		if response == "" {
			return defaultValue, nil
		}
		return response, nil
	}
	return "", errors.New("could not scan text input")
}

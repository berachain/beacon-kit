// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

type Reader interface {
	Read([]byte) (int, error)
}

type Prompt struct {
	Cmd *cobra.Command

	Text       string
	Default    string
	ValidateFn func(string) error
}

// Ask prompts the user and stores their response.
func (p *Prompt) Ask() (string, error) {
	au := aurora.NewAurora(true)
	prompt := au.Sprintf("%s:\n", p.Text)
	if p.Default != "" {
		prompt = au.Sprintf(
			"%s (%s: %s):\n", p.Text,
			au.BrightGreen("default"),
			p.Default,
		)
	}

	input := p.Default
	inputReader := p.Cmd.InOrStdin()
	scanner := bufio.NewScanner(inputReader)
	p.Cmd.Print(prompt)
	if scanner.Scan() {
		if text := scanner.Text(); text != "" {
			input = text
		}
	}

	return input, scanner.Err()
}

// AskAndValidate prompts the user and validates their response.
// Equivalent to Ask() if no validate function is specified.
func (p *Prompt) AskAndValidate() (string, error) {
	input, err := p.Ask()
	if err != nil {
		return "", err
	}

	// If validate function is specified, validate the input.
	if p.ValidateFn != nil {
		if err = p.ValidateFn(input); err != nil {
			return input, err
		}
	}
	return input, nil
}

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

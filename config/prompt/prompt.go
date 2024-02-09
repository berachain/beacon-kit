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
	"errors"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
)

// DefaultPrompt prompts the user and validates their response.
// Returns only when the user has provided a valid response.
func DefaultPrompt(
	cmd *cobra.Command, promptText, defaultValue string,
) error {
	au := aurora.NewAurora(true)
	if defaultValue != "" {
		promptText = au.Sprintf(
			"%s (%s: %s):\n", promptText,
			au.BrightGreen("default"),
			defaultValue,
		)
	} else {
		promptText = au.Sprintf("%s:\n", promptText)
	}

	var input string
	return survey.AskOne(
		&survey.Input{
			Message: promptText,
			Default: defaultValue,
		},
		&input,
		survey.WithValidator(func(val interface{}) error {
			input := val.(string)
			if !strings.EqualFold(input, "accept") {
				return errors.New("you have to accept Terms and Conditions in order to continue")
			}
			return nil
		}),
		survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = ""
		}),
	)
}

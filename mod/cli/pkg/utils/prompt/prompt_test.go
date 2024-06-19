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

package prompt_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/prompt"
	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/prompt/mocks"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	p = &prompt.Prompt{
		Cmd:  &cobra.Command{},
		Text: "test",
	}
)

func setupPromptWithInput(input string) {
	inputBuffer := bytes.NewReader([]byte(input + "\n"))
	p.Cmd.SetIn(inputBuffer)
}

func setupPromptWithErrorReader() {
	errReader := &mocks.Reader{}
	errReader.On("Read", mock.Anything).Return(0, errors.New("error"))
	p.Cmd.SetIn(errReader)
	p.Cmd.SetOut(os.NewFile(0, os.DevNull))
}

func TestAsk(t *testing.T) {
	tests := []struct {
		name        string
		setup       func()
		expected    string
		expectError bool
	}{
		{
			name: "Failure Case",
			setup: func() {
				setupPromptWithErrorReader()
			},
			expectError: true,
		},
		{
			name: "Success Case",
			setup: func() {
				setupPromptWithInput("response")
			},
			expected:    "response",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			input, err := p.Ask()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, input)
			}
		})
	}
}

func TestAskWithDefault(t *testing.T) {
	tests := []struct {
		name        string
		setup       func()
		expected    string
		expectError bool
	}{
		{
			name: "No Input",
			setup: func() {
				p.Default = "default"
			},
			expected:    "default",
			expectError: false,
		},
		{
			name: "With Input",
			setup: func() {
				p.Default = "default"
				setupPromptWithInput("response")
			},
			expected:    "response",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			input, err := p.Ask()
			require.NoError(t, err)
			require.Equal(t, tt.expected, input)
		})
	}
}

func TestAskAndValidate(t *testing.T) {
	failedValidationErr := errors.New("wrong response")
	p.ValidateFn = func(input string) error {
		if input == "response" {
			return nil
		}
		return failedValidationErr
	}

	tests := []struct {
		name        string
		setup       func()
		expectError bool
		expectedErr error
	}{
		{
			name: "Failed Validation",
			setup: func() {
				p.Default = "default"
			},
			expectError: true,
			expectedErr: failedValidationErr,
		},
		{
			name: "Success",
			setup: func() {
				setupPromptWithInput("response")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := p.AskAndValidate()
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

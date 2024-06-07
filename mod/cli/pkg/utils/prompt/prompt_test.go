// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	//nolint:gochecknoglobals // this is a test
	p = &prompt.Prompt{
		Cmd:  &cobra.Command{},
		Text: "test",
	}
)

func TestAsk(t *testing.T) {
	// Failure Case
	errReader := &mocks.Reader{}
	errReader.On("Read", mock.Anything).Return(0, errors.New("error"))

	p.Cmd.SetOut(os.NewFile(0, os.DevNull))
	p.Cmd.SetIn(errReader)
	_, err := p.Ask()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Success Case
	inputBuffer := bytes.NewReader([]byte("response\n"))
	p.Cmd.SetIn(inputBuffer)
	input, err := p.Ask()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if input != "response" {
		t.Errorf("Expected response, got %s", input)
	}
}

func TestAskWithDefault(t *testing.T) {
	p.Default = "default"
	input, err := p.Ask()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if input != "default" {
		t.Errorf("Expected default, got %s", input)
	}

	inputBuffer := bytes.NewReader([]byte("response\n"))
	p.Cmd.SetIn(inputBuffer)
	input, err = p.Ask()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if input != "response" {
		t.Errorf("Expected response, got %s", input)
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

	p.Default = "default"
	_, err := p.AskAndValidate()
	require.Error(t, err)
	require.Equal(t, failedValidationErr, err)

	inputBuffer := bytes.NewReader([]byte("response\n"))
	p.Cmd.SetIn(inputBuffer)
	_, err = p.AskAndValidate()
	require.NoError(t, err)
}

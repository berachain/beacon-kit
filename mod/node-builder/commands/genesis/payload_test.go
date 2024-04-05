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

package genesis

import (
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/stretchr/testify/require"
)

func TestAddExecutionPayloadCmd(t *testing.T) {
	// // Failure Case
	// errReader := &mocks.Reader{}
	// errReader.On("Read", mock.Anything).Return(0, errors.New("error"))

	// p.Cmd.SetOut(os.NewFile(0, os.DevNull))
	// p.Cmd.SetIn(errReader)
	// _, err := p.Ask()
	// if err == nil {
	// 	t.Errorf("Expected error, got nil")
	// }

	// // Success Case
	// inputBuffer := bytes.NewReader([]byte("response\n"))
	// p.Cmd.SetIn(inputBuffer)
	// input, err := p.Ask()
	// if err != nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }
	// if input != "response" {
	// 	t.Errorf("Expected response, got %s", input)
	// }
}

func TestExecutableDataToExecutionPayload(t *testing.T) {
	data := &ethtypes.ExecutableData{}
	payload := executableDataToExecutionPayload(data)
	require.NotNil(t, payload)
}

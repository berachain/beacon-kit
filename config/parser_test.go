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

package config_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/config/mocks"
)

func TestParser(t *testing.T) {
	var parser *config.AppOptionsParser
	var appOpts = new(mocks.AppOptions)

	parser = config.NewAppOptionsParser(appOpts)

	t.Run("should set and retrieve a string option", func(t *testing.T) {
		value := "testValue"
		runTest(t, appOpts, parser.GetString, value)
	})

	t.Run("should set and retrieve an integer option", func(t *testing.T) {
		value := int(42)
		runTest(t, appOpts, parser.GetInt, value)
	})

	t.Run("should handle an int64 option", func(t *testing.T) {
		value := int64(42)
		runTest(t, appOpts, parser.GetInt64, value)
	})

	t.Run("should set and retrieve a uint64 option", func(t *testing.T) {
		value := uint64(42)
		runTest(t, appOpts, parser.GetUint64, value)
	})

	t.Run("should set and retrieve a pointer to a uint64 option", func(t *testing.T) {
		value := uint64(42)
		runTestWithOutput(t, appOpts, parser.GetUint64Ptr, "42", &value)
	})

	t.Run("should set and retrieve a big.Int option", func(t *testing.T) {
		value := new(big.Int).SetInt64(42)
		runTestWithOutput(t, appOpts, parser.GetBigInt, "42", value)
	})

	t.Run("should set and retrieve a float64 option", func(t *testing.T) {
		value := 3.14159
		runTest(t, appOpts, parser.GetFloat64, value)
	})

	t.Run("should set and retrieve a boolean option", func(t *testing.T) {
		value := true
		runTest(t, appOpts, parser.GetBool, value)
	})

	t.Run("should set and retrieve a string slice option", func(t *testing.T) {
		value := []string{"apple", "banana", "cherry"}
		runTest(t, appOpts, parser.GetStringSlice, value)
	})

	t.Run("should set and retrieve a time.Duration option", func(t *testing.T) {
		value := time.Second * 10
		runTest(t, appOpts, parser.GetTimeDuration, value)
	})

	t.Run("should set and retrieve a common.Address option", func(t *testing.T) {
		addressStr := "0x18df82c7e422a42d47345ed86b0e935e9718ebda"
		runTestWithOutput(
			t, appOpts, parser.GetCommonAddress, addressStr, common.HexToAddress(addressStr))
	})

	t.Run("should set and retrieve a list of common.Address options", func(t *testing.T) {
		addressStrs := []string{
			"0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4",
			"0x18df82c7e422a42d47345ed86b0e935e9718ebda",
		}
		expectedAddresses := []common.Address{
			common.HexToAddress(addressStrs[0]),
			common.HexToAddress(addressStrs[1]),
		}
		runTestWithOutput(
			t, appOpts, parser.GetCommonAddressList, addressStrs, expectedAddresses)
	})

	t.Run("should set and retrieve a hexutil.Bytes option", func(t *testing.T) {
		bytesStr := "0x1234567890abcdef"
		expectedBytes := hexutil.MustDecode(bytesStr)
		runTest(t, appOpts, parser.GetHexutilBytes, expectedBytes)
	})
}

func runTest[A any](
	t *testing.T, appOpts *mocks.AppOptions, parser func(string) (A, error), value A,
) {
	runTestWithOutput(t, appOpts, parser, value, value)
}

func runTestWithOutput[A, B any](
	t *testing.T, appOpts *mocks.AppOptions, parser func(string) (B, error), value A, output B,
) {
	appOpts.On("Get", mock.Anything).Return(value).Once()

	retrievedValue, err := parser("myTestKey")

	require.NoError(t, err)
	require.Equal(t, output, retrievedValue)
}

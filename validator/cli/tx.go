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

package cli

import (
	"fmt"

	beacontypes "github.com/berachain/beacon-kit/runtime/modules/beacon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func ProvideMessageValidator() types.MessageValidator {
	return DefaultMessageValidator
}

func DefaultMessageValidator(msgs []sdk.Msg) error {
	if len(msgs) != 1 {
		return fmt.Errorf(
			"unexpected number of GenTx messages; got: %d, expected: 1",
			len(msgs),
		)
	}
	if _, ok := msgs[0].(*beacontypes.MsgCreateValidatorX); !ok {
		return fmt.Errorf(
			"unexpected GenTx message type; expected: MsgCreateValidator, got: %T",
			msgs[0],
		)
	}

	if m, ok := msgs[0].(sdk.HasValidateBasic); ok {
		if err := m.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid GenTx '%s': %w", msgs[0], err)
		}
	}

	return nil
}

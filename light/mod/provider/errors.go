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

package provider

import (
	abci "github.com/cometbft/cometbft/abci/types"
	legacyerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// abciErrorToSdkError converts an ABCI response not OK to a more
// informative SDK error.
func abciErrorToSdkError(resp abci.ResponseQuery) error {
	switch resp.Code {
	case legacyerrors.ErrInvalidRequest.ABCICode():
		return legacyerrors.ErrInvalidRequest.Wrap(resp.Log)
	case legacyerrors.ErrUnauthorized.ABCICode():
		return legacyerrors.ErrUnauthorized.Wrap(resp.Log)
	case legacyerrors.ErrKeyNotFound.ABCICode():
		return legacyerrors.ErrKeyNotFound.Wrap(resp.Log)
	default:
		// TODO: maybe change up this error message, as it's not very
		// informative and a default error case does not necessarily imply it's
		// a logic error.
		return legacyerrors.ErrLogic.Wrap(resp.Log)
	}
}

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

package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-builder/config/flags"
	"github.com/berachain/beacon-kit/mod/node-builder/utils/jwt"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// TrustedSetupInput is the input for the dep inject framework.
type JWTSecretInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

// TrustedSetupOutput is the output for the dep inject framework.
type JWTSecretOutput struct {
	depinject.Out
	JWTSecret *jwt.Secret
}

// ProvideJWTSecret is a function that provides the module to the application.
func ProvideJWTSecret(in JWTSecretInput) JWTSecretOutput {
	jwtSecret, err := jwt.LoadFromFile(
		cast.ToString(in.AppOpts.Get(flags.JWTSecretPath)))
	if err != nil {
		panic(err)
	}

	return JWTSecretOutput{
		JWTSecret: jwtSecret,
	}
}

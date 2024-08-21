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

package components

import (
	"strings"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/cli/pkg/flags"
	servertypes "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/server/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

// JWTSecretInput is the input for the dep inject framework.
type JWTSecretInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

// ProvideJWTSecret is a function that provides the module to the application.
func ProvideJWTSecret(in JWTSecretInput) (*jwt.Secret, error) {
	return LoadJWTFromFile(cast.ToString(in.AppOpts.Get(flags.JWTSecretPath)))
}

// LoadJWTFromFile reads the JWT secret from a file and returns it.
func LoadJWTFromFile(filepath string) (*jwt.Secret, error) {
	data, err := afero.ReadFile(afero.NewOsFs(), filepath)
	if err != nil {
		// Return an error if the file cannot be read.
		return nil, err
	}
	return jwt.NewFromHex(strings.TrimSpace(string(data)))
}

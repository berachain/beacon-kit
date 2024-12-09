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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/cli/flags"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

// TrustedSetupInput is the input for the dep inject framework.
type TrustedSetupInput struct {
	depinject.In
	AppOpts config.AppOptions
}

// ProvideTrustedSetup provides the module to the application.
func ProvideTrustedSetup(
	in TrustedSetupInput,
) (*gokzg4844.JSONTrustedSetup, error) {
	return ReadTrustedSetup(
		cast.ToString(in.AppOpts.Get(flags.KZGTrustedSetupPath)),
	)
}

// ReadTrustedSetup reads the trusted setup from the file system.
func ReadTrustedSetup(filePath string) (*gokzg4844.JSONTrustedSetup, error) {
	config, err := afero.ReadFile(afero.NewOsFs(), filePath)
	if err != nil {
		return nil, err
	}
	params := new(gokzg4844.JSONTrustedSetup)
	if err = json.Unmarshal(config, params); err != nil {
		return nil, err
	}
	if err = gokzg4844.CheckTrustedSetupIsWellFormed(params); err != nil {
		return nil, err
	}
	return params, nil
}

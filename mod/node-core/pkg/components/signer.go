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
	"path/filepath"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	clientFlags "github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// BlsSignerInput is the input for the dep inject framework.
type BlsSignerInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
	PrivKey LegacyKey `optional:"true"`
}

// ProvideBlsSigner is a function that provides the module to the application.
func ProvideBlsSigner(in BlsSignerInput) (crypto.BLSSigner, error) {
	if in.PrivKey == [constants.BLSSecretKeyLength]byte{} {
		// if no private key is provided, use privval signer
		homeDir := cast.ToString(in.AppOpts.Get(clientFlags.FlagHome))
		privValKeyFile := cast.ToString(
			in.AppOpts.Get("priv_validator_key_file"),
		)
		privValStateFile := cast.ToString(
			in.AppOpts.Get("priv_validator_state_file"),
		)
		// If privValKeyFile is not an absolute path, join with homeDir
		if !filepath.IsAbs(privValKeyFile) {
			privValKeyFile = filepath.Join(homeDir, privValKeyFile)
		}
		// If privValStateFile is not an absolute path, join with homeDir
		if !filepath.IsAbs(privValStateFile) {
			privValStateFile = filepath.Join(homeDir, privValStateFile)
		}
		return signer.NewBLSSigner(privValKeyFile, privValStateFile), nil
	}
	return signer.NewLegacySigner(in.PrivKey)
}

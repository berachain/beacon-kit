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

package components

import (
	"cosmossdk.io/depinject"
	blobkzg "github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/kzg"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/cast"
)

// TrustedSetupInput is the input for the dep inject framework.
type TrustedSetupInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

// ProvideTrustedSetup provides the trusted setup to the depinject framework.
func ProvideTrustedSetup(
	in TrustedSetupInput,
) (*gokzg4844.JSONTrustedSetup, error) {
	return kzg.ReadTrustedSetup(
		cast.ToString(in.AppOpts.Get(flags.KZGTrustedSetupPath)),
	)
}

type BlobProofVerifierInput struct {
	depinject.In
	// Cfg is the BeaconKit configuration.
	Cfg *config.Config
	// KZGTrustedSetup is the trusted setup.
	KZGTrustedSetup *gokzg4844.JSONTrustedSetup
}

// ProvideBlobProofVerifier provides the blob proof verifier to the depinject
// framework.
func ProvideBlobProofVerifier(
	in BlobProofVerifierInput,
) blobkzg.BlobProofVerifier {
	// #nosec:G703 // F
	verifier, _ := blobkzg.NewBlobProofVerifier(
		in.Cfg.KZG.Implementation, in.KZGTrustedSetup,
	)
	return verifier
}

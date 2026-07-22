// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"fmt"
	"path/filepath"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config"
	dablob "github.com/berachain/beacon-kit/da/blob"
	"github.com/berachain/beacon-kit/da/blobreactor"
	"github.com/berachain/beacon-kit/da/kzg"
	dastore "github.com/berachain/beacon-kit/da/store"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
)

// BlobReactorInput is the input for the blob reactor provider.
type BlobReactorInput struct {
	depinject.In

	Config            *config.Config
	ChainSpec         chain.Spec
	Logger            *phuslu.Logger
	AvailabilityStore *dastore.Store
	BlobProcessor     BlobProcessor
	TelemetrySink     *metrics.TelemetrySink
}

// ProvideBlobReactor provides the blob reactor, which distributes blob sidecars over a dedicated CometBFT p2p channel.
func ProvideBlobReactor(in BlobReactorInput) *blobreactor.BlobReactor {
	cfg := in.Config.BlobReactor
	defaults := blobreactor.DefaultConfig()
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = defaults.RequestTimeout
	}
	if cfg.FetchTimeout <= 0 {
		cfg.FetchTimeout = defaults.FetchTimeout
	}

	return blobreactor.NewBlobReactor(
		in.AvailabilityStore,
		in.BlobProcessor,
		in.Logger.With("service", "blob-reactor"),
		cfg,
		in.ChainSpec.MaxBlobsPerBlock(),
		in.TelemetrySink,
	)
}

// BlobReconstructorInput is the input for the blob reconstructor provider.
type BlobReconstructorInput struct {
	depinject.In

	EngineClient      *client.EngineClient
	SidecarFactory    *dablob.SidecarFactory
	BlobProofVerifier kzg.BlobProofVerifier
	Logger            *phuslu.Logger
}

// ProvideBlobReconstructor provides the reconstructor that rebuilds blob sidecars from blobs fetched off the local execution client.
func ProvideBlobReconstructor(in BlobReconstructorInput) (*dablob.Reconstructor, error) {
	prover, ok := in.BlobProofVerifier.(kzg.BlobProofProver)
	if !ok {
		return nil, fmt.Errorf("kzg implementation %s cannot compute blob proofs", in.BlobProofVerifier.GetImplementation())
	}
	return dablob.NewReconstructor(in.EngineClient, in.SidecarFactory, prover, in.Logger.With("service", "blob-reconstructor")), nil
}

// BlobFetcherInput is the input for the blob fetcher provider.
type BlobFetcherInput struct {
	depinject.In

	AppOpts        config.AppOptions
	ChainSpec      chain.Spec
	Logger         *phuslu.Logger
	BlobProcessor  BlobProcessor
	BlobReactor    *blobreactor.BlobReactor
	StorageBackend *storage.Backend
	TelemetrySink  *metrics.TelemetrySink
}

// ProvideBlobFetcher provides the background blob fetcher.
func ProvideBlobFetcher(in BlobFetcherInput) (blockchain.BlobFetcher, error) {
	return blockchain.NewBlobFetcher(
		filepath.Join(cast.ToString(in.AppOpts.Get(flags.FlagHome)), "data"),
		in.Logger.With("service", "blob-fetcher"),
		in.BlobProcessor,
		in.BlobReactor,
		in.StorageBackend,
		in.ChainSpec,
		blockchain.DefaultBlobFetcherConfig(),
		in.TelemetrySink,
	)
}

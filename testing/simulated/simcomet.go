//go:build simulated

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

package simulated

import (
	"context"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/validator"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/builder"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/types"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtcrypto "github.com/cometbft/cometbft/crypto"
	pvm "github.com/cometbft/cometbft/privval"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.ConsensusService = (*SimComet)(nil)

// SimComet is normal Comet under the hood, but we override the Start method to avoid starting the actual
// CometBFT core loop so that we can orchestrate it ourselves.
type SimComet struct {
	// We are forced to stutter here as we want to override the implementations of the original comet service.
	Comet *cometbft.Service
	// Used to initialize the node address.
	cmtCfg *cmtcfg.Config
}

func ProvideSimComet(
	logger *phuslu.Logger,
	blockchain blockchain.BlockchainI,
	blockBuilder validator.BlockBuilderI,
	db dbm.DB,
	cs chain.Spec,
	cmtCfg *cmtcfg.Config,
	appOpts config.AppOptions,
	telemetrySink *metrics.TelemetrySink) *SimComet {
	return &SimComet{
		Comet: cometbft.NewService(
			logger,
			db,
			blockchain,
			blockBuilder,
			cs,
			cmtCfg,
			telemetrySink,
			builder.DefaultServiceOptions(appOpts)...,
		),
		cmtCfg: cmtCfg,
	}
}

// Start sets the ctx and the node address for the SimComet service.
func (s *SimComet) Start(ctx context.Context) error {
	s.Comet.ResetAppCtx(ctx)
	return nil
}

// GetNodeAddress returns the node address for the SimComet service.
func (s *SimComet) GetNodeAddress() (cmtcrypto.Address, error) {
	privVal, err := pvm.LoadOrGenFilePV(
		s.cmtCfg.PrivValidatorKeyFile(),
		s.cmtCfg.PrivValidatorStateFile(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	pubKey, err := privVal.GetPubKey()
	if err != nil {
		return nil, err
	}
	return pubKey.Address(), nil
}

func (s *SimComet) Stop() error {
	return nil
}

func (s *SimComet) Name() string {
	return s.Comet.Name()
}

func (s *SimComet) IsAppReady() error {
	return s.Comet.IsAppReady()
}

func (s *SimComet) CreateQueryContext(height int64, prove bool) (sdk.Context, error) {
	return s.Comet.CreateQueryContext(height, prove)
}

func (s *SimComet) GetSyncData() (int64, int64) {
	panic("unimplemented")
}

func (s *SimComet) GetBlock(height int64) *cmttypes.Block {
	return s.Comet.GetBlock(height)
}

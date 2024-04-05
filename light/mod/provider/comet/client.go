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

package comet

import (
	"context"
	"strings"

	"github.com/berachain/beacon-kit/light/mod/provider/comet/types"
	"github.com/cockroachdb/errors"
	cometMath "github.com/cometbft/cometbft/libs/math"
	"github.com/cometbft/cometbft/light"
)

// NewClient returns a new comet light client with the given configuration.
func NewClient(cfg *Config) (*light.Client, error) {
	// Create a new local database
	db, err := types.NewDB(cfg.Directory, cfg.ChainID)
	if err != nil {
		return nil, err
	}

	// Check if the primary address is set, if not, check for existing providers
	if cfg.PrimaryAddr == "" {
		cfg.PrimaryAddr, cfg.WitnessesAddrs, err = db.CheckForExistingProviders()
		if err != nil {
			return nil, err
		}
		if cfg.PrimaryAddr == "" {
			return nil, errors.New(types.NoPrimaryAddress)
		}
	} else {
		err = db.SaveProviders(cfg.PrimaryAddr, strings.Join(cfg.WitnessesAddrs, ","))
		if err != nil {
			cfg.Logger.Error(types.FailedSaveProviders, err)
		}
	}

	// Set options for the light client
	verification := light.SequentialVerification()
	if !cfg.Sequential {
		// Parse the trust level from the input to the fraction
		var tl cometMath.Fraction
		tl, err = cometMath.ParseFraction(cfg.TrustLevel)
		if err != nil {
			return nil, err
		}

		verification = light.SkippingVerification(tl)
	}
	options := []light.Option{
		light.Logger(cfg.Logger),
		light.ConfirmationFunction(cfg.ConfirmationFunc),
		verification,
	}

	var client *light.Client
	if cfg.TrustedHeight > 0 && len(cfg.TrustedHash) > 0 { // fresh installation
		client, err = light.NewHTTPClient(
			context.Background(),
			cfg.ChainID,
			light.TrustOptions{
				Period: cfg.TrustingPeriod,
				Height: cfg.TrustedHeight,
				Hash:   cfg.TrustedHash,
			},
			cfg.PrimaryAddr,
			cfg.WitnessesAddrs,
			db.Store,
			options...,
		)
	} else { // continue from latest state
		client, err = light.NewHTTPClientFromTrustedStore(
			cfg.ChainID,
			cfg.TrustingPeriod,
			cfg.PrimaryAddr,
			cfg.WitnessesAddrs,
			db.Store,
			options...,
		)
	}
	if err != nil {
		return nil, errors.Wrap(err, types.FailedToCreateClient)
	}

	return client, nil
}

// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package light

import (
	"context"
	"strings"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	cometMath "github.com/cometbft/cometbft/libs/math"
	"github.com/cometbft/cometbft/light"
	"github.com/itsdevbear/bolaris/light/types"
)

func NewClient(
	logger log.Logger,
	chainID string,
	trustingPeriod time.Duration,
	trustedHeight int64,
	trustedHash []byte,
	tl string,
	sequential bool,
	primaryAddr string,
	witnesses []string,
	dir string,
	confirmationFunc func(string) bool,
) (*light.Client, error) {
	db, err := types.NewDB(dir, chainID)
	if err != nil {
		return nil, err
	}

	// primary address is the address of the primary node
	var (
		primaryAddress string
		witnessesAddrs []string
	)

	witnessesAddrs = witnesses
	// If no primary address is specified, check the db for existing providers
	if primaryAddr == "" {
		primaryAddress, witnessesAddrs, err = db.CheckForExistingProviders()
		if err != nil {
			return nil, types.NewError(types.FailedCheckForExistingProviders, err)
		}
		if primaryAddress == "" {
			return nil, types.NewError(types.NoPrimaryAddress, nil)
		}
	} else {
		// store the primary and witness addresses in the db
		err = db.SaveProviders(primaryAddr, strings.Join(witnesses, ","))
		if err != nil {
			// Should this error out?
			logger.Error(types.FailedSaveProviders, err)
		}
	}

	// parse the trust level from the input to the fraction
	trustLevel, err := cometMath.ParseFraction(tl)
	if err != nil {
		return nil, err
	}

	// set the options for the light client
	options := []light.Option{
		light.Logger(logger),
		light.ConfirmationFunction(confirmationFunc),
	}

	if sequential {
		options = append(options, light.SequentialVerification())
	} else {
		options = append(options, light.SkippingVerification(trustLevel))
	}

	var c *light.Client
	if trustedHeight > 0 && len(trustedHash) > 0 { // fresh installation
		c, err = light.NewHTTPClient(
			context.Background(),
			chainID,
			light.TrustOptions{
				Period: trustingPeriod,
				Height: trustedHeight,
				Hash:   trustedHash,
			},
			primaryAddr,
			witnessesAddrs,
			db.Store,
			options...,
		)
	} else { // continue from latest state
		c, err = light.NewHTTPClientFromTrustedStore(
			chainID,
			trustingPeriod,
			primaryAddr,
			witnessesAddrs,
			db.Store,
			options...,
		)
	}
	if err != nil {
		return nil, types.NewError(types.FailedToCreateClient, err)
	}

	return c, nil
}

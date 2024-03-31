package comet

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	cometMath "github.com/cometbft/cometbft/libs/math"
	"github.com/cometbft/cometbft/light"

	"github.com/berachain/beacon-kit/light/provider/comet/types"
)

func NewClient(
	logger log.Logger,
	chainID string,
	trustingPeriod time.Duration,
	trustedHeight int64,
	trustedHash []byte,
	trustLevel string,
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

	if primaryAddr == "" {
		primaryAddr, witnesses, err = db.CheckForExistingProviders()
		if err != nil {
			return nil, err
		}
		if primaryAddr == "" {
			return nil, errors.New(types.NoPrimaryAddress)
		}
	} else {
		err = db.SaveProviders(primaryAddr, strings.Join(witnesses, ","))
		if err != nil {
			logger.Error(types.FailedSaveProviders, err) // Should this error out?
		}
	}

	// set the options for the light client
	verification := light.SequentialVerification()
	if !sequential {
		// parse the trust level from the input to the fraction
		tl, err := cometMath.ParseFraction(trustLevel)
		if err != nil {
			return nil, err
		}

		verification = light.SkippingVerification(tl)
	}
	options := []light.Option{
		light.Logger(logger),
		light.ConfirmationFunction(confirmationFunc),
		verification,
	}

	var client *light.Client
	if trustedHeight > 0 && len(trustedHash) > 0 { // fresh installation
		client, err = light.NewHTTPClient(
			context.Background(),
			chainID,
			light.TrustOptions{
				Period: trustingPeriod,
				Height: trustedHeight,
				Hash:   trustedHash,
			},
			primaryAddr,
			witnesses,
			db.Store,
			options...,
		)
	} else { // continue from latest state
		client, err = light.NewHTTPClientFromTrustedStore(
			chainID,
			trustingPeriod,
			primaryAddr,
			witnesses,
			db.Store,
			options...,
		)
	}
	if err != nil {
		return nil, types.NewError(types.FailedToCreateClient, err)
	}

	return client, nil
}

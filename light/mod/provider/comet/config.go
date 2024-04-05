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
	"time"

	"github.com/cometbft/cometbft/libs/log"
)

// Config defines the configuration for the Comet client.
type Config struct {
	Logger             log.Logger
	ChainID            string
	TrustingPeriod     time.Duration
	TrustedHeight      int64
	TrustedHash        []byte
	TrustLevel         string
	ListeningAddr      string
	Sequential         bool
	PrimaryAddr        string
	WitnessesAddrs     []string
	Directory          string
	MaxOpenConnections int
	ConfirmationFunc   func(string) bool
}

// NewConfig returns a new Comet configuration.
func NewConfig(
	logger log.Logger,
	chainID string,
	trustingPeriod time.Duration,
	trustedHeight int64,
	trustedHash []byte,
	trustLevel string,
	listeningAddr string,
	sequential bool,
	primaryAddr string,
	witnessesAddrs []string,
	directory string,
	maxOpenConnections int,
	confirmationFunc func(string) bool,
) *Config {
	return &Config{
		Logger:             logger,
		ChainID:            chainID,
		TrustingPeriod:     trustingPeriod,
		TrustedHeight:      trustedHeight,
		TrustedHash:        trustedHash,
		TrustLevel:         trustLevel,
		ListeningAddr:      listeningAddr,
		Sequential:         sequential,
		PrimaryAddr:        primaryAddr,
		WitnessesAddrs:     witnessesAddrs,
		Directory:          directory,
		MaxOpenConnections: maxOpenConnections,
		ConfirmationFunc:   confirmationFunc,
	}
}

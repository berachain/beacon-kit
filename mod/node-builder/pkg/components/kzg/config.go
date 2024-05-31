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

package kzg

import (
	"encoding/json"

	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
)

const (
	// defaultTrustedSetupPath is the default path to the trusted setup.
	defaultTrustedSetupPath = "./testing/files/kzg-trusted-setup.json"
	// defaultImplementation is the default KZG implementation to use.
	// Options are `crate-crypto/go-kzg-4844` or `ethereum/c-kzg-4844`.
	defaultImplementation = "crate-crypto/go-kzg-4844"
)

type Config struct {
	// TrustedSetupPath is the path to the trusted setup.
	TrustedSetupPath string `mapstructure:"trusted-setup-path"`
	// Implementation is the KZG implementation to use.
	Implementation string `mapstructure:"implementation"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		TrustedSetupPath: defaultTrustedSetupPath,
		Implementation:   defaultImplementation,
	}
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

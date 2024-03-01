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

package config

import (
	"time"

	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/io/cli/parser"
)

// Beacon conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Beacon] = Beacon{}

// Builder is the configuration for the payload builder.
type Builder struct {
	// LocalBuilderEnabled determines if the local builder is enabled.
	LocalBuilderEnabled bool
	// LocalBuildPayloadTimeout is the timeout parameter for local build
	// payload. This should match, or be slightly less than the configured
	// timeout on your
	// execution client. It also must be less than timeout_proposal in the
	// CometBFT configuration.
	LocalBuildPayloadTimeout time.Duration
}

// DefaultBuilderConfig returns the default fork configuration.
func DefaultBuilderConfig() Builder {
	return Builder{
		LocalBuilderEnabled:      true,
		LocalBuildPayloadTimeout: 2 * time.Second, //nolint:gomnd // default config.
	}
}

// Parse parses the configuration.
func (c Builder) Parse(parser parser.AppOptionsParser) (*Builder, error) {
	localBuilderEnabled, err := parser.GetBool(flags.LocalBuilderEnabled)
	if err != nil {
		return nil, err
	}

	payloadTimeout, err := parser.GetTimeDuration(
		flags.LocalBuildPayloadTimeout,
	)
	if err != nil {
		return nil, err
	}

	return &Builder{
		LocalBuilderEnabled:      localBuilderEnabled,
		LocalBuildPayloadTimeout: payloadTimeout,
	}, nil
}

// Template returns the configuration template.
func (c Builder) Template() string {
	//nolint:lll
	return `
[beacon-kit.builder]
# LocalBuilderEnabled determines if the local payload builder is enabled.
local-builder-enabled = {{ .BeaconKit.Builder.LocalBuilderEnabled }}

# The timeout for local build payload. This should match, or be slightly less
# than the configured timeout on your execution client. It also must be less than
# timeout_proposal in the CometBFT configuration.
local-build-payload-timeout = "{{ .BeaconKit.Builder.LocalBuildPayloadTimeout }}"
`
}

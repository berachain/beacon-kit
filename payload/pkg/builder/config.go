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

package builder

import (
	"time"

	"github.com/berachain/beacon-kit/primitives/pkg/common"
)

const (
	// defaultPayloadTimeout is the default value for local build
	// payload timeout.
	defaultPayloadTimeout = 1200 * time.Millisecond
)

// Config is the configuration for the payload builder.
//
//nolint:lll // struct tags.
type Config struct {
	// Enabled determines if the local builder is enabled.
	Enabled bool `mapstructure:"enabled"`
	// SuggestedFeeRecipient is the address that will receive the transaction
	// fees produced by any blocks from this node.
	SuggestedFeeRecipient common.ExecutionAddress `mapstructure:"suggested-fee-recipient"`
	// PayloadTimeout is the timeout parameter for local build
	// payload. This should match, or be slightly less than the configured
	// timeout on your execution client. It also must be less than
	// timeout_proposal in the CometBFT configuration.
	PayloadTimeout time.Duration `mapstructure:"payload-timeout"`
}

// DefaultConfig returns the default fork configuration.
func DefaultConfig() Config {
	return Config{
		Enabled:               true,
		SuggestedFeeRecipient: common.ExecutionAddress{},
		PayloadTimeout:        defaultPayloadTimeout,
	}
}

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

package validator

const (
	// defaultGraffiti is the default graffiti string.
	defaultGraffiti = ""

	// defaultEnableOptimisticPayloadBuilds is the default
	// for enabling the optimistic payload builder.
	defaultEnableOptimisticPayloadBuilds = true
)

// Config is the validator configuration.
//
//nolint:lll // struct tags.
type Config struct {
	// Graffiti is the string that will be included in the
	// graffiti field of the beacon block.
	Graffiti string `mapstructure:"graffiti"`

	// EnableOptimisticPayloadBuilds is the optimistic block builder.
	EnableOptimisticPayloadBuilds bool `mapstructure:"enable-optimistic-payload-builds"`
}

// DefaultConfig returns the default fork configuration.
func DefaultConfig() Config {
	return Config{
		Graffiti:                      defaultGraffiti,
		EnableOptimisticPayloadBuilds: defaultEnableOptimisticPayloadBuilds,
	}
}

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
	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/io/cli/parser"
)

type SyncMode string

const (
	SyncModeRegular    = "regular"
	SyncModeOptimistic = "optimistic"
	SyncModeLight      = "light"
)

// Sync conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Sync] = &Sync{}

// DefaultSyncConfig returns the default Sync configuration.
func DefaultSyncConfig() Sync {
	return Sync{
		Mode: SyncModeOptimistic,
	}
}

// Config represents the configuration struct for the Sync.
type Sync struct {
	// Mode is used to determine the sync mode.
	Mode SyncMode
}

// Parse parses the configuration.
func (c Sync) Parse(parser parser.AppOptionsParser) (*Sync, error) {
	mode, err := parser.GetString(
		flags.SyncMode,
	)
	if err != nil {
		return nil, err
	}
	c.Mode = SyncMode(mode)

	return &c, nil
}

// Template returns the configuration template.
func (c Sync) Template() string {
	return `
[beacon-kit.beacon-config.sync]
# Determines the sync mode. 
# Possibilities: ["optimistic", "regular", "light"]
mode = {{.BeaconKit.Beacon.Sync.Mode}}
`
}

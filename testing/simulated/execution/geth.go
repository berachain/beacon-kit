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

package execution

import (
	"fmt"

	"github.com/ory/dockertest/docker"
)

// NewGethNode creates a new execution node configured for Geth by using the default Geth command builder.
func NewGethNode(homeDir string, image docker.PullImageOptions) *ExecNode {
	return NewExecNode(homeDir, image, defaultGethCmdStrBuilder)
}

// ValidGethImage returns the default Docker image options for the Geth node.
func ValidGethImage() docker.PullImageOptions {
	return docker.PullImageOptions{
		Repository: "ghcr.io/berachain/bera-geth",
		Tag:        "latest",
	}
}

// defaultGethCmdStrBuilder returns a command string tailored for running a Geth node.
func defaultGethCmdStrBuilder(genesisFile string) string {
	return fmt.Sprintf(`
		geth init --datadir /tmp/gethdata /testdata/%s && 
		geth --http --http.addr 0.0.0.0 --http.api eth,net,web3,debug \
			 --authrpc.addr 0.0.0.0 \
			 --authrpc.jwtsecret /testing/files/jwt.hex \
			 --authrpc.vhosts '*' \
			 --datadir /tmp/gethdata \
			 --ipcpath /tmp/gethdata/geth.ipc \
			 --syncmode full \
			 --verbosity 4 \
			 --nodiscover
	`, genesisFile)
}

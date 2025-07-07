#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2025 Berachain Foundation
#
# Permission is hereby granted, free of charge, to any person
# obtaining a copy of this software and associated documentation
# files (the "Software"), to deal in the Software without
# restriction, including without limitation the rights to use,
# copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following
# conditions:
#
# The above copyright notice and this permission notice shall be
# included in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
# OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
# WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

# function to resolve absolute path from relative
resolve_path() {
	if [[ "$1" =~ : ]]; then
        # treat as an address or url, return as is
        echo "$1"
        return
	fi
    cd "$(dirname "$1")"
    local abs_path
    abs_path="$(pwd -P)/$(basename "$1")"
    echo "$abs_path"
}

# Check if the chain spec is provided as an argument.
CHAIN_SPEC=""
CHAIN_SPEC_ARG=""
if [ -z "$1" ]; then
    echo "No chain spec provided, falling back on devnet"
    CHAIN_SPEC="devnet"
    CHAIN_SPEC_ARG="--beacon-kit.chain-spec $CHAIN_SPEC"
else
	CHAIN_SPEC="$1"
    CHAIN_SPEC_ARG="--beacon-kit.chain-spec $CHAIN_SPEC"
	if [ "$CHAIN_SPEC" == "file" ]; then
		CHAIN_SPEC_FILE=$(resolve_path "$2")
		CHAIN_SPEC_ARG="$CHAIN_SPEC_ARG --beacon-kit.chain-spec-file $CHAIN_SPEC_FILE"
	fi
fi

CHAINID="beacond-2061"
MONIKER="localtestnet"
LOGLEVEL="info"
HOMEDIR="${HOMEDIR:-./.tmp/beacond}"
NON_INTERACTIVE="${NON_INTERACTIVE:-}"

# Path variables
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json
ETH_GENESIS=$(resolve_path "./testing/files/eth-genesis.json")
ETH_NETHER_GENESIS=$(resolve_path "./testing/files/eth-nether-genesis.json")
KZG_PATH=$(resolve_path "./testing/files/kzg-trusted-setup.json")

# used to exit on first error (any non-zero exit code)
set -e

if [ -z "$NON_INTERACTIVE" ]; then
  # Reinstall daemon
  make build

  overwrite="N"
  if [ -d $HOMEDIR ]; then
    printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" $HOMEDIR
    echo "Overwrite the existing configuration and start a new local node? [y/n]"
    read -r overwrite
  else
    overwrite="Y"
  fi
else
  # Executing within the e2e framework

  if [ -f "$HOMEDIR/emulate-latency.sh" ]; then
      "$HOMEDIR/emulate-latency.sh"
  fi

  # Forcibly remove any stray UNIX sockets left behind from previous runs
  rm -rf /var/run/privval.sock /var/run/app.sock || true

  # If there is no configuration provided, create one (for one-node testing only)
  if [ ! -f "$GENESIS" ]; then
    overwrite="Y"
  fi
fi

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
	rm -rf $HOMEDIR || true
	./build/bin/beacond init $MONIKER --chain-id $CHAINID --home $HOMEDIR $CHAIN_SPEC_ARG

	if [ "$CHAIN_SPEC" == "testnet" ]; then
	    network_dir="testing/networks/80069"
		cp -f $network_dir/*.toml $network_dir/genesis.json ${HOMEDIR}/config
    	KZG_PATH=$network_dir/kzg-trusted-setup.json
	elif [ "$CHAIN_SPEC" == "mainnet" ]; then
		network_dir="testing/networks/80094"
		cp -f $network_dir/*.toml $network_dir/genesis.json ${HOMEDIR}/config
    	KZG_PATH=$network_dir/kzg-trusted-setup.json
	else
		./build/bin/beacond genesis add-premined-deposit --home $HOMEDIR \
			32000000000 0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4 $CHAIN_SPEC_ARG
		./build/bin/beacond genesis collect-premined-deposits --home $HOMEDIR $CHAIN_SPEC_ARG
		./build/bin/beacond genesis set-deposit-storage "$ETH_GENESIS" --home $HOMEDIR $CHAIN_SPEC_ARG
		./build/bin/beacond genesis set-deposit-storage "$ETH_NETHER_GENESIS" --nethermind --home $HOMEDIR $CHAIN_SPEC_ARG
		./build/bin/beacond genesis execution-payload "$HOMEDIR/eth-genesis.json" --home $HOMEDIR $CHAIN_SPEC_ARG
	fi
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
BEACON_START_CMD="./build/bin/beacond start $CHAIN_SPEC_ARG --pruning=nothing "$TRACE" \
--beacon-kit.logger.log-level $LOGLEVEL --home $HOMEDIR \
--beacon-kit.engine.jwt-secret-path ${JWT_SECRET_PATH} \
--beacon-kit.kzg.trusted-setup-path ${KZG_PATH}  \
--beacon-kit.block-store-service.enabled \
--beacon-kit.node-api.enabled --beacon-kit.node-api.logging"

# Conditionally add the rpc-dial-url flag if RPC_DIAL_URL is not empty
if [ -n "$RPC_DIAL_URL" ]; then
	# this will overwrite the default dial url
	RPC_DIAL_URL=$(resolve_path "$RPC_DIAL_URL")
	echo "Overwriting the default dial url with $RPC_DIAL_URL"
	BEACON_START_CMD="$BEACON_START_CMD --beacon-kit.engine.rpc-dial-url ${RPC_PREFIX}${RPC_DIAL_URL}"
fi

eval $BEACON_START_CMD

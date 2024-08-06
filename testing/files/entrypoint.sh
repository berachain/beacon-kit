#!/bin/bash
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2024 Berachain Foundation
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
	fi
    cd "$(dirname "$1")"
    local abs_path
    abs_path="$(pwd -P)/$(basename "$1")"
    echo "$abs_path"
}

CHAINID="beacond-2062"
MONIKER="localtestnet"
LOGLEVEL="info"
CONSENSUS_KEY_ALGO="bls12_381"
HOMEDIR="./.tmp/beacond"

# Path variables
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json
ETH_GENESIS=$(resolve_path "./testing/files/eth-genesis.json")
PERSISTENT_PEERS="834419a275ea8584d9e0d62be054517fb28d6b82@10.0.2.106:26656,d1a215e21aec12057b2f195f472f065326ce331e@10.0.13.254:26656,3987bb7be4024fbc04cbf21b54935be3fe90ecc2@10.0.1.173:26656,8a6ab87c7c550f5fa84e3b641371523b352db264@10.0.15.167:26656"

sudo chmod 777 -R ./.tmp

# used to exit on first error (any non-zero exit code)
set -e

# Reinstall daemon
make build

overwrite="N"
if [ -d $HOMEDIR ]; then
  if [ $1 == "1" ]; then
    printf "\nAn existing folder at '%s' was found. skip overwrite\n" $HOMEDIR
  else
    printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" $HOMEDIR
    echo "Overwrite the existing configuration and start a new local node? [y/n]"
    read -r overwrite
  fi
else
  overwrite="Y"
fi

export CHAIN_SPEC="devnet"

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" || $3 == "onlyInit" ]]; then
	rm -rf $HOMEDIR

	if [[ $2 == "validator" && $3 != "locally" ]]; then
    cp -rf "./testing/files/beacond-validator-$3" $HOMEDIR/
  else
    ./build/bin/beacond init $MONIKER \
    --chain-id $CHAINID \
    --home $HOMEDIR \
    --consensus-key-algo $CONSENSUS_KEY_ALGO
  fi

  if [[ $3 == "onlyInit" ]]; then
    cp -rf ./testing/files/genesis.json $HOMEDIR/config/genesis.json
    exit 0
  fi

  if [[ $3 == "locally" ]]; then
    ./build/bin/beacond genesis add-premined-deposit --home $HOMEDIR
  fi

  ./build/bin/beacond genesis collect-premined-deposits --home $HOMEDIR
  ./build/bin/beacond genesis execution-payload "$ETH_GENESIS" --home $HOMEDIR

  if [[ $3 != "locally" ]]; then
    cp -rf ./testing/files/genesis.json $HOMEDIR/config/genesis.json
  fi
fi


# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
BEACON_START_CMD="./build/bin/beacond start --pruning=nothing "$TRACE" \
--beacon-kit.logger.log-level $LOGLEVEL --api.enabled-unsafe-cors \
--api.enable --api.swagger --minimum-gas-prices=0.0001abgt \
--home $HOMEDIR --beacon-kit.engine.jwt-secret-path ${JWT_SECRET_PATH} \
--beacon-kit.block-store-service.enabled --beacon-kit.block-store-service.pruner-enabled \
--beacon-kit.node-api.enabled --beacon-kit.node-api.logging" 

# Conditionally add the rpc-dial-url flag if RPC_DIAL_URL is not empty
if [ -n "$RPC_DIAL_URL" ]; then
	# this will overwrite the default dial url
	RPC_DIAL_URL=$(resolve_path "$RPC_DIAL_URL")
	echo "Overwriting the default dial url with $RPC_DIAL_URL"
	BEACON_START_CMD="$BEACON_START_CMD --beacon-kit.engine.rpc-dial-url ${RPC_PREFIX}${RPC_DIAL_URL}"
fi

eval $BEACON_START_CMD

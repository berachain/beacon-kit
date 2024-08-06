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

MONIKER="localtestnet"
CONSENSUS_KEY_ALGO="bls12_381"
HOMEDIR="./.tmp/beacond"

# if LOGLEVEL exists then use it, otherwise set default to info
if [ -n "$LOGLEVEL" ]; then
	LOGLEVEL=$LOGLEVEL
else
	LOGLEVEL="info"
fi

# Path variables
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json
ETH_GENESIS=$(resolve_path $ETH_GENESIS)

# used to exit on first error (any non-zero exit code)
set -e

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

echo "CHAINID: $CHAINID"
echo CHAIN_SPEC: $CHAIN_SPEC

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
	rm -rf $HOMEDIR
	./build/bin/beacond init $MONIKER \
		--chain-id $CHAINID \
		--home $HOMEDIR \
		--consensus-key-algo $CONSENSUS_KEY_ALGO
	./build/bin/beacond genesis add-premined-deposit --home $HOMEDIR
	./build/bin/beacond genesis collect-premined-deposits --home $HOMEDIR 
	./build/bin/beacond genesis execution-payload "$ETH_GENESIS" --home $HOMEDIR
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
